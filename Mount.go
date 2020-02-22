package gorest

import (
	"net/http"
	"strings"
)

func Mount(multiplexer interface {
	Handle(pattern string, handler http.Handler)
}, pattern string, handler *Controller) {
	pattern = `/` + strings.TrimPrefix(pattern, `/`)
	pattern = strings.TrimSuffix(pattern, `/`)
	h := http.StripPrefix(pattern, handler)
	multiplexer.Handle(pattern, h)
	multiplexer.Handle(pattern+`/`, h)
}
