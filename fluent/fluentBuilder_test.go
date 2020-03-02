package fluent_test

import (
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestBuild(t *testing.T) {
	t.Run("scalar build should work", func(t *testing.T) {
		n := fluent.MustBuild(basicnode.Style__String{}, func(fna fluent.NodeAssembler) {
			fna.AssignString("fine")
		})
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_String)
		v2, err := n.AsString()
		Wish(t, err, ShouldEqual, nil)
		Wish(t, v2, ShouldEqual, "fine")
	})
	t.Run("map build should work", func(t *testing.T) {
		n := fluent.MustBuild(basicnode.Style__Map{}, func(fna fluent.NodeAssembler) {
			fna.CreateMap(3, func(fma fluent.MapNodeAssembler) {
				fma.AssembleDirectly("k1").AssignString("fine")
				fma.AssembleDirectly("k2").AssignString("super")
				fma.AssembleDirectly("k3").CreateMap(3, func(fma fluent.MapNodeAssembler) {
					fma.AssembleDirectly("k31").AssignString("thanks")
					fma.AssembleDirectly("k32").AssignString("for")
					fma.AssembleDirectly("k33").AssignString("asking")
				})
			})
		})
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Map)
		Wish(t, n.Length(), ShouldEqual, 3)
		Wish(t, must.String(must.Node(n.LookupString("k1"))), ShouldEqual, "fine")
		Wish(t, must.String(must.Node(n.LookupString("k2"))), ShouldEqual, "super")
		n = must.Node(n.LookupString("k3"))
		Wish(t, n.Length(), ShouldEqual, 3)
		Wish(t, must.String(must.Node(n.LookupString("k31"))), ShouldEqual, "thanks")
		Wish(t, must.String(must.Node(n.LookupString("k32"))), ShouldEqual, "for")
		Wish(t, must.String(must.Node(n.LookupString("k33"))), ShouldEqual, "asking")
	})
	t.Run("list build should work", func(t *testing.T) {
		n := fluent.MustBuild(basicnode.Style__List{}, func(fna fluent.NodeAssembler) {
			fna.CreateList(1, func(fla fluent.ListNodeAssembler) {
				fla.AssembleValue().CreateList(1, func(fla fluent.ListNodeAssembler) {
					fla.AssembleValue().CreateList(1, func(fla fluent.ListNodeAssembler) {
						fla.AssembleValue().CreateList(1, func(fla fluent.ListNodeAssembler) {
							fla.AssembleValue().AssignInt(2)
						})
					})
				})
			})
		})
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, n.Length(), ShouldEqual, 1)
		n = must.Node(n.LookupIndex(0))
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, n.Length(), ShouldEqual, 1)
		n = must.Node(n.LookupIndex(0))
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, n.Length(), ShouldEqual, 1)
		n = must.Node(n.LookupIndex(0))
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_List)
		Wish(t, n.Length(), ShouldEqual, 1)
		n = must.Node(n.LookupIndex(0))
		Wish(t, n.ReprKind(), ShouldEqual, ipld.ReprKind_Int)
		Wish(t, must.Int(n), ShouldEqual, 2)
	})
}
