package gorest

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type Handler struct {
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if cause := recover(); cause != nil {
			h.internalServerError(w, r)
		}
	}()
	switch r.URL.Path {
	case `/`, ``:
		switch r.Method {
		case http.MethodGet:
			h.serve(h.List, w, r)

		case http.MethodPost:
			h.serve(h.Create, w, r)

		}
	default: // dynamic path
		ctx := r.Context()

		r, resourceID := UnshiftPathParamFromRequest(r)
		ctx, found, err := h.handleResourceID(r.Context(), resourceID)
		if err != nil {
			h.internalServerError(w, r)
			return
		}
		if !found {
			h.notFound(w, r)
			return
		}

		r = r.WithContext(ctx)

		if h.routes != nil && h.routes.HasResource(r.URL.Path) {
			h.routes.ServeHTTP(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.serve(h.Show, w, r)
		case http.MethodPut, http.MethodPatch:
			h.serve(h.Update, w, r)
		case http.MethodDelete:
			h.serve(h.Delete, w, r)
		default:
			h.serve(nil, w, r)
		}
	}
}

func (h *Handler) serve(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	if handler != nil {
		handler.ServeHTTP(w, r)
		return
	}

	if h.routes != nil && h.routes.HasRoot() {
		h.routes.ServeHTTP(w, r)
		return
	}

	h.notFound(w, r)
}

func (h *Handler) internalServerError(w http.ResponseWriter, r *http.Request) {
	if h.InternalServerError == nil {
		h.defaultInternalServerError(w, r)
		return
	}
	defer func() {
		if cause := recover(); cause != nil {
			h.defaultInternalServerError(w, r)
		}
	}()
	h.InternalServerError.ServeHTTP(w, r)
}

func (h *Handler) defaultInternalServerError(w http.ResponseWriter, r *http.Request) {
	const code = http.StatusInternalServerError
	http.Error(w, http.StatusText(code), code)
	return
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	if h.NotFound == nil {
		http.NotFound(w, r)
		return
	}
	h.NotFound.ServeHTTP(w, r)
}

func (h *Handler) handleResourceID(ctx context.Context, resourceID string) (context.Context, bool, error) {
	if h.ContextHandler == nil {
		return ctx, true, nil
	}
	return h.ContextHandler.ContextWithResource(ctx, resourceID)
}

func (h *Handler) getRoutes() *routes {
	if h.routes == nil {
		h.routes = newRoutes()
	}
	return h.routes
}

func (h *Handler) Mount(name string, handler *Handler) error {
	if strings.Contains(name, `/`) {
		return errors.New(`resource should not include "/"`)
	}
	var (
		path    = `/` + name
		pattern = path + `/`
	)
	Mount(h.getRoutes(), pattern, handler)
	return nil
}

func (h *Handler) Handle(pattern string, handler http.Handler) {
	h.getRoutes().Handle(pattern, handler)
}
