// Apollo provides `net/context`-aware middleware chaining
package apollo

import (
	"net/http"

	"golang.org/x/net/context"
)

// Handler is a context-aware interface analagous to the `net/http` http.Handler interface
// The only difference is that a context.Context is required as the first parameter in ServeHTTP.
type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

// HandlerFunc, similar to http.HandlerFunc, is an adapter to convert ordinary functions
// into handlers.
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTP calls the wrapped function h(ctx, w, r)
func (h HandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h(ctx, w, r)
}

// addsContext is an adapter that wraps a Handler and implements the http.Handler interface.
// The resulting object is used as a bridge to integrate existing `net/http` functions with
// a context-aware chain or handler.
// Internally, it is used as an onramp to the chain in Then(), and as an adapter for
// injecting non-context-aware handlers with Wrap()
type addsContext struct {
	ctx     context.Context
	handler Handler
}

// ServeHTTP calls the stored handler with the stored context, passing through the received
// HTTP Response and Request
func (a *addsContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(a.ctx, w, r)
}

// stripsContext is an adapter that wraps a http.Handler and implements the Handler interface.
// The resulting object can be used to insert a non-context-aware function into a context-aware
// chain.
// Internally, it is used to in Wrap() to inject standard handlers.  It could also be used to
// link context-aware middleware to a standard handler.
type stripsContext struct {
	handler http.Handler
}

// ServeHTTP calls the stored handler, dropping the context and passing only the HTTP ResponseWriter and Request
func (s *stripsContext) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

// Wrap allows injection of normal http.Handler middleware into an
// apollo middleware chain
// The context will be preserved and passed through intact
func Wrap(h func(http.Handler) http.Handler) Constructor {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			stubHandler := &addsContext{
				ctx:     ctx,
				handler: next,
			}
			h(stubHandler).ServeHTTP(w, r)
		})
	}
}
