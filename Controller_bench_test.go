package gorest_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamluzsi/gorest"
)

func BenchmarkController_ServeHTTP(b *testing.B) {
	ctrl := &gorest.Controller{Show: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, `/resourceID`, &bytes.Buffer{})
		ctrl.ServeHTTP(w, r)
	}
}
