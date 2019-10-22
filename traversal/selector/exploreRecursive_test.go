package selector

import (
	"bytes"
	"fmt"
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	. "github.com/warpfork/go-wish"
)

func TestParseExploreRecursive(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := fnb.CreateInt(0)
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector body must be a map"))
	})
	t.Run("parsing map node without sequence field should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_LimitDepth), vnb.CreateInt(2))
			}))
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: sequence field must be present in ExploreRecursive selector"))
	})
	t.Run("parsing map node without limit field should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: limit field must be present in ExploreRecursive selector"))
	})
	t.Run("parsing map node with limit field that is not a map should fail", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateString("cheese"))
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: limit in ExploreRecursive is a keyed union and thus must be a map"))
	})
	t.Run("parsing map node with limit field that is not a single entry map should fail", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_LimitDepth), vnb.CreateInt(2))
				mb.Insert(knb.CreateString(SelectorKey_LimitNone), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: limit in ExploreRecursive is a keyed union and thus must be a single-entry map"))
	})
	t.Run("parsing map node with limit field that does not have a known key should fail", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString("applesauce"), vnb.CreateInt(2))
			}))
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: \"applesauce\" is not a known member of the limit union in ExploreRecursive"))
	})
	t.Run("parsing map node with limit field of type depth that is not an int should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_LimitDepth), vnb.CreateString("cheese"))
			}))
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: limit field of type depth must be a number in ExploreRecursive selector"))
	})
	t.Run("parsing map node with sequence field with invalid selector node should return child's error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_LimitDepth), vnb.CreateInt(2))
			}))
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateInt(0))
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map"))
	})
	t.Run("parsing map node with sequence field with valid selector w/o ExploreRecursiveEdge should not parse", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_LimitDepth), vnb.CreateInt(2))
			}))
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_ExploreAll), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
					}))
				}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: ExploreRecursive must have at least one ExploreRecursiveEdge"))
	})
	t.Run("parsing map node that is ExploreRecursiveEdge without ExploreRecursive parent should not parse", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {})
		_, err := ParseContext{}.ParseExploreRecursiveEdge(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: ExploreRecursiveEdge must be beneath ExploreRecursive"))
	})
	t.Run("parsing map node with sequence field with valid selector node should parse", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_LimitDepth), vnb.CreateInt(2))
			}))
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_ExploreAll), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(SelectorKey_ExploreRecursiveEdge), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
					}))
				}))
			}))
		})
		s, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreRecursive{ExploreAll{ExploreRecursiveEdge{}}, ExploreAll{ExploreRecursiveEdge{}}, RecursionLimit{RecursionLimit_Depth, 2}})
	})

	t.Run("parsing map node with sequence field with valid selector node and limit type none should parse", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Limit), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_LimitNone), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
			mb.Insert(knb.CreateString(SelectorKey_Sequence), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_ExploreAll), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
						mb.Insert(knb.CreateString(SelectorKey_ExploreRecursiveEdge), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
					}))
				}))
			}))
		})
		s, err := ParseContext{}.ParseExploreRecursive(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreRecursive{ExploreAll{ExploreRecursiveEdge{}}, ExploreAll{ExploreRecursiveEdge{}}, RecursionLimit{RecursionLimit_None, 0}})
	})

}

/*

{
	exploreRecursive: {
		maxDepth: 3
		sequence: {
			exploreFields: {
				fields: {
					Parents: {
						exploreAll: {
							exploreRecursiveEdge: {}
						}
					}
				}
			}
		}
	}
 }

*/

