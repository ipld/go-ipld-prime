package datamodel

// AmendFn takes a Node and returns a NodeAmender that stores any applied transformations. The returned NodeAmender
// allows further transformations to be applied to the Node under construction.
type AmendFn func(Node) (NodeAmender, error)

// NodeAmender adds to NodeBuilder the ability to transform all or part of a Node under construction.
type NodeAmender interface {
	NodeBuilder

	// Transform takes in a Node (or a child Node of a recursive node) along with a transformation function that returns
	// a new NodeAmender with the transformed results.
	//
	// Transform returns the previous state of the target Node.
	Transform(path Path, transform AmendFn) (Node, error)
}

// containerAmender is an internal type for representing the interface for amendable containers (like maps and lists)
type containerAmender interface {
	Empty() bool
	Length() int64
	Clear()
	Values() (Node, error) // returns a list Node with the values

	NodeAmender
}

// MapAmender adds a map-like interface to NodeAmender
type MapAmender interface {
	Put(key string, value Node) error
	Get(key string) (Node, error)
	Remove(key string) (bool, error)
	Keys() (Node, error) // returns a list Node with the keys

	containerAmender
}

// ListAmender adds a list-like interface to NodeAmender
type ListAmender interface {
	Get(idx int64) (Node, error)
	Remove(idx int64) error
	// Append will add Node(s) to the end of the list. It can accept a list Node with multiple values to append.
	Append(value Node) error
	// Insert will add Node(s) at the specified index and shift subsequent elements to the right. It can accept a list
	// Node with multiple values to insert.
	// Passing an index equal to the length of the list will add Node(s) to the end of the list like Append.
	Insert(idx int64, value Node) error
	// Set will add Node(s) at the specified index and shift subsequent elements to the right. It can accept a list Node
	// with multiple values to insert.
	// Passing an index equal to the length of the list will add Node(s) to the end of the list like Append.
	// Set is different from Insert in that it will start its insertion at the specified index, overwriting it in the
	// process, while Insert will only add the Node(s).
	Set(idx int64, value Node) error

	containerAmender
}
