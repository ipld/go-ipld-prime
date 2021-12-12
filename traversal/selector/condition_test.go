package selector

import (
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp"
	"github.com/ipfs/go-cid"

	"github.com/ipld/go-ipld-prime/fluent"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

var deepEqualsAllowAllUnexported = qt.CmpEquals(cmp.Exporter(func(reflect.Type) bool { return true }))

func TestParseCondition(t *testing.T) {
	t.Run("parsing non map node should error", func(t *testing.T) {
		sn := basicnode.NewInt(0)
		_, err := ParseContext{}.ParseCondition(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: condition body must be a map")
	})
	t.Run("parsing map node without field should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype.Map, 0, func(na fluent.MapAssembler) {})
		_, err := ParseContext{}.ParseCondition(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: condition is a keyed union and thus must be single-entry map")
	})

	t.Run("parsing map node keyed to invalid type should error", func(t *testing.T) {
		sn := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(string(ConditionMode_Link)).AssignInt(0)
		})
		_, err := ParseContext{}.ParseCondition(sn)
		qt.Check(t, err, qt.ErrorMatches, "selector spec parse rejected: condition_link must be a link")
	})
	t.Run("parsing map node with condition field with valid selector node should parse", func(t *testing.T) {
		lnk := cidlink.Link{Cid: cid.Undef}
		sn := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
			na.AssembleEntry(string(ConditionMode_Link)).AssignLink(lnk)
		})
		s, err := ParseContext{}.ParseCondition(sn)
		qt.Check(t, err, qt.IsNil)
		lnkNode := basicnode.NewLink(lnk)
		qt.Check(t, s, deepEqualsAllowAllUnexported, Condition{mode: ConditionMode_Link, match: lnkNode})
	})
}
