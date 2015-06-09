package apollo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestHandlerFunc(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		ctx := context.Background()
		r, _ := http.NewRequest("GET", "http://github.com/", nil)
		w := httptest.NewRecorder()

		handler := HandlerFunc(handlerOne)
		assert.Implements((*Handler)(nil), handler)

		handler.ServeHTTP(ctx, w, r)
		assert.Equal(w.Code, 200)
		assert.Equal(w.Body.String(), "h1\n")
	})
}

func TestAddsContextServe(t *testing.T) {
	assert := assert.New(t)
	adapter := AddsContext{
		ctx:     context.Background(),
		handler: HandlerFunc(handlerOne),
	}
	assert.NotPanics(func() {
		r, _ := http.NewRequest("GET", "http://github.com/", nil)
		w := httptest.NewRecorder()

		adapter.ServeHTTP(w, r)
		assert.Equal(w.Code, 200)
		assert.Equal(w.Body.String(), "h1\n")
	})
}

func TestStripsContextServe(t *testing.T) {
	assert := assert.New(t)
	adapter := StripsContext{http.HandlerFunc(handlerZero)}
	assert.NotPanics(func() {
		ctx := context.Background()
		r, _ := http.NewRequest("GET", "http://github.com/", nil)
		w := httptest.NewRecorder()

		adapter.ServeHTTP(ctx, w, r)
		assert.Equal(w.Code, 200)
		assert.Equal(w.Body.String(), "h0\n")
	})
}
