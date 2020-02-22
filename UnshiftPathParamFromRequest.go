package gorest

import (
	"net/http"
	"net/url"
	"strings"
)

func UnshiftPathParamFromRequest(r *http.Request) (*http.Request, string) {
	const separator = `/`

	path := r.URL.Path
	isRootPath := strings.HasPrefix(path, separator)
	path = strings.TrimFunc(path, func(r rune) bool { return r == '/' })
	parts := strings.Split(path, separator)

	param := parts[0]
	newPath := strings.Join(parts[1:], separator)
	if isRootPath {
		newPath = `/` + newPath
	}

	r2 := new(http.Request)
	*r2 = *r // copy
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = newPath
	return r2, param
}
