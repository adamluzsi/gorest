package controllers

import "net/http"

type ListControllerByHTTPHandler struct{ http.Handler }

func (ctrl ListControllerByHTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}

type CreateControllerByHTTPHandler struct{ http.Handler }

func (ctrl CreateControllerByHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}

type ShowControllerByHTTPHandler struct{ http.Handler }

func (ctrl ShowControllerByHTTPHandler) Show(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}

type UpdateControllerByHTTPHandler struct{ http.Handler }

func (ctrl UpdateControllerByHTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}

type DeleteControllerByHTTPHandler struct{ http.Handler }

func (ctrl DeleteControllerByHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}
