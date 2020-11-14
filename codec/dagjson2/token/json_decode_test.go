package jsontoken

import (
	"io"
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/codec/codectools"
	"github.com/ipld/go-ipld-prime/codec/codectools/scratch"
)

func makeReader(s string) *scratch.Reader {
	r := &scratch.Reader{}
	r.InitSlice([]byte(s))
	return r
}

var inf int = 1 << 31

func TestDecode(t *testing.T) {
	t.Run("SimpleString", func(t *testing.T) {
		var d Decoder
		d.Init(strings.NewReader(`"asdf"`))
		tok, err := d.Step(&inf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, tok.Kind, ShouldEqual, codectools.TokenKind_String)
		Wish(t, tok.Str, ShouldEqual, "asdf")
		tok, err = d.Step(&inf)
		Wish(t, err, ShouldEqual, io.EOF)
	})
	t.Run("SingleMap", func(t *testing.T) {
		var d Decoder
		d.Init(strings.NewReader(`{"a":"b","c":"d"}`))
		tok, err := d.Step(&inf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, tok.Kind, ShouldEqual, codectools.TokenKind_MapOpen)
		Wish(t, d.phase, ShouldEqual, decoderPhase_acceptMapKeyOrEnd)
		tok, err = d.Step(&inf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, tok.Kind, ShouldEqual, codectools.TokenKind_String)
		Wish(t, tok.Str, ShouldEqual, "a")
		Wish(t, d.phase, ShouldEqual, decoderPhase_acceptMapValue)
		tok, err = d.Step(&inf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, tok.Kind, ShouldEqual, codectools.TokenKind_String)
		Wish(t, tok.Str, ShouldEqual, "b")
		tok, err = d.Step(&inf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, tok.Kind, ShouldEqual, codectools.TokenKind_String)
		Wish(t, tok.Str, ShouldEqual, "c")
		tok, err = d.Step(&inf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, tok.Kind, ShouldEqual, codectools.TokenKind_String)
		Wish(t, tok.Str, ShouldEqual, "d")
		tok, err = d.Step(&inf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, tok.Kind, ShouldEqual, codectools.TokenKind_MapClose)
		tok, err = d.Step(&inf)
		Wish(t, err, ShouldEqual, io.EOF)
	})
}
