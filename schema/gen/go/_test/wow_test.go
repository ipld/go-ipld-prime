package whee

import (
	"testing"

	"github.com/polydawn/refmt/tok"
	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
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

func plz(n ipld.Node, e error) ipld.Node {
	if e != nil {
		panic(e)
	}
	return n
}

func TestScalarUnmarshal(t *testing.T) {
	t.Run("string node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TString, Str: "zooooom"},
		}}
		nb := String__NodeBuilder{}
		n, err := encoding.Unmarshal(nb, tb)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, fluent.WrapNode(n).AsString(), ShouldEqual, "zooooom")
		Wish(t, tb.read, ShouldEqual, 1)
	})
}

// n.b. testing unmarshal is something different (this doesn't exercise
// representation node/nodebuilder, just the semantic/type level one).
func TestStructBuilder(t *testing.T) {
	t.Run("stroct", func(t *testing.T) {
		t.Run("all fields set", func(t *testing.T) {
			mb, err := Stroct__NodeBuilder{}.CreateMap()
			Require(t, err, ShouldEqual, nil)
			mb.Insert(ipldfree.String("f1"), plz(String__NodeBuilder{}.CreateString("a")))
			mb.Insert(ipldfree.String("f2"), plz(String__NodeBuilder{}.CreateString("b")))
			mb.Insert(ipldfree.String("f3"), plz(String__NodeBuilder{}.CreateString("c")))
			mb.Insert(ipldfree.String("f4"), plz(String__NodeBuilder{}.CreateString("d")))
			n, err := mb.Build()

			Wish(t, err, ShouldEqual, nil)
			Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
			Wish(t, plz(n.LookupString("f1")), ShouldEqual, plz(String__NodeBuilder{}.CreateString("a")))
		})
	})
}

/*
soon...


func TestStructUnmarshal(t *testing.T) {
	t.Run("stroct", func(t *testing.T) {
		t.Run("all fields set", func(t *testing.T) {
			tb := &TokenSourceBucket{tokens: []tok.Token{
				{Type: tok.TMapOpen, Length: 4},
				{Type: tok.TString, Str: "f1"}, {Type: tok.TString, Str: "a"},
				{Type: tok.TString, Str: "f2"}, {Type: tok.TString, Str: "b"},
				{Type: tok.TString, Str: "f3"}, {Type: tok.TString, Str: "c"},
				{Type: tok.TString, Str: "f4"}, {Type: tok.TString, Str: "d"},
				{Type: tok.TMapClose},
			}}
			nb := Stroct__NodeBuilder{}
			n, err := encoding.Unmarshal(nb, tb)
			// ... asserts ...
		})
	})
}

*/
