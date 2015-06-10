// Apollo provides `net/context`-aware middleware chaining
package apollo

import (
	"net/http"
	"strconv"

	"golang.org/x/net/context"
)

func handlerZero(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("h0\n"))
}

func handlerOne(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("h1\n"))
}

func handlerContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if value, ok := FromContext(ctx); ok {
		contents := strconv.Itoa(value) + "\n"
		w.Write([]byte(contents))
	}
}

func middleZero(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("m0\n"))
		h.ServeHTTP(w, r)
	})
}

func middleOne(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("m1\n"))
		h.ServeHTTP(ctx, w, r)
	})
}

func middleTwo(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("m2\n"))
		h.ServeHTTP(ctx, w, r)
	})
}

// TestContext
type key int

const testKey key = 0

func NewTestContext(ctx context.Context, dummy int) context.Context {
	return context.WithValue(ctx, testKey, dummy)
}

func FromContext(ctx context.Context) (int, bool) {
	dummy, ok := ctx.Value(testKey).(int)
	return dummy, ok
}
