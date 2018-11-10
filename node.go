package ipld

import "github.com/ipfs/go-cid"

type Node interface {
	// GetField resolves a path in the the object and returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	TraverseField(path string) (Node, error)

	// GetIndex is the equivalent of GetField but for indexing into an array
	// (or a numerically-keyed map).  Like GetField, it returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	TraverseIndex(idx int) (Node, error)

	AsBool() (bool, error)
	AsString() (string, error)
	AsInt() (int, error)
	AsLink() (cid.Cid, error)
}

type SerializableNode interface {
	CID() cid.Cid
}

type MutableNode interface {
	// FUTURE: setter methods go here
}
