package traversal_test

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/traversal"
)

func TestSelectLinks(t *testing.T) {
	t.Run("Scalar", func(t *testing.T) {
		lnks, _ := traversal.SelectLinks(leafAlpha)
		Wish(t, lnks, ShouldEqual, []ipld.Link(nil))
	})
	t.Run("DeepMap", func(t *testing.T) {
		lnks, _ := traversal.SelectLinks(middleMapNode)
		Wish(t, lnks, ShouldEqual, []ipld.Link{leafAlphaLnk})
	})
	t.Run("List", func(t *testing.T) {
		lnks, _ := traversal.SelectLinks(rootNode)
		Wish(t, lnks, ShouldEqual, []ipld.Link{leafAlphaLnk, middleMapNodeLnk, middleListNodeLnk})
	})
}
