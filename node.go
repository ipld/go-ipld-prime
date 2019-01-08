package ipld

import "github.com/ipfs/go-cid"

type Node interface {
	// Kind returns a value from the ReprKind enum describing what the
	// essential serializable kind of this node is (map, list, int, etc).
	// Most other handling of a node requires first switching upon the kind.
	Kind() ReprKind

	// GetField resolves a path in the the object and returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	TraverseField(path string) (Node, error)

	// GetIndex is the equivalent of GetField but for indexing into an array
	// (or a numerically-keyed map).  Like GetField, it returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	TraverseIndex(idx int) (Node, error)

	// Keys returns instructions for traversing the node.
	// If the node kind is a map, the keys slice has content;
	// if it's a list, the length int will be positive
	// (and if it's a zero length list, there's not to traverse, right?);
	// and if it's a primitive type the returned values are nil and zero.
	Keys() ([]string, int)

	// Undefined nodes are returned when traversing a struct field that is
	// defined by a schema but unset in the data.  (Undefined nodes are not
	// possible otherwise; you'll only see them from `typed.Node`.)
	//
	// REVIEW: unsure if this should use ReprKind_Invalid or another enum.
	// Since ReprKind_Invalid is returned for UnionStyle_Kinded, confusing.
	//IsUndefined() bool

	IsNull() bool
	AsBool() (bool, error)
	AsString() (string, error)
	AsInt() (int, error)
	AsLink() (cid.Cid, error)
}

type SerializableNode interface {
	CID() cid.Cid
}

// MutableNode is an interface which allows setting its value.
// MutableNode can be all the same kinds as Node can, and can also
// be freely coerced between kinds.
//
// Using a method which coerces a Mutable node to a new kind can
// discard data; the SetField and SetIndex methods are concatenative,
// but using any of the Set{Scalar} methods will result in any map or
// array content being discarded(!).
type MutableNode interface {
	Node
	SetField(k string, v Node) // SetField coerces the node to a map kind and sets a key:val pair.
	SetIndex(k int, v Node)    // SetIndex coerces the node to an array kind and sets an index:val pair.  (It will implicitly increase the size to include the index.)
	SetBool(v bool)
	SetString(v string)
	SetInt(v int)
	SetLink(v cid.Cid)
}
