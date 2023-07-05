package selector_test

import (
	"fmt"
	"math"
	"regexp"
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
		{"bytesNode", testutil.NewSimpleBytes([]byte(expectedString))},
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
		from  int64
		to    int64
		exp   string
		match bool
	}{
		{0, math.MaxInt64, expectedString, true},
		{0, int64(len(expectedString)), expectedString, true},
		{0, 0, "", true},
		{0, 1, "f", true},
		{0, 2, "fo", true},
		{0, 3, "foo", true},
		{0, 4, "foob", true},
		{1, 4, "oob", true},
		{2, 4, "ob", true},
		{3, 4, "b", true},
		{4, 4, "", true},
		{4, math.MaxInt64, "arbaz!", true},
		{4, int64(len(expectedString)), "arbaz!", true},
		{4, int64(len(expectedString) - 1), "arbaz", true},
		{0, int64(len(expectedString) - 1), expectedString[0 : len(expectedString)-1], true},
		{0, int64(len(expectedString) - 2), expectedString[0 : len(expectedString)-2], true},
		{0, -1, expectedString[0 : len(expectedString)-1], true},
		{0, -2, expectedString[0 : len(expectedString)-2], true},
		{-2, -1, "z", true},
		{-1, math.MaxInt64, "!", true},
		{-int64(len(expectedString)), math.MaxInt64, expectedString, true},
		{math.MaxInt64 - 1, math.MaxInt64, "", false},
		{int64(len(expectedString)), math.MaxInt64, "", false},
		{-1, -2, "", false},                 // To < From, no match
		{-1, -1, "", true},                  // To==From, match zero bytes
		{-1000, -100, "", false},            // From undeflow, adjusted to 0, To underflow, not adjusted, To < From, no match
		{-100, -1000, "", false},            // From undeflow, adjusted to 0, To underflow, adjusted to 0, To < From, no match
		{-1000, 1000, expectedString, true}, // From undeflow, adjusted to 0, To overflow, adjusted to len, match all
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

				if tc.match {
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
				} else {
					qt.Assert(t, got, qt.IsNil)
				}
			})
		}
	}

	// when both are positive, we can validate ranges up-front
	t.Run("invalid range", func(t *testing.T) {
		selNode, err := mkRangeSelector(1000, 100)
		qt.Assert(t, err, qt.IsNil)
		re, err := regexp.Compile("from.*less than or equal to.*to")
		qt.Assert(t, err, qt.IsNil)
		ss, err := selector.ParseSelector(selNode)
		qt.Assert(t, ss, qt.IsNil)
		qt.Assert(t, err, qt.ErrorMatches, re)
	})
}
