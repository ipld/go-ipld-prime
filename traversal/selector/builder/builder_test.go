package builder_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	nodetests "github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	"github.com/ipld/go-ipld-prime/traversal/selector/builder"
)

func TestBuildingSelectors(t *testing.T) {
	np := basicnode.Prototype.Any
	ssb := builder.NewSelectorSpecBuilder(np)
	t.Run("Matcher builds matcher nodes", func(t *testing.T) {
		sn := ssb.Matcher().Node()
		esn := fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
	})
	t.Run("ExploreRecursiveEdge builds ExploreRecursiveEdge nodes", func(t *testing.T) {
		sn := ssb.ExploreRecursiveEdge().Node()
		esn := fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapAssembler) {})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
	})
	t.Run("ExploreAll builds ExploreAll nodes", func(t *testing.T) {
		sn := ssb.ExploreAll(ssb.Matcher()).Node()
		esn := fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
			})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
	})
	t.Run("ExploreIndex builds ExploreIndex nodes", func(t *testing.T) {
		sn := ssb.ExploreIndex(2, ssb.Matcher()).Node()
		esn := fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_ExploreIndex).CreateMap(2, func(na fluent.MapAssembler) {
				na.AssembleEntry(selector.SelectorKey_Index).AssignInt(2)
				na.AssembleEntry(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
			})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
	})
	t.Run("ExploreRange builds ExploreRange nodes", func(t *testing.T) {
		sn := ssb.ExploreRange(2, 3, ssb.Matcher()).Node()
		esn := fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_ExploreRange).CreateMap(3, func(na fluent.MapAssembler) {
				na.AssembleEntry(selector.SelectorKey_Start).AssignInt(2)
				na.AssembleEntry(selector.SelectorKey_End).AssignInt(3)
				na.AssembleEntry(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
			})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
	})
	t.Run("ExploreRecursive builds ExploreRecursive nodes", func(t *testing.T) {
		sn := ssb.ExploreRecursive(selector.RecursionLimitDepth(2), ssb.ExploreAll(ssb.ExploreRecursiveEdge())).Node()
		esn := fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_ExploreRecursive).CreateMap(2, func(na fluent.MapAssembler) {
				na.AssembleEntry(selector.SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_LimitDepth).AssignInt(2)
				})
				na.AssembleEntry(selector.SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleEntry(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
							na.AssembleEntry(selector.SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapAssembler) {})
						})
					})
				})
			})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
		sn = ssb.ExploreRecursive(selector.RecursionLimitNone(), ssb.ExploreAll(ssb.ExploreRecursiveEdge())).Node()
		esn = fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_ExploreRecursive).CreateMap(2, func(na fluent.MapAssembler) {
				na.AssembleEntry(selector.SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_LimitNone).CreateMap(0, func(na fluent.MapAssembler) {})
				})
				na.AssembleEntry(selector.SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleEntry(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
							na.AssembleEntry(selector.SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapAssembler) {})
						})
					})
				})
			})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
	})
	t.Run("ExploreUnion builds ExploreUnion nodes", func(t *testing.T) {
		sn := ssb.ExploreUnion(ssb.Matcher(), ssb.ExploreIndex(2, ssb.Matcher())).Node()
		esn := fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_ExploreUnion).CreateList(2, func(na fluent.ListAssembler) {
				na.AssembleValue().CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
				na.AssembleValue().CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(selector.SelectorKey_ExploreIndex).CreateMap(2, func(na fluent.MapAssembler) {
						na.AssembleEntry(selector.SelectorKey_Index).AssignInt(2)
						na.AssembleEntry(selector.SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
							na.AssembleEntry(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
						})
					})
				})
			})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
	})
	t.Run("ExploreFields builds ExploreFields nodes", func(t *testing.T) {
		sn := ssb.ExploreFields(func(efsb builder.ExploreFieldsSpecBuilder) { efsb.Insert("applesauce", ssb.Matcher()) }).Node()
		esn := fluent.MustBuildMap(np, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(selector.SelectorKey_ExploreFields).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(selector.SelectorKey_Fields).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry("applesauce").CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleEntry(selector.SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
					})
				})
			})
		})
		qt.Check(t, sn, nodetests.NodeContentEquals, esn)
	})
}
