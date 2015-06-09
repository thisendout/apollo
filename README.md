Apollo
======

Apollo is a middleware-chaining helper for Golang web applications using google's `net/context` package.  Apollo is a fork of [Alice](https://github.com/justinas/alice), modified to support passing the `ctx context.Context` param through middleware and HTTP handlers.

Apollo is meant to chain handler functions with this signature:
```
func (context.Context, http.ResponseWriter, *http.Request)
```

Relevant and influential articles:
 * https://blog.golang.org/context
 * https://joeshaw.org/net-context-and-http-handler/
 * https://elithrar.github.io/article/map-string-interface/
 * http://www.alexedwards.net/blog/making-and-using-middleware
 * http://laicos.com/writing-handsome-golang-middleware/

# Usage

```
apollo.New(Middleware1, Middlware2, Middleware3).With(ctx).Then(App)
```

# Motivation

Given these handlers:
```
func HandlerOne(w http.ResponseWriter, r *http.Request) {}
//FIXME
func HandlerTwo(h){
  return func(w http.ResponseWriter, r *http.Request) {
    h(w, r)
  }
}
```

We can serve them using the following:
```
http.HandleFunc("/one", HandlerOne)
// or http.Handle("/one", http.HandlerFunc(HandlerOne))
http.Handle("/two", HandlerTwo)
```

However, given a handler that expects a `net/context`:
```
func HandlerAlpha(ctx context.Context, w http.ResponseWriter, r *http.Request) {}
```

We would need to create a wrapper along the lines of:
```
func withContext(ctx context.Context, fn func(context.Context, http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fn(ctx, w, r)
  }
}
```
and serve with:
```
ctx := context.Background()
http.Handle("/alpha", withContext(ctx, HandlerAlpha))
```

With this pattern, we can build nested middleware/handler calls that can be used with any `net/http` compatible router/mux. However, we can't implement an Alice-style chain because we no longer conform to the http.Handler-chaining interface that Alice expects.

Apollo enables context-aware middleware and handlers by passing the context through the chain as a required parameter.
