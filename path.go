package ipld

import (
	"strings"
)

// Path is used in describing progress in a traversal;
// and can also be used as an instruction for a specific traverse.
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
// (Note: escaping may be specified and supported in the future; currently, it is not.)
//
type Path struct {
	segments []PathSegment
}

// ParsePath converts a string to an IPLD Path, parsing the string into a segmented Path.
//
// Each segment of the path string should be separated by a "/" character.
//
// Multiple subsequent "/" characters will be silently collapsed.
// E.g., `"foo///bar"` will be treated equivalently to `"foo/bar"`.
// Prefixed and suffixed extraneous "/" characters are also discarded.
//
// No "cleaning" of the path occurs.  See the documentation of the Path struct;
// in particular, note that ".." does not mean "go up", nor does "." mean "stay here" --
// correspondingly, there isn't anything to "clean".
func ParsePath(pth string) Path {
	// FUTURE: we should probably have some escaping mechanism which makes
	//  it possible to encode a slash in a segment.  Specification needed.
	ss := strings.FieldsFunc(pth, func(r rune) bool { return r == '/' })
	ssl := len(ss)
	p := Path{make([]PathSegment, ssl)}
	for i := 0; i < ssl; i++ {
		p.segments[i] = PathSegmentOfString(ss[i])
	}
	return p
}

// String representation of a Path is simply the join of each segment with '/'.
// It does not include a leading nor trailing slash.
func (p Path) String() string {
	l := len(p.segments)
	if l == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i := 0; i < l-1; i++ {
		sb.WriteString(p.segments[i].String())
		sb.WriteByte('/')
	}
	sb.WriteString(p.segments[l-1].String())
	return sb.String()
}

// Segements returns a slice of the path segment strings.
//
// It is not lawful to mutate nor append the returned slice.
func (p Path) Segments() []PathSegment {
	return p.segments
}

// Join creates a new path composed of the concatenation of this and the given path's segments.
func (p Path) Join(p2 Path) Path {
	combinedSegments := make([]PathSegment, len(p.segments)+len(p2.segments))
	copy(combinedSegments, p.segments)
	copy(combinedSegments[len(p.segments):], p2.segments)
	p.segments = combinedSegments
	return p
}

// AppendSegmentString is as per Join, but a shortcut when appending single segments using strings.
func (p Path) AppendSegment(ps PathSegment) Path {
	l := len(p.segments)
	combinedSegments := make([]PathSegment, l+1)
	copy(combinedSegments, p.segments)
	combinedSegments[l] = ps
	p.segments = combinedSegments
	return p
}

// AppendSegmentString is as per Join, but a shortcut when appending single segments using strings.
func (p Path) AppendSegmentString(ps string) Path {
	return p.AppendSegment(PathSegmentOfString(ps))
}

// Parent returns a path with the last of its segments popped off (or
// the zero path if it's already empty).
func (p Path) Parent() Path {
	if len(p.segments) == 0 {
		return Path{}
	}
	return Path{p.segments[0 : len(p.segments)-1]}
}

// Truncate returns a path with only as many segments remaining as requested.
func (p Path) Truncate(i int) Path {
	return Path{p.segments[0:i]}
}
