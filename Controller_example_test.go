package gorest_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleController() {
	mux := http.NewServeMux()

	var teapotHandler http.Handler = &gorest.Controller{
		ContextHandler: TeapotResourceHandler{},
		Show:           TeapotShowAction{},
	}

	// to mount into a serve multiplexer
	teapotHandler = http.StripPrefix(`/teapots`, teapotHandler)
	mux.Handle(`/teapots`, teapotHandler)
	mux.Handle(`/teapots/`, teapotHandler)

	// OR do the same, but with this function.
	// I often forget to mount my controller also to the path without slash suffix.
	gorest.Mount(mux, `/teapots`, &gorest.Controller{Show: TeapotShowAction{}})

}

type TeapotResourceHandler struct{}

func (t TeapotResourceHandler) ContextWithResource(ctx context.Context, teapotID string) (newCTX context.Context, found bool, err error) {
	teapot, found, err := lookupTeapotByID(ctx, teapotID)
	if err != nil {
		// teapot lookup encountered an unexpected error
		return ctx, false, err
	}
	if !found {
		// teapot not found by id
		return ctx, false, nil
	}
	// set teapot object in context so handlers can access the Resource easily
	return context.WithValue(ctx, `teapot`, teapot), true, nil
}

type TeapotShowAction struct{}

func (e TeapotShowAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	teapot := r.Context().Value(`teapot`).(Teapot)
	_, _ = fmt.Fprintf(w, `my teapot resource: %v`, teapot)
}

// business entity, probably in a different pkg
type Teapot struct{}

// oversimplified external Resource lookup
func lookupTeapotByID(ctx context.Context, id string) (Teapot, bool, error) {
	return Teapot{}, true, nil
}
