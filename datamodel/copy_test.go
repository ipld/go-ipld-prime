package datamodel_test

import (
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	basic "github.com/ipld/go-ipld-prime/node/basicnode"
)

var copyTests = []struct {
	name string
	na   datamodel.NodeBuilder
	n    datamodel.Node
	err  string
}{
	{name: "Null / Any", na: basic.Prototype.Any.NewBuilder(), n: datamodel.Null},
	{name: "Int / Any", na: basic.Prototype.Any.NewBuilder(), n: basic.NewInt(100)},
	{name: "Int / Int", na: basic.Prototype.Int.NewBuilder(), n: basic.NewInt(1000)},
	{name: "Bool / Any", na: basic.Prototype.Any.NewBuilder(), n: basic.NewBool(true)},
	{name: "Bool / Bool", na: basic.Prototype.Bool.NewBuilder(), n: basic.NewBool(false)},
	{name: "Float / Any", na: basic.Prototype.Any.NewBuilder(), n: basic.NewFloat(1.1)},
	{name: "Float / Float", na: basic.Prototype.Float.NewBuilder(), n: basic.NewFloat(1.2)},
	{name: "String / Any", na: basic.Prototype.Any.NewBuilder(), n: basic.NewString("mary had")},
	{name: "String / String", na: basic.Prototype.String.NewBuilder(), n: basic.NewString("a little lamb")},
	{name: "Bytes / Any", na: basic.Prototype.Any.NewBuilder(), n: basic.NewBytes([]byte("mary had"))},
	{name: "Bytes / Bytes", na: basic.Prototype.Bytes.NewBuilder(), n: basic.NewBytes([]byte("a little lamb"))},
	{name: "Link / Any", na: basic.Prototype.Any.NewBuilder(), n: basic.NewLink(globalLink)},
	{name: "Link / Link", na: basic.Prototype.Link.NewBuilder(), n: basic.NewLink(globalLink2)},
	{
		name: "List / Any",
		na:   basic.Prototype.Any.NewBuilder(),
		n: qpMust(qp.BuildList(basic.Prototype.Any, -1, func(am datamodel.ListAssembler) {
			qp.ListEntry(am, qp.Int(7))
			qp.ListEntry(am, qp.Int(8))
		})),
	},
	{
		name: "List / List",
		na:   basic.Prototype.List.NewBuilder(),
		n: qpMust(qp.BuildList(basic.Prototype.List, -1, func(am datamodel.ListAssembler) {
			qp.ListEntry(am, qp.String("yep"))
			qp.ListEntry(am, qp.Int(8))
			qp.ListEntry(am, qp.String("nope"))
		})),
	},
	{
		name: "Map / Any",
		na:   basic.Prototype.Any.NewBuilder(),
		n: qpMust(qp.BuildMap(basic.Prototype.Any, -1, func(am datamodel.MapAssembler) {
			qp.MapEntry(am, "foo", qp.Int(7))
			qp.MapEntry(am, "bar", qp.Int(8))
		})),
	},
	{
		name: "Map / Map",
		na:   basic.Prototype.Map.NewBuilder(),
		n: qpMust(qp.BuildMap(basic.Prototype.Map, -1, func(am datamodel.MapAssembler) {
			qp.MapEntry(am, "foo", qp.Int(7))
			qp.MapEntry(am, "bar", qp.Int(8))
			qp.MapEntry(am, "bang", qp.Link(globalLink))
		})),
	},
	{name: "nil", na: basic.Prototype.Any.NewBuilder(), n: nil, err: "cannot copy a nil node"},
	{name: "absent", na: basic.Prototype.Any.NewBuilder(), n: datamodel.Absent, err: "copying an absent node makes no sense"},
}

func TestCopy(t *testing.T) {
	for _, tt := range copyTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := datamodel.Copy(tt.n, tt.na)
			if err != nil {
				if tt.err != "" {
					if err.Error() != tt.err {
						t.Fatalf("expected error %q, got %q", tt.err, err.Error())
					}
				} else {
					t.Fatal(err)
				}
				return
			} else if tt.err != "" {
				t.Fatalf("expected error %q, got nil", tt.err)
				return
			}
			out := tt.na.Build()
			if !datamodel.DeepEqual(tt.n, out) {
				t.Fatalf("deep equal failed")
			}
		})
	}
}
