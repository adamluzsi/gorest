package gorest_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/adamluzsi/gorest"
	"github.com/adamluzsi/gorest/controllers"
)

var _ interface {
	http.Handler
	gorest.Multiplexer
} = &gorest.Handler{}

func TestHandler_ServeHTTP(t *testing.T) {
	s := testcase.NewSpec(t)

	s.Let(`controller`, func(t *testcase.T) interface{} { return nil })
	s.Let(`handler`, func(t *testcase.T) interface{} { return gorest.NewHandler(t.I(`controller`)) })
	var handler = func(t *testcase.T) *gorest.Handler { return t.I(`handler`).(*gorest.Handler) }

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

		s.When(`index handler is not yet set`, func(s *testcase.Spec) {
			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`controller with list action is provided`, func(s *testcase.Spec) {
			s.Let(`controller`, func(t *testcase.T) interface{} {
				return controllers.ListControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, 200, `index`)}
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

		s.When(`create handler is not yet set`, func(s *testcase.Spec) {
			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`controller with create action is provided`, func(s *testcase.Spec) {
			s.Let(`controller`, func(t *testcase.T) interface{} {
				return controllers.CreateControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, http.StatusCreated, `created`)}
			})

			s.Then(`it will use the create handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, http.StatusCreated, resp.Code)
				require.Equal(t, `created`, strings.TrimSpace(resp.Body.String()))
			})
		})
	})

	type ContextKeyTestResourceHandlerResourceID struct{}

	var andWhenResourceHandlerIs = func(s *testcase.Spec, sub func(s *testcase.Spec)) {
		s.When(`Resource handler is supplied`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				handler(t).ContextHandler = gorest.ContextHandlerFunc(func(ctx context.Context, id string) (context.Context, bool, error) {
					ctx = context.WithValue(ctx, ContextKeyTestResourceHandlerResourceID{}, id)
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

				s.And(`if a custom internal server controller with error action is provided`, func(s *testcase.Spec) {
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

				s.Context(``, sub)

				s.Then(`it yields no error`, func(t *testcase.T) {
					rr := serve(t)
					require.NotEqual(t, http.StatusInternalServerError, rr.Code)
					require.NotEqual(t, http.StatusNotFound, rr.Code)
				})
			})
		})
	}

	s.Describe(`GET /{resourceID} - show`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

		s.When(`show handler is not yet set`, func(s *testcase.Spec) {
			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`controller with show action is provided`, func(s *testcase.Spec) {
			const code = 201
			s.Let(`controller`, func(t *testcase.T) interface{} {
				return controllers.ShowControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, code, `show`)}
			})

			s.Then(`it will use the index handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, code, resp.Code)
				require.Equal(t, `show`, strings.TrimSpace(resp.Body.String()))
			})

			andWhenResourceHandlerIs(s, func(s *testcase.Spec) {
				s.Let(`controller`, func(t *testcase.T) interface{} {
					return controllers.ShowControllerByHTTPHandler{Handler: GenericResourceHandler{
						Message:    "show",
						ContextKey: ContextKeyTestResourceHandlerResourceID{},
					}}
				})

				s.Then(`stored values from the context can be retrieved`, func(t *testcase.T) {
					require.Contains(t, serve(t).Body.String(), fmt.Sprintf(`show:%s`, resourceID(t)))
				})
			})
		})
	})

	s.Describe(`PUT|PATCH /{resourceID} - update`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} {
			if rand.Intn(1) == 0 {
				return http.MethodPut
			}
			return http.MethodPatch
		})
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

		s.When(`update handler is not yet set`, func(s *testcase.Spec) {
			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`controller with update action is provided`, func(s *testcase.Spec) {
			const code = 203
			s.Let(`controller`, func(t *testcase.T) interface{} {
				return controllers.UpdateControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, code, `update`)}
			})

			s.Then(`it will use the index handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, code, resp.Code)
				require.Equal(t, `update`, strings.TrimSpace(resp.Body.String()))
			})

			andWhenResourceHandlerIs(s, func(s *testcase.Spec) {
				s.Let(`controller`, func(t *testcase.T) interface{} {
					return controllers.UpdateControllerByHTTPHandler{Handler: GenericResourceHandler{
						Message:    "update",
						ContextKey: ContextKeyTestResourceHandlerResourceID{},
					}}
				})

				s.Then(`stored values from the context can be retrieved`, func(t *testcase.T) {
					require.Contains(t, serve(t).Body.String(), fmt.Sprintf(`update:%s`, resourceID(t)))
				})
			})
		})
	})

	s.Describe(`DELETE /{resourceID} - delete`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodDelete })
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

		s.When(`delete handler is not set yet`, func(s *testcase.Spec) {
			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`controller with delete action is provided`, func(s *testcase.Spec) {
			const code = 204
			s.Let(`controller`, func(t *testcase.T) interface{} {
				return controllers.DeleteControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, code, `delete`)}
			})

			s.Then(`it will use the delete handler`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, code, resp.Code)
				require.Equal(t, `delete`, strings.TrimSpace(resp.Body.String()))
			})

			andWhenResourceHandlerIs(s, func(s *testcase.Spec) {
				s.Let(`controller`, func(t *testcase.T) interface{} {
					return controllers.DeleteControllerByHTTPHandler{Handler: GenericResourceHandler{
						Message:    "delete",
						ContextKey: ContextKeyTestResourceHandlerResourceID{},
					}}
				})

				s.Then(`stored values from the context can be retrieved`, func(t *testcase.T) {
					require.Contains(t, serve(t).Body.String(), fmt.Sprintf(`delete:%s`, resourceID(t)))
				})
			})
		})
	})

	s.Describe(`with Mount`, func(s *testcase.Spec) {
		s.Let(`controller`, func(t *testcase.T) interface{} {
			t.Log(`given the top level controller has all the action`)
			h := NewTestControllerMockHandler(t, http.StatusForbidden, `FORBIDDEN`)
			t.Log(`but none of those actions should be called`)
			return struct {
				controllers.CreateControllerByHTTPHandler
				controllers.ListControllerByHTTPHandler
				controllers.ShowControllerByHTTPHandler
				controllers.UpdateControllerByHTTPHandler
				controllers.DeleteControllerByHTTPHandler
			}{
				CreateControllerByHTTPHandler: controllers.CreateControllerByHTTPHandler{Handler: h},
				ListControllerByHTTPHandler:   controllers.ListControllerByHTTPHandler{Handler: h},
				ShowControllerByHTTPHandler:   controllers.ShowControllerByHTTPHandler{Handler: h},
				UpdateControllerByHTTPHandler: controllers.UpdateControllerByHTTPHandler{Handler: h},
				DeleteControllerByHTTPHandler: controllers.DeleteControllerByHTTPHandler{Handler: h},
			}
		})

		s.When(`valid sub controller successfully mounted`, func(s *testcase.Spec) {
			s.Let(`sub-handler`, func(t *testcase.T) interface{} {
				return gorest.NewHandler(
					struct {
						gorest.ContextHandler
						controllers.CreateControllerByHTTPHandler
						controllers.ListControllerByHTTPHandler
						controllers.ShowControllerByHTTPHandler
						controllers.UpdateControllerByHTTPHandler
						controllers.DeleteControllerByHTTPHandler
					}{
						ContextHandler:                gorest.DefaultContextHandler{ContextKey: `sub-id`},
						CreateControllerByHTTPHandler: controllers.CreateControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, http.StatusOK, `Create`)},
						ListControllerByHTTPHandler:   controllers.ListControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, http.StatusOK, `List`)},
						ShowControllerByHTTPHandler:   controllers.ShowControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, http.StatusOK, `Show`)},
						UpdateControllerByHTTPHandler: controllers.UpdateControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, http.StatusOK, `Update`)},
						DeleteControllerByHTTPHandler: controllers.DeleteControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, http.StatusOK, `Delete`)},
					},
				)
			})
			s.Before(func(t *testcase.T) { gorest.Mount(handler(t), `/subs/`, t.I(`sub-handler`).(*gorest.Handler)) })

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

		s.Test(`E2E`, func(t *testcase.T) {
			h := &gorest.Handler{}

			ch := gorest.DefaultContextHandler{ContextKey: `bookID`}
			books := gorest.NewHandler(struct {
				gorest.ContextHandler
				controllers.ShowControllerByHTTPHandler
			}{
				ContextHandler: ch,
				ShowControllerByHTTPHandler: controllers.ShowControllerByHTTPHandler{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = fmt.Fprintf(w, `%s`, ch.GetResourceID(r.Context()))
				})},
			})
			gorest.Mount(h, `/books/`, books)

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
						s.Let(`controller`, func(t *testcase.T) interface{} {
							return controllers.CreateControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, code, msg)}
						})

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not yet set`, func(s *testcase.Spec) {
						thenItWillUseTheAttachedHandler(s)
					})
				})

				s.Context(`list action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
					s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Let(`controller`, func(t *testcase.T) interface{} {
							return controllers.ListControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, code, msg)}
						})

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not yet set`, func(s *testcase.Spec) {
						thenItWillUseTheAttachedHandler(s)
					})
				})
				s.Context(`show action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
					s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Let(`controller`, func(t *testcase.T) interface{} {
							return controllers.ShowControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, code, msg)}
						})

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not yet set`, func(s *testcase.Spec) {
						thenItWillUseTheAttachedHandler(s)
					})
				})
				s.Context(`update action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodPut })
					s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Let(`controller`, func(t *testcase.T) interface{} {
							return controllers.UpdateControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, code, msg)}
						})

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not yet set`, func(s *testcase.Spec) {
						thenItWillUseTheAttachedHandler(s)
					})
				})
				s.Context(`delete action`, func(s *testcase.Spec) {
					s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodDelete })
					s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

					s.And(`action is set`, func(s *testcase.Spec) {
						s.Let(`controller`, func(t *testcase.T) interface{} {
							return controllers.DeleteControllerByHTTPHandler{Handler: NewTestControllerMockHandler(t, code, msg)}
						})

						thenItWillUseTheControllerHandler(s)
					})
					s.And(`action is not yet set`, func(s *testcase.Spec) {
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

	s.Describe(`CUSTOM / - unknown http method used`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return `CUSTOM` })
		s.Let(`path`, func(t *testcase.T) interface{} { return `/` })

		s.Before(func(t *testcase.T) {
			t.Log(`given we have no controller action defined regarding collection level operation`)
			_, ok := handler(t).LookupCollectionHandler(http.MethodGet, t.I(`path`).(string))
			require.False(t, ok)
			_, ok = handler(t).LookupCollectionHandler(http.MethodPost, t.I(`path`).(string))
			require.False(t, ok)
		})

		s.When(`nothing set to handle the request`, func(s *testcase.Spec) {
			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`a global handler is set as fallback solution`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				handler(t).Handle(`/`, NewTestControllerMockHandler(t, http.StatusTeapot, http.StatusText(http.StatusTeapot)))
			})

			s.Then(`it will use the attached`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, http.StatusTeapot, resp.Code)
				require.Equal(t, http.StatusText(http.StatusTeapot), strings.TrimSpace(resp.Body.String()))
			})
		})
	})

	s.Describe(`CUSTOM /{resourceID} - unknown http method used`, func(s *testcase.Spec) {
		s.Let(`method`, func(t *testcase.T) interface{} { return `CUSTOM` })
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

		s.When(`nothing set to handle the request`, func(s *testcase.Spec) {
			s.Then(`it will return with 404`, func(t *testcase.T) {
				require.Equal(t, http.StatusNotFound, serve(t).Code)
			})

			andWhenCustomNotFoundHandlerProvided(s)
		})

		s.When(`a global handler is set as fallback solution`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				handler(t).Handle(`/`, NewTestControllerMockHandler(t, http.StatusTeapot, http.StatusText(http.StatusTeapot)))
			})

			s.Then(`it will use the attached`, func(t *testcase.T) {
				resp := serve(t)
				require.Equal(t, http.StatusTeapot, resp.Code)
				require.Equal(t, http.StatusText(http.StatusTeapot), strings.TrimSpace(resp.Body.String()))
			})

			andWhenResourceHandlerIs(s, func(s *testcase.Spec) {})
		})
	})

	s.Describe(`#InternalServerError`, func(s *testcase.Spec) {
		const respBody = "a custom internal server error response"
		s.Before(func(t *testcase.T) {
			handler(t).InternalServerError = NewInternalServerErrorHandler(InternalServerErrorController{
				Code: http.StatusInternalServerError,
				Msg:  respBody,
			})
		})

		s.Let(`method`, func(t *testcase.T) interface{} { return http.MethodGet })
		s.Let(`path`, func(t *testcase.T) interface{} { return fmt.Sprintf(`/%s`, resourceID(t)) })

		s.When(`error occurs during context setup`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				handler(t).ContextHandler = ErrorContextHandler{Err: errors.New(`boom`)}
			})

			s.Then(`custom internal server error handler will be used`, func(t *testcase.T) {
				require.Contains(t, serve(t).Body.String(), respBody)
			})
		})

		s.When(`panic occurs during controller action`, func(s *testcase.Spec) {
			s.Let(`controller`, func(t *testcase.T) interface{} {
				return controllers.ShowControllerByHTTPHandler{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic(`boom`)
				})}
			})

			s.Then(`custom internal server error handler will be used`, func(t *testcase.T) {
				require.Contains(t, serve(t).Body.String(), respBody)
			})

			s.And(`internal server error also return with panic`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					handler(t).InternalServerError = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						panic(`boom`)
					})
				})

				s.Then(`generic internal server error is used as fallback`, func(t *testcase.T) {
					resp := serve(t)
					require.Equal(t, http.StatusInternalServerError, resp.Code)
					require.Contains(t, resp.Body.String(), http.StatusText(http.StatusInternalServerError))
				})
			})
		})
	})
}

func BenchmarkController_ServeHTTP(b *testing.B) {
	h := gorest.NewHandler(struct {
		gorest.ContextHandler
		controllers.ShowControllerByHTTPHandler
	}{
		ContextHandler:              gorest.DefaultContextHandler{ContextKey: `bench`},
		ShowControllerByHTTPHandler: controllers.ShowControllerByHTTPHandler{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})},
	})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, `/resourceID`, &bytes.Buffer{})
		h.ServeHTTP(w, r)
	}
}
