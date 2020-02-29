package gorest_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
)

func NewTestControllerMockHandler(t *testcase.T, code int, msg string) TestControllerMockHandler {
	m := TestControllerMockHandler{T: t, Code: code, Msg: msg}
	return m
}

type TestControllerMockHandler struct {
	T    *testcase.T
	Code int
	Msg  string
}

func (m TestControllerMockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.T != nil {
		m.T.Let(`request ctx`, r.Context())
		bs, err := ioutil.ReadAll(r.Body)
		require.Nil(m.T, err)
		require.Equal(m.T, m.T.I(`body.content`).(string), string(bs))
	}

	http.Error(w, m.Msg, m.Code)
}

type InternalServerErrorController struct {
	Code int
	Msg  string
}

func (h InternalServerErrorController) InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(h.Code)
	_, _ = fmt.Fprintf(w, h.Msg)
}

func NewInternalServerErrorHandler(ctrl interface {
	InternalServerError(w http.ResponseWriter, r *http.Request)
}) http.Handler {
	return http.HandlerFunc(ctrl.InternalServerError)
}

type ErrorContextHandler struct {
	Err error
}

func (ctrl ErrorContextHandler) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	return nil, false, ctrl.Err
}

func (ctrl ErrorContextHandler) Create(w http.ResponseWriter, r *http.Request) {}

func (ctrl ErrorContextHandler) List(w http.ResponseWriter, r *http.Request) {}

func (ctrl ErrorContextHandler) Show(w http.ResponseWriter, r *http.Request) {}

func (ctrl ErrorContextHandler) Update(w http.ResponseWriter, r *http.Request) {}

func (ctrl ErrorContextHandler) Delete(w http.ResponseWriter, r *http.Request) {}


type TestController struct{}

type ContextTestIDKey struct{}

func (d TestController) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	return context.WithValue(ctx, ContextTestIDKey{}, resourceID), true, nil
}

func (d TestController) Create(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `create`)
}

func (d TestController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `list`)
}

func (d TestController) Show(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `show:%s`, r.Context().Value(ContextTestIDKey{}))
}

func (d TestController) Update(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `update:%s`, r.Context().Value(ContextTestIDKey{}))
}

func (d TestController) Delete(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `delete:%s`, r.Context().Value(ContextTestIDKey{}))
}

func (d TestController) NotFound(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `not-found`)
}

func (d TestController) InternalServerError(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `internal-server-error`)
}

var _ interface {
	gorest.ListController
	gorest.CreateController

	gorest.ContextHandler
	gorest.ShowController
	gorest.UpdateController
	gorest.DeleteController

	gorest.WithNotFoundHandler
	gorest.WithInternalServerErrorHandler
} = TestController{}
