package gorest_test

import (
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
)

func TestCleanPath(t *testing.T) {
	s := testcase.NewSpec(t)

	var subject = func(t *testcase.T) string {
		return gorest.CanonicalPath(t.I(`path`).(string))
	}

	s.When(`path is a canonical non root path`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `/a/canonical/path` })

		s.Then(`it will leave it as is`, func(t *testcase.T) {
			require.Equal(t, `/a/canonical/path`, subject(t))
		})
	})

	s.When(`path is a canonical root path`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

		s.Then(`it will leave it as is`, func(t *testcase.T) {
			require.Equal(t, `/`, subject(t))
		})
	})

	s.When(`path is empty`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `` })

		s.Then(`it will `, func(t *testcase.T) {
			require.Equal(t, `/`, subject(t))
		})
	})

	s.When(`path is has no leading slash`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `test` })

		s.Then(`it will add the leading slash`, func(t *testcase.T) {
			require.Equal(t, `/test`, subject(t))
		})
	})

	s.When(`path is has multiple leading slash`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `//test` })

		s.Then(`it will remove the extra leading slash`, func(t *testcase.T) {
			require.Equal(t, `/test`, subject(t))
		})
	})

	s.When(`path is starting with leading dot`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `./test` })

		s.Then(`it will remove the leading dot`, func(t *testcase.T) {
			require.Equal(t, `/test`, subject(t))
		})
	})

	s.When(`path is has parent directory reference as double dot`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `/../test` })

		s.Then(`it will remove the parent directory reference double dot`, func(t *testcase.T) {
			require.Equal(t, `/test`, subject(t))
		})
	})

	s.When(`path has trailing slash`, func(s *testcase.Spec) {
		s.Let(`path`, func(t *testcase.T) interface{} { return `/test/` })

		s.Then(`it will preserve the trailing slash`, func(t *testcase.T) {
			require.Equal(t, `/test/`, subject(t))
		})
	})
}

func BenchmarkCanonicalPath(b *testing.B) {
	const path = `/canonical/path`
	for i := 0; i < b.N; i++ {
		gorest.CanonicalPath(path)
	}
}
