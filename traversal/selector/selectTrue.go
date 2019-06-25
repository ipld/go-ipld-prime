package selector

import (
	ipld "github.com/ipld/go-ipld-prime"
)

// SelectAll is a dummy selector that other selectors can return to say
// "the content at this path?  definitely this".
type SelectTrue struct{}

func (s SelectTrue) Interests() []PathSegment {
	return []PathSegment{}
}

func (s SelectTrue) Explore(n ipld.Node, p PathSegment) Selector {
	return nil
}

func (s SelectTrue) Decide(n ipld.Node) bool {
	return true
}
