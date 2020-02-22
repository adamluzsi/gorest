package gorest_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
)

func TestMount(t *testing.T) {
	ctxHandler := gorest.DefaultContextHandler{ContextKey: `resourceID`}

	resourcesShow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `%s`, ctxHandler.GetResourceID(r.Context()))
	})

	resources := &gorest.Controller{
		ContextHandler: ctxHandler,
		Show:           resourcesShow,
	}

	mux := http.NewServeMux()
	gorest.Mount(mux, `/routes`, resources)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, `/routes/resourceID/`, &bytes.Buffer{})
	mux.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `resourceID`, w.Body.String())
}
