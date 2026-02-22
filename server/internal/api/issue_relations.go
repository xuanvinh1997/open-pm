package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
)

var validRelationTypes = map[string]bool{
	"relates_to": true, "blocks": true, "blocked_by": true, "duplicate_of": true,
}

var inverseRelationType = map[string]string{
	"blocks":       "blocked_by",
	"blocked_by":   "blocks",
	"relates_to":   "relates_to",
	"duplicate_of": "duplicate_of",
}

type AddIssueRelationRequest struct {
	RelatedIssueID uuid.UUID `json:"related_issue_id" validate:"required"`
	RelationType   string    `json:"relation_type" validate:"required"`
}

// ListIssueRelations handles GET .../issues/{issueID}/relations
func (a *API) ListIssueRelations(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	relations, err := a.queries.ListIssueRelations(r.Context(), issueID)
	if err != nil {
		return internalServerError("failed to list relations")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": relations})
}

// AddIssueRelation handles POST .../issues/{issueID}/relations
func (a *API) AddIssueRelation(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)

	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	var req AddIssueRelationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if !validRelationTypes[req.RelationType] {
		return badRequestError("invalid relation type: %s", req.RelationType)
	}
	if issueID == req.RelatedIssueID {
		return badRequestError("cannot relate an issue to itself")
	}

	rel, err := a.queries.AddIssueRelation(ctx, issueID, req.RelatedIssueID, req.RelationType, &userID)
	if err != nil {
		return internalServerError("failed to add relation")
	}

	// Create inverse relation
	inverse := inverseRelationType[req.RelationType]
	if _, err := a.queries.AddIssueRelation(ctx, req.RelatedIssueID, issueID, inverse, &userID); err != nil {
		log.Warn().Err(err).Msg("failed to create inverse relation")
	}

	// Log activity
	if _, err := a.queries.CreateIssueActivity(ctx, CreateActivityParams{
		IssueID:       &issueID,
		ProjectID:     projectID,
		WorkspaceID:   workspaceID,
		Verb:          "created",
		Field:         strPtr("relation"),
		NewValue:      strPtr(req.RelationType),
		NewIdentifier: &req.RelatedIssueID,
		Comment:       "added a relation",
		ActorID:       &userID,
	}); err != nil {
		log.Warn().Err(err).Msg("failed to log relation activity")
	}

	return sendJSON(w, http.StatusCreated, rel)
}

// RemoveIssueRelation handles DELETE .../issues/{issueID}/relations/{relationID}
func (a *API) RemoveIssueRelation(w http.ResponseWriter, r *http.Request) error {
	relationID, err := uuid.FromString(chi.URLParam(r, "relationID"))
	if err != nil {
		return badRequestError("invalid relation ID")
	}

	// Fetch the relation first to know its direction for inverse cleanup
	rel, err := a.queries.GetIssueRelationByID(r.Context(), relationID)
	if err != nil {
		return notFoundError("relation not found")
	}

	// Remove the forward relation
	if err := a.queries.RemoveIssueRelation(r.Context(), relationID); err != nil {
		return internalServerError("failed to remove relation")
	}

	// Remove the inverse relation
	inverse := inverseRelationType[rel.RelationType]
	if err := a.queries.RemoveInverseIssueRelation(r.Context(), rel.RelatedIssueID, rel.IssueID, inverse); err != nil {
		log.Warn().Err(err).Msg("failed to remove inverse relation")
	}

	return sendEmpty(w, http.StatusNoContent)
}
