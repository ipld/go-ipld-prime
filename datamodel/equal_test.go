package datamodel_test

import (
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basic "github.com/ipld/go-ipld-prime/node/basicnode" // shorter name for the tests
)

var (
	globalNode = basic.NewString("global")
	globalLink = func() datamodel.Link {
		someCid, _ := cid.Cast([]byte{1, 85, 0, 5, 0, 1, 2, 3, 4})
		return cidlink.Link{Cid: someCid}
	}()
	globalLink2 = func() datamodel.Link {
		someCid, _ := cid.Cast([]byte{1, 85, 0, 5, 0, 5, 6, 7, 8})
		return cidlink.Link{Cid: someCid}
	}()
)

func qpMust(node datamodel.Node, err error) datamodel.Node {
	if err != nil {
		panic(err)
	}
	return node
}

var deepEqualTests = []struct {
	name        string
	left, right datamodel.Node
	want        bool
}{
	{"MismatchingKinds", basic.NewBool(true), basic.NewInt(3), false},

	{"SameNodeSamePointer", globalNode, globalNode, true},
	// Repeated basicnode.New invocations might return different pointers.
	{"SameNodeDiffPointer", basic.NewString("same"), basic.NewString("same"), true},

	{"NilVsNil", nil, nil, true},
	{"NilVsNull", nil, datamodel.Null, false},
	{"SameKindNull", datamodel.Null, datamodel.Null, true},
	{"DiffKindNull", datamodel.Null, datamodel.Absent, false},
	{"SameKindBool", basic.NewBool(true), basic.NewBool(true), true},
	{"DiffKindBool", basic.NewBool(true), basic.NewBool(false), false},
	{"SameKindInt", basic.NewInt(12), basic.NewInt(12), true},
	{"DiffKindInt", basic.NewInt(12), basic.NewInt(15), false},
	{"SameKindFloat", basic.NewFloat(1.25), basic.NewFloat(1.25), true},
	{"DiffKindFloat", basic.NewFloat(1.25), basic.NewFloat(1.75), false},
	{"SameKindString", basic.NewString("foobar"), basic.NewString("foobar"), true},
	{"DiffKindString", basic.NewString("foobar"), basic.NewString("baz"), false},
	{"SameKindBytes", basic.NewBytes([]byte{5, 2, 3}), basic.NewBytes([]byte{5, 2, 3}), true},
	{"DiffKindBytes", basic.NewBytes([]byte{5, 2, 3}), basic.NewBytes([]byte{5, 8, 3}), false},
	{"SameKindLink", basic.NewLink(globalLink), basic.NewLink(globalLink), true},
	{"DiffKindLink", basic.NewLink(globalLink), basic.NewLink(globalLink2), false},

	{
		"SameKindList",
		qpMust(qp.BuildList(basic.Prototype.Any, -1, func(am datamodel.ListAssembler) {
			qp.ListEntry(am, qp.Int(7))
			qp.ListEntry(am, qp.Int(8))
		})),
		qpMust(qp.BuildList(basic.Prototype.Any, -1, func(am datamodel.ListAssembler) {
			qp.ListEntry(am, qp.Int(7))
			qp.ListEntry(am, qp.Int(8))
		})),
		true,
	},
	{
		"DiffKindList_length",
		qpMust(qp.BuildList(basic.Prototype.Any, -1, func(am datamodel.ListAssembler) {
			qp.ListEntry(am, qp.Int(7))
			qp.ListEntry(am, qp.Int(8))
		})),
		qpMust(qp.BuildList(basic.Prototype.Any, -1, func(am datamodel.ListAssembler) {
			qp.ListEntry(am, qp.Int(7))
		})),
		false,
	},
	{
		"DiffKindList_elems",
		qpMust(qp.BuildList(basic.Prototype.Any, -1, func(am datamodel.ListAssembler) {
			qp.ListEntry(am, qp.Int(7))
			qp.ListEntry(am, qp.Int(8))
		})),
		qpMust(qp.BuildList(basic.Prototype.Any, -1, func(am datamodel.ListAssembler) {
			qp.ListEntry(am, qp.Int(3))
			qp.ListEntry(am, qp.Int(2))
		})),
		false,
	},

	{
		"SameKindMap",
		qpMust(qp.BuildMap(basic.Prototype.Any, -1, func(am datamodel.MapAssembler) {
			qp.MapEntry(am, "foo", qp.Int(7))
			qp.MapEntry(am, "bar", qp.Int(8))
		})),
		qpMust(qp.BuildMap(basic.Prototype.Any, -1, func(am datamodel.MapAssembler) {
			qp.MapEntry(am, "foo", qp.Int(7))
			qp.MapEntry(am, "bar", qp.Int(8))
		})),
		true,
	},
	{
		"DiffKindMap_length",
		qpMust(qp.BuildMap(basic.Prototype.Any, -1, func(am datamodel.MapAssembler) {
			qp.MapEntry(am, "foo", qp.Int(7))
			qp.MapEntry(am, "bar", qp.Int(8))
		})),
		qpMust(qp.BuildMap(basic.Prototype.Any, -1, func(am datamodel.MapAssembler) {
			qp.MapEntry(am, "foo", qp.Int(7))
		})),
		false,
	},
	{
		"DiffKindMap_elems",
		qpMust(qp.BuildMap(basic.Prototype.Any, -1, func(am datamodel.MapAssembler) {
			qp.MapEntry(am, "foo", qp.Int(7))
			qp.MapEntry(am, "bar", qp.Int(8))
		})),
		qpMust(qp.BuildMap(basic.Prototype.Any, -1, func(am datamodel.MapAssembler) {
			qp.MapEntry(am, "foo", qp.Int(3))
			qp.MapEntry(am, "baz", qp.Int(8))
		})),
		false,
	},

	// TODO: tests involving different implementations, once bindnode is ready

}

func TestDeepEqual(t *testing.T) {
	t.Parallel()
	for _, tc := range deepEqualTests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := datamodel.DeepEqual(tc.left, tc.right)
			if got != tc.want {
				t.Fatalf("DeepEqual got %v, want %v", got, tc.want)
			}
		})
	}
}
