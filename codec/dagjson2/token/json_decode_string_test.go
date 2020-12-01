package jsontoken

import (
	"errors"
	"io"
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestDecodeString(t *testing.T) {
	t.Run("SimpleString", func(t *testing.T) {
		s, err := DecodeString(makeReader(`"asdf"`))
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, "asdf")
	})
	t.Run("NonString", func(t *testing.T) {
		s, err := DecodeString(makeReader(`not prefixed right`))
		Wish(t, err, ShouldEqual, errors.New(`not a string: strings must begin with '"', not 'n'`))
		Wish(t, s, ShouldEqual, "")
	})
	t.Run("UnterminatedString", func(t *testing.T) {
		s, err := DecodeString(makeReader(`"ohno`))
		Wish(t, err, ShouldEqual, io.ErrUnexpectedEOF)
		Wish(t, s, ShouldEqual, "")
	})
	t.Run("StringWithEscapes", func(t *testing.T) {
		s, err := DecodeString(makeReader(`"as\tdf\bwow"`))
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, "as\tdf\bwow")
	})
}
