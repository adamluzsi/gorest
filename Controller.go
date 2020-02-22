package gorest

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type Controller struct {
	ContextHandler ContextHandler

	Create http.Handler
	List   http.Handler
	Show   http.Handler
	Update http.Handler
	Delete http.Handler

	NotFound            http.Handler
	InternalServerError http.Handler

	routes *routes
}

func (ctrl *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case `/`, ``:
		switch r.Method {
		case http.MethodGet:
			ctrl.serve(ctrl.List, w, r)

		case http.MethodPost:
			ctrl.serve(ctrl.Create, w, r)

		}
	default: // dynamic path
		ctx := r.Context()

		r, resourceID := UnshiftPathParamFromRequest(r)
		ctx, found, err := ctrl.handleResourceID(r.Context(), resourceID)
		if err != nil {
			ctrl.internalServerError(w, r)
			return
		}
		if !found {
			ctrl.notFound(w, r)
			return
		}

		r = r.WithContext(ctx)

		if ctrl.routes != nil && ctrl.routes.HasResource(r.URL.Path) {
			ctrl.routes.ServeHTTP(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			ctrl.serve(ctrl.Show, w, r)
		case http.MethodPut, http.MethodPatch:
			ctrl.serve(ctrl.Update, w, r)
		case http.MethodDelete:
			ctrl.serve(ctrl.Delete, w, r)
		default:
			ctrl.serve(nil, w, r)
		}
	}
}

func (ctrl *Controller) serve(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	if handler != nil {
		handler.ServeHTTP(w, r)
		return
	}

	if ctrl.routes != nil && ctrl.routes.HasRoot() {
		ctrl.routes.ServeHTTP(w, r)
		return
	}

	ctrl.notFound(w, r)
}

func (ctrl *Controller) internalServerError(w http.ResponseWriter, r *http.Request) {
	if ctrl.InternalServerError == nil {
		const code = http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}
	ctrl.InternalServerError.ServeHTTP(w, r)
}

func (ctrl *Controller) notFound(w http.ResponseWriter, r *http.Request) {
	if ctrl.NotFound == nil {
		http.NotFound(w, r)
		return
	}
	ctrl.NotFound.ServeHTTP(w, r)
}

func (ctrl *Controller) handleResourceID(ctx context.Context, resourceID string) (context.Context, bool, error) {
	if ctrl.ContextHandler == nil {
		return ctx, true, nil
	}
	return ctrl.ContextHandler.WithResource(ctx, resourceID)
}

func (ctrl *Controller) getRoutes() *routes {
	if ctrl.routes == nil {
		ctrl.routes = newRoutes()
	}
	return ctrl.routes
}

func (ctrl *Controller) Mount(name string, handler *Controller) error {
	if strings.Contains(name, `/`) {
		return errors.New(`resource should not include "/"`)
	}
	var (
		path    = `/` + name
		pattern = path + `/`
	)
	Mount(ctrl.getRoutes(), pattern, handler)
	return nil
}

func (ctrl *Controller) Handle(pattern string, handler http.Handler) {
	ctrl.getRoutes().Handle(pattern, handler)
}
