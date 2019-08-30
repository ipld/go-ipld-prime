package ipld

import (
	"strconv"
)

// PathSegment can describe either a key in a map, or an index in a list.
//
// Path segments are "stringly typed" -- they may be interpreted as either strings or ints depending on context.
// A path segment of "123" will be used as a string when traversing a node of map kind;
// and it will be converted to an integer when traversing a node of list kind.
// (If a path segment string cannot be parsed to an int when traversing a node of list kind, then traversal will error.)
// It is not possible to ask which kind (string or integer) a PathSegment is, because that is not defined -- this is *only* intepreted contextually.
//
// Internally, PathSegment will store either a string or an integer,
// depending on how it was constructed,
// and will automatically convert to the other on request.
// (This means if two pieces of code communicate using PathSegment, one producing ints and the other expecting ints, they will work together efficiently.)
// PathSegment in a Path produced by ParsePath generally have all strings internally,
// because there is distinction possible when parsing a Path string
// (and attempting to pre-parse all strings into ints "in case" would waste time in almost all cases).
type PathSegment struct {
	/*
		A quick implementation note about the Go compiler and "union" semantics:

		There are roughly two ways to do "union" semantics in Go.

		The first is to make a struct with each of the values.

		The second is to make an interface and use an unexported method to keep it closed.

		The second tactic provides somewhat nicer semantics to the programmer.
		(Namely, it's clearly impossible to have two inhabitants, which is... the point.)
		The downside is... putting things in interfaces generally incurs an allocation
		(grep your assembly output for "runtime.conv*").

		The first tactic looks kludgier, and would seem to waste memory
		(the struct reserves space for each possible value, even though the semantic is that only one may be non-zero).
		However, in most cases, more *bytes* are cheaper than more *allocs* --
		garbage collection costs are domininated by alloc count, not alloc size.

		Because PathSegment is something we expect to put in fairly "hot" paths,
		we're using the first tactic.

		(We also currently get away with having no extra discriminator bit
		because empty string is not considered a valid segment,
		and thus we can use it as a sentinel value.
		This may change if the IPLD Path spec comes to other conclusions about this.)
	*/

	s string
	i int
}

// ParsePathSegment parses a string into a PathSegment,
// handling any escaping if present.
// (Note: there is currently no escaping specified for PathSegments,
// so this is currently functionally equivalent to PathSegmentOfString.)
func ParsePathSegment(s string) PathSegment {
	return PathSegment{s: s}
}

// PathSegmentOfString boxes a string into a PathSegement.
// It does not attempt to parse any escaping; use ParsePathSegment for that.
func PathSegmentOfString(s string) PathSegment {
	return PathSegment{s: s}
}

// PathSegmentOfString boxes an int into a PathSegement.
func PathSegmentOfInt(i int) PathSegment {
	return PathSegment{i: i}
}

// containsString is unexported because we use it to see what our *storage* form is,
// but this is considered an implementation detail that's non-semantic.
// If it returns false, it implicitly means "containsInt", as these are the only options.
func (ps PathSegment) containsString() bool {
	return ps.s != ""
}

// String returns the PathSegment as a string.
func (ps PathSegment) String() string {
	switch ps.containsString() {
	case true:
		return ps.s
	case false:
		return strconv.Itoa(ps.i)
	}
	panic("unreachable")
}

// Index returns the PathSegment as an int,
// or returns an error if the segment is a string that can't be parsed as an int.
func (ps PathSegment) Index() (int, error) {
	switch ps.containsString() {
	case true:
		return strconv.Atoi(ps.s)
	case false:
		return ps.i, nil
	}
	panic("unreachable")
}

// Equals checks if two PathSegment values are equal.
// This is equivalent to checking if their strings are equal --
// if one of the PathSegment values is backed by an int and the other is a string,
// they may still be "equal".
func (x PathSegment) Equals(o PathSegment) bool {
	if !x.containsString() && !o.containsString() {
		return x.i == o.i
	}
	return x.String() == o.String()
}
