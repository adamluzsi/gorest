package gorest_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleHandler() {
	mux := http.NewServeMux()
	teapotHandler := gorest.NewHandler(TeapotController{})

	// to mount into a serve multiplexer
	h := http.StripPrefix(`/teapots`, teapotHandler)
	mux.Handle(`/teapots`, h)
	mux.Handle(`/teapots/`, h)

	// or do the same, but with this function.
	gorest.Mount(mux, `/teapots`, teapotHandler)
}

type TeapotController struct{}

type ContextKeyTeapot struct{}

func (ctrl TeapotController) ContextWithResource(ctx context.Context, teapotID string) (newCTX context.Context, found bool, err error) {
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
	return context.WithValue(ctx, ContextKeyTeapot{}, teapot), true, nil
}

func (ctrl TeapotController) Show(w http.ResponseWriter, r *http.Request) {
	teapot := r.Context().Value(ContextKeyTeapot{}).(Teapot)
	_, _ = fmt.Fprintf(w, `my teapot resource: %v`, teapot)
}

// business entity, probably in a different pkg
type Teapot struct{}

// oversimplified external Resource lookup
func lookupTeapotByID(ctx context.Context, id string) (Teapot, bool, error) {
	return Teapot{}, true, nil
}
