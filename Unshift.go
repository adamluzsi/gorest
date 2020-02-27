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
	const separator = `/`

	isRootPath := strings.HasPrefix(path, separator)
	path = strings.TrimFunc(path, func(r rune) bool { return r == '/' })
	parts := strings.Split(path, separator)

	param := parts[0]
	newPath := strings.Join(parts[1:], separator)

	if isRootPath {
		newPath = `/` + newPath
	}

	return param, newPath
}
