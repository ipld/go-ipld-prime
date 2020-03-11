package selector

import (
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestParseExploreFields(t *testing.T) {
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := basicnode.NewInt(0)
		_, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector body must be a map"))
	})
	t.Run("parsing map node without fields value should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 0, func(na fluent.MapAssembler) {})
		_, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: fields in ExploreFields selector must be present"))
	})
	t.Run("parsing map node with fields value that is not a map should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(SelectorKey_Fields).AssignString("cheese")
		})
		_, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: fields in ExploreFields selector must be a map"))
	})
	t.Run("parsing map node with selector node in fields that is invalid should return child's error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(SelectorKey_Fields).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleDirectly("applesauce").AssignInt(0)
			})
		})
		_, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, fmt.Errorf("selector spec parse rejected: selector is a keyed union and thus must be a map"))
	})
	t.Run("parsing map node with fields value that is map of only valid selector node should parse", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Style__Map{}, 1, func(na fluent.MapAssembler) {
			na.AssembleDirectly(SelectorKey_Fields).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleDirectly("applesauce").CreateMap(1, func(na fluent.MapAssembler) {
					na.AssembleDirectly(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
				})
			})
		})
		s, err := ParseContext{}.ParseExploreFields(sn)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, ExploreFields{map[string]Selector{"applesauce": Matcher{}}, []ipld.PathSegment{ipld.PathSegmentOfString("applesauce")}})
	})
}
