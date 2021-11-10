package datamodel

import (
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp"
)

func TestParsePath(t *testing.T) {
	// Allow equality checker to check all unexported fields in PathSegment.
	pathSegmentEquals := qt.CmpEquals(cmp.Exporter(func(reflect.Type) bool { return true }))
	t.Run("parsing one segment", func(t *testing.T) {
		qt.Check(t, ParsePath("0").segments, pathSegmentEquals, []PathSegment{{s: "0", i: -1}})
	})
	t.Run("parsing three segments", func(t *testing.T) {
		qt.Check(t, ParsePath("0/foo/2").segments, pathSegmentEquals, []PathSegment{{s: "0", i: -1}, {s: "foo", i: -1}, {s: "2", i: -1}})
	})
	t.Run("eliding leading slashes", func(t *testing.T) {
		qt.Check(t, ParsePath("/0/2").segments, pathSegmentEquals, []PathSegment{{s: "0", i: -1}, {s: "2", i: -1}})
	})
	t.Run("eliding trailing", func(t *testing.T) {
		qt.Check(t, ParsePath("0/2/").segments, pathSegmentEquals, []PathSegment{{s: "0", i: -1}, {s: "2", i: -1}})
	})
	t.Run("eliding empty segments", func(t *testing.T) { // NOTE: a spec for string encoding might cause this to change in the future!
		qt.Check(t, ParsePath("0//2").segments, pathSegmentEquals, []PathSegment{{s: "0", i: -1}, {s: "2", i: -1}})
	})
	t.Run("escaping segments", func(t *testing.T) { // NOTE: a spec for string encoding might cause this to change in the future!
		qt.Check(t, ParsePath(`0/\//2`).segments, pathSegmentEquals, []PathSegment{{s: "0", i: -1}, {s: `\`, i: -1}, {s: "2", i: -1}})
	})
}

func TestPathSegmentZeroValue(t *testing.T) {
	qt.Check(t, PathSegment{}.String(), qt.Equals, "0")
	i, err := PathSegment{}.Index()
	qt.Check(t, err, qt.IsNil)
	qt.Check(t, i, qt.Equals, int64(0))
}
