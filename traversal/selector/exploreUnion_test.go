package selector

import (
	"fmt"
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	. "github.com/warpfork/go-wish"
)

func TestParseExploreUnion(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	t.Run("parsing non list node should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {})
		_, err := ParseExploreUnion(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: explore union selector must be a list"))
	})
	t.Run("parsing list node where one node is invalid should return child's error", func(t *testing.T) {
		sn := fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
			lb.Append(vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
			lb.Append(vnb.CreateInt(2))
		})
		_, err := ParseExploreUnion(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map"))
	})

	t.Run("parsing map node with next field with valid selector node should parse", func(t *testing.T) {
		sn := fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
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
		})
		s, err := ParseExploreUnion(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreUnion{[]Selector{Matcher{}, ExploreIndex{Matcher{}, []PathSegment{PathSegmentInt{I: 2}}}}})
	})
}

func TestExploreUnionExplore(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	n := fnb.CreateList(func(lb fluent.ListBuilder, vnb fluent.NodeBuilder) {
		lb.AppendAll([]ipld.Node{fnb.CreateInt(0), fnb.CreateInt(1), fnb.CreateInt(2), fnb.CreateInt(3)})
	})
	t.Run("exploring should return nil if all member selectors return nil when explored", func(t *testing.T) {
		s := ExploreUnion{[]Selector{Matcher{}, ExploreIndex{Matcher{}, []PathSegment{PathSegmentInt{I: 2}}}}}
		returnedSelector := s.Explore(n, PathSegmentInt{I: 3})
		Wish(t, returnedSelector, ShouldEqual, nil)
	})

	t.Run("if exactly one member selector returns a non-nil selector when explored, exploring should return that value", func(t *testing.T) {
		s := ExploreUnion{[]Selector{Matcher{}, ExploreIndex{Matcher{}, []PathSegment{PathSegmentInt{I: 2}}}}}

		returnedSelector := s.Explore(n, PathSegmentInt{I: 2})
		Wish(t, returnedSelector, ShouldEqual, Matcher{})
	})
	t.Run("exploring should return a new union selector if more than one member selector returns a non nil selector when explored", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			Matcher{},
			ExploreIndex{Matcher{}, []PathSegment{PathSegmentInt{I: 2}}},
			ExploreRange{Matcher{}, map[int]struct{}{2: struct{}{}}, []PathSegment{PathSegmentInt{I: 2}}},
			ExploreFields{map[string]Selector{"applesauce": Matcher{}}, []PathSegment{PathSegmentString{S: "applesauce"}}},
		}}

		returnedSelector := s.Explore(n, PathSegmentInt{I: 2})
		Wish(t, returnedSelector, ShouldEqual, ExploreUnion{[]Selector{Matcher{}, Matcher{}}})
	})
}

func TestExploreUnionInterests(t *testing.T) {
	t.Run("if any member selector is high-cardinality, interests should be high-cardinality", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			ExploreAll{Matcher{}},
			Matcher{},
			ExploreIndex{Matcher{}, []PathSegment{PathSegmentInt{I: 2}}},
		}}
		Wish(t, s.Interests(), ShouldEqual, []PathSegment(nil))
	})
	t.Run("if no member selector is high-cardinality, interests should be combination of member selectors interests", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			ExploreFields{map[string]Selector{"applesauce": Matcher{}}, []PathSegment{PathSegmentString{S: "applesauce"}}},
			Matcher{},
			ExploreIndex{Matcher{}, []PathSegment{PathSegmentInt{I: 2}}},
		}}
		Wish(t, s.Interests(), ShouldEqual, []PathSegment{PathSegmentString{S: "applesauce"}, PathSegmentInt{I: 2}})
	})
}

func TestExploreUnionDecide(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	n := fnb.CreateInt(2)
	t.Run("if any member selector returns true, decide should be true", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			ExploreAll{Matcher{}},
			Matcher{},
			ExploreIndex{Matcher{}, []PathSegment{PathSegmentInt{I: 2}}},
		}}
		Wish(t, s.Decide(n), ShouldEqual, true)
	})
	t.Run("if no member selector returns true, decide should be false", func(t *testing.T) {
		s := ExploreUnion{[]Selector{
			ExploreFields{map[string]Selector{"applesauce": Matcher{}}, []PathSegment{PathSegmentString{S: "applesauce"}}},
			ExploreAll{Matcher{}},
			ExploreIndex{Matcher{}, []PathSegment{PathSegmentInt{I: 2}}},
		}}
		Wish(t, s.Decide(n), ShouldEqual, false)
	})
}
