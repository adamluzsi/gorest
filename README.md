[![Build Status](https://travis-ci.org/adamluzsi/gorest.svg?branch=master)](https://travis-ci.org/adamluzsi/gorest)
[![GoDoc](https://godoc.org/github.com/adamluzsi/gorest?status.png)](https://godoc.org/github.com/adamluzsi/gorest)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamluzsi/gorest)](https://goreportcard.com/report/github.com/adamluzsi/gorest)
[![codecov](https://codecov.io/gh/adamluzsi/gorest/branch/master/graph/badge.svg)](https://codecov.io/gh/adamluzsi/gorest)
# gorest

gorest is a minimalist approach to build restful API designs through composition.

## What problem it solves?

The reason I made this package, because when I work with restful API design,
I prefer the pipeline pattern to create steps for each resource validation.
This convention led me to the pattern where I have a `http.Handler` that act as a controller,
a `http.ServeMux` that composite the controllers, and middlewares which are setup the resource values in the `http.Request`'s context.

For example if I have the following path: `/stores/mystore/foods/cucumber`
then the `stores` resource `mystore` resource id first being validated,
and store objects that represent `mystore` set in the context.
After that, the same will happen with the `foods` resource `cucumber` id,
but for this, I need the `mystore` instance as well from the context.

This pattern worked nicely so far, as it allows guard clauses for handling cases
when a resource is not found or should not be returned to the requester.
But it was also kinda boilerplate to setup this with `http.ServeMux`.

Using a `router` that would allow me to have the path params would be less efficient,
as my controllers under a certain resource use the assumption that the resource exists and can be used already,
to remove a lot of repetition from each controller code.

## How does it solve it?

The above-mentioned problem is solved by introducing a convention,
by having a well-tested controller, that can have actions such as `List`, `Create`, `Show` and so on.

Also, the controller has the convention to take an object that knows what to do with the resourceID it received
and how to setup the request context to include this resource instance.

## Example

* [godoc examples](https://godoc.org/github.com/adamluzsi/gorest#pkg-examples)

### Using a controller struct

A controller that only implement certain actions

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func main() {
    mux := http.NewServeMux()
    gorest.Mount(mux, `/xys/`, gorest.NewHandler(XYController{}))

    if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err.Error())
	}
}

type XYController struct{}

func (ctrl XYController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `list`)
}

func (ctrl XYController) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	return context.WithValue(ctx, `id`, resourceID), true, nil
}

func (ctrl XYController) Show(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `show:%s`, r.Context().Value(`id`))
}


```

Or if the controller implements all resource, then the handler will use the other methods as well.
It is documented by the `gorest.Controller` interface.

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

type XYController struct{}

func (ctrl XYController) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	return context.WithValue(ctx, `id`, resourceID), true, nil
}

func (ctrl XYController) Create(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `create`)
}

func (ctrl XYController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `list`)
}

func (ctrl XYController) Show(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `show:%s`, r.Context().Value(`id`))
}

func (ctrl XYController) Update(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `update:%s`, r.Context().Value(`id`))
}

func (ctrl XYController) Delete(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `delete:%s`, r.Context().Value(`id`))
}

func (ctrl XYController) NotFound(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `not-found`)
}

func (ctrl XYController) InternalServerError(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `internal-server-error`)
}
```

### Using gorest.Handler directly

```go
teapotController := &gorest.Controller{
    ContextHandler: TeapotResourceHandler{},
    Show:           TeapotShowHandler{},
}

// GET /teapots/:teapotID/droplets/:dropletID -> SHOW
teapotController.Mount(`droplets`, &gorest.Controller{
    ContextHandler: DropletResourceHandler{},
    Show:           DropletShowHandler{},
})

mux := http.NewServeMux()
gorest.Mount(mux, `/teapots`, teapotController)
```

## Q&A

### And what to do when I need an Endpoints that is outside of the restful path convention?

You can use [Handle](https://godoc.org/github.com/adamluzsi/gorest#Controller.Handle) to mount a `http.Handler` to any `http.ServeMux` supported pattern,

### How hard is the testing?

I usually use BDD approach and setup the testing context,
so I don't have to care too much about implementation details from that perspective,
but if you need to mock, the `ContextHandler` can help a lot, as you can test the implementation separately,
and supply a mock for the controller action tests.

IMO, and E2E test with context setup is less of a risk, but that is just a matter of testing taste.
