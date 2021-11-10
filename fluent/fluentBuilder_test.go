package fluent_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestBuild(t *testing.T) {
	t.Run("scalar build should work", func(t *testing.T) {
		n := fluent.MustBuild(basicnode.Prototype__String{}, func(fna fluent.NodeAssembler) {
			fna.AssignString("fine")
		})
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_String)
		v2, err := n.AsString()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, v2, qt.Equals, "fine")
	})
	t.Run("map build should work", func(t *testing.T) {
		n := fluent.MustBuild(basicnode.Prototype.Map, func(fna fluent.NodeAssembler) {
			fna.CreateMap(3, func(fma fluent.MapAssembler) {
				fma.AssembleEntry("k1").AssignString("fine")
				fma.AssembleEntry("k2").AssignString("super")
				fma.AssembleEntry("k3").CreateMap(3, func(fma fluent.MapAssembler) {
					fma.AssembleEntry("k31").AssignString("thanks")
					fma.AssembleEntry("k32").AssignString("for")
					fma.AssembleEntry("k33").AssignString("asking")
				})
			})
		})
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Map)
		qt.Check(t, n.Length(), qt.Equals, int64(3))
		qt.Check(t, must.String(must.Node(n.LookupByString("k1"))), qt.Equals, "fine")
		qt.Check(t, must.String(must.Node(n.LookupByString("k2"))), qt.Equals, "super")
		n = must.Node(n.LookupByString("k3"))
		qt.Check(t, n.Length(), qt.Equals, int64(3))
		qt.Check(t, must.String(must.Node(n.LookupByString("k31"))), qt.Equals, "thanks")
		qt.Check(t, must.String(must.Node(n.LookupByString("k32"))), qt.Equals, "for")
		qt.Check(t, must.String(must.Node(n.LookupByString("k33"))), qt.Equals, "asking")
	})
	t.Run("list build should work", func(t *testing.T) {
		n := fluent.MustBuild(basicnode.Prototype.List, func(fna fluent.NodeAssembler) {
			fna.CreateList(1, func(fla fluent.ListAssembler) {
				fla.AssembleValue().CreateList(1, func(fla fluent.ListAssembler) {
					fla.AssembleValue().CreateList(1, func(fla fluent.ListAssembler) {
						fla.AssembleValue().CreateList(1, func(fla fluent.ListAssembler) {
							fla.AssembleValue().AssignInt(2)
						})
					})
				})
			})
		})
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_List)
		qt.Check(t, n.Length(), qt.Equals, int64(1))
		n = must.Node(n.LookupByIndex(0))
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_List)
		qt.Check(t, n.Length(), qt.Equals, int64(1))
		n = must.Node(n.LookupByIndex(0))
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_List)
		qt.Check(t, n.Length(), qt.Equals, int64(1))
		n = must.Node(n.LookupByIndex(0))
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_List)
		qt.Check(t, n.Length(), qt.Equals, int64(1))
		n = must.Node(n.LookupByIndex(0))
		qt.Check(t, n.Kind(), qt.Equals, datamodel.Kind_Int)
		qt.Check(t, must.Int(n), qt.Equals, int64(2))
	})
}
