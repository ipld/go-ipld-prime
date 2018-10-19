package ipld

type Node interface {
	// GetField resolves a merklepath against the object and returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	//
	// If a Node is returned, it will be a unrooted node -- that is,
	// it can be used to view the fields below it, but since it was not
	// originally stored as a full node, you cannot immediately take
	// a link to it for embedding in other objects (you'd have to make
	// a new RootNode with the same content first, then store that).
	GetField(path []string) (interface{}, error)

	// Distinctly suspecting a GetField(string)(iface,err) API would be better.
	// And a GetIndex(int)(iface,err) for arrays and int-keyed maps.
	// Much easier to code.
	// And traversals are apparently going to require a type schema parameter!
	// Main counterargument: might be more efficient to let nodes take part,
	// so they can do work without generating intermediate nodes.
	// That logic applies to some but not all impls.  e.g. freenode always
	// already has the intermediates.  ipldcbor.Node prob will too.
	// Hrm.  I guess ipldbind.Node is the only one that can get usefully fancy there.
}

type SerializableNode interface {
}

type MutableNode interface {
}
