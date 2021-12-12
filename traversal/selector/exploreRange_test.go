package selector

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestParseExploreRange(t *testing.T) {
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := basicnode.NewInt(0)
		_, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: selector body must be a map")
	})
	t.Run("parsing map node without next field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Start).AssignInt(2)
			na.AssembleEntry(SelectorKey_End).AssignInt(3)
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: next field must be present in ExploreRange selector")
	})
	t.Run("parsing map node without start field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_End).AssignInt(3)
			na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: start field must be present in ExploreRange selector")
	})
	t.Run("parsing map node with start field that is not an int should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Start).AssignString("cheese")
			na.AssembleEntry(SelectorKey_End).AssignInt(3)
			na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: start field must be a number in ExploreRange selector")
	})
	t.Run("parsing map node without end field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 2, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Start).AssignInt(2)
			na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: end field must be present in ExploreRange selector")
	})
	t.Run("parsing map node with end field that is not an int should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Start).AssignInt(2)
			na.AssembleEntry(SelectorKey_End).AssignString("cheese")
			na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: end field must be a number in ExploreRange selector")
	})
	t.Run("parsing map node where end is not greater than start should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Start).AssignInt(3)
			na.AssembleEntry(SelectorKey_End).AssignInt(2)
			na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: end field must be greater than start field in ExploreRange selector")
	})
	t.Run("parsing map node with next field with invalid selector node should return child's error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Start).AssignInt(2)
			na.AssembleEntry(SelectorKey_End).AssignInt(3)
			na.AssembleEntry(SelectorKey_Next).AssignInt(0)
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: selector is a keyed union and thus must be a map")
	})

	t.Run("parsing map node with next field with valid selector node should parse", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype__Map{}, 3, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Start).AssignInt(2)
			na.AssembleEntry(SelectorKey_End).AssignInt(3)
			na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		s, err := ParseContext{}.ParseExploreRange(sn)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, s, deepEqualsAllowAllUnexported, ExploreRange{Matcher{}, 2, 3, []datamodel.PathSegment{datamodel.PathSegmentOfInt(2)}})
	})
}

func TestExploreRangeExplore(t *testing.T) {
	s := ExploreRange{Matcher{}, 3, 4, []datamodel.PathSegment{datamodel.PathSegmentOfInt(3)}}
	t.Run("exploring should return nil unless node is a list", func(t *testing.T) {
		n := fluent.MustBuildMap(basicnode.Prototype__Map{}, 0, func(na fluent.MapAssembler) {})
		returnedSelector, _ := s.Explore(n, datamodel.PathSegmentOfInt(3))
		qt.Check(t, returnedSelector, qt.IsNil)
	})
	n := fluent.MustBuildList(basicnode.Prototype__List{}, 4, func(na fluent.ListAssembler) {
		na.AssembleValue().AssignInt(0)
		na.AssembleValue().AssignInt(1)
		na.AssembleValue().AssignInt(2)
		na.AssembleValue().AssignInt(3)
	})
	t.Run("exploring should return nil when given a path segment out of range", func(t *testing.T) {
		returnedSelector, _ := s.Explore(n, datamodel.PathSegmentOfInt(2))
		qt.Check(t, returnedSelector, qt.IsNil)
	})
	t.Run("exploring should return nil when given a path segment that isn't an index", func(t *testing.T) {
		returnedSelector, _ := s.Explore(n, datamodel.PathSegmentOfString("cheese"))
		qt.Check(t, returnedSelector, qt.IsNil)
	})
	t.Run("exploring should return the next selector when given a path segment with index in range", func(t *testing.T) {
		returnedSelector, _ := s.Explore(n, datamodel.PathSegmentOfInt(3))
		qt.Check(t, returnedSelector, qt.Equals, Matcher{})
	})
}
