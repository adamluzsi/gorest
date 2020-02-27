package gorest_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleMount() {
	subresource := &gorest.Handler{
		ContextHandler: ContextHandlerForSubResource{},
		Show: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// have access to the top resource because the top resource handler set it for us
			res := r.Context().Value(`resource`).(Resource)
			// have access to the sub resource because the sub resource handler set it for us
			subres := r.Context().Value(`subresource`).(SubResource)
			// print it, because why not
			fmt.Fprintf(w, `resource: %v | subresource: %v`, res, subres)
		}),
	}

	resource := &gorest.Handler{
		ContextHandler: ContextHandlerForResource{},
	}

	gorest.Mount(resource, `/subresources`, subresource)

	mux := http.NewServeMux()

	// this will cause http.ServeMux to have endpoints like:
	//	GET /resources
	//	POST /resources
	//	GET /resources/{resourceID}
	//	GET /resources/{resourceID}/subresources
	//	GET /resources/{resourceID}/subresources/{subresourceID}
	// and so on
	gorest.Mount(mux, `/resources`, resource)
}

type Resource struct{ ID string }

type ContextHandlerForResource struct{}

func (ContextHandlerForResource) ContextWithResource(ctx context.Context, resourceID string) (context.Context, bool, error) {
	// lookup Resource by id
	// err out if lookup failed
	// return false if not found
	return context.WithValue(ctx, `resource`, Resource{ID: resourceID}), true, nil
}

type SubResource struct{ ID string }

type ContextHandlerForSubResource struct{}

func (ContextHandlerForSubResource) ContextWithResource(ctx context.Context, subResourceID string) (context.Context, bool, error) {
	// lookup Resource by id
	// err out if lookup failed
	// return false if not found
	return context.WithValue(ctx, `subresource`, SubResource{ID: subResourceID}), true, nil
}
