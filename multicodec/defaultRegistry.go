package multicodec

import (
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/codecregistry"
)

type Registry = codecregistry.Registry

var DefaultRegistry = codecregistry.DefaultRegistry

func RegisterEncoder(indicator uint64, encodeFunc codec.Encoder) {
	codecregistry.RegisterEncoder(indicator, encodeFunc)
}

func LookupEncoder(indicator uint64) (codec.Encoder, error) {
	return codecregistry.LookupEncoder(indicator)
}

func ListEncoders() []uint64 {
	return codecregistry.ListEncoders()
}

func RegisterDecoder(indicator uint64, decodeFunc codec.Decoder) {
	codecregistry.RegisterDecoder(indicator, decodeFunc)
}

func LookupDecoder(indicator uint64) (codec.Decoder, error) {
	return codecregistry.LookupDecoder(indicator)
}

func ListDecoders() []uint64 {
	return codecregistry.ListDecoders()
}
