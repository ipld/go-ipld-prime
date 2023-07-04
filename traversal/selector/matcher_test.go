package selector_test

import (
	"fmt"
	"io"
	"math"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
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
		{"largeBytesNode", &MultiByteNode{
			Bytes: [][]byte{
				[]byte("foo"),
				[]byte("bar"),
				[]byte("baz"),
				[]byte("!"),
			},
		}},
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

func TestMultiByteNode_Sanity(t *testing.T) {
	mbn := &MultiByteNode{
		Bytes: [][]byte{
			[]byte("foo"),
			[]byte("bar"),
			[]byte("baz"),
			[]byte("!"),
		},
	}
	// Sanity check that the readseeker works.
	// (This is a test of the test, not the code under test.)

	for _, rl := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		t.Run("readseeker works with read length "+qt.Format(rl), func(t *testing.T) {
			rs, err := mbn.AsLargeBytes()
			qt.Assert(t, err, qt.IsNil)
			acc := make([]byte, 0, mbn.size())
			buf := make([]byte, rl)
			for {
				n, err := rs.Read(buf)
				if err == io.EOF {
					qt.Check(t, n, qt.Equals, 0)
					break
				}
				qt.Assert(t, err, qt.IsNil)
				acc = append(acc, buf[0:n]...)
			}
			qt.Assert(t, string(acc), qt.DeepEquals, "foobarbaz!")
		})
	}

	t.Run("readseeker can seek and read middle bytes", func(t *testing.T) {
		rs, err := mbn.AsLargeBytes()
		qt.Assert(t, err, qt.IsNil)
		_, err = rs.Seek(2, io.SeekStart)
		qt.Assert(t, err, qt.IsNil)
		buf := make([]byte, 2)
		acc := make([]byte, 0, 5)
		for len(acc) < 5 {
			n, err := rs.Read(buf)
			qt.Assert(t, err, qt.IsNil)
			acc = append(acc, buf[0:n]...)
		}
		qt.Assert(t, string(acc), qt.DeepEquals, "obarba")
	})

	t.Run("readseeker can seek and read last byte", func(t *testing.T) {
		rs, err := mbn.AsLargeBytes()
		qt.Assert(t, err, qt.IsNil)
		_, err = rs.Seek(-1, io.SeekEnd)
		qt.Assert(t, err, qt.IsNil)
		buf := make([]byte, 1)
		n, err := rs.Read(buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, n, qt.Equals, 1)
		qt.Check(t, string(buf[0]), qt.Equals, "!")
	})
}

var _ datamodel.Node = (*MultiByteNode)(nil)
var _ datamodel.LargeBytesNode = (*MultiByteNode)(nil)

// MultiByteNode is a node that is a concatenation of multiple byte slices.
// It's not particularly sophisticated but lets us exercise LargeBytesNode as a
// path through the selectors. The novel behaviour of Read() and Seek() on the
// AsLargeBytes is similar to that which would be expected from a LBN ADL, such
// as UnixFS sharded files.
type MultiByteNode struct {
	Bytes [][]byte
}

func (mbn *MultiByteNode) Kind() datamodel.Kind {
	return datamodel.Kind_Bytes
}

func (mbn *MultiByteNode) AsBytes() ([]byte, error) {
	ret := make([]byte, 0, mbn.size())
	for _, b := range mbn.Bytes {
		ret = append(ret, b...)
	}
	return ret, nil
}

func (mbn *MultiByteNode) size() int {
	var size int
	for _, b := range mbn.Bytes {
		size += len(b)
	}
	return size
}

func (mbn *MultiByteNode) AsLargeBytes() (io.ReadSeeker, error) {
	return &mbnReadSeeker{node: mbn}, nil
}

func (mbn *MultiByteNode) AsBool() (bool, error) {
	return false, datamodel.ErrWrongKind{TypeName: "bool", MethodName: "AsBool", AppropriateKind: datamodel.KindSet_JustBytes}
}

func (mbn *MultiByteNode) AsInt() (int64, error) {
	return 0, datamodel.ErrWrongKind{TypeName: "int", MethodName: "AsInt", AppropriateKind: datamodel.KindSet_JustBytes}
}

func (mbn *MultiByteNode) AsFloat() (float64, error) {
	return 0, datamodel.ErrWrongKind{TypeName: "float", MethodName: "AsFloat", AppropriateKind: datamodel.KindSet_JustBytes}
}

func (mbn *MultiByteNode) AsString() (string, error) {
	return "", datamodel.ErrWrongKind{TypeName: "string", MethodName: "AsString", AppropriateKind: datamodel.KindSet_JustBytes}
}

func (mbn *MultiByteNode) AsLink() (datamodel.Link, error) {
	return nil, datamodel.ErrWrongKind{TypeName: "link", MethodName: "AsLink", AppropriateKind: datamodel.KindSet_JustBytes}
}

func (mbn *MultiByteNode) AsNode() (datamodel.Node, error) {
	return nil, nil
}

func (mbn *MultiByteNode) Size() int {
	return 0
}

func (mbn *MultiByteNode) IsAbsent() bool {
	return false
}

func (mbn *MultiByteNode) IsNull() bool {
	return false
}

func (mbn *MultiByteNode) Length() int64 {
	return 0
}

func (mbn *MultiByteNode) ListIterator() datamodel.ListIterator {
	return nil
}

func (mbn *MultiByteNode) MapIterator() datamodel.MapIterator {
	return nil
}

func (mbn *MultiByteNode) LookupByIndex(idx int64) (datamodel.Node, error) {
	return nil, datamodel.ErrWrongKind{}
}

func (mbn *MultiByteNode) LookupByString(key string) (datamodel.Node, error) {
	return nil, datamodel.ErrWrongKind{}
}

func (mbn *MultiByteNode) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	return nil, datamodel.ErrWrongKind{}
}

func (mbn *MultiByteNode) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	return nil, datamodel.ErrWrongKind{}
}

func (mbn *MultiByteNode) Prototype() datamodel.NodePrototype {
	return basicnode.Prototype.Bytes // not really ... but it'll do for this test
}

type mbnReadSeeker struct {
	node   *MultiByteNode
	offset int
}

func (mbnrs *mbnReadSeeker) Read(p []byte) (int, error) {
	var acc int
	for _, byts := range mbnrs.node.Bytes {
		if mbnrs.offset-acc >= len(byts) {
			acc += len(byts)
			continue
		}
		n := copy(p, byts[mbnrs.offset-acc:])
		mbnrs.offset += n
		return n, nil
	}
	return 0, io.EOF
}

func (mbnrs *mbnReadSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		mbnrs.offset = int(offset)
	case io.SeekCurrent:
		mbnrs.offset += int(offset)
	case io.SeekEnd:
		mbnrs.offset = mbnrs.node.size() + int(offset)
	}
	return int64(mbnrs.offset), nil
}
