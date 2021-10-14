package linking_test

import (
	"fmt"

	"github.com/ipfs/go-cid"

	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/storage/memstore"
)

// storage is a map where we'll store serialized IPLD data.
//
// ExampleLinkSystem_Store will put data into this;
// ExampleLinkSystem_Load will read out from it.
//
// In a real program, you'll probably make functions to load and store from disk,
// or some network storage, or... whatever you want, really :)
var store = memstore.Store{}

// TODO: These examples are really heavy on CIDs and the multicodec and multihash magic tables.
// It would be good to have examples that create and use less magical LinkSystem constructions, too.

func ExampleLinkSystem_Store() {
	// Creating a Link is done by choosing a concrete link implementation (typically, CID),
	//  getting a LinkSystem that knows how to work with that, and then using the LinkSystem methods.

	// Let's get a LinkSystem.  We're going to be working with CID links,
	//  so let's get the default LinkSystem that's ready to work with those.
	lsys := cidlink.DefaultLinkSystem()

	// We want to store the serialized data somewhere.
	//  We'll use an in-memory store for this.  (It's a package scoped variable.)
	//  You can use any kind of storage system here;
	//   or if you need even more control, you could also write a function that conforms to the linking.BlockWriteOpener interface.
	lsys.SetWriteStorage(&store)

	// To create any links, first we need a LinkPrototype.
	// This gathers together any parameters that might be needed when making a link.
	// (For CIDs, the version, the codec, and the multihash type are all parameters we'll need.)
	// Often, you can probably make this a constant for your whole application.
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,    // Usually '1'.
		Codec:    0x71, // 0x71 means "dag-cbor" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhType:   0x13, // 0x20 means "sha2-512" -- See the multicodecs table: https://github.com/multiformats/multicodec/
		MhLength: 64,   // sha2-512 hash has a 64-byte sum.
	}}

	// And we need some data to link to!  Here's a quick piece of example data:
	n := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
		na.AssembleEntry("hello").AssignString("world")
	})

	// Before we use the LinkService, NOTE:
	//  There's a side-effecting import at the top of the file.  It's for the dag-cbor codec.
	//  The CID LinkSystem defaults use a global registry called the multicodec table;
	//  and the multicodec table is populated in part by the dag-cbor package when it's first imported.
	// You'll need that side-effecting import, too, to copy this example.
	//  It can happen anywhere in your program; once, in any package, is enough.
	//  If you don't have this import, the codec will not be registered in the multicodec registry,
	//  and when you use the LinkSystem we got from the cidlink package, it will return an error of type ErrLinkingSetup.
	// If you initialize a custom LinkSystem, you can control this more directly;
	//  these registry systems are only here as defaults.

	// Now: time to apply the LinkSystem, and do the actual store operation!
	lnk, err := lsys.Store(
		linking.LinkContext{}, // The zero value is fine.  Configure it it you want cancellability or other features.
		lp,                    // The LinkPrototype says what codec and hashing to use.
		n,                     // And here's our data.
	)
	if err != nil {
		panic(err)
	}

	// That's it!  We got a link.
	fmt.Printf("link: %s\n", lnk)
	fmt.Printf("concrete type: `%T`\n", lnk)

	// Remember: the serialized data was also stored to the 'store' variable as a side-effect.
	//  (We set this up back when we customized the LinkSystem.)
	//  We'll pick this data back up again in the example for loading.

	// Output:
	// link: bafyrgqhai26anf3i7pips7q22coa4sz2fr4gk4q4sqdtymvvjyginfzaqewveaeqdh524nsktaq43j65v22xxrybrtertmcfxufdam3da3hbk
	// concrete type: `cidlink.Link`
}

func ExampleLinkSystem_Load() {
	// Let's say we want to load this link (it's the same one we created in ExampleLinkSystem_Store).
	cid, _ := cid.Decode("bafyrgqhai26anf3i7pips7q22coa4sz2fr4gk4q4sqdtymvvjyginfzaqewveaeqdh524nsktaq43j65v22xxrybrtertmcfxufdam3da3hbk")
	lnk := cidlink.Link{Cid: cid}

	// Let's get a LinkSystem.  We're going to be working with CID links,
	//  so let's get the default LinkSystem that's ready to work with those.
	// (This is the same as we did in ExampleLinkSystem_Store.)
	lsys := cidlink.DefaultLinkSystem()

	// We need somewhere to go looking for any of the data we might want to load!
	//  We'll use an in-memory store for this.  (It's a package scoped variable.)
	//   (This particular memory store was filled with the data we'll load earlier, during ExampleLinkSystem_Store.)
	//  You can use any kind of storage system here;
	//   or if you need even more control, you could also write a function that conforms to the linking.BlockReadOpener interface.
	lsys.SetReadStorage(&store)

	// We'll need to decide what in-memory implementation of datamodel.Node we want to use.
	//  Here, we'll use the "basicnode" implementation.  This is a good getting-started choice.
	//   But you could also use other implementations, or even a code-generated type with special features!
	np := basicnode.Prototype.Any

	// Before we use the LinkService, NOTE:
	//  There's a side-effecting import at the top of the file.  It's for the dag-cbor codec.
	//  See the comments in ExampleLinkSystem_Store for more discussion of this and why it's important.

	// Apply the LinkSystem, and ask it to load our link!
	n, err := lsys.Load(
		linking.LinkContext{}, // The zero value is fine.  Configure it it you want cancellability or other features.
		lnk,                   // The Link we want to load!
		np,                    // The NodePrototype says what kind of Node we want as a result.
	)
	if err != nil {
		panic(err)
	}

	// Tada!  We have the data as node that we can traverse and use as desired.
	fmt.Printf("we loaded a %s with %d entries\n", n.Kind(), n.Length())

	// Output:
	// we loaded a map with 1 entries
}
