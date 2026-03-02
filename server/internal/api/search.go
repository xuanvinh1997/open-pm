package api

import (
	"net/http"
)

type SearchResponse struct {
	Issues []*Issue `json:"issues"`
	Pages  []*Page  `json:"pages"`
}

// GlobalSearch handles GET .../workspaces/{slug}/search?q=text
func (a *API) GlobalSearch(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	workspaceID := getWorkspaceID(ctx)

	q := r.URL.Query().Get("q")
	if len(q) < 2 {
		return sendJSON(w, http.StatusOK, SearchResponse{
			Issues: []*Issue{},
			Pages:  []*Page{},
		})
	}

	issues, err := a.queries.SearchIssues(ctx, workspaceID, q, 20)
	if err != nil {
		issues = []*Issue{}
	}

	pages, err := a.queries.SearchPages(ctx, workspaceID, q, 10)
	if err != nil {
		pages = []*Page{}
	}

	return sendJSON(w, http.StatusOK, SearchResponse{
		Issues: issues,
		Pages:  pages,
	})
}
