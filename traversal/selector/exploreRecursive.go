package selector

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
)

// ExploreRecursive traverses some structure recursively.
// To guide this exploration, it uses a "sequence", which is another Selector
// tree; some leaf node in this sequence should contain an ExploreRecursiveEdge
// selector, which denotes the place recursion should occur.
//
// In implementation, whenever evaluation reaches an ExploreRecursiveEdge marker
// in the recursion sequence's Selector tree, the implementation logically
// produces another new Selector which is a copy of the original
// ExploreRecursive selector, but with a decremented maxDepth parameter, and
// continues evaluation thusly.
//
// It is not valid for an ExploreRecursive selector's sequence to contain
// no instances of ExploreRecursiveEdge; it *is* valid for it to contain
// more than one ExploreRecursiveEdge.
//
// ExploreRecursive can contain a nested ExploreRecursive!
// This is comparable to a nested for-loop.
// In these cases, any ExploreRecursiveEdge instance always refers to the
// nearest parent ExploreRecursive (in other words, ExploreRecursiveEdge can
// be thought of like the 'continue' statement, or end of a for-loop body;
// it is *not* a 'goto' statement).
//
// Be careful when using ExploreRecursive with a large maxDepth parameter;
// it can easily cause very large traversals (especially if used in combination
// with selectors like ExploreAll inside the sequence).
type ExploreRecursive struct {
	sequence Selector // selector for element we're interested in
	current  Selector // selector to apply to the current node
	maxDepth int
}

// Interests for ExploreRecursive is empty (meaning traverse everything)
func (s ExploreRecursive) Interests() []PathSegment {
	return s.current.Interests()
}

// Explore returns the node's selector for all fields
func (s ExploreRecursive) Explore(n ipld.Node, p PathSegment) Selector {
	nextSelector := s.current.Explore(n, p)
	if nextSelector == nil {
		return nil
	}
	_, ok := nextSelector.(ExploreRecursiveEdge)
	if !ok {
		return ExploreRecursive{s.sequence, nextSelector, s.maxDepth}
	}
	if s.maxDepth < 2 {
		return nil
	}
	return ExploreRecursive{s.sequence, s.sequence, s.maxDepth - 1}
}

// Decide always returns false because this is not a matcher
func (s ExploreRecursive) Decide(n ipld.Node) bool {
	return s.current.Decide(n)
}

// ParseExploreRecursive assembles a Selector from a ExploreRecursive selector node
func ParseExploreRecursive(n ipld.Node) (Selector, error) {
	if n.ReprKind() != ipld.ReprKind_Map {
		return nil, fmt.Errorf("selector spec parse rejected: selector body must be a map")
	}

	maxDepthNode, err := n.TraverseField(maxDepthKey)
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: maxDepth field must be present in ExploreRecursive selector")
	}
	maxDepthValue, err := maxDepthNode.AsInt()
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: maxDepth field must be a number in ExploreRecursive selector")
	}
	sequence, err := n.TraverseField(sequenceKey)
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: sequence field must be present in ExploreRecursive selector")
	}
	selector, err := ParseSelector(sequence)
	if err != nil {
		return nil, err
	}
	return ExploreRecursive{selector, selector, maxDepthValue}, nil
}
