package gorest_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
)

var (
	_ gorest.Multiplexer = http.NewServeMux()
	_ gorest.Multiplexer = gorest.NewHandler(nil)
)

func TestMount(t *testing.T) {
	s := testcase.NewSpec(t)

	var subject = func(t *testcase.T) {
		gorest.Mount(
			t.I(`multiplexer`).(gorest.Multiplexer),
			t.I(`pattern`).(string),
			t.I(`handler`).(http.Handler),
		)
	}

	var multiplexer = func(t *testcase.T) gorest.Multiplexer { return t.I(`multiplexer`).(gorest.Multiplexer) }
	s.Let(`multiplexer`, func(t *testcase.T) interface{} { return http.NewServeMux() })
	s.Let(`handler`, func(t *testcase.T) interface{} { return gorest.NewHandler(TestController{}) })

	s.When(`pattern lack trailing slash`, func(s *testcase.Spec) {
		s.Let(`pattern`, func(t *testcase.T) interface{} { return `/path0` })

		s.Then(`it will be still available to call even for the under paths`, func(t *testcase.T) {
			subject(t)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, `/path0/123`, nil)
			multiplexer(t).ServeHTTP(w, r)
			require.Contains(t, w.Body.String(), `show:123`)
		})
	})

	s.When(`pattern lack leading slash`, func(s *testcase.Spec) {
		s.Let(`pattern`, func(t *testcase.T) interface{} { return `path1/` })

		s.Then(`it will be still available to call even for the under paths`, func(t *testcase.T) {
			subject(t)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, `/path1/123`, nil)
			multiplexer(t).ServeHTTP(w, r)
			require.Contains(t, w.Body.String(), `show:123`)
		})
	})

	s.When(`pattern lack leading and trailing slash`, func(s *testcase.Spec) {
		s.Let(`pattern`, func(t *testcase.T) interface{} { return `path2` })

		s.Then(`it will be still available to call even for the under paths`, func(t *testcase.T) {
			subject(t)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, `/path2/123`, nil)
			multiplexer(t).ServeHTTP(w, r)
			require.Contains(t, w.Body.String(), `show:123`)
		})
	})

	s.When(`pattern includes nested path`, func(s *testcase.Spec) {
		s.Let(`pattern`, func(t *testcase.T) interface{} { return `/test/this/out/` })

		s.Then(`it will be still available to call even for the under paths`, func(t *testcase.T) {
			subject(t)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, `/test/this/out/123`, nil)
			multiplexer(t).ServeHTTP(w, r)
			require.Contains(t, w.Body.String(), `show:123`)
		})
	})

	s.Test(`E2E`, func(t *testcase.T) {
		ctxHandler := gorest.DefaultContextHandler{ContextKey: `resourceID`}

		resourcesShow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, `%s`, ctxHandler.GetResourceID(r.Context()))
		})

		resources := &gorest.Handler{
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
	})
}
