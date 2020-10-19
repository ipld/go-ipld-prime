package ipld_test

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ipfs/go-cid"

	ipld "github.com/ipld/go-ipld-prime"
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

// storage is a map where we'll store serialized IPLD data.
//
// ExampleCreatingLink will put data into this;
// ExampleLoadingLink will read out from it.
//
// In a real program, you'll probably make functions to load and store from disk,
// or some network storage, or... whatever you want, really :)
var storage = make(map[ipld.Link][]byte)

func ExampleCreatingLink() {
	// Creating a link is done by choosing a concrete link implementation (typically, CID),
	//  importing that package, and using its functions to create the link.

	// First, create a LinkBuilder.  This gathers together any parameters that might be needed when making a link.
	// (For CIDs, the version, the codec, and the multihash type are all parameters we'll need.)
	lb := cidlink.LinkBuilder{cid.Prefix{
		Version:  1,    // Usually '1'.
		Codec:    0x71, // 0x71 means "dag-cbor" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhType:   0x15, // 0x15 means "sha3-384" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhLength: 48,   // sha3-224 hash has a 48-byte sum.
	}}

	// And we need some data to link to!  Here's a quick piece of example data:
	n := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
		na.AssembleEntry("hello").AssignString("world")
	})

	// Building a link takes a bunch of arguments:
	// - a `context.Context` -- this is a standard way to support cancellability in long-running tasks in golang.
	//    (Hashing to form a link is fast -- but you might be writing to a slow storage medium at the same time.)
	// - an `ipld.LinkContext` -- this can provide additional info (like a path -- the traversal package will do this), but can also be empty.
	// - the `ipld.Node` to serialize and create the link for!
	// - an `ipld.Storer` -- this is a function that defines where the serialized Node is written to.
	lnk, err := lb.Build(
		context.Background(),
		ipld.LinkContext{},
		n,
		func(ipld.LinkContext) (io.Writer, ipld.StoreCommitter, error) {
			buf := bytes.Buffer{}
			return &buf, func(lnk ipld.Link) error {
				storage[lnk] = buf.Bytes()
				return nil
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}

	// That's it!  We got a link.
	fmt.Printf("link: %s\n", lnk)
	fmt.Printf("concrete type: `%T`\n", lnk)

	// Output:
	// link: bafyrkmbukvrgzcs6qlsh4wvkvbe5wp7sclcblfnapnb2xfznisbykpbnlocet2qzley3cpxofoxqrnqgm3ta
	// concrete type: `*cidlink.Link`
}

func ExampleLoadingLink() {
	// Let's say we want to load this link (it's the same one we just created in the example above).
	cid, _ := cid.Decode("bafyrkmbukvrgzcs6qlsh4wvkvbe5wp7sclcblfnapnb2xfznisbykpbnlocet2qzley3cpxofoxqrnqgm3ta")
	lnk := &cidlink.Link{cid}

	// First, we'll need a Loader.  This function has to take a link as a parameter,
	//  then decides where to get the referenced raw data from,
	//   and returns that as a standard `io.Reader`.
	var loader ipld.Loader = func(lnk ipld.Link, _ ipld.LinkContext) (io.Reader, error) {
		return bytes.NewReader(storage[lnk]), nil
	}

	// Second, we'll need to decide what in-memory implementation of ipld.Node we want to use.
	//  Here, we'll use the "basicnode" implementation.
	//   But you could also use other implementations, or even a code-generated type with special features!
	// To encapsulate this decision, we create a NodeBuilder for the implementation we want.
	//  (If you are building a library and want to expose this choice, though, you'd probably want to accept a NodePrototype as the configuration for this.)
	nb := basicnode.Prototype.Any.NewBuilder()

	// Tell the link to load itself!
	//  This only returns an error...
	//  the data itself gets loaded into the NodeBuilder.
	//   (It's kinda like how if you use stdlib's `json.Unmarshal`, you give it `&something` as a parameter, and it fills that in.)
	err := lnk.Load(
		context.Background(), // As with creating links, a context here is so you can interrupt the process if it's slow.
		ipld.LinkContext{},   // The LinkContext can provide more info, but it's also fine if it's empty.
		nb,                   // Here's the NodeBuilder we'll pour the unmarshalled data into.
		loader,               // The loader is called to get the io.Reader for the raw data.
	)
	if err != nil {
		panic(err)
	}

	// We can get the reified data from the NodeBuilder:
	n := nb.Build()

	// Tada!  We have the data as node that we can traverse and use as desired.
	fmt.Printf("we loaded a %s with %d entries\n", n.ReprKind(), n.Length())

	// Output:
	// we loaded a map with 1 entries
}