func TestExploreRecursiveExplore(t *testing.T) {
	recursiveEdge := ExploreRecursiveEdge{}
	maxDepth := 3
	var rs Selector
	t.Run("exploring should traverse until we get to maxDepth", func(t *testing.T) {
		parentsSelector := ExploreAll{recursiveEdge}
		subTree := ExploreFields{map[string]Selector{"Parents": parentsSelector}, []ipld.PathSegment{ipld.PathSegmentOfString("Parents")}}
		rs = ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth}}
		nodeString := `{
			"Parents": [
				{
					"Parents": [
						{
							"Parents": [
								{
									"Parents": []
								}
							]
						}
					]
				}
			]
		}
		`
		rn, err := dagjson.Decoder(ipldfree.NodeBuilder(), bytes.NewBufferString(nodeString))
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))

		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, nil)
		Wish(t, err, ShouldEqual, nil)
	})

	t.Run("exploring should traverse indefinitely if no depth specified", func(t *testing.T) {
		parentsSelector := ExploreAll{recursiveEdge}
		subTree := ExploreFields{map[string]Selector{"Parents": parentsSelector}, []ipld.PathSegment{ipld.PathSegmentOfString("Parents")}}
		rs = ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_None, 0}}
		nodeString := `{
			"Parents": [
				{
					"Parents": [
						{
							"Parents": [
								{
									"Parents": []
								}
							]
						}
					]
				}
			]
		}
		`
		rn, err := dagjson.Decoder(ipldfree.NodeBuilder(), bytes.NewBufferString(nodeString))
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_None, 0}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_None, 0}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_None, 0}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_None, 0}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_None, 0}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_None, 0}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_None, 0}})
		Wish(t, err, ShouldEqual, nil)
	})

	t.Run("exploring should continue till we get to selector that returns nil on explore", func(t *testing.T) {
		parentsSelector := ExploreIndex{recursiveEdge, [1]ipld.PathSegment{ipld.PathSegmentOfInt(1)}}
		subTree := ExploreFields{map[string]Selector{"Parents": parentsSelector}, []ipld.PathSegment{ipld.PathSegmentOfString("Parents")}}
		rs = ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth}}
		nodeString := `{
			"Parents": {
			}
		}
		`
		rn, err := dagjson.Decoder(ipldfree.NodeBuilder(), bytes.NewBufferString(nodeString))
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		Wish(t, rs, ShouldEqual, nil)
	})
	t.Run("exploring should work when there is nested recursion", func(t *testing.T) {
		parentsSelector := ExploreAll{recursiveEdge}
		sideSelector := ExploreAll{recursiveEdge}
		subTree := ExploreFields{map[string]Selector{
			"Parents": parentsSelector,
			"Side":    ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}},
		}, []ipld.PathSegment{
			ipld.PathSegmentOfString("Parents"),
			ipld.PathSegmentOfString("Side"),
		},
		}
		s := ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth}}
		nodeString := `{
			"Parents": [
				{
					"Parents": [],
					"Side": {
						"cheese": {
							"whiz": {
							}
						}
					}
				}
			],
			"Side": {
				"real": {
					"apple": {
						"sauce": {
						}
					}
				}
			}
		}
		`
		n, err := dagjson.Decoder(ipldfree.NodeBuilder(), bytes.NewBufferString(nodeString))
		Wish(t, err, ShouldEqual, nil)

		// traverse down Parent nodes
		rn := n
		rs = s
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)

		// traverse down top level Side tree (nested recursion)
		rn = n
		rs = s
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Side"))
		rn, err = rn.LookupString("Side")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}}, RecursionLimit{RecursionLimit_Depth, maxDepth}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("real"))
		rn, err = rn.LookupString("real")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}}, RecursionLimit{RecursionLimit_Depth, maxDepth}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("apple"))
		rn, err = rn.LookupString("apple")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}}, RecursionLimit{RecursionLimit_Depth, maxDepth}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("sauce"))
		rn, err = rn.LookupString("sauce")
		Wish(t, rs, ShouldEqual, nil)
		Wish(t, err, ShouldEqual, nil)

		// traverse once down Parent (top level recursion) then down Side tree (nested recursion)
		rn = n
		rs = s
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Side"))
		rn, err = rn.LookupString("Side")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}}, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("cheese"))
		rn, err = rn.LookupString("cheese")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}}, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("whiz"))
		rn, err = rn.LookupString("whiz")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}}, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
	})
	t.Run("exploring should work with explore union and recursion", func(t *testing.T) {
		parentsSelector := ExploreUnion{[]Selector{ExploreAll{Matcher{}}, ExploreIndex{recursiveEdge, [1]ipld.PathSegment{ipld.PathSegmentOfInt(0)}}}}
		subTree := ExploreFields{map[string]Selector{"Parents": parentsSelector}, []ipld.PathSegment{ipld.PathSegmentOfString("Parents")}}
		rs = ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth}}
		nodeString := `{
			"Parents": [
				{
					"Parents": [
						{
							"Parents": [
								{
									"Parents": []
								}
							]
						}
					]
				}
			]
		}
		`
		rn, err := dagjson.Decoder(ipldfree.NodeBuilder(), bytes.NewBufferString(nodeString))
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))
		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, ExploreUnion{[]Selector{Matcher{}, subTree}}, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfString("Parents"))

		rn, err = rn.LookupString("Parents")
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}})
		Wish(t, err, ShouldEqual, nil)
		rs = rs.Explore(rn, ipld.PathSegmentOfInt(0))
		rn, err = rn.LookupIndex(0)
		Wish(t, rs, ShouldEqual, ExploreRecursive{subTree, ExploreUnion{[]Selector{Matcher{}, subTree}}, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}})
		Wish(t, err, ShouldEqual, nil)
	})
}
