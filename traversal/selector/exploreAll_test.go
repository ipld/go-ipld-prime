package selector

import (
	"fmt"
	"testing"

	"github.com/ipld/go-ipld-prime/fluent"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	. "github.com/warpfork/go-wish"
)

func TestParseExploreAll(t *testing.T) {
	fnb := fluent.WrapNodeBuilder(ipldfree.NodeBuilder()) // just for the other fixture building
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := fnb.CreateInt(0)
		_, err := ParseContext{}.ParseExploreAll(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector body must be a map"))
	})
	t.Run("parsing map node without next field should error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {})
		_, err := ParseContext{}.ParseExploreAll(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: next field must be present in ExploreAll selector"))
	})

	t.Run("parsing map node without next field with invalid selector node should return child's error", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateInt(0))
		})
		_, err := ParseContext{}.ParseExploreAll(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map"))
	})
	t.Run("parsing map node with next field with valid selector node should parse", func(t *testing.T) {
		sn := fnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
			mb.Insert(knb.CreateString(nextSelectorKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {
				mb.Insert(knb.CreateString(matcherKey), vnb.CreateMap(func(mb fluent.MapBuilder, knb fluent.NodeBuilder, vnb fluent.NodeBuilder) {}))
			}))
		})
		s, err := ParseContext{}.ParseExploreAll(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreAll{Matcher{}})
	})
}
