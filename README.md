Apollo
======
[![GoDoc](https://godoc.org/github.com/cyclopsci/apollo?status.svg)](https://godoc.org/github.com/cyclopsci/apollo)

Apollo is a middleware-chaining helper for Golang web applications using google's `net/context` package.  Apollo is a fork of [Alice](https://github.com/justinas/alice), modified to support passing the `ctx context.Context` param through middleware and HTTP handlers.

Apollo is meant to chain handler functions with this signature:
```go
func (context.Context, http.ResponseWriter, *http.Request)
```

Relevant and influential articles:
 * https://blog.golang.org/context
 * https://joeshaw.org/net-context-and-http-handler/
 * https://elithrar.github.io/article/map-string-interface/
 * http://www.alexedwards.net/blog/making-and-using-middleware
 * http://laicos.com/writing-handsome-golang-middleware/
 * http://nicolasmerouze.com/share-values-between-middlewares-context-golang/
 * https://elithrar.github.io/article/custom-handlers-avoiding-globals/
 * http://www.jerf.org/iri/post/2929

Apollo was originally hosted at [cyclopsci/apollo](https://github.com/cyclopsci/apollo), which is now a fork of this repo.

# Usage

```go
apollo.New(Middleware1, Middlware2, Middleware3).With(ctx).Then(App)
```

# Integration with http.Handler middleware

Apollo provides a `Wrap` function to inject normal http.Handler-based middleware into the chain.  The context will skip over the injected middleware and pass unharmed to the next context-aware handler in the chain.
```go
apollo.New(ContextMW1, apollo.Wrap(NormalMiddlware), ContextMW2).With(ctx).Then(App)
```
# Motivation

Given a handler:
```go
func HandlerOne(w http.ResponseWriter, r *http.Request) {}
```

We can serve it using the following:
```go
http.HandleFunc("/one", HandlerOne)
// or http.Handle("/one", http.HandlerFunc(HandlerOne))
```

However, given a handler that expects a `net/context`:
```go
func HandlerAlpha(ctx context.Context, w http.ResponseWriter, r *http.Request) {}
```

We would need to create a wrapper along the lines of:
```go
func withContext(ctx context.Context, fn func(context.Context, http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fn(ctx, w, r)
  }
}
```
and serve with:
```go
ctx := context.Background()
http.Handle("/alpha", withContext(ctx, HandlerAlpha))
```

With this pattern, we can build nested middleware/handler calls that can be used with any `net/http` compatible router/mux. However, we can't use Alice for chaining because we no longer conform to the http.Handler interface that Alice expects.

Apollo enables Alice-style chaining of context-aware middleware and handlers.

