package selector

import (
	ipld "github.com/ipld/go-ipld-prime"
)

// implementation note: union selectors can be generated at selector evaluation time!
// for example, globstar selectors do this: `**/foo` implicitly generates a union
//  after the first depth in order to check deeper star as well as the 'foo' match.
//
// imagine this example:
//
// a selector like `**/?oo/bar` is applied...
//
// ```
// ./zot/ -- prefix matches **
// ./zot/zoo/ -- prefix matches **, prefix matches **/?oo
// ./zot/zoo/foo -- prefix matches **, prefix matches **/?oo
// ./zot/zoo/foo/bar -- prefix matches **, FULL MATCH **/?oo/bar
// ./zot/zoo/foo/bar/baz -- prefix matches **
// ```
//
// as you can see, a union selector reasonably expresses the intermediate state
// needed during handling several of these paths.

// SelectUnion combines two or more other selectors and aggregates their behavior;
// if something is matched by any of the composed selectors, it's matched by the union.
type SelectUnion struct {
	Members []Selector
}

func (s SelectUnion) Interests() []PathSegment {
	// Check for any high-cardinality selectors first; if so, shortcircuit.
	//  (n.b. we're assuming the 'Interests' method is cheap here.)
	for _, m := range s.Members {
		if m.Interests() == nil {
			return nil
		}
	}
	// Accumulate the whitelist of interesting path segments.
	v := []PathSegment{}
	for _, m := range s.Members {
		v = append(v, m.Interests()...)
	}
	return v
}

func (s SelectUnion) Explore(n ipld.Node, p PathSegment) Selector {
	// this needs to call Explore for each member,
	//  and if more than one member returns a selector,
	//   we compose them into a new union automatically and return that.
	panic("TODO")
}

func (s SelectUnion) Decide(n ipld.Node) bool {
	for _, m := range s.Members {
		if m.Decide(n) {
			return true
		}
	}
	return false
}
