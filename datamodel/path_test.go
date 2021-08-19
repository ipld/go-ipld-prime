package datamodel

import (
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestParsePath(t *testing.T) {
	t.Run("parsing one segment", func(t *testing.T) {
		Wish(t, ParsePath("0").segments, ShouldEqual, []PathSegment{{s: "0", i: -1}})
	})
	t.Run("parsing three segments", func(t *testing.T) {
		Wish(t, ParsePath("0/foo/2").segments, ShouldEqual, []PathSegment{{s: "0", i: -1}, {s: "foo", i: -1}, {s: "2", i: -1}})
	})
	t.Run("eliding leading slashes", func(t *testing.T) {
		Wish(t, ParsePath("/0/2").segments, ShouldEqual, []PathSegment{{s: "0", i: -1}, {s: "2", i: -1}})
	})
	t.Run("eliding trailing", func(t *testing.T) {
		Wish(t, ParsePath("0/2/").segments, ShouldEqual, []PathSegment{{s: "0", i: -1}, {s: "2", i: -1}})
	})
	t.Run("eliding empty segments", func(t *testing.T) { // NOTE: a spec for string encoding might cause this to change in the future!
		Wish(t, ParsePath("0//2").segments, ShouldEqual, []PathSegment{{s: "0", i: -1}, {s: "2", i: -1}})
	})
	t.Run("escaping segments", func(t *testing.T) { // NOTE: a spec for string encoding might cause this to change in the future!
		Wish(t, ParsePath(`0/\//2`).segments, ShouldEqual, []PathSegment{{s: "0", i: -1}, {s: `\`, i: -1}, {s: "2", i: -1}})
	})
}

func TestPathSegmentZeroValue(t *testing.T) {
	Wish(t, PathSegment{}.String(), ShouldEqual, "0")
	i, err := PathSegment{}.Index()
	Wish(t, err, ShouldEqual, nil)
	Wish(t, i, ShouldEqual, int64(0))
}
