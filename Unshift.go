package gorest

import (
	"net/http"
	"net/url"
	"strings"
)

func UnshiftPathParamFromRequest(r *http.Request) (*http.Request, string) {
	id, path := Unshift(r.URL.Path)
	r2 := new(http.Request)
	*r2 = *r // copy
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = path
	return r2, id
}

func Unshift(path string) (id string, remainingPath string) {
	pathParam := strings.Split(CanonicalPath(path), `/`)[1]
	newPath := path[(strings.Index(path, pathParam) + len(pathParam)):]
	if len(newPath) == 0 || newPath[0] != '/' {
		newPath = `/` + newPath
	}
	return pathParam, newPath
}
