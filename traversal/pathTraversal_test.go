package traversal

import (
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"

	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
)

func TestPathTraversal(t *testing.T) {
	t.Run("traversing list", func(t *testing.T) {
		n := &ipldfree.Node{}
		n0 := &ipldfree.Node{}
		n0.SetString("asdf")
		n.SetIndex(0, n0)

		tp, nn, e := ParsePath("0").traverse(TraversalProgress{}, n)

		Wish(t, nn, ShouldEqual, n0)
		Wish(t, tp.Path, ShouldEqual, ParsePath("0"))
		Wish(t, e, ShouldEqual, nil)
	})
	t.Run("traversing map", func(t *testing.T) {
		n := &ipldfree.Node{}
		n0 := &ipldfree.Node{}
		n0.SetString("asdf")
		n.SetField("foo", n0)

		tp, nn, e := ParsePath("foo").traverse(TraversalProgress{}, n)

		Wish(t, nn, ShouldEqual, n0)
		Wish(t, tp.Path, ShouldEqual, ParsePath("foo"))
		Wish(t, e, ShouldEqual, nil)
	})
	t.Run("traversing deeper", func(t *testing.T) {
		n := &ipldfree.Node{}
		n0 := &ipldfree.Node{}
		n01 := &ipldfree.Node{}
		n010 := &ipldfree.Node{}
		n010.SetString("asdf")
		n01.SetField("bar", n010)
		n0.SetIndex(1, n01)
		n.SetField("foo", n0)

		tp, nn, e := ParsePath("foo/1/bar").traverse(TraversalProgress{}, n)

		Wish(t, nn, ShouldEqual, n010)
		Wish(t, tp.Path, ShouldEqual, ParsePath("foo/1/bar"))
		Wish(t, e, ShouldEqual, nil)
	})
	t.Run("traversal error on unexpected terminals", func(t *testing.T) {
		n := &ipldfree.Node{}
		n0 := &ipldfree.Node{}
		n01 := &ipldfree.Node{}
		n010 := &ipldfree.Node{}
		n010.SetString("asdf")
		n01.SetField("bar", n010)
		n0.SetIndex(1, n01)
		n.SetField("foo", n0)

		t.Run("deep terminal", func(t *testing.T) {
			tp, nn, e := ParsePath("foo/1/bar/drat").traverse(TraversalProgress{}, n)

			Wish(t, nn, ShouldEqual, nil)
			Wish(t, tp.Path, ShouldEqual, Path{})
			Wish(t, e, ShouldEqual, fmt.Errorf(`error traversing node at "foo/1/bar": cannot traverse terminals`))
		})
		t.Run("immediate terminal", func(t *testing.T) {
			tp, nn, e := ParsePath("drat").traverse(TraversalProgress{}, n010)

			Wish(t, nn, ShouldEqual, nil)
			Wish(t, tp.Path, ShouldEqual, Path{})
			Wish(t, e, ShouldEqual, fmt.Errorf(`error traversing node at "": cannot traverse terminals`))
		})
	})
	t.Run("traversal error and partial progress on missing members", func(t *testing.T) {
		n := &ipldfree.Node{}
		n0 := &ipldfree.Node{}
		n01 := &ipldfree.Node{}
		n010 := &ipldfree.Node{}
		n010.SetString("asdf")
		n01.SetField("bar", n010)
		n0.SetIndex(1, n01)
		n.SetField("foo", n0)

		tp, nn, e := ParsePath("foo/1/drat").traverse(TraversalProgress{}, n)

		Wish(t, nn, ShouldEqual, nil)
		Wish(t, tp.Path, ShouldEqual, Path{})
		Wish(t, e, ShouldEqual, fmt.Errorf(`error traversing node at "foo/1": 404`))
	})
}
