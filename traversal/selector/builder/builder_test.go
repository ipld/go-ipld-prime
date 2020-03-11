package builder

import (
	"testing"

	"github.com/ipld/go-ipld-prime/fluent"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	. "github.com/warpfork/go-wish"
)

func TestBuildingSelectors(t *testing.T) {
	ns := basicnode.Style__Any{}
	ssb := NewSelectorSpecBuilder(ns)
	t.Run("Matcher builds matcher nodes", func(t *testing.T) {
		sn := ssb.Matcher().Node()
		esn := fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRecursiveEdge builds ExploreRecursiveEdge nodes", func(t *testing.T) {
		sn := ssb.ExploreRecursiveEdge().Node()
		esn := fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapAssembler) {})
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreAll builds ExploreAll nodes", func(t *testing.T) {
		sn := ssb.ExploreAll(ssb.Matcher()).Node()
		esn := fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
			})
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreIndex builds ExploreIndex nodes", func(t *testing.T) {
		sn := ssb.ExploreIndex(2, ssb.Matcher()).Node()
		esn := fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreIndex).CreateMap(2, func(na fluent.MapAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Index).AssignInt(2)
				na.AssembleDirectly(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
			})
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRange builds ExploreRange nodes", func(t *testing.T) {
		sn := ssb.ExploreRange(2, 3, ssb.Matcher()).Node()
		esn := fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreRange).CreateMap(3, func(na fluent.MapAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Start).AssignInt(2)
				na.AssembleDirectly(selector.SelectorKey_End).AssignInt(3)
				na.AssembleDirectly(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
			})
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRecursive builds ExploreRecursive nodes", func(t *testing.T) {
		sn := ssb.ExploreRecursive(selector.RecursionLimitDepth(2), ssb.ExploreAll(ssb.ExploreRecursiveEdge())).Node()
		esn := fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreRecursive).CreateMap(2, func(na fluent.MapAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_LimitDepth).AssignInt(2)
				})
				na.AssembleDirectly(selector.SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleDirectly(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
							na.AssembleDirectly(selector.SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapAssembler) {})
						})
					})
				})
			})
		})
		Wish(t, sn, ShouldEqual, esn)
		sn = ssb.ExploreRecursive(selector.RecursionLimitNone(), ssb.ExploreAll(ssb.ExploreRecursiveEdge())).Node()
		esn = fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreRecursive).CreateMap(2, func(na fluent.MapAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_LimitNone).CreateMap(0, func(na fluent.MapAssembler) {})
				})
				na.AssembleDirectly(selector.SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleDirectly(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
							na.AssembleDirectly(selector.SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapAssembler) {})
						})
					})
				})
			})
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreUnion builds ExploreUnion nodes", func(t *testing.T) {
		sn := ssb.ExploreUnion(ssb.Matcher(), ssb.ExploreIndex(2, ssb.Matcher())).Node()
		esn := fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreUnion).CreateList(2, func(na fluent.ListAssembler) {
				na.AssembleValue().CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
				na.AssembleValue().CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(selector.SelectorKey_ExploreIndex).CreateMap(2, func(na fluent.MapAssembler) {
						na.AssembleDirectly(selector.SelectorKey_Index).AssignInt(2)
						na.AssembleDirectly(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
							na.AssembleDirectly(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
						})
					})
				})
			})
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreFields builds ExploreFields nodes", func(t *testing.T) {
		sn := ssb.ExploreFields(func(efsb ExploreFieldsSpecBuilder) { efsb.Insert("applesauce", ssb.Matcher()) }).Node()
		esn := fluent.MustBuildMap(ns, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(selector.SelectorKey_ExploreFields).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleDirectly(selector.SelectorKey_Fields).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly("applesauce").CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleDirectly(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
					})
				})
			})
		})
		Wish(t, sn, ShouldEqual, esn)
	})
}
