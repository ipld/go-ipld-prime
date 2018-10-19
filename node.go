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
}

type SerializableNode interface {
}

type MutableNode interface {
}
