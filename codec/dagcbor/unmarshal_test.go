package dagcbor

import (
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestFunBlocks(t *testing.T) {
	t.Run("zero length link", func(t *testing.T) {
		// This fixture has a zero length link -- not even the multibase byte (which dag-cbor insists must be zero) is there.
		buf := strings.NewReader("\x8d\x8d\x97\xd8*@")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decoder(nb, buf)
		Require(t, err, ShouldEqual, ErrInvalidMultibase)
	})
}
