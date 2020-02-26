package gorest_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleController() {
	if err := http.ListenAndServe(`:8080`, gorest.NewHandler(ExampleTestController{})); err != nil {
		panic(err.Error())
	}
}

type ExampleTestController struct{}

func (d ExampleTestController) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	return context.WithValue(ctx, `id`, resourceID), true, nil
}

func (d ExampleTestController) Create(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `create`)
}

func (d ExampleTestController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `list`)
}

func (d ExampleTestController) Show(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `show:%s`, r.Context().Value(`id`))
}

func (d ExampleTestController) Update(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `update:%s`, r.Context().Value(`id`))
}

func (d ExampleTestController) Delete(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `delete:%s`, r.Context().Value(`id`))
}

func (d ExampleTestController) NotFound(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `not-found`)
}

func (d ExampleTestController) InternalServerError(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `internal-server-error`)
}

var _ interface {
	gorest.ListController
	gorest.CreateController
	
	gorest.ContextHandler
	gorest.ShowController
	gorest.UpdateController
	gorest.DeleteController

	gorest.WithNotFoundHandler
	gorest.WithInternalServerErrorHandler
} = ExampleTestController{}
