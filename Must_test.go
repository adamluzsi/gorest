package gorest_test

import (
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
)

func TestMust(t *testing.T) {
	s := testcase.NewSpec(t)

	var subject = func(t *testcase.T) {
		h := &gorest.Handler{}
		gorest.Must(h.Mount(t.I(`name`).(string), &gorest.Handler{}))
	}

	s.When(`no error received as argument`, func(s *testcase.Spec) {
		s.Let(`name`, func(t *testcase.T) interface{} { return `good` })

		s.Then(`it will do nothing`, func(t *testcase.T) {
			subject(t)
		})
	})

	s.When(`error passed as an argument`, func(s *testcase.Spec) {
		s.Let(`name`, func(t *testcase.T) interface{} { return `/this/is/bad/name/for/mounting/a/sub/resource` })

		s.Then(`it will panics`, func(t *testcase.T) {
			require.Panics(t, func() { subject(t) })
		})
	})
}
