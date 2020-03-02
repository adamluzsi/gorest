package gorest_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleMount() {

	subResourceHandler := gorest.NewHandler(SubResourceController{})
	resourceHandler := gorest.NewHandler(ResourceController{})

	mux := http.NewServeMux()
	gorest.Mount(resourceHandler, `/sub-resources`, subResourceHandler)
	gorest.Mount(mux, `/resources`, resourceHandler)

	// this will cause http.ServeMux to have handlers by the controller structures in hierarchy:
	//	GET /resources/{resourceID}/sub-resources/{sub-resourceID}
}

type ResourceController struct{}

type ContextKeyResource struct{}

func (ctrl ResourceController) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	// lookup Resource by id
	// err out if lookup failed
	// return false if not found
	return context.WithValue(ctx, ContextKeyResource{}, Resource{ID: resourceID}), true, nil
}

type SubResourceController struct{}

type ContextKeySubResource struct{}

func (ctrl SubResourceController) ContextWithResource(ctx context.Context, subResourceID string) (context.Context, bool, error) {
	// lookup Resource by id
	// err out if lookup failed
	// return false if not found
	return context.WithValue(ctx, ContextKeySubResource{}, SubResource{ID: subResourceID}), true, nil
}

func (ctrl SubResourceController) Show(w http.ResponseWriter, r *http.Request) {
	// have access to the top resource because the top resource handler set it for us
	res := r.Context().Value(ContextKeyResource{}).(Resource)
	// have access to the sub resource because the sub resource handler set it for us
	subres := r.Context().Value(ContextKeySubResource{}).(SubResource)
	// print it, because why not
	_, _ = fmt.Fprintf(w, `resource: %v | subresource: %v`, res, subres)
}

type (
	Resource    struct{ ID string }
	SubResource struct{ ID string }
)
