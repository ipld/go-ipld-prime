package tests

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

func TestScalarUnmarshal(t *testing.T, nb ipld.NodeBuilder) {
	t.Run("null node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TNull},
		}}
		n, err := encoding.Unmarshal(nb, tb)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.ReprKind_Null)
		Wish(t, tb.read, ShouldEqual, 1)
	})
	t.Run("int node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TInt, Int: 1400},
		}}
		n, err := encoding.Unmarshal(nb, tb)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.ReprKind_Int)
		Wish(t, fluent.WrapNode(n).AsInt(), ShouldEqual, 1400)
		Wish(t, tb.read, ShouldEqual, 1)
	})
	t.Run("string node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TString, Str: "zooooom"},
		}}
		n, err := encoding.Unmarshal(nb, tb)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n.Kind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, fluent.WrapNode(n).AsString(), ShouldEqual, "zooooom")
		Wish(t, tb.read, ShouldEqual, 1)
	})
}

func TestRecursiveUnmarshal(t *testing.T, nb ipld.NodeBuilder) {
	t.Run("short list node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TArrOpen, Length: 1},
			{Type: tok.TString, Str: "asdf"},
			{Type: tok.TArrClose},
		}}
		n, err := encoding.Unmarshal(nb, tb)
		Require(t, err, ShouldEqual, nil)
		Require(t, n.Kind(), ShouldEqual, ipld.ReprKind_List)
		Require(t, n.Length(), ShouldEqual, 1)
		Wish(t, fluent.WrapNode(n).TraverseIndex(0).Kind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, fluent.WrapNode(n).TraverseIndex(0).AsString(), ShouldEqual, "asdf")
		Wish(t, tb.read, ShouldEqual, 3)
	})
	t.Run("nested list node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TArrOpen, Length: 2},
			{Type: tok.TArrOpen, Length: 1},
			{Type: tok.TString, Str: "asdf"},
			{Type: tok.TArrClose},
			{Type: tok.TString, Str: "quux"},
			{Type: tok.TArrClose},
		}}
		n, err := encoding.Unmarshal(nb, tb)
		Require(t, err, ShouldEqual, nil)
		Require(t, n.Kind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, n.Length(), ShouldEqual, 2)
		Require(t, fluent.WrapNode(n).TraverseIndex(0).Kind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, fluent.WrapNode(n).TraverseIndex(0).Length(), ShouldEqual, 1)
		Wish(t, fluent.WrapNode(n).TraverseIndex(0).TraverseIndex(0).Kind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, fluent.WrapNode(n).TraverseIndex(0).TraverseIndex(0).AsString(), ShouldEqual, "asdf")
		Wish(t, tb.read, ShouldEqual, 6)
	})
	t.Run("short map node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TMapOpen, Length: 1},
			{Type: tok.TString, Str: "asdf"},
			{Type: tok.TString, Str: "zomzom"},
			{Type: tok.TMapClose},
		}}
		n, err := encoding.Unmarshal(nb, tb)
		Require(t, err, ShouldEqual, nil)
		Require(t, n.Kind(), ShouldEqual, ipld.ReprKind_Map)
		Require(t, n.Length(), ShouldEqual, 1)
		Require(t, fluent.WrapNode(n).KeysImmediate(), ShouldEqual, []string{"asdf"})
		Wish(t, fluent.WrapNode(n).TraverseField("asdf").Kind(), ShouldEqual, ipld.ReprKind_String)
		Wish(t, fluent.WrapNode(n).TraverseField("asdf").AsString(), ShouldEqual, "zomzom")
		Wish(t, tb.read, ShouldEqual, 4)
	})
	t.Run("nested map node", func(t *testing.T) {
		tb := &TokenSourceBucket{tokens: []tok.Token{
			{Type: tok.TMapOpen, Length: 2},
			{Type: tok.TString, Str: "asdf"},
			{Type: tok.TMapOpen, Length: 1},
			{Type: tok.TString, Str: "awoo"},
			{Type: tok.TString, Str: "gah"},
			{Type: tok.TMapClose},
			{Type: tok.TString, Str: "zyzzy"},
			{Type: tok.TInt, Int: 9},
			{Type: tok.TMapClose},
		}}
		n, err := encoding.Unmarshal(nb, tb)
		Require(t, err, ShouldEqual, nil)
		Require(t, n.Kind(), ShouldEqual, ipld.ReprKind_Map)
		Require(t, n.Length(), ShouldEqual, 2)
		Require(t, fluent.WrapNode(n).KeysImmediate(), ShouldEqual, []string{"asdf", "zyzzy"})
		Wish(t, fluent.WrapNode(n).TraverseField("asdf").Kind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, fluent.WrapNode(n).TraverseField("asdf").TraverseField("awoo").AsString(), ShouldEqual, "gah")
		Wish(t, fluent.WrapNode(n).TraverseField("zyzzy").AsInt(), ShouldEqual, 9)
		Wish(t, tb.read, ShouldEqual, 9)
	})
}
