package ipld

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	_ Traversal = Path{}.Traverse // (type assertion)
)

// Path represents a MerklePath.  TODO:standards-doc-link.
//
// IPLD Paths can only go down: that is, each segment must traverse one node.
// There is no ".." which means "go up";
// and there is no "." which means "stay here";
// and it is not valid to have an empty path segment.
//
// (Note: path strings as interpreted by UnixFS may certainly have concepts
// of ".." and "."!  But UnixFS is built upon IPLD; IPLD has no idea of this.)
//
// Paths are representable as strings.  When represented as a string, each
// segment is separated by a "/" character.
// (It follows that path segments may not themselves contain a "/" character.)
//
// Path segments are stringly typed.  A path segment of "123" will be used
// as a string when traversing a node of map kind; and it will be converted
// to an integer when traversing a node of list kind.
// (If a path segment string cannot be parsed to an int when traversing a node
// of list kind, then traversal will error.)
type Path struct {
	segments []string
}

// ParsePath converts a string to an IPLD Path, parsing the string
// into a segemented Path.
//
// Each segment of the path string should be separated by a "/" character.
//
// Multiple subsequent "/" characters will be silently collapsed.
// E.g., `"foo///bar"` will be treated equivalently to `"foo/bar"`.
//
// No "cleaning" of the path occurs.  See the documentation of the Path
// struct; in particular, note that ".." does not mean "go up" -- so
// correspondingly, there is nothing to "clean".
func ParsePath(pth string) Path {
	// FUTURE: we should probably have some escaping mechanism which makes
	//  it possible to encode a slash in a segment.  Specification needed.
	return Path{strings.FieldsFunc(pth, func(r rune) bool { return r == '/' })}
}

// String representation of a Path is simply the join of each segment with '/'.
func (p Path) String() string {
	return strings.Join(p.segments, "/")
}

// Segements returns a slice of the path segment strings.
//
// It is not lawful to mutate the returned slice.
func (p Path) Segments() []string {
	return p.segments
}

// Join creates a new path composed of the concatenation of this and the
// given path's segments.
func (p Path) Join(p2 Path) Path {
	combinedSegments := make([]string, len(p.segments)+len(p2.segments))
	copy(combinedSegments, p.segments)
	copy(combinedSegments[len(p.segments):], p2.segments)
	p.segments = combinedSegments
	return p
}

// Path.Traverse is an implementation of Traversal that makes a simple
// direct walk over a sequence of nodes, using each segment of the path
// to get the next node until all path segments have been consumed.
//
// If one of the node traverse steps returns an error, that node and the
// path so far including that node will be returned, as well as the error.
func (p Path) Traverse(tp TraversalProgress, start Node) (_ TraversalProgress, reached Node, err error) {
	for i, seg := range p.segments {
		switch start.Kind() {
		case ReprKind_Invalid:
			return TraversalProgress{}, nil, fmt.Errorf("cannot traverse node at %q: it is undefined", Path{p.segments[0:i]})
		case ReprKind_Map:
			next, err := start.TraverseField(seg)
			if err != nil {
				return TraversalProgress{}, nil, fmt.Errorf("error traversing node at %q: %s", Path{p.segments[0:i]}, err)
			}
			start = next
		case ReprKind_List:
			intSeg, err := strconv.Atoi(seg)
			if err != nil {
				return TraversalProgress{}, nil, fmt.Errorf("cannot traverse node at %q: the next path segment (%q) cannot be parsed as a number and the node is a list", Path{p.segments[0:i]}, seg)
			}
			next, err := start.TraverseIndex(intSeg)
			if err != nil {
				return TraversalProgress{}, nil, fmt.Errorf("error traversing node at %q: %s", Path{p.segments[0:i]}, err)
			}
			start = next
		default:
			return TraversalProgress{}, nil, fmt.Errorf("error traversing node at %q: %s", Path{p.segments[0:i]}, fmt.Errorf("cannot traverse terminals"))
		}
	}
	tp.Path = tp.Path.Join(p)
	return tp, start, nil
}
