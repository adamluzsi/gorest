package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func main() {
	m := http.NewServeMux()
	h := gorest.NewHandler(MyController{})

	// GET http://localhost:8080/collection-name-in-plural
	// GET http://localhost:8080/collection-name-in-plural/
	// GET http://localhost:8080/collection-name-in-plural/resource-id
	gorest.Mount(m, `collection-name-in-plural`, h)

	if err := http.ListenAndServe(`:8080`, m); err != nil {
		log.Fatal(err.Error())
	}
}

type MyController struct{}

func (ctrl MyController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `list`)
}

type ContextKeyMyResource struct{}

type resource struct {
	ID string
}

func (ctrl MyController) ContextWithResource(ctx context.Context, resourceID string) (context.Context, bool, error) {
	return context.WithValue(ctx, ContextKeyMyResource{}, resource{ID: resourceID}), true, nil
}

func (ctrl MyController) Show(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `show; id: %v`, r.Context().Value(ContextKeyMyResource{}))
}
