package selector

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestParseExploreUnion(t *testing.T) {
	t.Run("parsing non list node should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 0, func(na fluent.MapAssembler) {})
		_, err := ParseContext{}.ParseExploreUnion(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: explore union selector must be a list")
	})
	t.Run("parsing list node where one node is invalid should return child's error", func(t *testing.T) {
		sn := fluent.MustBuildList(basicnode.Prototype__List{}, 2, func(na fluent.ListAssembler) {
			na.AssembleValue().CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
			na.AssembleValue().AssignInt(2)
		})
		_, err := ParseContext{}.ParseExploreUnion(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: selector is a keyed union and thus must be a map")
	})

	t.Run("parsing map node with next field with valid selector node should parse", func(t *testing.T) {
		sn := fluent.MustBuildList(basicnode.Prototype__List{}, 2, func(na fluent.ListAssembler) {
			na.AssembleValue().CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
			na.AssembleValue().CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_ExploreIndex).CreateMap(2, func(na fluent.MapAssembler) {
					na.AssembleEntry(SelectorKey_Index).AssignInt(2)
					na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
						na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
					})
				})
			})
		})
		s, err := ParseContext{}.ParseExploreUnion(sn)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, s, deepEqualsAllowAllUnexported, ExploreUnion{[]Selector{Matcher{}, ExploreIndex{Matcher{}, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}}}})
	})
}

func TestExploreUnionExplore(t *testing.T) {
	n := fluent.MustBuildList(basicnode.Prototype__List{}, 4, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignInt(0)
		na.AssembleValue().AssignInt(1)
		na.AssembleValue().AssignInt(2)
		na.AssembleValue().AssignInt(3)
	})
	t.Run("exploring should return nil if all member selectors return nil when explored", func(t *testing.T) {
		s := ExploreUnion{[]Selector{Matcher{}, ExploreIndex{Matcher{}, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}}}}
		returnedSelector, _ := s.Explore(n, datamodel.PathSegmentOfInt(3))
		qt.Check(t, returnedSelector, qt.IsNil)
	})

	t.Run("if exactly one member selector returns a non-nil selector when explored, exploring should return that value", func(t *testing.T) {
		s := ExploreUnion{[]Selector{Matcher{}, ExploreIndex{Matcher{}, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}}}}

		returnedSelector, _ := s.Explore(n, datamodel.PathSegmentOfInt(2))
		qt.Check(t, returnedSelector, qt.Equals, Matcher{})
	})
	t.Run("exploring should return a new union selector if more than one member selector returns a non nil selector when explored", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			Matcher{},
			ExploreIndex{Matcher{}, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}},
			ExploreRange{Matcher{}, 2, 3, []datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}},
			ExploreFields{map[string]Selector{"applesauce": Matcher{}}, []datamodel.PathSegment{datamodel.PathSegmentOfString("applesauce")}},
		}}

		returnedSelector, _ := s.Explore(n, datamodel.PathSegmentOfInt(2))
		qt.Check(t, returnedSelector, deepEqualsAllowAllUnexported, ExploreUnion{[]Selector{Matcher{}, Matcher{}}})
	})
}

func TestExploreUnionInterests(t *testing.T) {
	t.Run("if any member selector is high-cardinality, interests should be high-cardinality", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			ExploreAll{Matcher{}},
			Matcher{},
			ExploreIndex{Matcher{}, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}},
		}}
		qt.Check(t, s.Interests(), deepEqualsAllowAllUnexported, []datamodel.PathSegment(nil))
	})
	t.Run("if no member selector is high-cardinality, interests should be combination of member selectors interests", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			ExploreFields{map[string]Selector{"applesauce": Matcher{}}, []datamodel.PathSegment{datamodel.PathSegmentOfString("applesauce")}},
			Matcher{},
			ExploreIndex{Matcher{}, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}},
		}}
		qt.Check(t, s.Interests(), deepEqualsAllowAllUnexported, []datamodel.PathSegment{datamodel.PathSegmentOfString("applesauce"), datamodel.PathSegmentOfInt(2)})
	})
}

func TestExploreUnionDecide(t *testing.T) {
	n := basicnode.NewInt(2)
	t.Run("if any member selector returns true, decide should be true", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			ExploreAll{Matcher{}},
			Matcher{},
			ExploreIndex{Matcher{}, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}},
		}}
		qt.Check(t, s.Decide(n), qt.IsTrue)
	})
	t.Run("if no member selector returns true, decide should be false", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			ExploreFields{map[string]Selector{"applesauce": Matcher{}}, []datamodel.PathSegment{datamodel.PathSegmentOfString("applesauce")}},
			ExploreAll{Matcher{}},
			ExploreIndex{Matcher{}, [1]datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}},
		}}
		qt.Check(t, s.Decide(n), qt.IsFalse)
	})
}
