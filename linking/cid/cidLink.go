package cidlink

import (
	"fmt"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	mh "github.com/multiformats/go-multihash"
	multihash "github.com/multiformats/go-multihash"
)

var (
	_ ipld.Link          = Link{}
	_ ipld.LinkPrototype = LinkPrototype{}
)

// Link implements the ipld.Link interface using a CID.
// See https://github.com/ipfs/go-cid for more information about CIDs.
//
// When using this value, typically you'll use it as `Link`, and not `*Link`.
// This includes when handling the value as an `ipld.Link` interface -- the non-pointer form is typically preferable.
// This is because the ipld.Link inteface is often desirable to be able to use as a golang map key,
// and in that context, pointers would not result in the desired behavior.
type Link struct {
	cid.Cid
}

func (lnk Link) Prototype() ipld.LinkPrototype {
	return LinkPrototype{lnk.Cid.Prefix()}
}
func (lnk Link) String() string {
	return lnk.Cid.String()
}

type LinkPrototype struct {
	cid.Prefix
}

func (lp LinkPrototype) length() int {
	length := lp.MhLength
	if lp.MhType == multihash.ID {
		length = -1
	}
	return length
}

func (lp LinkPrototype) validate() error {
	if lp.Version == 0 && (lp.MhType != mh.SHA2_256 ||
		(lp.MhLength != 32 && lp.MhLength != -1)) {

		return fmt.Errorf("invalid v0 prefix")
	}

	if !mh.ValidCode(lp.MhType) {
		return fmt.Errorf("invalid code")
	}
	return nil
}

func (lp LinkPrototype) BuildLink(hashsum []byte) ipld.Link {
	if lp.length() != -1 {
		hashsum = hashsum[:lp.MhLength]
	}
	c, err := lp.Prefix.FromDigest(hashsum)
	if err != nil {
		panic(err)
	}
	return Link{Cid: c}
}
