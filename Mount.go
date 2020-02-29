package gorest

import (
	"net/http"
	"strings"
)

// Multiplexer represents a http request multiplexer.
type Multiplexer interface {
	Handle(pattern string, handler http.Handler)
}

// Mount will help to register a handler on a request multiplexer in both as the concrete path to the handler and as a prefix match.
// example:
//	if pattern -> "/something"
//	registered as "/something" for exact match
//	registered as "/something/" for prefix match
//
func Mount(multiplexer Multiplexer, pattern string, handler http.Handler) {
	pattern = `/` + strings.TrimPrefix(pattern, `/`)
	pattern = strings.TrimSuffix(pattern, `/`)
	h := http.StripPrefix(pattern, handler)
	multiplexer.Handle(pattern, h)
	multiplexer.Handle(pattern+`/`, h)
}
