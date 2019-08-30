package ipld

import (
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestParsePath(t *testing.T) {
	t.Run("parsing one segment", func(t *testing.T) {
		Wish(t, ParsePath("0").segments, ShouldEqual, []PathSegment{{s: "0"}})
	})
	t.Run("parsing three segments", func(t *testing.T) {
		Wish(t, ParsePath("0/foo/2").segments, ShouldEqual, []PathSegment{{s: "0"}, {s: "foo"}, {s: "2"}})
	})
	t.Run("eliding empty segments", func(t *testing.T) {
		Wish(t, ParsePath("0//2").segments, ShouldEqual, []PathSegment{{s: "0"}, {s: "2"}})
	})
	t.Run("eliding leading slashes", func(t *testing.T) {
		Wish(t, ParsePath("/0/2").segments, ShouldEqual, []PathSegment{{s: "0"}, {s: "2"}})
	})
	t.Run("eliding trailing", func(t *testing.T) {
		Wish(t, ParsePath("0/2/").segments, ShouldEqual, []PathSegment{{s: "0"}, {s: "2"}})
	})
}
