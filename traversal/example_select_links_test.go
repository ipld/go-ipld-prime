package traversal_test

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"

	// Import all the codecs so we can use them in our example; each of these will
	// set themselves up in our multicodec registry so the LinkSystem can use
	// them for encoding (when a LinkPrototype says so) and decoding (when a CID's
	// codec code says so).
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	_ "github.com/ipld/go-ipld-prime/codec/dagjson"
	_ "github.com/ipld/go-ipld-prime/codec/raw"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/multiformats/go-multihash"
)

func ExampleSelectLinks() {
	// Setup: make some blocks, store them in our memory store

	blocks := make([]cid.Cid, 0)

	// make 3 raw blocks
	c1 := encodeAndStore(rawlp, basicnode.NewBytes([]byte{0xca, 0xfe, 0xbe, 0xef}))
	blocks = append(blocks, c1)
	c2 := encodeAndStore(rawlp, basicnode.NewBytes([]byte{0xde, 0xad, 0xbe, 0xef}))
	blocks = append(blocks, c2)
	c3 := encodeAndStore(rawlp, basicnode.NewBytes([]byte{0xba, 0xad, 0xf0, 0x0d}))
	blocks = append(blocks, c3)

	// Pretend we're doing a dag-pb encode here but since we don't want to pull in
	// the dagpb package we'll do it as dag-json. This should use
	// dagpb.Type.PBNode if it were real.
	pbn, err := qp.BuildMap(basicnode.Prototype.Map, 1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "Links", qp.List(3, func(la datamodel.ListAssembler) {
			qp.ListEntry(la, qp.Map(2, func(ma datamodel.MapAssembler) {
				qp.MapEntry(ma, "Name", qp.String("01"))
				qp.MapEntry(ma, "Hash", qp.Link(cidlink.Link{Cid: c1}))
			}))
			qp.ListEntry(la, qp.Map(2, func(ma datamodel.MapAssembler) {
				qp.MapEntry(ma, "Name", qp.String("02"))
				qp.MapEntry(ma, "Hash", qp.Link(cidlink.Link{Cid: c2}))
			}))
			qp.ListEntry(la, qp.Map(2, func(ma datamodel.MapAssembler) {
				qp.MapEntry(ma, "Name", qp.String("03"))
				qp.MapEntry(ma, "Hash", qp.Link(cidlink.Link{Cid: c3}))
			}))
		}))
	})
	if err != nil {
		panic(err)
	}
	cpb := encodeAndStore(pblp, pbn)
	blocks = append(blocks, cpb)

	// make a dag-cbor block with a bunch of links in it, stored in various ways
	cbn, err := qp.BuildList(basicnode.Prototype.List, -1, func(la datamodel.ListAssembler) {
		qp.ListEntry(la, qp.String("not a link!"))
		qp.ListEntry(la, qp.Link(cidlink.Link{Cid: c1}))
		qp.ListEntry(la, qp.Link(cidlink.Link{Cid: c2}))
		qp.ListEntry(la, qp.Int(42))
		qp.ListEntry(la, qp.Link(cidlink.Link{Cid: c3}))
		qp.ListEntry(la, qp.Map(-1, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, "Boop!", qp.String("boop!"))
			qp.MapEntry(ma, "This is a dag-pb link:", qp.Link(cidlink.Link{Cid: cpb}))
			qp.MapEntry(ma, "And this one the same link:", qp.Link(cidlink.Link{Cid: cpb}))
			qp.MapEntry(ma, "But thus one is a raw link:", qp.Link(cidlink.Link{Cid: c2}))
		}))
	})
	if err != nil {
		panic(err)
	}
	ccb := encodeAndStore(dclp, cbn)
	blocks = append(blocks, ccb)

	// Example code: load the blocks as datamodel.Node form and traverse
	// them using the SelectLinks function to find all the links

	for _, c := range blocks {
		// load Node form of the block
		n := loadNode(c)
		// Select all links from the node
		links, err := traversal.SelectLinks(n)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s (%d links)\n", c, len(links))
		// Print the links
		for _, l := range links {
			fmt.Printf("\t➜ %s\n", l)
		}
	}

	// Demonstrating how we might do this if we don't have a LinkSystem but do
	// have the bytes of a block we want to traverse.

	// load bytes of the block
	byts := loadBytes(ccb)
	// decode the block into Node form
	n, err := ipld.Decode(byts, dagcbor.Decode)
	if err != nil {
		panic(err)
	}
	// Select all links from the node
	links, err := traversal.SelectLinks(n)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Manually decoded %s has %d links too, surprise!\n", ccb, len(links))

	// Output:
	// bafkreicwcm4sqhux7ipwbwro2inf4vbjoprzwoqnl3t4zta4gphofot7du (0 links)
	// bafkreic7pdbte5heh6u54vszezob3el6exadoiw4wc4ne7ny2x7kvajzkm (0 links)
	// bafkreiafc5f36dkaocd6iwysxkxboelue2cs745j4wgrfihlxgqqwqexim (0 links)
	// baguqeeraawfdfutq7b4dgmcelypy4r36d7ndngbzpbejzhjaxffcrbgq3wka (3 links)
	// 	➜ bafkreicwcm4sqhux7ipwbwro2inf4vbjoprzwoqnl3t4zta4gphofot7du
	// 	➜ bafkreic7pdbte5heh6u54vszezob3el6exadoiw4wc4ne7ny2x7kvajzkm
	// 	➜ bafkreiafc5f36dkaocd6iwysxkxboelue2cs745j4wgrfihlxgqqwqexim
	// bafyreietdfd5sh743y6c4f4zrhwkqdadwtcuvziot3mvwhyqojjtvtqoxi (6 links)
	// 	➜ bafkreicwcm4sqhux7ipwbwro2inf4vbjoprzwoqnl3t4zta4gphofot7du
	// 	➜ bafkreic7pdbte5heh6u54vszezob3el6exadoiw4wc4ne7ny2x7kvajzkm
	// 	➜ bafkreiafc5f36dkaocd6iwysxkxboelue2cs745j4wgrfihlxgqqwqexim
	// 	➜ baguqeeraawfdfutq7b4dgmcelypy4r36d7ndngbzpbejzhjaxffcrbgq3wka
	// 	➜ baguqeeraawfdfutq7b4dgmcelypy4r36d7ndngbzpbejzhjaxffcrbgq3wka
	// 	➜ bafkreic7pdbte5heh6u54vszezob3el6exadoiw4wc4ne7ny2x7kvajzkm
	// Manually decoded bafyreietdfd5sh743y6c4f4zrhwkqdadwtcuvziot3mvwhyqojjtvtqoxi has 6 links too, surprise!
}

