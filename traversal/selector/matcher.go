package selector

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
)

// Matcher marks a node to be included in the "result" set.
// (All nodes traversed by a selector are in the "covered" set (which is a.k.a.
// "the merkle proof"); the "result" set is a subset of the "covered" set.)
//
// In libraries using selectors, the "result" set is typically provided to
// some user-specified callback.
//
// A selector tree with only "explore*"-type selectors and no Matcher selectors
// is valid; it will just generate a "covered" set of nodes and no "result" set.
// TODO: From spec: implement conditions and labels
type Matcher struct{}

// Interests are empty for a matcher (for now) because
// It is always just there to match, not explore further
func (s Matcher) Interests() []datamodel.PathSegment {
	return []datamodel.PathSegment{}
}

// Explore will return nil because a matcher is a terminal selector
func (s Matcher) Explore(n datamodel.Node, p datamodel.PathSegment) (Selector, error) {
	return nil, nil
}

// Decide is always true for a match cause it's in the result set
// TODO: Implement boolean logic for conditionals
func (s Matcher) Decide(n datamodel.Node) bool {
	return true
}

// ParseMatcher assembles a Selector
// from a matcher selector node
// TODO: Parse labels and conditions
func (pc ParseContext) ParseMatcher(n datamodel.Node) (Selector, error) {
	if n.Kind() != datamodel.Kind_Map {
		return nil, fmt.Errorf("selector spec parse rejected: selector body must be a map")
	}
	return Matcher{}, nil
}
