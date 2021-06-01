package tests

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/schema"
)

func SchemaTestString(t *testing.T, engine Engine) {
	ts := schema.TypeSystem{}
	ts.Init()

	ts.Accumulate(schema.SpawnString("String"))
	engine.Init(t, ts)

	np := engine.PrototypeByName("String")
	t.Run("create string", func(t *testing.T) {
		nb := np.NewBuilder()
		Wish(t, nb.AssignString("woiu"), ShouldEqual, nil)
		n := nb.Build().(schema.TypedNode)
		t.Run("read string", func(t *testing.T) {
			Wish(t, n.Kind(), ShouldEqual, ipld.Kind_String)
		})
		t.Run("read representation", func(t *testing.T) {
			nr := n.Representation()
			Wish(t, nr.Kind(), ShouldEqual, ipld.Kind_String)
			Wish(t, str(nr), ShouldEqual, "woiu")
		})
	})
	t.Run("create null is rejected", func(t *testing.T) {
		nb := np.NewBuilder()
		Wish(t, nb.AssignNull(), ShouldBeSameTypeAs, ipld.ErrWrongKind{})
	})
}
