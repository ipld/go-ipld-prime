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
	Values() (Node, error) // returns a list node with the values

	NodeAmender
}

// MapAmender adds a map-like interface to NodeAmender
type MapAmender interface {
	Put(key Node, value Node) error
	Get(key Node) (Node, error)
	Remove(key Node) (bool, error)
	Keys() (Node, error) // returns a list node with the keys

	containerAmender
}

// ListAmender adds a list-like interface to NodeAmender
type ListAmender interface {
	Get(idx int64) (Node, error)
	Remove(idx int64) error
	Append(values Node) error            // accepts a list node
	Insert(idx int64, values Node) error // accepts a list node
	Set(idx int64, value Node) error

	containerAmender
}
