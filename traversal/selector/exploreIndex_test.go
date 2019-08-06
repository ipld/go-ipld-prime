package selector

import (
	"fmt"
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	. "github.com/warpfork/go-wish"
)

func TestParseExploreIndex(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := fnb.CreateInt(0)
		_, err := ParseContext{}.ParseExploreIndex(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector body must be a map"))
	})
	t.Run("parsing map node without next field should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(indexKey), vnb.CreateInt(2))
		})
		_, err := ParseContext{}.ParseExploreIndex(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: next field must be present in ExploreIndex selector"))
	})
	t.Run("parsing map node without index field should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreIndex(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: index field must be present in ExploreIndex selector"))
	})
	t.Run("parsing map node with index field that is not an int should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(indexKey), vnb.CreateString("cheese"))
			mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		_, err := ParseContext{}.ParseExploreIndex(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: index field must be a number in ExploreIndex selector"))
	})
	t.Run("parsing map node with next field with invalid selector node should return child's error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(indexKey), vnb.CreateInt(2))
			mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateInt(0))
		})
		_, err := ParseContext{}.ParseExploreIndex(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map"))
	})
	t.Run("parsing map node with next field with valid selector node should parse", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(indexKey), vnb.CreateInt(2))
			mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		s, err := ParseContext{}.ParseExploreIndex(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreIndex{Matcher{}, [1]PathSegment{PathSegmentInt{I: 2}}})
	})
}

func TestExploreIndexExplore(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	s := ExploreIndex{Matcher{}, [1]PathSegment{PathSegmentInt{I: 3}}}
	t.Run("exploring should return nil unless node is a list", func(t *testing.T) {
		n := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {})
		returnedSelector := s.Explore(n, PathSegmentInt{I: 3})
		Wish(t, returnedSelector, ShouldEqual, nil)
	})
	t.Run("exploring should return nil when given a path segment with a different index", func(t *testing.T) {
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
	t.Run("exploring should return the next selector when given a path segment with the right index", func(t *testing.T) {
		n := fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.AppendAll([]ipld.Node{fnb.CreateInt(0), fnb.CreateInt(1), fnb.CreateInt(2), fnb.CreateInt(3)})
		})
		returnedSelector := s.Explore(n, PathSegmentInt{I: 3})
		Wish(t, returnedSelector, ShouldEqual, Matcher{})
	})
}
