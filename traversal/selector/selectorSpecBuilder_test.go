package selector

import (
	"testing"

	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	. "github.com/warpfork/go-wish"
)

func TestBuildingSelectors(t *testing.T) {
	nb := ipldfree.NodeBuilder()
	fnb := fluent.WrapNodeBuilder(nb)
	ssb := NewSelectorSpecBuilder(nb)
	t.Run("Matcher builds matcher nodes", func(t *testing.T) {
		sn := ssb.Matcher().Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRecursiveEdge builds ExploreRecursiveEdge nodes", func(t *testing.T) {
		sn := ssb.ExploreRecursiveEdge().Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreRecursiveEdgeKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreAll builds ExploreAll nodes", func(t *testing.T) {
		sn := ssb.ExploreAll(ssb.Matcher()).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreAllKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreIndex builds ExploreIndex nodes", func(t *testing.T) {
		sn := ssb.ExploreIndex(2, ssb.Matcher()).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreIndexKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(indexKey), vnb.CreateInt(2))
				mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRange builds ExploreRange nodes", func(t *testing.T) {
		sn := ssb.ExploreRange(2, 3, ssb.Matcher()).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreRangeKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(startKey), vnb.CreateInt(2))
				mb.Insert(knb.CreateString(endKey), vnb.CreateInt(3))
				mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
	t.Run("ExploreRecursive builds ExploreRecursive nodes", func(t *testing.T) {
		sn := ssb.ExploreRecursive(2, ssb.ExploreAll(ssb.ExploreRecursiveEdge())).Node()
		esn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(exploreRecursiveKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(maxDepthKey), vnb.CreateInt(2))
				mb.Insert(knb.CreateString(sequenceKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(exploreAllKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
							mb.Insert(knb.CreateString(exploreRecursiveEdgeKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
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
			mb.Insert(knb.CreateString(exploreUnionKey), vnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
				lb.Append(vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
				lb.Append(vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(exploreIndexKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(indexKey), vnb.CreateInt(2))
						mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
							mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
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
			mb.Insert(knb.CreateString(exploreFieldsKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(fieldsKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString("applesauce"), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
					}))
				}))
			}))
		})
		Wish(t, sn, ShouldEqual, esn)
	})
}
