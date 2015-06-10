// Apollo provides `net/context`-aware middleware chaining
package apollo

import (
	"io/ioutil"
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
	adapter := addsContext{
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
	adapter := stripsContext{http.HandlerFunc(handlerZero)}
	assert.NotPanics(func() {
		ctx := context.Background()
		r, _ := http.NewRequest("GET", "http://github.com/", nil)
		w := httptest.NewRecorder()

		adapter.ServeHTTP(ctx, w, r)
		assert.Equal(w.Code, 200)
		assert.Equal(w.Body.String(), "h0\n")
	})
}

func TestWrap(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		con := Wrap(middleZero)
		assert.IsType(con, *new(Constructor))
	})
}

func TestWrapChains(t *testing.T) {
	assert := assert.New(t)
	ctx := NewTestContext(context.Background(), 10)
	value, _ := FromContext(ctx)
	assert.Equal(value, 10)

	chain := New(middleOne, Wrap(middleZero), middleTwo).With(ctx).ThenFunc(handlerContext)

	ts := httptest.NewServer(chain)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	assert.Equal(200, res.StatusCode)
	assert.Equal("m1\nm0\nm2\n10\n", string(body))
}
