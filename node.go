package ipld

type Node interface {
	// Kind returns a value from the ReprKind enum describing what the
	// essential serializable kind of this node is (map, list, int, etc).
	// Most other handling of a node requires first switching upon the kind.
	ReprKind() ReprKind

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

	// MapIterator returns an iterator which yields key-value pairs
	// traversing the node.
	// If the node kind is anything other than a map, the iterator will
	// yield error values.
	//
	// The iterator will yield every entry in the map; that is, it
	// can be expected that itr.Next will be called node.Length times
	// before itr.Done becomes true.
	MapIterator() MapIterator

	// ListIterator returns an iterator which yields key-value pairs
	// traversing the node.
	// If the node kind is anything other than a list, the iterator will
	// yield error values.
	//
	// The iterator will yield every entry in the list; that is, it
	// can be expected that itr.Next will be called node.Length times
	// before itr.Done becomes true.
	ListIterator() ListIterator

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

// MapIterator is an interface for traversing map nodes.
// Sequential calls to Next() will yield key-value pairs;
// Done() describes whether iteration should continue.
//
// Iteration order is defined to be stable: two separate MapIterator
// created to iterate the same Node will yield the same key-value pairs
// in the same order.
// The order itself may be defined by the Node implementation: some
// Nodes may retain insertion order, and some may return iterators which
// always yield data in sorted order, for example.
type MapIterator interface {
	// Next returns the next key-value pair.
	//
	// An error value can also be returned at any step: in the case of advanced
	// data structures with incremental loading, it's possible to encounter
	// cancellation or I/O errors at any point in iteration.
	// If an error is returned, the boolean will always be false (so it's
	// correct to check the bool first and short circuit to continuing if true).
	// If an error is returned, the key and value may be nil.
	Next() (key Node, value Node, err error)

	// Done returns false as long as there's at least one more entry to iterate.
	// When Done returns true, iteration can stop.
	//
	// Implementers of iterators for advanced data layouts (e.g. more than
	// one chunk of backing data, which is loaded incrementally), if your
	// implementation does any I/O during the Done method, and it encounters
	// an error, it must return 'true', so that the following Next call
	// has an opportunity to return the error.
	Done() bool
}

// ListIterator is an interface for traversing list nodes.
// Sequential calls to Next() will yield index-value pairs;
// Done() describes whether iteration should continue.
//
// A loop which iterates from 0 to Node.Length is a valid
// alternative to using a ListIterator.
type ListIterator interface {
	// Next returns the next index and value.
	//
	// An error value can also be returned at any step: in the case of advanced
	// data structures with incremental loading, it's possible to encounter
	// cancellation or I/O errors at any point in iteration.
	// If an error is returned, the boolean will always be false (so it's
	// correct to check the bool first and short circuit to continuing if true).
	// If an error is returned, the key and value may be nil.
	Next() (idx int, value Node, err error)

	// Done returns false as long as there's at least one more entry to iterate.
	// When Done returns false, iteration can stop.
	//
	// Implementers of iterators for advanced data layouts (e.g. more than
	// one chunk of backing data, which is loaded incrementally), if your
	// implementation does any I/O during the Done method, and it encounters
	// an error, it must return 'true', so that the following Next call
	// has an opportunity to return the error.
	Done() bool
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
