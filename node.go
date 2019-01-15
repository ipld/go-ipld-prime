package ipld

import (
	"github.com/polydawn/refmt/shared"

	"github.com/ipfs/go-cid"
)

type Node interface {
	// Kind returns a value from the ReprKind enum describing what the
	// essential serializable kind of this node is (map, list, int, etc).
	// Most other handling of a node requires first switching upon the kind.
	Kind() ReprKind

	// GetField resolves a path in the the object and returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	//
	// If the Kind of this Node is not ReprKind_Map, a nil node and an error
	// will be returned.
	//
	// If the key does not exist, a nil node and an error will be returned.
	TraverseField(key string) (Node, error)

	// GetIndex is the equivalent of GetField but for indexing into an array
	// (or a numerically-keyed map).  Like GetField, it returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	//
	// If the Kind of this Node is not ReprKind_List, a nil node and an error
	// will be returned.
	//
	// If idx is out of range, a nil node and an error will be returned.
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
	// (But since this is only relevant for `typed.Node`, we can make that
	// choice locally to that package.)
	//IsUndefined() bool

	IsNull() bool
	AsBool() (bool, error)
	AsInt() (int, error)
	AsFloat() (float64, error)
	AsString() (string, error)
	AsBytes() ([]byte, error)
	AsLink() (cid.Cid, error)

	// PushTokens converts this node and its children into a stream of refmt Tokens
	// and push them sequentially into the given TokenSink.
	// This is useful for serializing, or hashing, or feeding to any other
	// TokenSink (for example, converting to other ipld.Node implementations
	// which can construct themselves from a token stream).
	PushTokens(sink shared.TokenSink) error
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
	SetNull()
	SetBool(v bool)
	SetInt(v int)
	SetFloat(v float64)
	SetString(v string)
	SetBytes(v []byte)
	SetLink(v cid.Cid)
}

// REVIEW: having a an immediate-mode Keys() method rather than iterator
// might actually be a bad idea.  We're aiming to reuse this interface
// for *advanced layouts as well*, and those can be *large*.
//
// Similar goes for AsBytes().
//
// Probable solution is having both immediate and iterator return methods.
// Returning a reader for bytes when you know you want a slice already
// is going to be high friction without purpose in many common uses.
//
// Unclear what SetByteStream() would look like for advanced layouts.
// One could try to encapsulate the chunking entirely within the advlay
// node impl... but would it be graceful?  Not sure.  Maybe.  Hopefully!
// Yes?  The advlay impl would still tend to use SetBytes for the raw
// data model layer nodes its composing, so overall, it shakes out nicely.
