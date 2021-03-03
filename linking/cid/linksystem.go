package cidlink

import (
	"fmt"
	"hash"

	cid "github.com/ipfs/go-cid"
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
		DecoderChooser: func(lp ipld.LinkPrototype) (ipld.Decoder, error) {
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
		Prototype: func(lnk ipld.Link) ipld.LinkPrototype {
			c, ok := lnk.(cid.Cid)
			if ok {
				return LinkPrototype{c.Prefix()}
			}
			return nil
		},
	}
}
