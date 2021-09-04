package codec_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	_ "github.com/ipld/go-ipld-prime/codec/cbor"
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	_ "github.com/ipld/go-ipld-prime/codec/dagjson"
	_ "github.com/ipld/go-ipld-prime/codec/json"
	mcregistry "github.com/ipld/go-ipld-prime/multicodec"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/multiformats/go-multicodec"
)

func TestDecodeZero(t *testing.T) {
	for _, code := range []multicodec.Code{
		multicodec.Cbor,
		multicodec.DagCbor,
		multicodec.Json,
		multicodec.DagJson,
	} {
		t.Run(code.String(), func(t *testing.T) {
			nb := basicnode.Prototype.Any.NewBuilder()
			decode, err := mcregistry.LookupDecoder(uint64(code))
			if err != nil {
				t.Fatal(err)
			}

			err = decode(nb, strings.NewReader(""))
			if !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
