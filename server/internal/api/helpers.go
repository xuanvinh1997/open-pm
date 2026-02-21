package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func sendJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func sendEmpty(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	return nil
}

type PaginatedResponse struct {
	Results    interface{} `json:"results"`
	TotalCount int64       `json:"total_count"`
	NextPage   int         `json:"next_page,omitempty"`
	PrevPage   int         `json:"prev_page,omitempty"`
}

type PaginationParams struct {
	Page    int
	PerPage int
}

func parsePagination(r *http.Request) PaginationParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}

	return PaginationParams{Page: page, PerPage: perPage}
}

func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PerPage
}
