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
