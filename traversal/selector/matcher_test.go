package selector_test

import (
	"fmt"
	"math"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/testutil"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func TestSubsetMatch(t *testing.T) {
	expectedString := "foobarbaz!"
	nodes := []struct {
		name string
		node datamodel.Node
	}{
		{"stringNode", basicnode.NewString(expectedString)},
		{"bytesNode", basicnode.NewBytes([]byte(expectedString))},
		{"largeBytesNode", testutil.NewMultiByteNode(
			[]byte("foo"),
			[]byte("bar"),
			[]byte("baz"),
			[]byte("!"),
		)},
	}

	// selector for a slice of the value of the "bipbop" field within a map
	mkRangeSelector := func(from int64, to int64) (datamodel.Node, error) {
		return qp.BuildMap(basicnode.Prototype.Map, 1, func(na datamodel.MapAssembler) {
			qp.MapEntry(na, selector.SelectorKey_ExploreFields, qp.Map(1, func(na datamodel.MapAssembler) {
				qp.MapEntry(na, selector.SelectorKey_Fields, qp.Map(1, func(na datamodel.MapAssembler) {
					qp.MapEntry(na, "bipbop", qp.Map(1, func(na datamodel.MapAssembler) {
						qp.MapEntry(na, selector.SelectorKey_Matcher, qp.Map(1, func(na datamodel.MapAssembler) {
							qp.MapEntry(na, selector.SelectorKey_Subset, qp.Map(1, func(na datamodel.MapAssembler) {
								qp.MapEntry(na, selector.SelectorKey_From, qp.Int(from))
								qp.MapEntry(na, selector.SelectorKey_To, qp.Int(to))
							}))
						}))
					}))
				}))
			}))
		})
	}

	for _, tc := range []struct {
		from int64
		to   int64
		exp  string
	}{
		{0, math.MaxInt64, expectedString},
		{0, int64(len(expectedString)), expectedString},
		{0, 0, ""},
		{0, 1, "f"},
		{0, 2, "fo"},
		{0, 3, "foo"},
		{0, 4, "foob"},
		{1, 4, "oob"},
		{2, 4, "ob"},
		{3, 4, "b"},
		{4, 4, ""},
		{4, math.MaxInt64, "arbaz!"},
		{4, int64(len(expectedString)), "arbaz!"},
		{4, int64(len(expectedString) - 1), "arbaz"},
		{0, int64(len(expectedString) - 1), expectedString[0 : len(expectedString)-1]},
		{0, int64(len(expectedString) - 2), expectedString[0 : len(expectedString)-2]},
	} {
		for _, variant := range nodes {
			t.Run(fmt.Sprintf("%s[%d:%d]", variant.name, tc.from, tc.to), func(t *testing.T) {
				selNode, err := mkRangeSelector(tc.from, tc.to)
				qt.Assert(t, err, qt.IsNil)
				ss, err := selector.ParseSelector(selNode)
				qt.Assert(t, err, qt.IsNil)

				// node that the selector will match, with our variant node embedded in the "bipbop" field
				n, err := qp.BuildMap(basicnode.Prototype.Map, 1, func(na datamodel.MapAssembler) {
					qp.MapEntry(na, "bipbop", qp.Node(variant.node))
				})

				var got datamodel.Node
				qt.Assert(t, err, qt.IsNil)
				err = traversal.WalkMatching(n, ss, func(prog traversal.Progress, n datamodel.Node) error {
					qt.Assert(t, got, qt.IsNil)
					got = n
					return nil
				})
				qt.Assert(t, err, qt.IsNil)

				qt.Assert(t, got, qt.IsNotNil)
				qt.Assert(t, got.Kind(), qt.Equals, variant.node.Kind())
				var gotString string
				switch got.Kind() {
				case datamodel.Kind_String:
					gotString, err = got.AsString()
					qt.Assert(t, err, qt.IsNil)
				case datamodel.Kind_Bytes:
					byts, err := got.AsBytes()
					qt.Assert(t, err, qt.IsNil)
					gotString = string(byts)
				}
				qt.Assert(t, gotString, qt.DeepEquals, tc.exp)
			})
		}
	}
}
