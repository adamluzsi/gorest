package gorest

import (
	"net/http"
	"strings"
)

func newRoutes() *routes {
	return &routes{
		mux:         http.NewServeMux(),
		hasResource: make(map[string]struct{}),
	}
}

type routes struct {
	mux         *http.ServeMux
	hasResource map[string]struct{}
	hasRoot     bool
}

func (m *routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (m *routes) Handle(pattern string, handler http.Handler) {
	if pattern == `/` {
		m.hasRoot = true
	} else {
		name := m.name(pattern)
		m.hasResource[name] = struct{}{}
	}

	m.mux.Handle(pattern, handler)
}

func (m *routes) HasResource(path string) bool {
	_, ok := m.hasResource[m.name(path)]
	return ok
}

func (m *routes) HasRoot() bool {
	return m.hasRoot
}

func (m *routes) name(path string) string {
	for _, part := range strings.Split(path, `/`) {
		if part != `` {
			return part
		}
	}

	return ``
}
