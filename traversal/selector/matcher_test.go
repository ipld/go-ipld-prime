package selector_test

import (
	"fmt"
	"math"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/mixins"
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
		{"bytesNode", simpleBytes([]byte(expectedString))},
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
		{-1000, -100, "", false},
		{-1, -2, "", false},
		{-1, -1, "", true}, // matches empty node
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
}

var _ datamodel.Node = simpleBytes{}

// simpleBytes is like basicnode's plainBytes but it doesn't implement
// LargeBytesNode so we can exercise the non-LBN case.
type simpleBytes []byte

// -- Node interface methods -->

func (simpleBytes) Kind() datamodel.Kind {
	return datamodel.Kind_Bytes
}
func (simpleBytes) LookupByString(string) (datamodel.Node, error) {
	return mixins.Bytes{TypeName: "bytes"}.LookupByString("")
}
func (simpleBytes) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	return mixins.Bytes{TypeName: "bytes"}.LookupByNode(nil)
}
func (simpleBytes) LookupByIndex(idx int64) (datamodel.Node, error) {
	return mixins.Bytes{TypeName: "bytes"}.LookupByIndex(0)
}
func (simpleBytes) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	return mixins.Bytes{TypeName: "bytes"}.LookupBySegment(seg)
}
func (simpleBytes) MapIterator() datamodel.MapIterator {
	return nil
}
func (simpleBytes) ListIterator() datamodel.ListIterator {
	return nil
}
func (simpleBytes) Length() int64 {
	return -1
}
func (simpleBytes) IsAbsent() bool {
	return false
}
func (simpleBytes) IsNull() bool {
	return false
}
func (simpleBytes) AsBool() (bool, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsBool()
}
func (simpleBytes) AsInt() (int64, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsInt()
}
func (simpleBytes) AsFloat() (float64, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsFloat()
}
func (simpleBytes) AsString() (string, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsString()
}
func (n simpleBytes) AsBytes() ([]byte, error) {
	return []byte(n), nil
}
func (simpleBytes) AsLink() (datamodel.Link, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsLink()
}
func (simpleBytes) Prototype() datamodel.NodePrototype {
	return basicnode.Prototype__Bytes{}
}
