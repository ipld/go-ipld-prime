package selector

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestParseExploreRecursive(t *testing.T) {
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := basicnode.NewInt(0)
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: selector body must be a map")
	})
	t.Run("parsing map node without sequence field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_LimitDepth).AssignInt(2)
			})
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: sequence field must be present in ExploreRecursive selector")
	})
	t.Run("parsing map node without limit field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: limit field must be present in ExploreRecursive selector")
	})
	t.Run("parsing map node with limit field that is not a map should fail", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).AssignString("cheese")
			na.AssembleEntry(SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: limit in ExploreRecursive is a keyed union and thus must be a map")
	})
	t.Run("parsing map node with limit field that is not a single entry map should fail", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).CreateMap(2, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_LimitDepth).AssignInt(2)
				na.AssembleEntry(SelectorKey_LimitNone).CreateMap(0, func(na fluent.MapAssembler) {})
			})
			na.AssembleEntry(SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: limit in ExploreRecursive is a keyed union and thus must be a single-entry map")
	})
	t.Run("parsing map node with limit field that does not have a known key should fail", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry("applesauce").AssignInt(2)
			})
			na.AssembleEntry(SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: \"applesauce\" is not a known member of the limit union in ExploreRecursive")
	})
	t.Run("parsing map node with limit field of type depth that is not an int should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_LimitDepth).AssignString("cheese")
			})
			na.AssembleEntry(SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: limit field of type depth must be a number in ExploreRecursive selector")
	})
	t.Run("parsing map node with sequence field with invalid selector node should return child's error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_LimitDepth).AssignInt(2)
			})
			na.AssembleEntry(SelectorKey_Sequence).AssignInt(0)
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: selector is a keyed union and thus must be a map")
	})
	t.Run("parsing map node with sequence field with valid selector w/o ExploreRecursiveEdge should not parse", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_LimitDepth).AssignInt(2)
			})
			na.AssembleEntry(SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
					})
				})
			})
		})
		_, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: ExploreRecursive must have at least one ExploreRecursiveEdge")
	})
	t.Run("parsing map node that is ExploreRecursiveEdge without ExploreRecursive parent should not parse", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 0, func(na fluent.MapAssembler) {})
		_, err := ParseContext{}.ParseExploreRecursiveEdge(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: ExploreRecursiveEdge must be beneath ExploreRecursive")
	})
	t.Run("parsing map node with sequence field with valid selector node should parse", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_LimitDepth).AssignInt(2)
			})
			na.AssembleEntry(SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleEntry(SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapAssembler) {})
					})
				})
			})
		})
		s, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, s, qt.Equals, ExploreRecursive{ExploreAll{ExploreRecursiveEdge{}}, ExploreAll{ExploreRecursiveEdge{}}, RecursionLimit{RecursionLimit_Depth, 2}, nil})
	})

	t.Run("parsing map node with sequence field with valid selector node and limit type none should parse", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Limit).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_LimitNone).CreateMap(0, func(na fluent.MapAssembler) {})
			})
			na.AssembleEntry(SelectorKey_Sequence).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_ExploreAll).CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleEntry(SelectorKey_ExploreRecursiveEdge).CreateMap(0, func(na fluent.MapAssembler) {})
					})
				})
			})
		})
		s, err := ParseContext{}.ParseExploreRecursive(sn)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, s, qt.Equals, ExploreRecursive{ExploreAll{ExploreRecursiveEdge{}}, ExploreAll{ExploreRecursiveEdge{}}, RecursionLimit{RecursionLimit_None, 0}, nil})
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
	maxDepth := int64(3)
	var rs Selector
	t.Run("exploring should traverse until we get to maxDepth", func(t *testing.T) {
		parentsSelector := ExploreAll{recursiveEdge}
		subTree := ExploreFields{map[string]Selector{"Parents": parentsSelector}, []datamodel.PathSegment{datamodel.PathSegmentOfString("Parents")}}
		rs = ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil}
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
		nb := basicnode.Prototype__Any{}.NewBuilder()
		err := dagjson.Decode(nb, strings.NewReader(nodeString))
		qt.Check(t, err, qt.IsNil)
		rn := nb.Build()
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		rn, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))

		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		rn, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		_, err = rn.LookupByIndex(0)
		qt.Check(t, rs, qt.IsNil)
		qt.Check(t, err, qt.IsNil)
	})

	t.Run("exploring should traverse indefinitely if no depth specified", func(t *testing.T) {
		parentsSelector := ExploreAll{recursiveEdge}
		subTree := ExploreFields{map[string]Selector{"Parents": parentsSelector}, []datamodel.PathSegment{datamodel.PathSegmentOfString("Parents")}}
		rs = ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_None, 0}, nil}
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
		nb := basicnode.Prototype__Any{}.NewBuilder()
		err := dagjson.Decode(nb, strings.NewReader(nodeString))
		qt.Check(t, err, qt.IsNil)
		rn := nb.Build()
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_None, 0}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		rn, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_None, 0}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_None, 0}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		rn, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_None, 0}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_None, 0}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		rn, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_None, 0}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		_, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_None, 0}, nil})
		qt.Check(t, err, qt.IsNil)
	})

	t.Run("exploring should continue till we get to selector that returns nil on explore", func(t *testing.T) {
		parentsSelector := ExploreIndex{recursiveEdge, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(1)}}
		subTree := ExploreFields{map[string]Selector{"Parents": parentsSelector}, []datamodel.PathSegment{datamodel.PathSegmentOfString("Parents")}}
		rs = ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil}
		nodeString := `{
			"Parents": {
			}
		}
		`
		nb := basicnode.Prototype__Any{}.NewBuilder()
		err := dagjson.Decode(nb, strings.NewReader(nodeString))
		qt.Check(t, err, qt.IsNil)
		rn := nb.Build()
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		qt.Check(t, rs, qt.IsNil)
	})
	t.Run("exploring should work when there is nested recursion", func(t *testing.T) {
		parentsSelector := ExploreAll{recursiveEdge}
		sideSelector := ExploreAll{recursiveEdge}
		subTree := ExploreFields{map[string]Selector{
			"Parents": parentsSelector,
			"Side":    ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil},
		}, []datamodel.PathSegment{
			datamodel.PathSegmentOfString("Parents"),
			datamodel.PathSegmentOfString("Side"),
		},
		}
		s := ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil}
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
		nb := basicnode.Prototype__Any{}.NewBuilder()
		err := dagjson.Decode(nb, strings.NewReader(nodeString))
		qt.Check(t, err, qt.IsNil)
		n := nb.Build()

		// traverse down Parent nodes
		rn := n
		rs = s
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		rn, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		_, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)

		// traverse down top level Side tree (nested recursion)
		rn = n
		rs = s
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Side"))
		rn, err = rn.LookupByString("Side")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil}, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("real"))
		rn, err = rn.LookupByString("real")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil}, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("apple"))
		rn, err = rn.LookupByString("apple")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}, nil}, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("sauce"))
		_, err = rn.LookupByString("sauce")
		qt.Check(t, rs, qt.IsNil)
		qt.Check(t, err, qt.IsNil)

		// traverse once down Parent (top level recursion) then down Side tree (nested recursion)
		rn = n
		rs = s
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		rn, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Side"))
		rn, err = rn.LookupByString("Side")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil}, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("cheese"))
		rn, err = rn.LookupByString("cheese")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil}, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("whiz"))
		_, err = rn.LookupByString("whiz")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, ExploreRecursive{sideSelector, sideSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}, nil}, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
	})
	t.Run("exploring should work with explore union and recursion", func(t *testing.T) {
		parentsSelector := ExploreUnion{[]Selector{ExploreAll{Matcher{}}, ExploreIndex{recursiveEdge, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(0)}}}}
		subTree := ExploreFields{map[string]Selector{"Parents": parentsSelector}, []datamodel.PathSegment{datamodel.PathSegmentOfString("Parents")}}
		rs = ExploreRecursive{subTree, subTree, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil}
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
		nb := basicnode.Prototype__Any{}.NewBuilder()
		err := dagjson.Decode(nb, strings.NewReader(nodeString))
		qt.Check(t, err, qt.IsNil)
		rn := nb.Build()
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))
		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		rn, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, ExploreUnion{[]Selector{Matcher{}, subTree}}, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfString("Parents"))

		rn, err = rn.LookupByString("Parents")
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, parentsSelector, RecursionLimit{RecursionLimit_Depth, maxDepth - 1}, nil})
		qt.Check(t, err, qt.IsNil)
		rs, _ = rs.Explore(rn, datamodel.PathSegmentOfInt(0))
		_, err = rn.LookupByIndex(0)
		qt.Check(t, rs, deepEqualsAllowAllUnexported, ExploreRecursive{subTree, ExploreUnion{[]Selector{Matcher{}, subTree}}, RecursionLimit{RecursionLimit_Depth, maxDepth - 2}, nil})
		qt.Check(t, err, qt.IsNil)
	})
}
