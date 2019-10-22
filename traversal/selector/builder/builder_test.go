package builder

import (
	"testing"

	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	. "github.com/warpfork/go-wish"
)

func TestBuildingSelectors(t *testing.T) {
	nb := ipldfree.NodeBuilder()
	fnb := fluent.WrapNodeBuilder(nb)
	ssb := NewSelectorSpecBuilder(nb)
	t.Run("Matcher builds matcher nodes", func(t *testing.T) {
		sn := ssb.Matcher().Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRecursiveEdge builds ExploreRecursiveEdge nodes", func(t *testing.T) {
		sn := ssb.ExploreRecursiveEdge().Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_ExploreRecursiveEdge), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreAll builds ExploreAll nodes", func(t *testing.T) {
		sn := ssb.ExploreAll(ssb.Matcher()).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_ExploreAll), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(selector.SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreIndex builds ExploreIndex nodes", func(t *testing.T) {
		sn := ssb.ExploreIndex(2, ssb.Matcher()).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_ExploreIndex), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(selector.SelectorKey_Index), vnb.CreateInt(2))
				mb.Insert(knb.CreateString(selector.SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRange builds ExploreRange nodes", func(t *testing.T) {
		sn := ssb.ExploreRange(2, 3, ssb.Matcher()).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_ExploreRange), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(selector.SelectorKey_Start), vnb.CreateInt(2))
				mb.Insert(knb.CreateString(selector.SelectorKey_End), vnb.CreateInt(3))
				mb.Insert(knb.CreateString(selector.SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRecursive builds ExploreRecursive nodes", func(t *testing.T) {
		sn := ssb.ExploreRecursive(selector.RecursionLimitDepth(2), ssb.ExploreAll(ssb.ExploreRecursiveEdge())).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_ExploreRecursive), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(selector.SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_LimitDepth), vnb.CreateInt(2))
				}))
				mb.Insert(knb.CreateString(selector.SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_ExploreAll), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(selector.SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
							mb.Insert(knb.CreateString(selector.SelectorKey_ExploreRecursiveEdge), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
						}))
					}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
		sn = ssb.ExploreRecursive(selector.RecursionLimitNone(), ssb.ExploreAll(ssb.ExploreRecursiveEdge())).Node()
		esn = fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_ExploreRecursive), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(selector.SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_LimitNone), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
				mb.Insert(knb.CreateString(selector.SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_ExploreAll), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(selector.SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
							mb.Insert(knb.CreateString(selector.SelectorKey_ExploreRecursiveEdge), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
						}))
					}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreUnion builds ExploreUnion nodes", func(t *testing.T) {
		sn := ssb.ExploreUnion(ssb.Matcher(), ssb.ExploreIndex(2, ssb.Matcher())).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_ExploreUnion), vnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
				lb.Append(vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
				lb.Append(vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(selector.SelectorKey_ExploreIndex), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(selector.SelectorKey_Index), vnb.CreateInt(2))
						mb.Insert(knb.CreateString(selector.SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
							mb.Insert(knb.CreateString(selector.SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
						}))
					}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreFields builds ExploreFields nodes", func(t *testing.T) {
		sn := ssb.ExploreFields(func(efsb ExploreFieldsSpecBuilder) { efsb.Insert("applesauce", ssb.Matcher()) }).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(selector.SelectorKey_ExploreFields), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(selector.SelectorKey_Fields), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString("applesauce"), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(selector.SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
					}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
}
