package apollo

import (
	"net/http"

	"golang.org/x/net/context"
)

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h HandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h(ctx, w, r)
}

type addsContext struct {
	ctx     context.Context
	handler Handler
}

func (a *addsContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(a.ctx, w, r)
}

type stripsContext struct {
	handler http.Handler
}

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
