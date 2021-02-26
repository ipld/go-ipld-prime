package cidlink

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/multihash"
)

// EncoderChooser is the default encoder chooser for the cidlink system
func EncoderChooser(lp ipld.LinkPrototype) (ipld.Encoder, error) {
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
}

// DecoderChooser is the default decoder chooser for a cidlink Link system
func DecoderChooser(lnk ipld.Link) (ipld.Decoder, error) {
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
}

// HasherChooser is the default choser for a cidlink LinkSystem that uses the new hashing system
func HasherChooser(lp ipld.LinkPrototype) (ipld.Hasher, error) {
	switch lp2 := lp.(type) {
	case LinkPrototype:
		err := lp2.validate()
		if err != nil {
			return nil, err
		}
		fn, ok := multihash.Registry[lp2.MhType]
		if !ok {
			return nil, fmt.Errorf("no hasher registered for multihash indicator 0x%x", lp2.MhType)
		}
		return fn(), nil
	default:
		return nil, fmt.Errorf("this decoderChooser can only handle cidlink.LinkPrototype; got %T", lp)
	}
}

// LegacyHasherChooser uses go-multihash to build hash digests, insuring compatibility with the codec table
func LegacyHasherChooser(lp ipld.LinkPrototype) (ipld.Hasher, error) {
	switch lp2 := lp.(type) {
	case LinkPrototype:
		err := lp2.validate()
		if err != nil {
			return nil, err
		}
		return newMHHash(lp2.MhType, lp2.length())
	default:
		return nil, fmt.Errorf("this decoderChooser can only handle cidlink.LinkPrototype; got %T", lp)
	}
}

func DefaultLinkSystem() ipld.LinkSystem {
	return ipld.LinkSystem{
		EncoderChooser: EncoderChooser,
		DecoderChooser: DecoderChooser,
		HasherChooser:  HasherChooser,
	}
}
