package dagjson

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/fluent"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

var byteNode = fluent.MustBuildMap(basicnode.Prototype__Map{}, 4, func(na fluent.MapAssembler) {
	na.AssembleEntry("plain").AssignString("olde string")
	na.AssembleEntry("bytes").AssignBytes([]byte("deadbeef"))
})
var byteSerial = `{
	"plain": "olde string",
	"bytes": {
		"/": {
			"bytes": "ZGVhZGJlZWY="
		}
	}
}
`

func TestRoundtripBytes(t *testing.T) {
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := Encode(byteNode, &buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, buf.String(), ShouldEqual, byteSerial)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := strings.NewReader(byteSerial)
		nb := basicnode.Prototype__Map{}.NewBuilder()
		err := Decode(nb, buf)
		Require(t, err, ShouldEqual, nil)
		Wish(t, nb.Build(), ShouldEqual, byteNode)
	})
}
