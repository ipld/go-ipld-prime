package selector

import (
	ipld "github.com/ipld/go-ipld-prime"
)

// SelectAll is a non-recursive kleene-star match (e.g., it's `./*`).
// If SelectAll is a leaf in a Selector tree, it will match all content;
// if it has a 'next' selector (e.g., it's like `./*/foo`), it'll yield
// that next selector for explore of any and all pathsegments.
type SelectAll struct {
	next Selector // set to SelectTrue at parse time if appropriate.
}

func (s SelectAll) Interests() []PathSegment {
	return nil
}

func (s SelectAll) Explore(n ipld.Node, p PathSegment) Selector {
	return s.next
}

func (s SelectAll) Decide(n ipld.Node) bool {
	return false // this is an intermediate selector: it doesn't itself call for a thing, only indirectly does so by sometimes returning SelectTrue.
}
