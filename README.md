[![Build Status](https://travis-ci.org/adamluzsi/gorest.svg?branch=master)](https://travis-ci.org/adamluzsi/gorest)
[![GoDoc](https://godoc.org/github.com/adamluzsi/gorest?status.png)](https://godoc.org/github.com/adamluzsi/gorest)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamluzsi/gorest)](https://goreportcard.com/report/github.com/adamluzsi/gorest)
[![codecov](https://codecov.io/gh/adamluzsi/gorest/branch/master/graph/badge.svg)](https://codecov.io/gh/adamluzsi/gorest)
# gorest

`gorest` is a minimalist approach to build restful API designs through composition.

## What problem it solves?

The reason I made this package, because when I design restful APIs,
I prefer to decouple the resource operations from the resource retrieve aspect.
The reason for this is that when I need to represent relationship of resources,
I often found myself testing the same logic over and over again in each `http.Handler`
about which input value is used used to retrieve a certain resource.
This then caused leaky abstraction between the `http.Handler` implementations,
since the handler had to know about how a resource is specified on the API,
and also how to retrieve this resource.

A simple example to this is, imagine you have `users`, `organizations` and `permissions` collections.
Given your API provides `/users/...` path where the API user can retrieve information about other users.
Also the API requester may retrieve information about a certain user publicly listed organizations under `/users/{user-id}/organizations`.

Now if you build your handler in a way that your `organizations` collection handler is aware that user-id must be a path parameter,
and must be used to retrieve the user entity, you couple the two heavily together.
You also bind the knowledge that the API requester must be checked if the API requester has permission to know this information.

Now imagine if your API need a `/organizations` collection that represent the API requester's organization.
Maybe a `/organizations/{organization-id}/permissions` as well.
Now you need to somehow inject the user-id in a way that it will be compatible with your `organizations` collection handler implementation.

I prefer the pipeline pattern to create steps where I represent such aspects as who is the user in the current context.
This allow me to reduce the need to test this responsibility in all the collection handler that depends on a resource existence.
In the case of the `organizations` handler, this would be the dependency on a `current context's user`. 

This convention led me to the pattern where I have a `http.Handler` that act as a controller,
a `http.ServeMux` that composite the controllers, and a couple of middleware.
Using a middleware allowed me to do permission validations 
and resource retrieval based on a defined input parameter.
Then this resource can be stored in the context,
and in the collection handler tests, I can define these resources as dependencies 
that must be present in the request context in order to use the collection handler.

This pattern worked nicely so far, as it allows guard clauses for handling cases
when a resource is not found or should not be returned to the requester.
As a side effect, it caused a lot of boilerplate while I used `http.ServeMux` purely. 

Using a `router` that would allow me to have the path params would be less efficient,
as my controllers under a certain resource use the assumption that the resource exists and can be used already,
to remove a lot of repetition from each controller code.

## How does this package solve the problem?

The above-mentioned problem is solved by introducing a convention.
By having a well-tested controller package, that can have actions such as `List`, `Create`, `Show` and so on
it becomes easier to focus on the operation aspects and the retrieval aspects.

The resource retrieval and visibility validation is solved by having a controller function ([ContextWithResource](https://godoc.org/github.com/adamluzsi/gorest#ContextHandler)),
that focus only on this, and the rest of the resource oriented actions like `Show`, `Update` and `Delete`
no longer have to cover in they specification this.

More about this between the examples.

## Resource Oriented Design

The architectural style of REST was introduced primarily to work well with HTTP/1.1.
It also helps to reduce the learning curve a developer need to do in order to understand how to use a new API.
Its core principle is to define named resources that can be manipulated using a small number of methods.
The resources and methods are known as nouns and verbs of APIs.
With the HTTP protocol, the resource names naturally map to URLs,
and methods naturally map to HTTP methods POST, GET, PUT, PATCH, and DELETE.
This results in much fewer things to learn, since developers can focus on the resources and their relationship,
and assume that they have the same small number of standard methods.

## What is a REST API?

A REST API is modeled as collections of individually-addressable resources (the nouns of the API).
Resources are referenced with their resource names and manipulated via a small set of methods (also known as verbs or operations).

Standard methods for REST APIs (also known as REST methods) are List, Show, Create, Update, and Delete.

You can create a controller simply as:

```go
package myhttpapi

import (
	"context"
	"net/http"

	"github.com/adamluzsi/gorest"
)

func NewMyCollectionHandler() http.Handler {
    return gorest.NewHandler(MyCollectionController{})
}

type MyCollectionController struct{}

func (ctrl MyCollectionController) List(w http.ResponseWriter, r *http.Request) {}
```

Custom methods (also known as custom verbs or custom operations) are also available to API designers
for functionality that doesn't easily map to one of the standard methods,
such as database transactions.

You can apply them using the multiplexer interface of the `gorest.Handler`

```go
handler := gorest.NewHandler(TestController{})
var myCustomOperationHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
handler.Handle(`/my-custom-operation`, myCustomOperationHandler)
```

Note: Custom verbs does not mean creating custom HTTP verbs to support custom methods.
For HTTP-based APIs, they simply map to the most suitable HTTP verbs.

## Design flow

we suggests taking the following steps when designing resource-oriented APIs:
- Determine what types of resources an API provides.
- Determine the relationships between resources.
- Decide the resource name schemes based on types and relationships.
- Decide the resource schemas.
- Attach minimum set of methods to resources.

### Resources

A resource-oriented API is generally modeled as a resource hierarchy, 
where each node is either a simple resource or a collection resource.
For convenience, they are often called as a resource and a collection, respectively.
A `gorest.Controller` represents operations on a collection resources.
Not all function must be implemented, if the given collection doesn't need it.   

A collection contains a list of resources of the same type. 
For example, a user has a collection of contacts.
A resource has some state and zero or more sub-resources. 
Each sub-resource can be either a simple resource or a collection resource.
For example, an API may have a collection of users, each user has a collection of messages, a profile resource, and several setting resources.

While there is some conceptual alignment between storage systems and REST APIs, 
a service with a resource-oriented API is not necessarily a database,
and has enormous flexibility in how it interprets resources and methods.
For example, creating a calendar event (resource) may create additional events for attendees,
send email invitations to attendees, reserve conference rooms, and update video conference schedules.

### Methods

The key characteristic of a resource-oriented API is that it emphasizes resources (data model) over the methods performed on the resources (functionality).
A typical resource-oriented API exposes a large number of resources with a small number of methods.
The methods can be either the standard methods or custom methods. 
For `gorest`, the standard methods are: List, Show, Create, Update, and Delete.

Where API functionality naturally maps to one of the standard methods, that method should be used in the API design.
For functionality that does not naturally map to one of the standard methods, custom methods may be used.
Custom methods offer the same design freedom as traditional RPC APIs, 
which can be used to implement common programming patterns, such as database transactions or data analysis.

## Examples

[You can find examples regarding the usage of the package between the godoc examples.](https://godoc.org/github.com/adamluzsi/gorest#pkg-examples)

The following sections present a few examples on how to use a `gorest` controller.

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
    gorest.Mount(mux, `/my-collection-id-name-in-plural/`, gorest.NewHandler(MyCollectionController{}))

    if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err.Error())
	}
}

type MyCollectionController struct{}

func (ctrl MyCollectionController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `list`)
}

func (ctrl MyCollectionController) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	return context.WithValue(ctx, `id`, resourceID), true, nil
}

func (ctrl MyCollectionController) Show(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `show:%s`, r.Context().Value(`id`))
}
```

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamluzsi/gorest"
)

type MyCollectionController struct{}

func (ctrl MyCollectionController) ContextWithResource(ctx context.Context, resourceID string) (newContext context.Context, found bool, err error) {
	return context.WithValue(ctx, `id`, resourceID), true, nil
}

func (ctrl MyCollectionController) Create(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `create`)
}

func (ctrl MyCollectionController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `list`)
}

func (ctrl MyCollectionController) Show(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `show:%s`, r.Context().Value(`id`))
}

func (ctrl MyCollectionController) Update(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `update:%s`, r.Context().Value(`id`))
}

func (ctrl MyCollectionController) Delete(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `delete:%s`, r.Context().Value(`id`))
}

func (ctrl MyCollectionController) NotFound(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `not-found`)
}

func (ctrl MyCollectionController) InternalServerError(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, `internal-server-error`)
}
```