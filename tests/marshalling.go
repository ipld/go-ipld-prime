package tests

import (
	"testing"

	"github.com/polydawn/refmt/tok"
	. "github.com/warpfork/go-wish"
)

// TokenBucket acts as a TokenSink; you can dump data into it and then
// do test assertions on it with go-wish.
type TokenBucket []tok.Token

func (tb *TokenBucket) Step(consume *tok.Token) (done bool, err error) {
	if tb == nil {
		*tb = make(TokenBucket, 0, 10)
	}
	*tb = append(*tb, *consume)
	return false, nil
}

// This should really be a utility func in refmt tok.  -.-
func (tb TokenBucket) Minimalize() TokenBucket {
	for i, v := range tb {
		switch v.Type {
		case tok.TMapOpen:
			tb[i] = tok.Token{Type: v.Type, Length: v.Length, Tagged: v.Tagged, Tag: v.Tag}
		case tok.TMapClose:
			tb[i] = tok.Token{Type: v.Type}
		case tok.TArrOpen:
			tb[i] = tok.Token{Type: v.Type, Length: v.Length, Tagged: v.Tagged, Tag: v.Tag}
		case tok.TArrClose:
			tb[i] = tok.Token{Type: v.Type}
		case tok.TNull:
			tb[i] = tok.Token{Type: v.Type, Tagged: v.Tagged, Tag: v.Tag}
		case tok.TString:
			tb[i] = tok.Token{Type: v.Type, Str: v.Str, Tagged: v.Tagged, Tag: v.Tag}
		case tok.TBytes:
			tb[i] = tok.Token{Type: v.Type, Bytes: v.Bytes, Tagged: v.Tagged, Tag: v.Tag}
		case tok.TBool:
			tb[i] = tok.Token{Type: v.Type, Bool: v.Bool, Tagged: v.Tagged, Tag: v.Tag}
		case tok.TInt:
			tb[i] = tok.Token{Type: v.Type, Int: v.Int, Tagged: v.Tagged, Tag: v.Tag}
		case tok.TUint:
			tb[i] = tok.Token{Type: v.Type, Uint: v.Uint, Tagged: v.Tagged, Tag: v.Tag}
		case tok.TFloat64:
			tb[i] = tok.Token{Type: v.Type, Float64: v.Float64, Tagged: v.Tagged, Tag: v.Tag}
		}
	}
	return tb
}

func TestScalarMarshal(t *testing.T, newNode MutableNodeFactory) {
	t.Run("null node", func(t *testing.T) {
		n0 := newNode()
		n0.SetNull()
		var tb TokenBucket
		n0.PushTokens(&tb)
		Wish(t, tb, ShouldEqual, TokenBucket{
			{Type: tok.TNull},
		})
	})
}

func TestRecursiveMarshal(t *testing.T, newNode MutableNodeFactory) {
	t.Run("short list node", func(t *testing.T) {
		n0 := newNode()
		n00 := newNode()
		n00.SetString("asdf")
		n0.SetIndex(0, n00)
		var tb TokenBucket
		n0.PushTokens(&tb)
		Wish(t, tb.Minimalize(), ShouldEqual, TokenBucket{
			{Type: tok.TArrOpen, Length: 1},
			{Type: tok.TString, Str: "asdf"},
			{Type: tok.TArrClose},
		})
	})
	t.Run("nested list node", func(t *testing.T) {
		n0 := newNode()
		n00 := newNode()
		n0.SetIndex(0, n00)
		n000 := newNode()
		n000.SetString("asdf")
		n00.SetIndex(0, n000)
		n01 := newNode()
		n01.SetString("quux")
		n0.SetIndex(1, n01)
		var tb TokenBucket
		n0.PushTokens(&tb)
		Wish(t, tb.Minimalize(), ShouldEqual, TokenBucket{
			{Type: tok.TArrOpen, Length: 2},
			{Type: tok.TArrOpen, Length: 1},
			{Type: tok.TString, Str: "asdf"},
			{Type: tok.TArrClose},
			{Type: tok.TString, Str: "quux"},
			{Type: tok.TArrClose},
		})
	})
}
