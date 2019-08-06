package selector

import (
	"fmt"
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	. "github.com/warpfork/go-wish"
)

func TestParseExploreRange(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := fnb.CreateInt(0)
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector body must be a map"))
	})
	t.Run("parsing map node without next field should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Start), vnb.CreateInt(2))
			mb.Insert(knb.CreateString(SelectorKey_End), vnb.CreateInt(3))
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: next field must be present in ExploreRange selector"))
	})
	t.Run("parsing map node without start field should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_End), vnb.CreateInt(3))
			mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: start field must be present in ExploreRange selector"))
	})
	t.Run("parsing map node with start field that is not an int should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Start), vnb.CreateString("cheese"))
			mb.Insert(knb.CreateString(SelectorKey_End), vnb.CreateInt(3))
			mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: start field must be a number in ExploreRange selector"))
	})
	t.Run("parsing map node without end field should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Start), vnb.CreateInt(2))
			mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: end field must be present in ExploreRange selector"))
	})
	t.Run("parsing map node with end field that is not an int should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Start), vnb.CreateInt(2))
			mb.Insert(knb.CreateString(SelectorKey_End), vnb.CreateString("cheese"))
			mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: end field must be a number in ExploreRange selector"))
	})
	t.Run("parsing map node where end is not greater than start should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Start), vnb.CreateInt(3))
			mb.Insert(knb.CreateString(SelectorKey_End), vnb.CreateInt(2))
			mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: end field must be greater than start field in ExploreRange selector"))
	})
	t.Run("parsing map node with next field with invalid selector node should return child's error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Start), vnb.CreateInt(2))
			mb.Insert(knb.CreateString(SelectorKey_End), vnb.CreateInt(3))
			mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateInt(0))
		})
		_, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map"))
	})

	t.Run("parsing map node with next field with valid selector node should parse", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Start), vnb.CreateInt(2))
			mb.Insert(knb.CreateString(SelectorKey_End), vnb.CreateInt(3))
			mb.Insert(knb.CreateString(SelectorKey_Next), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		s, err := ParseContext{}.ParseExploreRange(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreRange{Matcher{}, 2, 3, []PathSegment{PathSegmentInt{I: 2}}})
	})
}

func TestExploreRangeExplore(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	s := ExploreRange{Matcher{}, 3, 4, []PathSegment{PathSegmentInt{I: 3}}}
	t.Run("exploring should return nil unless node is a list", func(t *testing.T) {
		n := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {})
		returnedSelector := s.Explore(n, PathSegmentInt{I: 3})
		Wish(t, returnedSelector, ShouldEqual, nil)
	})
	t.Run("exploring should return nil when given a path segment out of range", func(t *testing.T) {
		n := fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.AppendAll([]ipld.Node{fnb.CreateInt(0), fnb.CreateInt(1), fnb.CreateInt(2), fnb.CreateInt(3)})
		})
		returnedSelector := s.Explore(n, PathSegmentInt{I: 2})
		Wish(t, returnedSelector, ShouldEqual, nil)
	})
	t.Run("exploring should return nil when given a path segment that isn't an index", func(t *testing.T) {
		n := fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.AppendAll([]ipld.Node{fnb.CreateInt(0), fnb.CreateInt(1), fnb.CreateInt(2), fnb.CreateInt(3)})
		})
		returnedSelector := s.Explore(n, PathSegmentString{S: "cheese"})
		Wish(t, returnedSelector, ShouldEqual, nil)
	})
	t.Run("exploring should return the next selector when given a path segment with index in range", func(t *testing.T) {
		n := fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.AppendAll([]ipld.Node{fnb.CreateInt(0), fnb.CreateInt(1), fnb.CreateInt(2), fnb.CreateInt(3)})
		})
		returnedSelector := s.Explore(n, PathSegmentInt{I: 3})
		Wish(t, returnedSelector, ShouldEqual, Matcher{})
	})
}
