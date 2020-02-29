package gorest_test

import (
	"net/http"

	"github.com/adamluzsi/gorest"
)

func ExampleHandler_Handle() {
	handler := gorest.NewHandler(TestController{})
	var myCustomOperationHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler.Handle(`/my-custom-operation`, myCustomOperationHandler)
}
