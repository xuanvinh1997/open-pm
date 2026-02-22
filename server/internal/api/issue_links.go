package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateIssueLinkRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
	URL   string `json:"url" validate:"required,url"`
}

type UpdateIssueLinkRequest struct {
	Title *string `json:"title,omitempty"`
	URL   *string `json:"url,omitempty"`
}

// ListIssueLinks handles GET .../issues/{issueID}/links
func (a *API) ListIssueLinks(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	links, err := a.queries.ListIssueLinksByIssue(r.Context(), issueID)
	if err != nil {
		return internalServerError("failed to list links")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": links})
}

// CreateIssueLink handles POST .../issues/{issueID}/links
func (a *API) CreateIssueLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := getUserID(ctx)

	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	var req CreateIssueLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	link, err := a.queries.CreateIssueLink(ctx, issueID, req.Title, req.URL, &userID)
	if err != nil {
		return internalServerError("failed to create link")
	}
	return sendJSON(w, http.StatusCreated, link)
}

// UpdateIssueLink handles PUT .../issues/{issueID}/links/{linkID}
func (a *API) UpdateIssueLink(w http.ResponseWriter, r *http.Request) error {
	linkID, err := uuid.FromString(chi.URLParam(r, "linkID"))
	if err != nil {
		return badRequestError("invalid link ID")
	}

	var req UpdateIssueLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	link, err := a.queries.UpdateIssueLink(r.Context(), linkID, req.Title, req.URL)
	if err != nil {
		return internalServerError("failed to update link")
	}
	return sendJSON(w, http.StatusOK, link)
}

// DeleteIssueLink handles DELETE .../issues/{issueID}/links/{linkID}
func (a *API) DeleteIssueLink(w http.ResponseWriter, r *http.Request) error {
	linkID, err := uuid.FromString(chi.URLParam(r, "linkID"))
	if err != nil {
		return badRequestError("invalid link ID")
	}

	if err := a.queries.DeleteIssueLink(r.Context(), linkID); err != nil {
		return internalServerError("failed to delete link")
	}
	return sendEmpty(w, http.StatusNoContent)
}
