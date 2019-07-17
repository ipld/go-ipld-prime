package selector

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
)

// ExploreRecursiveEdge is a special sentinel value which is used to mark
// the end of a sequence started by an ExploreRecursive selector: the recursion
// goes back to the initial state of the earlier ExploreRecursive selector,
// and proceeds again (with a decremented maxDepth value).
//
// An ExploreRecursive selector that doesn't contain an ExploreRecursiveEdge
// is nonsensical.  Containing more than one ExploreRecursiveEdge is valid.
// An ExploreRecursiveEdge without an enclosing ExploreRecursive is an error.
type ExploreRecursiveEdge struct{}

// Interests should ultimately never get called for an ExploreRecursiveEdge selector
func (s ExploreRecursiveEdge) Interests() []PathSegment {
	return []PathSegment{}
}

// Explore should ultimately never get called for an ExploreRecursiveEdge selector
func (s ExploreRecursiveEdge) Explore(n ipld.Node, p PathSegment) Selector {
	return nil
}

// Decide should ultimately never get called for an ExploreRecursiveEdge selector
func (s ExploreRecursiveEdge) Decide(n ipld.Node) bool {
	return false
}

// ParseExploreRecursiveEdge assembles a Selector
// from a exploreRecursiveEdge selector node
func ParseExploreRecursiveEdge(n ipld.Node) (Selector, error) {
	if n.ReprKind() != ipld.ReprKind_Map {
		return nil, fmt.Errorf("selector spec parse rejected: selector body must be a map")
	}
	return ExploreRecursiveEdge{}, nil
}
