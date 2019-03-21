// go-ipld-prime is a series of go interfaces for manipulating IPLD data.
//
// See https://github.com/ipld/specs for more information about the basics
// of "What is IPLD?".
//
// See https://github.com/ipld/go-ipld-prime/tree/master/doc/README.md
// for more documentation about go-ipld-prime's architecture and usage.
//
// Here in the godoc, the first couple of types to look at should be:
//
//   - Node
//   - NodeBuilder
//
// These types provide a generic description of the data model.
//
// If working with linked data (data which is split into multiple
// trees of Nodes, loaded separately, and connected by some kind of
// "link" reference), the next types you should look at are:
//
//   - Link
//   - LinkBuilder
//   - Loader
//   - Storer
//
// All of these types are interfaces.  There are several implementations you
// can choose; we've provided some in subpackages, or you can bring your own.
//
// Particularly interesting subpackages include:
//
//   - impl/* -- various Node + NodeBuilder implementations
//   - encoding/* -- functions for serializing and deserializing Nodes
//   - linking/* -- various Link + LinkBuilder implementation
//   - traversal -- functions for walking Node graphs (including
//        automatic link loading) and visiting
//   - typed -- Node implementations with constraints
//   - fluent -- Node interfaces with streamlined error handling
//
package ipld
