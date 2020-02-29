package gorest_test

import (
	"context"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleController_creatingHandlerFromController() {
	NewMyHandler()
}

func NewMyHandler() http.Handler {
	return gorest.NewHandler(MyController{})
}

type MyController struct{}

func (ctrl MyController) List(w http.ResponseWriter, r *http.Request) {}

func (ctrl MyController) Create(w http.ResponseWriter, r *http.Request) {}

type ContextKeyMyResource struct{}

func (ctrl MyController) ContextWithResource(ctx context.Context, resourceID string) (context.Context, bool, error) {
	var resource interface{} // fetch resource by id
	return context.WithValue(ctx, ContextKeyMyResource{}, resource), true, nil
}

func (ctrl MyController) Show(w http.ResponseWriter, r *http.Request) {}

func (ctrl MyController) Update(w http.ResponseWriter, r *http.Request) {}

func (ctrl MyController) Delete(w http.ResponseWriter, r *http.Request) {}
