package gengo

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func TestString(t *testing.T) {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &AdjunctCfg{
		maybeUsesPtr: map[schema.TypeName]bool{},
	}

	ts.Accumulate(schema.SpawnString("String"))

	prefix := "justString"
	pkgName := "main" // has to be 'main' for plugins to work.  this stricture makes little sense to me, but i didn't write the rules.
	genAndCompilerAndTest(t, prefix, pkgName, ts, adjCfg, func(t *testing.T, getStyleByName func(string) ipld.NodeStyle) {
		ns := getStyleByName("String")
		t.Run("create string", func(t *testing.T) {
			nb := ns.NewBuilder()
			Wish(t, nb.AssignString("woiu"), ShouldEqual, nil)
			n := nb.Build().(schema.TypedNode)
			t.Run("read string", func(t *testing.T) {
				Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_String)
			})
			t.Run("read representation", func(t *testing.T) {
				nr := n.Representation()
				Wish(t, nr.ReprKind(), ShouldEqual, ipld.ReprKind_String)
				Wish(t, str(nr), ShouldEqual, "woiu")
			})
		})
		t.Run("create null is rejected", func(t *testing.T) {
			nb := ns.NewBuilder()
			Wish(t, nb.AssignNull(), ShouldBeSameTypeAs, ipld.ErrWrongKind{})
		})
	})
}
