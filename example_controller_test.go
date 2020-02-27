package gorest_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleController_listenAndServe() {
	if err := http.ListenAndServe(`:8080`, gorest.NewHandler(TestController{})); err != nil {
		panic(err.Error())
	}
}

type TestController struct{}

type ContextTestIDKey struct{}

func (d TestController) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	return context.WithValue(ctx, ContextTestIDKey{}, resourceID), true, nil
}

func (d TestController) Create(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `create`)
}

func (d TestController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `list`)
}

func (d TestController) Show(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `show:%s`, r.Context().Value(ContextTestIDKey{}))
}

func (d TestController) Update(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `update:%s`, r.Context().Value(ContextTestIDKey{}))
}

func (d TestController) Delete(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `delete:%s`, r.Context().Value(ContextTestIDKey{}))
}

func (d TestController) NotFound(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `not-found`)
}

func (d TestController) InternalServerError(w http.ResponseWriter, r *http.Request) {
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
} = TestController{}
