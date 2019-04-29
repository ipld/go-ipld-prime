package selector

import (
	"fmt"
	"strconv"

	ipld "github.com/ipld/go-ipld-prime"
)

type Selector interface {
	Interests() []PathSegment                // returns the segments we're likely interested in **or nil** if we're a high-cardinality or expression based matcher and need all segments proposed to us.
	Explore(ipld.Node, PathSegment) Selector // explore one step -- iteration comes from outside (either whole node, or by following suggestions of Interests).  returns nil if no interest.  you have to traverse to the next node yourself (the selector doesn't do it for you because you might be considering multiple selection reasons at the same time).
	Decide(ipld.Node) bool
}

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
	case "f":
		return ParseSelectFields(v)
	// FUTURE:
	// case "a":
	//	return ParseSelectAll(v)
	// case "i":
	//	return ParseSelectIndexes(v)
	// case "r":
	//	return ParseSelectRange(v)
	// case "+":
	//	return ParseSelectTrue(v)
	default:
		return nil, fmt.Errorf("selector spec parse rejected: %q is not a known member of the selector union", kstr)
	}
}

type PathSegment interface {
	String() string
	Index() (int, error)
}

type PathSegmentString struct {
	S string
}
type PathSegmentInt struct {
	I int
}

func (ps PathSegmentString) String() string {
	return ps.S
}
func (ps PathSegmentString) Index() (int, error) {
	return strconv.Atoi(ps.S)
}
func (ps PathSegmentInt) String() string {
	return strconv.Itoa(ps.I)
}
func (ps PathSegmentInt) Index() (int, error) {
	return ps.I, nil
}
