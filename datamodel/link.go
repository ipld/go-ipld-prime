package datamodel

// Link is a special kind of value in IPLD which can be "loaded" to access more nodes.
//
// Nodes can be a Link: "link" is one of the kinds in the IPLD Data Model;
// and accordingly there is an `ipld.Kind_Link` enum value, and Node has an `AsLink` method.
//
// Links are considered a scalar value in the IPLD Data Model,
// but when "loaded", the result can be any other IPLD kind:
// maps, lists, strings, etc.
//
// Link is an interface in the go-ipld implementation,
// but the most common instantiation of it comes from the `linking/cid` package,
// and represents CIDs (see https://github.com/multiformats/cid).
//
// The Link interface says very little by itself; it's generally necessary to
// use type assertions to unpack more specific forms of data.
// The only real contract is that the Link must be able to return a LinkPrototype,
// which must be able to produce new Link values of a similar form.
// (In practice: if you're familiar with CIDs: Link.Prototype is analogous to cid.Prefix.)
//
// The traversal package contains powerful features for walking through large graphs of Nodes
// while automatically loading and traversing links as the walk goes.
//
// Note that the Link interface should typically be inhabited by a struct or string, as opposed to a pointer.
// This is because Link is often desirable to be able to use as a golang map key,
// and in that context, pointers would not result in the desired behavior.
type Link interface {
	// Prototype should return a LinkPrototype which carries the information
	// to make more Link values similar to this one (but with different hashes).
	Prototype() LinkPrototype

	// String should return a reasonably human-readable debug-friendly representation the Link.
	// There is no contract that requires that the string be able to be parsed back into a Link value,
	// but the string should be unique (e.g. not elide any parts of the hash).
	String() string
}

// LinkPrototype encapsulates any implementation details and parameters
// necessary for creating a Link, expect for the hash result itself.
//
// LinkPrototype, like Link, is an interface in go-ipld,
// but the most common instantiation of it comes from the `linking/cid` package,
// and represents CIDs (see https://github.com/multiformats/cid).
// If using CIDs as an implementation, LinkPrototype will encapsulate information
// like multihashType, multicodecType, and cidVersion, for example.
// (LinkPrototype is analogous to cid.Prefix.)
type LinkPrototype interface {
	// BuildLink should return a new Link value based on the given hashsum.
	// The hashsum argument should typically be a value returned from a
	// https://golang.org/pkg/hash/#Hash.Sum call.
	//
	// The hashsum reference must not be retained (the caller is free to reuse it).
	BuildLink(hashsum []byte) Link
}
