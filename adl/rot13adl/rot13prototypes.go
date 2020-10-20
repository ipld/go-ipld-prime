package rot13adl

// Prototype embeds a NodePrototype for every kind of Node implementation in this package.
// This includes both the synthesized node as well as the substrate root node
// (other substrate interior node prototypes are not exported here;
// it's unlikely they'll be useful outside of the scope of the ADL's implementation package.)
//
// You can use it like this:
//
// 		rot13adl.Prototype.Node.NewBuilder() //...
//
// and:
//
// 		rot13adl.Prototype.SubstrateRoot.NewBuilder() // ...
//
var Prototype prototype

// This may seem a tad mundane, but we do it for consistency with the way
// other Node implementation packages (like basicnode) do this:
// it's a convention for "pkgname.Prototype.SpecificThing.NewBuilder()" to be defined.

type prototype struct {
	Node          _R13String__Prototype
	SubstrateRoot _Substrate__Prototype
}

// REVIEW: In ADLs that use codegenerated substrate types defined by an IPLD Schema, the `Prototype.SubstrateRoot` field...
// should it be a `_SubstrateRoot__Prototype`, or a `_SubstrateRoot__ReprPrototype`?
//  Probably the latter, because the only thing an external user of this package should be interested in is how to read data into memory in an optimal path.
// But does this answer all questions?  Codegen ReprPrototypes currently still return the type-level node from their Build method!
//  This probably would functionally work -- we could make the Reify methods check for either the type-level or repr-level types -- but would it be weird/surprising/hard-to-use?
