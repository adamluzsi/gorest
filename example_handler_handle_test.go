package gorest_test

import (
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleHandler_Handle() {
	handler := gorest.NewHandler(TestController{})
	var myCustomOperationHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	//
	// this will register a handler for /{resource-id}/my-custom-operation
	handler.Handle(`/my-custom-operation`, myCustomOperationHandler)
}
