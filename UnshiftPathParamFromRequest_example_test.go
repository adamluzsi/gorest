package gorest_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleUnshiftPathParam() {
	var withResourceID = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r, param := gorest.UnshiftPathParamFromRequest(r)
			// verify Resource id is valid and present
			// verify authority of the requester to this resourceID
			r = r.WithContext(context.WithValue(r.Context(), exampleUnshiftCtxKeyResourceID{}, param)) // add request context
			next.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()
	mux.Handle(`/routes/`, http.StripPrefix(`/routes`, withResourceID(NewResourceHandler())))
}

type exampleUnshiftCtxKeyResourceID struct{}

type ExampleUnshiftHTTPHandler struct {
	ServeMux *http.ServeMux
}

func NewResourceHandler() *ExampleUnshiftHTTPHandler {
	h := &ExampleUnshiftHTTPHandler{ServeMux: http.NewServeMux()}
	h.ServeMux.HandleFunc(`/show`, h.show)
	return h
}

// http path for this is /routes/:resource_id/show
// but anything is fine as long the context has the resourceID
func (h *ExampleUnshiftHTTPHandler) show(w http.ResponseWriter, r *http.Request) {
	resourceID := r.Context().Value(exampleUnshiftCtxKeyResourceID{}).(string)
	_, _ = fmt.Fprintf(w, `path param is: %s`, resourceID)
}

func (h *ExampleUnshiftHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeMux.ServeHTTP(w, r)
}
