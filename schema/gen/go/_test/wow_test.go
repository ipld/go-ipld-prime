package whee

import (
	"testing"

	"github.com/polydawn/refmt/tok"
	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/fluent"
)

// TokenSourceBucket acts like a TokenSource by yielding tokens from a pre-made
// slice; and also keeps track of how far it's been read.
type TokenSourceBucket struct {
	tokens []tok.Token
	read   int
}

func (tb *TokenSourceBucket) Step(yield *tok.Token) (done bool, err error) {
	*yield = tb.tokens[tb.read]
	tb.read++
	return tb.read > len(tb.tokens), nil
}

func TestScalarUnmarshal(t *testing.T) {
	t.Run("string node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TString, Str: "zooooom"},
		}}
		nb := Strang__NodeBuilder{}
		n, err := encoding.Unmarshal(nb, tb)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, fluent.WrapNode(n).AsString(), ShouldEqual, "zooooom")
		Wish(t, tb.read, ShouldEqual, 1)
	})
}
