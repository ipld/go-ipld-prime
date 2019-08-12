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
	maxDepth := s.maxDepth
	if nextSelector == nil {
		return nil
	}
	if !s.hasRecursiveEdge(nextSelector) {
		return ExploreRecursive{s.sequence, nextSelector, maxDepth}
	}
	if maxDepth < 2 {
		return s.replaceRecursiveEdge(nextSelector, nil)
	}
	return ExploreRecursive{s.sequence, s.replaceRecursiveEdge(nextSelector, s.sequence), s.maxDepth - 1}
}

func (s ExploreRecursive) hasRecursiveEdge(nextSelector Selector) bool {
	_, isRecursiveEdge := nextSelector.(ExploreRecursiveEdge)
	if isRecursiveEdge {
		return true
	}
	exploreUnion, isUnion := nextSelector.(ExploreUnion)
	if isUnion {
		for _, selector := range exploreUnion.Members {
			if s.hasRecursiveEdge(selector) {
				return true
			}
		}
	}
	return false
}

func (s ExploreRecursive) replaceRecursiveEdge(nextSelector Selector, replacement Selector) Selector {
	_, isRecursiveEdge := nextSelector.(ExploreRecursiveEdge)
	if isRecursiveEdge {
		return replacement
	}
	exploreUnion, isUnion := nextSelector.(ExploreUnion)
	if isUnion {
		replacementMembers := make([]Selector, 0, len(exploreUnion.Members))
		for _, selector := range exploreUnion.Members {
			newSelector := s.replaceRecursiveEdge(selector, replacement)
			if newSelector != nil {
				replacementMembers = append(replacementMembers, newSelector)
			}
		}
		if len(replacementMembers) == 0 {
			return nil
		}
		if len(replacementMembers) == 1 {
			return replacementMembers[0]
		}
		return ExploreUnion{replacementMembers}
	}
	return nextSelector
}

// Decide always returns false because this is not a matcher
func (s ExploreRecursive) Decide(n ipld.Node) bool {
	return s.current.Decide(n)
}

type exploreRecursiveContext struct {
	edgesFound int
}

func (erc *exploreRecursiveContext) Link(s Selector) bool {
	_, ok := s.(ExploreRecursiveEdge)
	if ok {
		erc.edgesFound++
	}
	return ok
}

// ParseExploreRecursive assembles a Selector from a ExploreRecursive selector node
func (pc ParseContext) ParseExploreRecursive(n ipld.Node) (Selector, error) {
	if n.ReprKind() != ipld.ReprKind_Map {
		return nil, fmt.Errorf("selector spec parse rejected: selector body must be a map")
	}

	maxDepthNode, err := n.LookupString(SelectorKey_MaxDepth)
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: maxDepth field must be present in ExploreRecursive selector")
	}
	maxDepthValue, err := maxDepthNode.AsInt()
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: maxDepth field must be a number in ExploreRecursive selector")
	}
	sequence, err := n.LookupString(SelectorKey_Sequence)
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: sequence field must be present in ExploreRecursive selector")
	}
	erc := &exploreRecursiveContext{}
	selector, err := pc.PushParent(erc).ParseSelector(sequence)
	if err != nil {
		return nil, err
	}
	if erc.edgesFound == 0 {
		return nil, fmt.Errorf("selector spec parse rejected: ExploreRecursive must have at least one ExploreRecursiveEdge")
	}
	return ExploreRecursive{selector, selector, maxDepthValue}, nil
}
