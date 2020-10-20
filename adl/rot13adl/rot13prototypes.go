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
