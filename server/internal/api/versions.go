package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type CreateVersionRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	ReleaseDate *string `json:"release_date,omitempty"`
}

type UpdateVersionRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty"`
	StartDate   *string  `json:"start_date,omitempty"`
	ReleaseDate *string  `json:"release_date,omitempty"`
	Released    *bool    `json:"released,omitempty"`
	SortOrder   *float64 `json:"sort_order,omitempty"`
}

type VersionIssueRequest struct {
	IssueID uuid.UUID `json:"issue_id" validate:"required"`
}

// ListVersions handles GET .../projects/{projectID}/versions
func (a *API) ListVersions(w http.ResponseWriter, r *http.Request) error {
	projectID := getProjectID(r.Context())
	versions, err := a.queries.ListVersionsByProject(r.Context(), projectID)
	if err != nil {
		return internalServerError("failed to list versions")
	}
	return sendJSON(w, http.StatusOK, map[string]interface{}{"results": versions})
}

// CreateVersion handles POST .../projects/{projectID}/versions
func (a *API) CreateVersion(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	projectID := getProjectID(ctx)
	workspaceID := getWorkspaceID(ctx)
	userID := getUserID(ctx)

	var req CreateVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	var startDate, releaseDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return badRequestError("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = &t
	}
	if req.ReleaseDate != nil {
		t, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err != nil {
			return badRequestError("invalid release_date format, expected YYYY-MM-DD")
		}
		releaseDate = &t
	}

	version, err := a.queries.CreateVersion(ctx, projectID, workspaceID, req.Name, req.Description, startDate, releaseDate, &userID)
	if err != nil {
		return internalServerError("failed to create version")
	}
	return sendJSON(w, http.StatusCreated, version)
}

// GetVersion handles GET .../projects/{projectID}/versions/{versionID}
func (a *API) GetVersion(w http.ResponseWriter, r *http.Request) error {
	versionID, err := uuid.FromString(chi.URLParam(r, "versionID"))
	if err != nil {
		return badRequestError("invalid version ID")
	}

	version, err := a.queries.GetVersionByID(r.Context(), versionID)
	if err != nil {
		return notFoundError("version not found")
	}
	return sendJSON(w, http.StatusOK, version)
}

// UpdateVersion handles PUT .../projects/{projectID}/versions/{versionID}
func (a *API) UpdateVersion(w http.ResponseWriter, r *http.Request) error {
	versionID, err := uuid.FromString(chi.URLParam(r, "versionID"))
	if err != nil {
		return badRequestError("invalid version ID")
	}

	var req UpdateVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}
	if err := a.validator.Struct(req); err != nil {
		return validationError(err)
	}

	var startDate, releaseDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return badRequestError("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = &t
	}
	if req.ReleaseDate != nil {
		t, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err != nil {
			return badRequestError("invalid release_date format, expected YYYY-MM-DD")
		}
		releaseDate = &t
	}
	version, err := a.queries.UpdateVersion(r.Context(), versionID, req.Name, req.Description, startDate, releaseDate, req.Released, req.SortOrder)
	if err != nil {
		return internalServerError("failed to update version")
	}
	return sendJSON(w, http.StatusOK, version)
}

// DeleteVersion handles DELETE .../projects/{projectID}/versions/{versionID}
func (a *API) DeleteVersion(w http.ResponseWriter, r *http.Request) error {
	versionID, err := uuid.FromString(chi.URLParam(r, "versionID"))
	if err != nil {
		return badRequestError("invalid version ID")
	}

	if err := a.queries.DeleteVersion(r.Context(), versionID); err != nil {
		return internalServerError("failed to delete version")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// AddIssueFixVersion handles POST .../issues/{issueID}/fix-versions
func (a *API) AddIssueFixVersion(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	var req VersionIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	if err := a.queries.AddIssueFixVersion(r.Context(), issueID, req.IssueID); err != nil {
		return internalServerError("failed to add fix version")
	}
	return sendEmpty(w, http.StatusCreated)
}

// RemoveIssueFixVersion handles DELETE .../issues/{issueID}/fix-versions/{versionID}
func (a *API) RemoveIssueFixVersion(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}
	versionID, err := uuid.FromString(chi.URLParam(r, "versionID"))
	if err != nil {
		return badRequestError("invalid version ID")
	}

	if err := a.queries.RemoveIssueFixVersion(r.Context(), issueID, versionID); err != nil {
		return internalServerError("failed to remove fix version")
	}
	return sendEmpty(w, http.StatusNoContent)
}

// AddIssueAffectsVersion handles POST .../issues/{issueID}/affects-versions
func (a *API) AddIssueAffectsVersion(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}

	var req VersionIssueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError("invalid request body")
	}

	if err := a.queries.AddIssueAffectsVersion(r.Context(), issueID, req.IssueID); err != nil {
		return internalServerError("failed to add affects version")
	}
	return sendEmpty(w, http.StatusCreated)
}

// RemoveIssueAffectsVersion handles DELETE .../issues/{issueID}/affects-versions/{versionID}
func (a *API) RemoveIssueAffectsVersion(w http.ResponseWriter, r *http.Request) error {
	issueID, err := uuid.FromString(chi.URLParam(r, "issueID"))
	if err != nil {
		return badRequestError("invalid issue ID")
	}
	versionID, err := uuid.FromString(chi.URLParam(r, "versionID"))
	if err != nil {
		return badRequestError("invalid version ID")
	}

	if err := a.queries.RemoveIssueAffectsVersion(r.Context(), issueID, versionID); err != nil {
		return internalServerError("failed to remove affects version")
	}
	return sendEmpty(w, http.StatusNoContent)
}
