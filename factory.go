package gorest

import "net/http"

func NewHandler(ctrl interface{}) *Handler {
	h := &Handler{}
	if i, ok := ctrl.(ContextHandler); ok {
		h.ContextHandler = i
	}
	if i, ok := ctrl.(interface{ Create(w http.ResponseWriter, r *http.Request) }); ok {
		h.Create = http.HandlerFunc(i.Create)
	}
	if i, ok := ctrl.(interface{ List(w http.ResponseWriter, r *http.Request) }); ok {
		h.List = http.HandlerFunc(i.List)
	}
	if i, ok := ctrl.(interface{ Show(w http.ResponseWriter, r *http.Request) }); ok {
		h.Show = http.HandlerFunc(i.Show)
	}
	if i, ok := ctrl.(interface{ Update(w http.ResponseWriter, r *http.Request) }); ok {
		h.Update = http.HandlerFunc(i.Update)
	}
	if i, ok := ctrl.(interface{ Delete(w http.ResponseWriter, r *http.Request) }); ok {
		h.Delete = http.HandlerFunc(i.Delete)
	}
	if i, ok := ctrl.(interface{ NotFound(w http.ResponseWriter, r *http.Request) }); ok {
		h.NotFound = http.HandlerFunc(i.NotFound)
	}
	if i, ok := ctrl.(interface{ InternalServerError(w http.ResponseWriter, r *http.Request) }); ok {
		h.InternalServerError = http.HandlerFunc(i.InternalServerError)
	}
	return h
}
