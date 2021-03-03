package cidlink

import (
	"fmt"
	"hash"

	cid "github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/multicodec"
	"github.com/ipld/go-ipld-prime/multihash"
	mh "github.com/multiformats/go-multihash"
)

func DefaultLinkSystem() ipld.LinkSystem {
	return ipld.LinkSystem{
		EncoderChooser: func(lp ipld.LinkPrototype) (ipld.Encoder, error) {
			switch lp2 := lp.(type) {
			case cid.Prefix:
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
			case cid.Prefix:
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
			case cid.Prefix:
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
			if !ok {
				panic("cannot work with non-cid links")
			}
			return c.Prefix()
		},
		BuildLink: func(lp ipld.LinkPrototype, hashsum []byte) ipld.Link {
			// Does this method body look surprisingly complex?  I agree.
			//  We actually have to do all this work.  The go-cid package doesn't expose a constructor that just lets us directly set the bytes and the prefix numbers next to each other.
			//  No, `cid.Prefix.Sum` is not the method you are looking for: that expects the whole data body.
			//  Most of the logic here is the same as the body of `cid.Prefix.Sum`; we just couldn't get at the relevant parts without copypasta.
			//  There is also some logic that's sort of folded in from the go-multihash module.  This is really a mess.
			//  Note that there are also several things that error, hard, even though that may not be reasonable.  For example, multihash.Encode rejects multihash indicator numbers it doesn't explicitly know about, which is not great for forward compatibility.
			//  The go-cid package needs review.  So does go-multihash.  Their responsibilies are not well compartmentalized and they don't play well with other stdlib golang interfaces.
			p, ok := lp.(cid.Prefix)
			if !ok {
				panic("cannot work with non-cid links")
			}

			length := p.MhLength
			if p.MhType == mh.ID {
				length = -1
			}
			if p.Version == 0 && (p.MhType != mh.SHA2_256 ||
				(p.MhLength != 32 && p.MhLength != -1)) {
				panic(fmt.Errorf("invalid cid v0 prefix"))
			}

			if length != -1 {
				hashsum = hashsum[:p.MhLength]
			}

			mh, err := mh.Encode(hashsum, p.MhType)
			if err != nil {
				panic(err)
			}

			switch p.Version {
			case 0:
				return cid.NewCidV0(mh)
			case 1:
				return cid.NewCidV1(p.Codec, mh)
			default:
				panic(fmt.Errorf("invalid cid version"))
			}
		},
	}
}
