package cbor

import (
	"io"

	"github.com/polydawn/refmt/cbor"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/multicodec"
)

var (
	_ ipld.Decoder = Decode
	_ ipld.Encoder = Encode
)

func init() {
	multicodec.RegisterEncoder(0x51, Encode)
	multicodec.RegisterDecoder(0x51, Decode)
}

func Decode(na ipld.NodeAssembler, r io.Reader) error {
	return dagcbor.Unmarshal(na, cbor.NewDecoder(cbor.DecodeOptions{}, r),
		dagcbor.UnmarshalOptions{AllowLinks: false})
}

func Encode(n ipld.Node, w io.Writer) error {
	return dagcbor.Marshal(n, cbor.NewEncoder(w),
		dagcbor.MarshalOptions{AllowLinks: false})
}