// LinkPrototypes to tell the LinkSystem how to store the Nodes in our memory
// store and how to generate the CIDs (i.e. codec to store, multihash to use in
// the CID).

var dclp = cidlink.LinkPrototype{
	Prefix: cid.Prefix{
		Version:  1,
		Codec:    cid.DagCBOR,
		MhType:   multihash.SHA2_256,
		MhLength: 32,
	},
}

var pblp = cidlink.LinkPrototype{
	Prefix: cid.Prefix{
		Version:  1,
		Codec:    cid.DagJSON, // cid.DagProtobuf, but we're pretending
		MhType:   multihash.SHA2_256,
		MhLength: 32,
	},
}

var rawlp = cidlink.LinkPrototype{
	Prefix: cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   multihash.SHA2_256,
		MhLength: 32,
	},
}

// Given a LinkPrototype and a Node, encode the Node and store it in our memory
// store, returning the CID of the stored Node.
func encodeAndStore(lp cidlink.LinkPrototype, n datamodel.Node) cid.Cid {
	lsys := cidlink.DefaultLinkSystem()
	// note this is in our package: var store = memstore.Store{}
	lsys.SetWriteStorage(&store)
	lsys.SetReadStorage(&store)
	lnk := lsys.MustStore(linking.LinkContext{}, lp, n)
	return lnk.(cidlink.Link).Cid
}

// Given a CID, load the Node from our memory store and return it.
func loadNode(c cid.Cid) datamodel.Node {
	lsys := cidlink.DefaultLinkSystem()
	// note this is in our package: var store = memstore.Store{}
	lsys.SetReadStorage(&store)
	nb := lsys.MustLoad(linking.LinkContext{}, cidlink.Link{Cid: c}, basicnode.Prototype.Any)
	return nb
}

// Given a CID, load the bytes from our memory store and return them.
func loadBytes(c cid.Cid) []byte {
	lsys := cidlink.DefaultLinkSystem()
	// note this is in our package: var store = memstore.Store{}
	lsys.SetReadStorage(&store)
	byts, err := lsys.LoadRaw(linking.LinkContext{}, cidlink.Link{Cid: c})
	if err != nil {
		panic(err)
	}
	return byts
}
