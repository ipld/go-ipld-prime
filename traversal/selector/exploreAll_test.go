package selector

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestParseExploreAll(t *testing.T) {
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := basicnode.NewInt(0)
		_, err := ParseContext{}.ParseExploreAll(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: selector body must be a map")
	})
	t.Run("parsing map node without next field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype.Map, 0, func(na fluent.MapAssembler) {})
		_, err := ParseContext{}.ParseExploreAll(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: next field must be present in ExploreAll selector")
	})

	t.Run("parsing map node without next field with invalid selector node should return child's error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Next).AssignInt(0)
		})
		_, err := ParseContext{}.ParseExploreAll(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: selector is a keyed union and thus must be a map")
	})
	t.Run("parsing map node with next field with valid selector node should parse", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(SelectorKey_Next).CreateMap(1, func(na fluent.MapAssembler) {
				na.AssembleEntry(SelectorKey_Matcher).CreateMap(0, func(na fluent.MapAssembler) {})
			})
		})
		s, err := ParseContext{}.ParseExploreAll(sn)
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, s, qt.Equals, ExploreAll{Matcher{}})
	})
}
