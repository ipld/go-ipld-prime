package selector

import (
	"fmt"
	"strconv"

	ipld "github.com/ipld/go-ipld-prime"
)

// Selector is the programmatic representation of an IPLD Selector Node
// and can be applied to traverse a given IPLD DAG
type Selector interface {
	Interests() []PathSegment                // returns the segments we're likely interested in **or nil** if we're a high-cardinality or expression based matcher and need all segments proposed to us.
	Explore(ipld.Node, PathSegment) Selector // explore one step -- iteration comes from outside (either whole node, or by following suggestions of Interests).  returns nil if no interest.  you have to traverse to the next node yourself (the selector doesn't do it for you because you might be considering multiple selection reasons at the same time).
	Decide(ipld.Node) bool
}

// ParseSelector creates a Selector that can be traversed from an IPLD Selector node
func ParseSelector(n ipld.Node) (Selector, error) {
	if n.ReprKind() != ipld.ReprKind_Map {
		return nil, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map")
	}
	if n.Length() != 1 {
		return nil, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be single-entry map")
	}
	kn, v, _ := n.MapIterator().Next()
	kstr, _ := kn.AsString()
	// Switch over the single key to determine which selector body comes next.
	//  (This switch is where the keyed union discriminators concretely happen.)
	switch kstr {
	case exploreFieldsKey:
		return ParseExploreFields(v)
	case exploreAllKey:
		return ParseExploreAll(v)
	case exploreIndexKey:
		return ParseExploreIndex(v)
	case exploreRangeKey:
		return ParseExploreRange(v)
	case exploreUnionKey:
		return ParseExploreUnion(v)
	case matcherKey:
		return ParseMatcher(v)
	default:
		return nil, fmt.Errorf("selector spec parse rejected: %q is not a known member of the selector union", kstr)
	}
}

// PathSegment can describe either an index in a list or a key in a map, as either int or a string
type PathSegment interface {
	String() string
	Index() (int, error)
}

// PathSegmentString represents a PathSegment with an underlying string
type PathSegmentString struct {
	S string
}

// PathSegmentInt represents a PathSegment with an underlying int
type PathSegmentInt struct {
	I int
}

func (ps PathSegmentString) String() string {
	return ps.S
}

// Index attempts to parse a string as an int for a PathSegmentString
func (ps PathSegmentString) Index() (int, error) {
	return strconv.Atoi(ps.S)
}

func (ps PathSegmentInt) String() string {
	return strconv.Itoa(ps.I)
}

// Index is always just the underlying int for a PathSegmentInt
func (ps PathSegmentInt) Index() (int, error) {
	return ps.I, nil
}
