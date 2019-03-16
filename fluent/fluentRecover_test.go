package fluent

import (
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/impl/free"
)

func TestRecover(t *testing.T) {
	t.Run("simple traversal error should capture", func(t *testing.T) {
		Wish(t,
			Recover(func() {
				WrapNode(&ipldfree.Node{}).TraverseIndex(0).AsString()
				t.Fatal("should not be reached")
			}),
			ShouldEqual,
			Error{fmt.Errorf("cannot traverse a node that is undefined")},
		)
	})
	t.Run("correct traversal should return nil", func(t *testing.T) {
		Wish(t,
			Recover(func() {
				n, _ := ipldfree.NodeBuilder().CreateString("foo")
				WrapNode(n).AsString()
			}),
			ShouldEqual,
			nil,
		)
	})
	t.Run("other panics should continue to rise", func(t *testing.T) {
		Wish(t,
			func() (r interface{}) {
				defer func() { r = recover() }()
				Recover(func() {
					panic("fuqawds")
				})
				return
			}(),
			ShouldEqual,
			"fuqawds",
		)
	})
}
