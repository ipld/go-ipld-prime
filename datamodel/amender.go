package datamodel

// NodeAmender layers onto NodeBuilder the ability to transform all or part of a Node under construction.
type NodeAmender interface {
	NodeBuilder

	// Transform takes in a Node (or a child Node of a recursive node) along with a transformation function that returns
	// a new NodeAmender with the transformed results.
	//
	// Transform returns the previous state of the target Node.
	Transform(path Path, transform func(Node) (NodeAmender, error)) (Node, error)
}
