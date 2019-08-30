package selector

import (
	"fmt"
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	. "github.com/warpfork/go-wish"
)

func TestParseExploreFields(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := fnb.CreateInt(0)
		_, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector body must be a map"))
	})
	t.Run("parsing map node without fields value should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {})
		_, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: fields in ExploreFields selector must be present"))
	})
	t.Run("parsing map node with fields value that is not a map should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Fields), vnb.CreateString("cheese"))
		})
		_, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: fields in ExploreFields selector must be a map"))
	})
	t.Run("parsing map node with selector node in fields that is invalid should return child's error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Fields), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString("applesauce"), vnb.CreateInt(0))
			}))
		})
		_, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map"))
	})
	t.Run("parsing map node with fields value that is map of only valid selector node should parse", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(SelectorKey_Fields), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString("applesauce"), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
					mb.Insert(knb.CreateString(SelectorKey_Matcher), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
				}))
			}))
		})
		s, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreFields{map[string]Selector{"applesauce": Matcher{}}, []ipld.PathSegment{ipld.PathSegmentOfString("applesauce")}})
	})
}
