package selector

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
)

// SelectFields selects some fields by name (or index),
// and may contain more nested selectors per field.
//
// If you're familiar with GraphQL queries, you can thing of SelectFields
// as similar to the basic unit of composition in GraphQL queries.
//
// SelectFields also works for selecting specific elements out of a list;
// if the "field" is a base-10 int, it will be coerced and do the right thing.
// SelectIndexes is more appropriate, however, and should be preferred.
type SelectFields struct {
	selections map[string]Selector
	interests  []PathSegment // keys of above; already boxed as that's the only way we consume them
}

func (s SelectFields) Interests() []PathSegment {
	return s.interests
}

func (s SelectFields) Explore(n ipld.Node, p PathSegment) Selector {
	return s.selections[p.String()]
}

func (s SelectFields) Decide(n ipld.Node) bool {
	return false // this is an intermediate selector: it doesn't itself call for a thing, only indirectly does so by sometimes returning SelectTrue.
}

func ParseSelectFields(n ipld.Node) (Selector, error) {
	if n.ReprKind() != ipld.ReprKind_Map {
		return nil, fmt.Errorf("selector spec parse rejected: selector body must be a map")
	}
	x := SelectFields{
		make(map[string]Selector, n.Length()),
		make([]PathSegment, 0, n.Length()),
	}
	for itr := n.MapIterator(); !itr.Done(); {
		kn, v, err := itr.Next()
		if err != nil {
			return nil, fmt.Errorf("error during selector spec parse: %s", err)
		}

		kstr, _ := kn.AsString()
		x.interests = append(x.interests, PathSegmentString{kstr})
		switch v.ReprKind() {
		case ipld.ReprKind_Map: // deeper!
			x.selections[kstr], err = ParseSelector(v)
			if err != nil {
				return nil, err
			}
		case ipld.ReprKind_Bool:
			b, _ := v.AsBool()
			if !b {
				// FUTURE: boolean-as-unit is not currently expressible in the schema spec; might be something we want, just for human ergonomics.
				return nil, fmt.Errorf("selector spec parse rejected: entries in selectFields must be either a nested selector or the value 'true'")
			}
			x.selections[kstr] = SelectTrue{}
		}
	}
	return x, nil
}
