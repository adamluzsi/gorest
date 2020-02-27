package gorest

import "net/http"

// NewHandler builds a new *Handler instance and try to setup the handler parameters with the passed controller.
func NewHandler(ctrl interface{}) *Handler {
	h := &Handler{}
	if i, ok := ctrl.(ContextHandler); ok {
		h.ContextHandler = i
	}
	if i, ok := ctrl.(CreateController); ok {
		h.Create = http.HandlerFunc(i.Create)
	}
	if i, ok := ctrl.(ListController); ok {
		h.List = http.HandlerFunc(i.List)
	}
	if i, ok := ctrl.(ShowController); ok {
		h.Show = http.HandlerFunc(i.Show)
	}
	if i, ok := ctrl.(UpdateController); ok {
		h.Update = http.HandlerFunc(i.Update)
	}
	if i, ok := ctrl.(DeleteController); ok {
		h.Delete = http.HandlerFunc(i.Delete)
	}
	if i, ok := ctrl.(WithNotFoundHandler); ok {
		h.NotFound = http.HandlerFunc(i.NotFound)
	}
	if i, ok := ctrl.(WithInternalServerErrorHandler); ok {
		h.InternalServerError = http.HandlerFunc(i.InternalServerError)
	}
	return h
}
