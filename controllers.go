package gorest

import (
	"fmt"
	"net/http"
)

func AsListController(i interface{}) ListController {
	switch i := i.(type) {
	case http.Handler:
		return httpHandlerAsListController{Handler: i}
	default:
		panic(fmt.Sprintf(`unknown type: %T`, i))
	}
}

func AsCreateController(i interface{}) CreateController {
	switch i := i.(type) {
	case http.Handler:
		return httpHandlerAsCreateController{Handler: i}
	default:
		panic(fmt.Sprintf(`unknown type: %T`, i))
	}
}

func AsShowController(i interface{}) ShowController {
	switch i := i.(type) {
	case http.Handler:
		return httpHandlerAsShowController{Handler: i}
	default:
		panic(fmt.Sprintf(`unknown type: %T`, i))
	}
}

func AsUpdateController(i interface{}) UpdateController {
	switch i := i.(type) {
	case http.Handler:
		return httpHandlerAsUpdateController{Handler: i}
	default:
		panic(fmt.Sprintf(`unknown type: %T`, i))
	}
}

func AsDeleteController(i interface{}) DeleteController {
	switch i := i.(type) {
	case http.Handler:
		return httpHandlerAsDeleteController{Handler: i}
	default:
		panic(fmt.Sprintf(`unknown type: %T`, i))
	}
}

//--------------------------------------------------- http.Handler ---------------------------------------------------//

type httpHandlerAsListController struct{ http.Handler }

func (ctrl httpHandlerAsListController) List(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}

type httpHandlerAsCreateController struct{ http.Handler }

func (ctrl httpHandlerAsCreateController) Create(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}

type httpHandlerAsShowController struct{ http.Handler }

func (ctrl httpHandlerAsShowController) Show(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}

type httpHandlerAsUpdateController struct{ http.Handler }

func (ctrl httpHandlerAsUpdateController) Update(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}

type httpHandlerAsDeleteController struct{ http.Handler }

func (ctrl httpHandlerAsDeleteController) Delete(w http.ResponseWriter, r *http.Request) {
	ctrl.Handler.ServeHTTP(w, r)
}
