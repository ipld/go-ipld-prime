package ipld

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
	AsLink() (Link, error)

	// NodeBuilder returns a NodeBuilder which can be used to build
	// new nodes of the same implementation type as this one.
	//
	// For map and list nodes, the NodeBuilder's append-oriented methods
	// will work using this node's values as a base.
	// If this is a typed node, the NodeBuilder will carry the same
	// typesystem constraints as this Node.
	//
	// (This feature is used by the traversal package, especially in
	// e.g. traversal.Transform, for doing tree updates while keeping the
	// existing implementation preferences and doing as many operations
	// in copy-on-write fashions as possible.)
	NodeBuilder() NodeBuilder
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
