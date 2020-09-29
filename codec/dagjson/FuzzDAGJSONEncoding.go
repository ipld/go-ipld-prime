package dagjson

import (
	"bytes"
	"strings"

	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	fleece "github.com/leastauthority/fleece/fuzzing"
)

// FuzzDAGJSONEncoding fuzz tests the encode and decode flow
// of dagjson
func FuzzDAGJSONEncoding(data []byte) int {
	buf := strings.NewReader(string(data))
	nb := basicnode.Prototype__Map{}.NewBuilder()
	err := Decoder(nb, buf)
	if err != nil {
		return fleece.FuzzNormal
	}

	// attempt to build the node to force processing
	n := nb.Build()
	buf2 := new(bytes.Buffer)
	err = Encoder(n, buf2)
	if err != nil {
		return fleece.FuzzInteresting
	}

	return fleece.FuzzInteresting
}
