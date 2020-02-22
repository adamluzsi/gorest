package gorest_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
)

func TestHandler_ServeHTTP(t *testing.T) {
	s := testcase.NewSpec(t)

	s.Let(`handler`, func(t *testcase.T) interface{} { return &gorest.Controller{} })
	var handler = func(t *testcase.T) *gorest.Controller { return t.I(`handler`).(*gorest.Controller) }

	var serve = func(t *testcase.T) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(
			t.I(`method`).(string),
			t.I(`path`).(string),
			t.I(`body`).(io.Reader),
		)
		handler(t).ServeHTTP(w, r)
		return w
	}

	s.Let(`body.content`, func(t *testcase.T) interface{} { return strconv.Itoa(rand.Int()) })
	s.Let(`body`, func(t *testcase.T) interface{} { return strings.NewReader(t.I(`body.content`).(string)) })

	var andWhenCustomNotFoundHandlerProvided = func(s *testcase.Spec) {
		s.And(`custom not found handler provided`, func(s *testcase.Spec) {
			const (
				code = http.StatusTeapot
				msg  = `I'm a teapot`
			)
			s.Before(func(t *testcase.T) {
				handler(t).NotFound = NewTestControllerMockHandler(t, code, msg)
			})

			s.Then(`the custom handler will be used`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, code, resp.Code)
				require.Equal(t, msg, strings.TrimSpace(resp.Body.String()))
			})
		})
	}

	s.Let(`resourceID`, func(t *testcase.T) interface{} { return strconv.Itoa(rand.Int()) })
	var resourceID = func(t *testcase.T) string { return t.I(`resourceID`).(string) }

	s.Describe(`GET / - list`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

		s.When(`index handler is not set`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) { handler(t).List = nil })

			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`index handler provided`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				handler(t).List = NewTestControllerMockHandler(t, 200, `index`)
			})

			s.Then(`it will use the index handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, 200, resp.Code)
				require.Equal(t, `index`, strings.TrimSpace(resp.Body.String()))
			})
		})
	})

	s.Describe(`POST / - create`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodPost })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

		s.When(`create handler is not set`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) { handler(t).Create = nil })

			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`create handler provided`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				handler(t).Create = NewTestControllerMockHandler(t, http.StatusCreated, `created`)
			})

			s.Then(`it will use the create handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Equal(t, `created`, strings.TrimSpace(resp.Body.String()))
			})
		})
	})

	var andWhenResourceHandlerIs = func(s *testcase.Spec, sub func(s *testcase.Spec)) {
		var mw = func(next http.Handler, t *testcase.T) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Let(`request context`, r.Context())

				if next != nil {
					next.ServeHTTP(w, r)
				}
			})
		}
		s.When(`Resource handler is supplied`, func(s *testcase.Spec) {
			s.Let(`id key`, func(t *testcase.T) interface{} { return rand.Int() })

			s.Before(func(t *testcase.T) {
				handler(t).ContextHandler = gorest.ContextHandlerFunc(func(ctx context.Context, id string) (context.Context, bool, error) {
					ctx = context.WithValue(ctx, t.I(`id key`), id)
					err, _ := t.I(`Resource handler error`).(error)
					return ctx, t.I(`Resource found`).(bool), err
				})
			})

			s.And(`Resource handler report an error`, func(s *testcase.Spec) {
				const errMsg = `boom`
				s.Let(`Resource handler error`, func(t *testcase.T) interface{} { return errors.New(errMsg) })
				s.Let(`Resource found`, func(t *testcase.T) interface{} { return false })

				s.Then(`then internal server error reported`, func(t *testcase.T) {
					require.Equal(t, http.StatusInternalServerError, serve(t).Code)
				})

				s.And(`if a custom internal server error handler provided`, func(s *testcase.Spec) {
					s.Before(func(t *testcase.T) {
						handler(t).InternalServerError = NewTestControllerMockHandler(t, 518, `boom-pot`)
					})

					s.Then(`the custom handler will be used`, func(t *testcase.T) {
						resp := serve(t)
						require.Equal(t, 518, resp.Code)
						require.Equal(t, `boom-pot`, strings.TrimSpace(resp.Body.String()))
					})
				})
			})

			s.And(`Resource handler states that Resource lookup yield no result`, func(s *testcase.Spec) {
				s.Let(`Resource found`, func(t *testcase.T) interface{} { return false })
				s.Let(`Resource handler error`, func(t *testcase.T) interface{} { return nil })

				s.Then(`it will return with 404`, func(t *testcase.T) {
					require.Equal(t, http.StatusNotFound, serve(t).Code)
				})

				andWhenCustomNotFoundHandlerProvided(s)
			})

			s.And(`Resource found without an error`, func(s *testcase.Spec) {
				s.Let(`Resource found`, func(t *testcase.T) interface{} { return true })
				s.Let(`Resource handler error`, func(t *testcase.T) interface{} { return nil })
				s.Before(func(t *testcase.T) {
					// MITM
					h := handler(t)
					h.Create = mw(h.Create, t)
					h.List = mw(h.List, t)
					h.Show = mw(h.Show, t)
					h.Update = mw(h.Update, t)
					h.Delete = mw(h.Delete, t)
				})

				s.Then(`context will be updated by the context that the Resource handler returned`, func(t *testcase.T) {
					serve(t)

					require.Equal(t, resourceID(t), t.I(`request ctx`).(context.Context).Value(t.I(`id key`)))
				})

				sub(s)
			})
		})
	}

	s.Describe(`GET /{resourceID} - show`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

		s.When(`show handler is not set`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) { handler(t).Show = nil })

			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`show handler provided`, func(s *testcase.Spec) {
			const code = 201
			s.Before(func(t *testcase.T) {
				handler(t).Show = NewTestControllerMockHandler(t, code, `show`)
			})

			s.Then(`it will use the index handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, code, resp.Code)
				require.Equal(t, `show`, strings.TrimSpace(resp.Body.String()))
			})

			andWhenResourceHandlerIs(s, func(s *testcase.Spec) {})
		})
	})

	s.Describe(`PUT|PATCH /{resourceID} - update`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} {
			if rand.Intn(1) == 0 {
				return http.MethodPut
			} else {
				return http.MethodPatch
			}
		})
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

		s.When(`update handler is not set`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) { handler(t).Update = nil })

			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`update handler provided`, func(s *testcase.Spec) {
			const code = 203
			s.Before(func(t *testcase.T) {
				handler(t).Update = NewTestControllerMockHandler(t, code, `update`)
			})

			s.Then(`it will use the index handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, code, resp.Code)
				require.Equal(t, `update`, strings.TrimSpace(resp.Body.String()))
			})

			andWhenResourceHandlerIs(s, func(s *testcase.Spec) {})
		})
	})

	s.Describe(`DELETE /{resourceID} - delete`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodDelete })
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

		s.When(`delete handler is not set`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) { handler(t).Delete = nil })

			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`delete handler provided`, func(s *testcase.Spec) {
			const code = 204
			s.Before(func(t *testcase.T) {
				handler(t).Delete = NewTestControllerMockHandler(t, code, `delete`)
			})

			s.Then(`it will use the delete handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, code, resp.Code)
				require.Equal(t, `delete`, strings.TrimSpace(resp.Body.String()))
			})

			andWhenResourceHandlerIs(s, func(s *testcase.Spec) {})
		})
	})

	s.Describe(`#Mount`, func(s *testcase.Spec) {
		s.Before(func(t *testcase.T) {
			t.Log(`given the top level controller has all the action`)
			handler(t).Create = NewTestControllerMockHandler(t, http.StatusForbidden, `FORBIDDEN`)
			handler(t).List = NewTestControllerMockHandler(t, http.StatusForbidden, `FORBIDDEN`)
			handler(t).Show = NewTestControllerMockHandler(t, http.StatusForbidden, `FORBIDDEN`)
			handler(t).Update = NewTestControllerMockHandler(t, http.StatusForbidden, `FORBIDDEN`)
			handler(t).Delete = NewTestControllerMockHandler(t, http.StatusForbidden, `FORBIDDEN`)
			t.Log(`but none of those actions should be called`)
		})
		s.When(`valid sub controller successfuly mounted`, func(s *testcase.Spec) {
			s.Let(`sub-ctrl`, func(t *testcase.T) interface{} {
				return &gorest.Controller{
					ContextHandler: gorest.DefaultContextHandler{ContextKey: `sub-id`},
					Create:         NewTestControllerMockHandler(t, http.StatusOK, `Create`),
					List:           NewTestControllerMockHandler(t, http.StatusOK, `List`),
					Show:           NewTestControllerMockHandler(t, http.StatusOK, `Show`),
					Update:         NewTestControllerMockHandler(t, http.StatusOK, `Update`),
					Delete:         NewTestControllerMockHandler(t, http.StatusOK, `Delete`),
				}
			})
			s.Before(func(t *testcase.T) { require.Nil(t, handler(t).Mount(`subs`, t.I(`sub-ctrl`).(*gorest.Controller))) })

			var thenActionWillReply = func(s *testcase.Spec, expectedCode int, expectedMessage string) {
				s.Then(`action will reply`, func(t *testcase.T) {
					resp := serve(t)
					require.Equal(t, expectedCode, resp.Code)
					require.Equal(t, expectedMessage, strings.TrimSpace(resp.Body.String()))
				})
			}

			s.And(`request aim create action`, func(s *testcase.Spec) {
				s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodPost })
				s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s/subs`, resourceID(t)) })
				thenActionWillReply(s, http.StatusOK, `Create`)
			})

			s.And(`request aim list action`, func(s *testcase.Spec) {
				s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
				s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s/subs`, resourceID(t)) })
				thenActionWillReply(s, http.StatusOK, `List`)
			})

			s.And(`request aim show action`, func(s *testcase.Spec) {
				s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
				s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s/subs/123`, resourceID(t)) })
				thenActionWillReply(s, http.StatusOK, `Show`)
			})

			s.And(`request aim update action`, func(s *testcase.Spec) {
				s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodPut })
				s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s/subs/123`, resourceID(t)) })
				thenActionWillReply(s, http.StatusOK, `Update`)
			})

			s.And(`request aim delete action`, func(s *testcase.Spec) {
				s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodDelete })
				s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s/subs/123`, resourceID(t)) })
				thenActionWillReply(s, http.StatusOK, `Delete`)
			})
		})

		s.When(`invalid resource name is provided`, func(s *testcase.Spec) {
			s.Then(`it will yield an error`, func(t *testcase.T) {
				require.Error(t, (&gorest.Controller{}).Mount(`/books/`, &gorest.Controller{}),
					`path is not accepted as resource identifier`)
			})
		})

		s.Test(`E2E`, func(t *testcase.T) {
			h := &gorest.Controller{}

			ch := gorest.DefaultContextHandler{ContextKey: `bookID`}
			booksShow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprintf(w, `%s`, ch.GetResourceID(r.Context()))
			})

			require.Nil(t, h.Mount(`books`, &gorest.Controller{ContextHandler: ch, Show: booksShow}))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, `/:topResourceID/books/42`, &bytes.Buffer{})
			h.ServeHTTP(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, `42`, w.Body.String())
		})
	})

	s.Describe(`#Handle`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`pattern`, func(t *testcase.T) interface{} { return `/foo/bar/baz/` })

		s.Before(func(t *testcase.T) {
			handler(t).Handle(t.I(`pattern`).(string), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				const code = http.StatusTeapot
				http.Error(w, http.StatusText(code), code)
			}))
		})

		s.When(`http.ServeMux pattern for the http.Handler is a slash which is basically the match anything`, func(s *testcase.Spec) {
			s.Let(`pattern`, func(t *testcase.T) interface{} { return `/` })

			const (
				code = http.StatusOK
				msg  = `it's fine`
			)

			var thenItWillUseTheControllerHandler = func(s *testcase.Spec) {
				s.Then(`it will use the controller handler`, func(t *testcase.T) {
					resp := serve(t)
					require.Equal(t, code, resp.Code)
					require.Equal(t, msg, strings.TrimSpace(resp.Body.String()))
				})
			}

			var thenItWillUseTheAttachedHandler = func(s *testcase.Spec) {
				s.Then(`it will use the attached handler`, func(t *testcase.T) {
					require.Equal(t, http.StatusTeapot, serve(t).Code)
				})
			}

			s.And(`the path points to`, func(s *testcase.Spec) {
				s.Context(`create action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodPost })
					s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).Create = NewTestControllerMockHandler(t, code, msg) })

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).Create = nil })

						thenItWillUseTheAttachedHandler(s)
					})
				})
				s.Context(`list action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
					s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).List = NewTestControllerMockHandler(t, code, msg) })

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).List = nil })

						thenItWillUseTheAttachedHandler(s)
					})
				})
				s.Context(`show action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
					s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).Show = NewTestControllerMockHandler(t, code, msg) })

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).Show = nil })

						thenItWillUseTheAttachedHandler(s)
					})
				})
				s.Context(`update action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodPut })
					s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).Update = NewTestControllerMockHandler(t, code, msg) })

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).Update = nil })

						thenItWillUseTheAttachedHandler(s)
					})
				})
				s.Context(`delete action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodDelete })
					s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).Delete = NewTestControllerMockHandler(t, code, msg) })

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not set`, func(s *testcase.Spec) {
						s.Before(func(t *testcase.T) { handler(t).Delete = nil })

						thenItWillUseTheAttachedHandler(s)
					})
				})
			})
		})

		s.When(`matching path called`, func(s *testcase.Spec) {
			s.Let(`path`, func(t *testcase.T) interface{} {
				pattern := t.I(`pattern`).(string)
				pattern = strings.TrimPrefix(pattern, `/`)
				pattern = strings.TrimSuffix(pattern, `/`)
				return fmt.Sprintf(`/%s/%s/this`, resourceID(t), pattern)
			})

			s.Then(`it will forward the request to the attached handler`, func(t *testcase.T) {
				require.Equal(t, http.StatusTeapot, serve(t).Code)
			})

			s.And(`multiple handler attached to the controller and one matches the path more precisely`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					pattern := fmt.Sprintf(`%s/this`, strings.TrimSuffix(t.I(`pattern`).(string), `/`))
					handler(t).Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						const code = http.StatusOK
						http.Error(w, http.StatusText(code), code)
					}))
				})

				s.Then(`the most matching one will be used for serving`, func(t *testcase.T) {
					require.Equal(t, http.StatusOK, serve(t).Code)
				})
			})

			s.And(`the path includes multiple slash in the beginning`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) { t.Let(`path`, `///`+t.I(`path`).(string)) })

				s.Then(`it will forward the request to the attached handler`, func(t *testcase.T) {
					require.Equal(t, http.StatusTeapot, serve(t).Code)
				})
			})
		})

		s.When(`non matching path is called`, func(s *testcase.Spec) {
			s.Let(`path`, func(t *testcase.T) interface{} {
				return fmt.Sprintf(`/%s/something/else/than/%s`, resourceID(t), t.I(`pattern`))
			})

			s.Then(`it will run into not found situation`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})
	})
}

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
