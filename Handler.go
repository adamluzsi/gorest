package gorest

import (
	"context"
	"net/http"
	"strings"
)

type Handler struct {
	ContextHandler      ContextHandler
	NotFound            http.Handler
	InternalServerError http.Handler
	operations          struct {
		collection operations
		resource   operations
	}
	handlers handlers
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if cause := recover(); cause != nil {
			h.internalServerError(w, r)
		}
	}()

	var method = r.Method

	switch r.URL.Path {
	case `/`, ``:
		ch, ok := h.lookupCollectionHandler(method, r.URL.Path)
		if !ok {
			h.notFound(w, r)
			return
		}

		ch.ServeHTTP(w, r)

	default: // dynamic path
		ctx := r.Context()
		r, resourceID := UnshiftPathParamFromRequest(r)
		ctx, found, err := h.handleResourceID(ctx, resourceID)

		if err != nil {
			h.internalServerError(w, r)
			return
		}

		if !found {
			h.notFound(w, r)
			return
		}

		r = r.WithContext(ctx)

		rh, ok := h.lookupResourceHandler(method, r.URL.Path)
		if !ok {
			h.notFound(w, r)
			return
		}

		rh.ServeHTTP(w, r)
	}
}

func (h *Handler) Handle(pattern string, handler http.Handler) {
	if pattern == `/` {
		h.handlers.hasRootHandler = true
	} else {
		if h.handlers.prefixes == nil {
			h.handlers.prefixes = make(map[string]struct{})
		}
		h.handlers.prefixes[h.prefix(pattern)] = struct{}{}
	}

	if h.handlers.ServeMux == nil {
		h.handlers.ServeMux = http.NewServeMux()
	}
	h.handlers.ServeMux.Handle(pattern, handler)
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

func (h *Handler) defaultInternalServerError(w http.ResponseWriter, _ *http.Request) {
	const code = http.StatusInternalServerError
	http.Error(w, http.StatusText(code), code)
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

func (h *Handler) hasHandlerWithPrefixThatMatch(path string) bool {
	if h.handlers.prefixes == nil {
		return false
	}
	_, ok := h.handlers.prefixes[h.prefix(path)]
	return ok
}

func (h *Handler) prefix(path string) string {
	for _, part := range strings.Split(path, `/`) {
		if part != `` {
			return part
		}
	}

	return ``
}

func (h *Handler) lookupCollectionHandler(method, path string) (http.Handler, bool) {
	handler, ok := h.operations.collection.Lookup(method, path)
	if !ok && h.handlers.hasRootHandler {
		return h.handlers, true
	}
	return handler, ok
}

func (h *Handler) lookupResourceHandler(method, path string) (http.Handler, bool) {
	if h.hasHandlerWithPrefixThatMatch(path) {
		return h.handlers, true
	}
	handler, ok := h.operations.resource.Lookup(method, path)
	if !ok && h.handlers.hasRootHandler {
		return h.handlers, true
	}
	return handler, ok
}

type operations struct {
	routes map[string]http.Handler
}

func (o operations) Lookup(method, path string) (http.Handler, bool) {
	if o.routes == nil {
		return nil, false
	}
	h, ok := o.routes[method]
	return h, ok
}

func (o *operations) Set(httpMethod string, handler http.Handler) {
	if o.routes == nil {
		o.routes = make(map[string]http.Handler)
	}
	o.routes[httpMethod] = handler
}

type handlers struct {
	*http.ServeMux
	prefixes       map[string]struct{}
	hasRootHandler bool
}
