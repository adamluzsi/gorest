package gorest

import "net/http"

type Controller interface {
	ListController
	CreateController
	ContextHandler
	ShowController
	UpdateController
	DeleteController
}

type ListController interface {
	// List -- GET /
	// List is the endpoint that responsible to list available resources to the requester.
	List(w http.ResponseWriter, r *http.Request)
}

type CreateController interface {
	// Create -- POST /
	// Create is expected to add a new element to the given collection.
	Create(w http.ResponseWriter, r *http.Request)
}

type ShowController interface {
	// Show -- GET /{resourceID}
	// Show expected to retrieve a specific resource from a collection by ID
	Show(w http.ResponseWriter, r *http.Request)
}

type UpdateController interface {
	// Update -- PUT /{resourceID}
	// Update expected to update the properties of a received resource that is identified by id.
	Update(w http.ResponseWriter, r *http.Request)
}

type DeleteController interface {
	// Delete -- DELETE /{resourceID}
	// Delete is expected to make the resource unavailable one way or an another.
	Delete(w http.ResponseWriter, r *http.Request)
}

type WithNotFoundHandler interface {
	// NotFound is expected to represent a resource not found response to the requester.
	NotFound(w http.ResponseWriter, r *http.Request)
}

type WithInternalServerErrorHandler interface {
	// InternalServerError is expected to represent an unexpected error occurrence in the request.
	InternalServerError(w http.ResponseWriter, r *http.Request)
}
