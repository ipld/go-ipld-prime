package selector

import (
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestParseExploreRange(t *testing.T) {
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := basicnode.NewInt(0)
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector body must be a map"))
	})
	t.Run("parsing map node without next field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 2, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(SelectorKey_Start).AssignInt(2)
			na.AssembleDirectly(SelectorKey_End).AssignInt(3)
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: next field must be present in ExploreRange selector"))
	})
	t.Run("parsing map node without start field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 2, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(SelectorKey_End).AssignInt(3)
			na.AssembleDirectly(SelectorKey_Next).CreateMap(1, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapNodeAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: start field must be present in ExploreRange selector"))
	})
	t.Run("parsing map node with start field that is not an int should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 3, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(SelectorKey_Start).AssignString("cheese")
			na.AssembleDirectly(SelectorKey_End).AssignInt(3)
			na.AssembleDirectly(SelectorKey_Next).CreateMap(1, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapNodeAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: start field must be a number in ExploreRange selector"))
	})
	t.Run("parsing map node without end field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 2, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(SelectorKey_Start).AssignInt(2)
			na.AssembleDirectly(SelectorKey_Next).CreateMap(1, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapNodeAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: end field must be present in ExploreRange selector"))
	})
	t.Run("parsing map node with end field that is not an int should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 3, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(SelectorKey_Start).AssignInt(2)
			na.AssembleDirectly(SelectorKey_End).AssignString("cheese")
			na.AssembleDirectly(SelectorKey_Next).CreateMap(1, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapNodeAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: end field must be a number in ExploreRange selector"))
	})
	t.Run("parsing map node where end is not greater than start should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 3, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(SelectorKey_Start).AssignInt(3)
			na.AssembleDirectly(SelectorKey_End).AssignInt(2)
			na.AssembleDirectly(SelectorKey_Next).CreateMap(1, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapNodeAssembler) {})
			})
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: end field must be greater than start field in ExploreRange selector"))
	})
	t.Run("parsing map node with next field with invalid selector node should return child's error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 3, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(SelectorKey_Start).AssignInt(2)
			na.AssembleDirectly(SelectorKey_End).AssignInt(3)
			na.AssembleDirectly(SelectorKey_Next).AssignInt(0)
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map"))
	})

	t.Run("parsing map node with next field with valid selector node should parse", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 3, func(na fluent.MapNodeAssembler) {
			na.AssembleDirectly(SelectorKey_Start).AssignInt(2)
			na.AssembleDirectly(SelectorKey_End).AssignInt(3)
			na.AssembleDirectly(SelectorKey_Next).CreateMap(1, func(na fluent.MapNodeAssembler) {
				na.AssembleDirectly(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapNodeAssembler) {})
			})
		})
		s, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreRange{Matcher{}, 2, 3, []ipld.PathSegment{ipld.PathSegmentOfInt(2)}})
	})
}

func TestExploreRangeExplore(t *testing.T) {
	s := ExploreRange{Matcher{}, 3, 4, []ipld.PathSegment{ipld.PathSegmentOfInt(3)}}
	t.Run("exploring should return nil unless node is a list", func(t *testing.T) {
		n := fluent.MustBuildMap(basicnode.Style__Map{}, 0, func(na fluent.MapNodeAssembler) {})
		returnedSelector := s.Explore(n, ipld.PathSegmentOfInt(3))
		Wish(t, returnedSelector, ShouldEqual, nil)
	})
	n := fluent.MustBuildList(basicnode.Style__List{}, 4, func(na fluent.ListNodeAssembler) {
		na.AssembleValue().AssignInt(0)
		na.AssembleValue().AssignInt(1)
		na.AssembleValue().AssignInt(2)
		na.AssembleValue().AssignInt(3)
	})
	t.Run("exploring should return nil when given a path segment out of range", func(t *testing.T) {
		returnedSelector := s.Explore(n, ipld.PathSegmentOfInt(2))
		Wish(t, returnedSelector, ShouldEqual, nil)
	})
	t.Run("exploring should return nil when given a path segment that isn't an index", func(t *testing.T) {
		returnedSelector := s.Explore(n, ipld.PathSegmentOfString("cheese"))
		Wish(t, returnedSelector, ShouldEqual, nil)
	})
	t.Run("exploring should return the next selector when given a path segment with index in range", func(t *testing.T) {
		returnedSelector := s.Explore(n, ipld.PathSegmentOfInt(3))
		Wish(t, returnedSelector, ShouldEqual, Matcher{})
	})
}
