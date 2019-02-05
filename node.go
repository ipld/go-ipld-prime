package ipld

import (
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

	// Keys returns an iterator which will yield keys for traversing the node.
	// If the node kind is anything other than a map, the iterator will
	// yield error values.
	Keys() KeyIterator

	// KeysImmediate returns a slice containing all keys for traversing the node.
	// The semantics are otherwise identical to using the Keys() iterator.
	//
	// KeysImmediate is for convenience of usage; callers should prefer to use
	// the iterator approach where possible, as it continues to behave well
	// even when using collections of extremely large size (and even when
	// the collection is split between multiple serial nodes, as with
	// Advanced Layouts, etc).
	KeysImmediate() ([]string, error)

	// Length returns the length of a list, or the number of entries in a map,
	// or -1 if the node is not of list nor map kind.
	Length() int

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
}

// KeyIterator is an interface for traversing nodes of kind map.
// Sequential calls to Next() will yield keys; HasNext() describes whether
// iteration should continue.
//
// Iteration order is defined to be stable.
//
// REVIEW: should Next return error?
// Other parts of the Node interface use that for kind mismatch rejection;
// so on those grounds, I'd say "no", because we know what the key kind is
// (but then Node.Keys should return error).
// In big nodes (composites using an AdvLayout), where do we return errors?
// Since we might be streaming, there are questions here.
type KeyIterator interface {
	Next() (string, error)
	HasNext() bool
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

// REVIEW: immediate-mode AsBytes() method (as opposed to e.g. returning
// an io.Reader instance) might be problematic, esp. if we introduce
// AdvancedLayouts which support large bytes natively.
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
