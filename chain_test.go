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

func TestDefaultNew(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		chain := New()
		assert.Len(chain.constructors, 0)
		assert.Equal(chain.context, context.Background())
	})
}

func TestThenNoMiddleware(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		chain := New()
		final := chain.Then(HandlerFunc(handlerOne))
		assert.Implements((*http.Handler)(nil), final)
	})
}

func TestThenFuncNoMiddleware(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		chain := New()
		final := chain.ThenFunc(handlerOne)
		assert.Implements((*http.Handler)(nil), final)
	})
}

func TestThenNil(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		final := New().Then(nil)
		assert.Implements((*http.Handler)(nil), final)
	})
}

func TestThenFuncNil(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		final := New().ThenFunc(nil)
		assert.Implements((*http.Handler)(nil), final)
	})
}

func TestAppend(t *testing.T) {
	assert := assert.New(t)
	chain := New(middleOne)
	newChain := chain.Append(middleTwo)
	assert.Len(chain.constructors, 1)
	assert.Len(newChain.constructors, 2)
}

func TestAppendContext(t *testing.T) {
	assert := assert.New(t)
	chain := New(middleOne)
	newChain := chain.Append(middleTwo)
	assert.Equal(chain.context, context.Background())
	assert.Equal(newChain.context, context.Background())
}

func TestAppendAfterWith(t *testing.T) {
	assert := assert.New(t)
	ctx := NewTestContext(context.Background(), 10)

	chain := New(middleOne)
	withChain := chain.With(ctx)
	newChain := withChain.Append(middleTwo)
	assert.Equal(chain.context, context.Background())
	assert.Equal(withChain.context, ctx)
	assert.Equal(newChain.context, ctx)
}

func TestWithAfterAppend(t *testing.T) {
	assert := assert.New(t)
	ctx := NewTestContext(context.Background(), 10)

	chain := New(middleOne)
	newChain := chain.Append(middleTwo)
	withChain := newChain.With(ctx)
	assert.Equal(chain.context, context.Background())
	assert.Equal(newChain.context, context.Background())
	assert.Equal(withChain.context, ctx)
}

func TestWithInPlace(t *testing.T) {
	assert := assert.New(t)
	ctx := NewTestContext(context.Background(), 10)

	chain := New(middleOne)
	chain.With(ctx)
	newChain := chain.Append(middleTwo)
	assert.Equal(chain.context, context.Background())
	assert.Equal(newChain.context, context.Background())
}

func TestWithReassigned(t *testing.T) {
	assert := assert.New(t)
	ctx := NewTestContext(context.Background(), 10)

	chain := New(middleOne)
	chain = chain.With(ctx)
	newChain := chain.Append(middleTwo)
	assert.Equal(chain.context, ctx)
	assert.Equal(newChain.context, ctx)
}

func TestWithMultiples(t *testing.T) {
	assert := assert.New(t)
	ctx := NewTestContext(context.Background(), 10)

	chain := New(middleOne)
	assert.Equal(chain.context, context.Background())
	chain = chain.With(ctx)
	assert.Equal(chain.context, ctx)
	chain = chain.With(context.TODO()).With(ctx)
	assert.Equal(chain.context, ctx)
}

func TestChains(t *testing.T) {
	assert := assert.New(t)
	ctx := NewTestContext(context.Background(), 10)
	value, _ := FromContext(ctx)
	assert.Equal(value, 10)

	chain := New(middleOne, middleTwo).With(ctx).ThenFunc(handlerContext)

	ts := httptest.NewServer(chain)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	assert.Equal(res.StatusCode, 200)
	assert.Equal(string(body), "m1\nm2\n10\n")
}
