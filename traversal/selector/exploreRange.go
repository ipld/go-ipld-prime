package selector

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
)

// ExploreRange traverses a list, and for each element in the range specified,
// will apply a next selector to those reached nodes.
type ExploreRange struct {
	next        Selector         // selector for element we're interested in
	hasSelector map[int]struct{} // a quick check whether to return the selector
	interest    []PathSegment    // index of element we're interested in
}

// Interests for ExploreRange are all path segments within the iteration range
func (s ExploreRange) Interests() []PathSegment {
	return s.interest
}

// Explore returns the node's selector if
// the path matches an index in the range of this selector
func (s ExploreRange) Explore(n ipld.Node, p PathSegment) Selector {
	index, err := p.Index()
	if err != nil {
		return nil
	}
	if _, ok := s.hasSelector[index]; !ok {
		return nil
	}
	return s.next
}

// Decide always returns false because this is not a matcher
func (s ExploreRange) Decide(n ipld.Node) bool {
	return false
}

// ParseExploreRange assembles a Selector
// from a ExploreRange selector node
func ParseExploreRange(n ipld.Node) (Selector, error) {
	if n.ReprKind() != ipld.ReprKind_Map {
		return nil, fmt.Errorf("selector spec parse rejected: selector body must be a map")
	}
	startNode, err := n.TraverseField(startKey)
	if err != nil || startNode.ReprKind() != ipld.ReprKind_Int {
		return nil, fmt.Errorf("selector spec parse rejected: start field must an int in ExploreRange selector")
	}
	startValue, err := startNode.AsInt()
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: start field must an int in ExploreRange selector")
	}
	endNode, err := n.TraverseField(endKey)
	if err != nil || endNode.ReprKind() != ipld.ReprKind_Int {
		return nil, fmt.Errorf("selector spec parse rejected: end field must an int in ExploreRange selector")
	}
	endValue, err := endNode.AsInt()
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: end field must an int in ExploreRange selector")
	}
	next, err := n.TraverseField(nextSelectorKey)
	if err != nil {
		return nil, fmt.Errorf("selector spec parse rejected: next field must be present in ExploreRange selector")
	}
	selector, err := ParseSelector(next)
	if err != nil {
		return nil, err
	}
	x := ExploreRange{
		selector,
		make(map[int]struct{}, endValue-startValue),
		make([]PathSegment, 0, endValue-startValue),
	}
	for i := startValue; i < endValue; i++ {
		x.interest = append(x.interest, PathSegmentInt{I: i})
		x.hasSelector[i] = struct{}{}
	}
	return x, nil
}
