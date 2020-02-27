package gorest

import (
	"net/http"
	"strings"
)

type Multiplexer interface {
	Handle(pattern string, handler http.Handler)
	http.Handler
}

func Mount(multiplexer Multiplexer, pattern string, handler http.Handler) {
	pattern = `/` + strings.TrimPrefix(pattern, `/`)
	pattern = strings.TrimSuffix(pattern, `/`)
	h := http.StripPrefix(pattern, handler)
	multiplexer.Handle(pattern, h)
	multiplexer.Handle(pattern+`/`, h)
}
