package dagjson

import (
	"bytes"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

// TestDecodeOptions_MaxDepth verifies that the configurable nesting-depth
// limit is respected, both with defaults and custom values.
func TestDecodeOptions_MaxDepth(t *testing.T) {
	nested := func(depth int) []byte {
		return []byte(strings.Repeat("[", depth) + "null" + strings.Repeat("]", depth))
	}

	t.Run("default depth rejects deeply nested structure", func(t *testing.T) {
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(nested(2000)))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("structure at default depth decodes", func(t *testing.T) {
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(nested(1024)))
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, nb.Build().Kind(), qt.Equals, datamodel.Kind_List)
	})

	t.Run("custom depth rejects when exceeded", func(t *testing.T) {
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{MaxDepth: 5}.Decode(nb, bytes.NewReader(nested(10)))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("custom depth accepts within limit", func(t *testing.T) {
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{MaxDepth: 10}.Decode(nb, bytes.NewReader(nested(5)))
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, nb.Build().Kind(), qt.Equals, datamodel.Kind_List)
	})

	t.Run("nested maps also limited", func(t *testing.T) {
		const depth = 2000
		buf := strings.Repeat(`{"x":`, depth) + "null" + strings.Repeat("}", depth)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader([]byte(buf)))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("zero MaxDepth resolves to default", func(t *testing.T) {
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{MaxDepth: 0}.Decode(nb, bytes.NewReader(nested(2000)))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("ParseLinks lookahead does not bypass depth", func(t *testing.T) {
		// A valid DAG-JSON link ({"/":"..."}) wrapped in deep list nesting.
		// The lookahead path for ParseLinks must still honour the depth limit.
		const depth = 2000
		buf := strings.Repeat("[", depth) + `{"/":"bafkqaaa"}` + strings.Repeat("]", depth)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{ParseLinks: true}.Decode(nb, bytes.NewReader([]byte(buf)))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("ParseBytes lookahead does not bypass depth", func(t *testing.T) {
		// A valid DAG-JSON bytes object wrapped in deep list nesting.
		const depth = 2000
		buf := strings.Repeat("[", depth) + `{"/":{"bytes":"aGVsbG8"}}` + strings.Repeat("]", depth)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{ParseBytes: true}.Decode(nb, bytes.NewReader([]byte(buf)))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("ParseLinks within limit resolves link correctly", func(t *testing.T) {
		// Depth 5 well within the default limit; ensure the lookahead path
		// still yields a Link node when not overflowing.
		buf := strings.Repeat("[", 5) + `{"/":"bafkqaaa"}` + strings.Repeat("]", 5)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{ParseLinks: true}.Decode(nb, bytes.NewReader([]byte(buf)))
		qt.Assert(t, err, qt.IsNil)
	})
}
