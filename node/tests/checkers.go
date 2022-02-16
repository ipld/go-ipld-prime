package tests

import (
	"errors"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/printer"
)

// NodeContentEquals checks whether two nodes have equal content by first encoding them via
// printer.Sprint, then checking that the generated encodings are identical.
//
// See: printer.Sprint.
var NodeContentEquals = &nodeContentEqualsChecker{}

type nodeContentEqualsChecker struct{}

func (n *nodeContentEqualsChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) error {
	want := args[0]
	if want == nil {
		return qt.IsNil.Check(got, args, note)
	}
	if got == nil {
		return qt.IsNotNil.Check(got, args, note)
	}
	wantNode, ok := want.(datamodel.Node)
	if !ok {
		return errors.New("this checker only supports checking datamodel.Node values")
	}
	wantPrint := printer.Sprint(wantNode)

	gotNode, ok := got.(datamodel.Node)
	if !ok {
		return errors.New("this checker only supports checking datamodel.Node values")
	}
	gotPrint := printer.Sprint(gotNode)
	return qt.Equals.Check(gotPrint, []interface{}{wantPrint}, note)
}

func (n *nodeContentEqualsChecker) ArgNames() []string {
	return []string{"got", "want node"}
}
