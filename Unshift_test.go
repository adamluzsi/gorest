package gorest_test

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/adamluzsi/frameless/fixtures"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
)

func TestUnshiftPathParamFromRequest(t *testing.T) {
	s := testcase.NewSpec(t)

	var subject = func(t *testcase.T) (*http.Request, string) {
		return gorest.UnshiftPathParamFromRequest(t.I(`request`).(*http.Request))
	}

	s.Let(`request`, func(t *testcase.T) interface{} {
		r, err := http.NewRequest(
			t.I(`method`).(string),
			t.I(`path`).(string),
			t.I(`body`).(io.Reader),
		)
		require.Nil(t, err)
		return r
	})

	s.Let(`method`, func(t *testcase.T) interface{} {
		return fixtures.RandomElementFromSlice([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodPatch,
			http.MethodConnect,
			http.MethodHead,
			http.MethodTrace,
		}).(string)
	})

	s.Let(`payload`, func(t *testcase.T) interface{} { return fixtures.RandomString(42) })
	s.Let(`body`, func(t *testcase.T) interface{} { return strings.NewReader(t.I(`payload`).(string)) })

	s.When(`request path has value but without slash prefix`, func(s *testcase.Spec) {
		const value = `value`
		s.Let(`path`, func(t *testcase.T) interface{} { return value })

		s.Then(`it unshift the parameter`, func(t *testcase.T) {
			r, param := subject(t)
			require.Equal(t, value, param)
			require.Equal(t, `/`, r.URL.Path)
		})
	})

	s.When(`request path has value but with slash prefix`, func(s *testcase.Spec) {
		const value = `value`
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, value) })

		s.Then(`it will unshift the parameter`, func(t *testcase.T) {
			r, param := subject(t)
			require.Equal(t, value, param)
			require.Equal(t, `/`, r.URL.Path)
		})

		s.And(`not just one but multiple slashes`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				t.Let(`path`, strings.Repeat(`/`, rand.Intn(40)+2)+t.I(`path`).(string))
			})

			s.Then(`it will unshift the parameter`, func(t *testcase.T) {
				r, param := subject(t)
				require.Equal(t, value, param)
				require.Equal(t, `/`, r.URL.Path)
			})
		})
	})

	s.When(`request path has multiple parts`, func(s *testcase.Spec) {
		const value = `not so random value`
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s/etc`, value) })

		s.Then(`it will unshift the parameter`, func(t *testcase.T) {
			r, param := subject(t)
			require.Equal(t, value, param)
			require.Equal(t, `/etc`, r.URL.Path)
		})
	})

	s.When(`request path is empty`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `` })

		s.Then(`it will unshift the parameter`, func(t *testcase.T) {
			r, param := subject(t)
			require.Equal(t, ``, param)
			require.Equal(t, `/`, r.URL.Path)
		})
	})
}
