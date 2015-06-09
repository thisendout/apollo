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

type AddsContext struct {
	ctx     context.Context
	handler Handler
}

func (a *AddsContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(a.ctx, w, r)
}

type StripsContext struct {
	handler http.Handler
}

func (s *StripsContext) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
