package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// apiHandler is a handler that returns an error (GoTrue pattern).
// The error is converted to an HTTP response by the wrapper.
type apiHandler func(w http.ResponseWriter, r *http.Request) error

// router wraps chi.Mux with error-returning handler support.
type router struct {
	chi chi.Router
}

func newRouter() *router {
	return &router{chi: chi.NewRouter()}
}

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.chi.ServeHTTP(w, r)
}

func (rt *router) Route(pattern string, fn func(r *router)) {
	rt.chi.Route(pattern, func(cr chi.Router) {
		sub := &router{chi: cr}
		fn(sub)
	})
}

func (rt *router) Group(fn func(r *router)) {
	rt.chi.Group(func(cr chi.Router) {
		sub := &router{chi: cr}
		fn(sub)
	})
}

func (rt *router) Use(middlewares ...func(http.Handler) http.Handler) {
	rt.chi.Use(middlewares...)
}

func (rt *router) Get(pattern string, h apiHandler) {
	rt.chi.Get(pattern, wrapHandler(h))
}

func (rt *router) Post(pattern string, h apiHandler) {
	rt.chi.Post(pattern, wrapHandler(h))
}

func (rt *router) Put(pattern string, h apiHandler) {
	rt.chi.Put(pattern, wrapHandler(h))
}

func (rt *router) Patch(pattern string, h apiHandler) {
	rt.chi.Patch(pattern, wrapHandler(h))
}

func (rt *router) Delete(pattern string, h apiHandler) {
	rt.chi.Delete(pattern, wrapHandler(h))
}

// wrapHandler converts an apiHandler to a standard http.HandlerFunc.
// Errors returned by the handler are written as JSON error responses.
func wrapHandler(h apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, r, err)
		}
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	var code int
	var body map[string]interface{}

	switch e := err.(type) {
	case *HTTPError:
		code = e.Code
		body = map[string]interface{}{
			"error":   http.StatusText(e.Code),
			"message": e.Message,
		}
		log.Error().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", code).
			Msg(e.Message)
	default:
		code = http.StatusInternalServerError
		body = map[string]interface{}{
			"error":   "Internal Server Error",
			"message": "An unexpected error occurred",
		}
		log.Error().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Err(err).
			Msg("internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
