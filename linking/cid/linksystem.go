package cidlink

import (
	"fmt"
	"hash"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/multihash"
)

func DefaultLinkSystem() ipld.LinkSystem {
	return ipld.LinkSystem{
		EncoderChooser: func(lp ipld.LinkPrototype) (ipld.Encoder, error) {
			switch lp2 := lp.(type) {
			case LinkPrototype:
				fn, ok := multicodec.EncoderRegistry[lp2.GetCodec()]
				if !ok {
					return nil, fmt.Errorf("no encoder registered for multicodec indicator 0x%x", lp2.GetCodec())
				}
				return fn, nil
			default:
				return nil, fmt.Errorf("this encoderChooser can only handle cidlink.LinkPrototype; got %T", lp)
			}
		},
		DecoderChooser: func(lnk ipld.Link) (ipld.Decoder, error) {
			lp := lnk.Prototype()
			switch lp2 := lp.(type) {
			case LinkPrototype:
				fn, ok := multicodec.DecoderRegistry[lp2.GetCodec()]
				if !ok {
					return nil, fmt.Errorf("no decoder registered for multicodec indicator 0x%x", lp2.GetCodec())
				}
				return fn, nil
			default:
				return nil, fmt.Errorf("this decoderChooser can only handle cidlink.LinkPrototype; got %T", lp)
			}
		},
		HasherChooser: func(lp ipld.LinkPrototype) (hash.Hash, error) {
			switch lp2 := lp.(type) {
			case LinkPrototype:
				fn, ok := multihash.Registry[lp2.MhType]
				if !ok {
					return nil, fmt.Errorf("no hasher registered for multihash indicator 0x%x", lp2.MhType)
				}
				return fn(), nil
			default:
				return nil, fmt.Errorf("this decoderChooser can only handle cidlink.LinkPrototype; got %T", lp)
			}
		},
	}
}
