package gorest

import "net/http"

type Controller interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)

	ContextHandler
	Show(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewHandler(i Controller) *Handler {
	return &Handler{
		ContextHandler: i,
		Create:         http.HandlerFunc(i.Create),
		List:           http.HandlerFunc(i.List),
		Show:           http.HandlerFunc(i.Show),
		Update:         http.HandlerFunc(i.Update),
		Delete:         http.HandlerFunc(i.Delete),
	}
}
