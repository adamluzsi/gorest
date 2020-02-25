package gorest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
)

func TestNewHandler(t *testing.T) {
	s := testcase.NewSpec(t)

	var request = func(t *testcase.T) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(t.I(`method`).(string), t.I(`path`).(string), nil)
		gorest.NewHandler(t.I(`controller`)).ServeHTTP(w, r)
		return w
	}

	s.Let(`controller`, func(t *testcase.T) interface{} { return ExampleTestController{} })

	s.Describe(`#List`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

		s.Then(`it will use List method to reply`, func(t *testcase.T) {
			require.Contains(t, request(t).Body.String(), `list`)
		})
	})

	s.Describe(`#Create`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodPost })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

		s.Then(`it will use Create method to reply`, func(t *testcase.T) {
			require.Contains(t, request(t).Body.String(), `create`)
		})
	})

	s.Describe(`#Show`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/42` })

		s.Then(`it will use Show method to reply`, func(t *testcase.T) {
			require.Contains(t, request(t).Body.String(), `show:42`)
		})
	})

	s.Describe(`#Update`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodPut })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/42` })

		s.Then(`it will use Update method to reply`, func(t *testcase.T) {
			require.Contains(t, request(t).Body.String(), `update:42`)
		})
	})

	s.Describe(`#Delete`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodDelete })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/42` })

		s.Then(`it will use Delete method to reply`, func(t *testcase.T) {
			require.Contains(t, request(t).Body.String(), `delete:42`)
		})
	})

	s.Describe(`#NotFound`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/not-found` })

		s.Then(`it will use the not found method to reply`, func(t *testcase.T) {
			require.Contains(t, request(t).Body.String(), `not-found`)
		})
	})

	s.Describe(`#InternalServerError`, func(s *testcase.Spec) {
		s.Let(`controller`, func(t *testcase.T) interface{} {
			return struct {
				InternalServerErrorController
				ErrorContextHandler
			}{
				InternalServerErrorController: InternalServerErrorController{
					Code: 500,
					Msg:  "custom-internal-server-error",
				},
				ErrorContextHandler: ErrorContextHandler{Err: errors.New(`boom`)},
			}
		})
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/42` })

		s.Then(`it will use the not found method to reply`, func(t *testcase.T) {
			require.Contains(t, request(t).Body.String(), `custom-internal-server-error`)
		})
	})

}
