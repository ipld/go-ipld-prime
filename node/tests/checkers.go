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
// Use DeepNodeContentsEquals if you want a less strict comparison that does not
// require map keys to be in the same order.
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

// DeepNodeContentsEquals checks whether two nodes have equal content by
// walking the nodes recursively and comparing their contents. This is similar
// to datamodel.DeepEquals except that map keys DO NOT need to be strictly
// the same order to be considered equal.
//
// Use NodeContentEquals if you want a more strict comparison.
var DeepNodeContentsEquals = &deepNodeContentsEqualsChecker{}

type deepNodeContentsEqualsChecker struct{}

func (n *deepNodeContentsEqualsChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) error {
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

	gotNode, ok := got.(datamodel.Node)
	if !ok {
		return errors.New("this checker only supports checking datamodel.Node values")
	}

	return deepEqualCheck(gotNode, wantNode, note)
}

func (n *deepNodeContentsEqualsChecker) ArgNames() []string {
	return []string{"got", "want node"}
}

func deepEqualCheck(x, y datamodel.Node, note func(key string, value interface{})) error {
	if x == nil || y == nil {
		if err := qt.Equals.Check(x, []interface{}{y}, note); err != nil {
			return err
		}
	}
	xk, yk := x.Kind(), y.Kind()
	if err := qt.Equals.Check(xk, []interface{}{yk}, note); err != nil {
		return err
	}

	switch xk {

	// Scalar kinds.
	case datamodel.Kind_Null:
		return qt.Equals.Check(x.IsNull(), []interface{}{y.IsNull()}, note)
	case datamodel.Kind_Bool:
		xv, err := x.AsBool()
		if err != nil {
			panic(err)
		}
		yv, err := y.AsBool()
		if err != nil {
			panic(err)
		}
		return qt.Equals.Check(xv, []interface{}{yv}, note)
	case datamodel.Kind_Int:
		xv, err := x.AsInt()
		if err != nil {
			panic(err)
		}
		yv, err := y.AsInt()
		if err != nil {
			panic(err)
		}
		return qt.Equals.Check(xv, []interface{}{yv}, note)
	case datamodel.Kind_Float:
		xv, err := x.AsFloat()
		if err != nil {
			panic(err)
		}
		yv, err := y.AsFloat()
		if err != nil {
			panic(err)
		}
		return qt.Equals.Check(xv, []interface{}{yv}, note)
	case datamodel.Kind_String:
		xv, err := x.AsString()
		if err != nil {
			panic(err)
		}
		yv, err := y.AsString()
		if err != nil {
			panic(err)
		}
		return qt.Equals.Check(xv, []interface{}{yv}, note)
	case datamodel.Kind_Bytes:
		xv, err := x.AsBytes()
		if err != nil {
			panic(err)
		}
		yv, err := y.AsBytes()
		if err != nil {
			panic(err)
		}
		return qt.Equals.Check(string(xv), []interface{}{string(yv)}, note)
	case datamodel.Kind_Link:
		xv, err := x.AsLink()
		if err != nil {
			panic(err)
		}
		yv, err := y.AsLink()
		if err != nil {
			panic(err)
		}
		// Links are just compared via ==.
		// This requires the types to exactly match,
		// and the values to be equal as per == too.
		// This will generally work,
		// as ipld-prime assumes link types to be consistent.
		return qt.Equals.Check(xv, []interface{}{yv}, note)

	// Recursive kinds.
	case datamodel.Kind_Map:
		if err := qt.Equals.Check(x.Length(), []interface{}{y.Length()}, note); err != nil {
			return err
		}
		ykeys := make([]string, y.Length())
		yitr := y.MapIterator()
		for !yitr.Done() {
			ykey, _, err := yitr.Next()
			if err != nil {
				panic(err)
			}
			yk, err := ykey.AsString()
			if err != nil {
				panic(err)
			}
			ykeys = append(ykeys, yk)
		}
		xitr := x.MapIterator()
		for !xitr.Done() {
			xkey, xval, err := xitr.Next()
			if err != nil {
				panic(err)
			}
			xk, err := xkey.AsString()
			if err != nil {
				panic(err)
			}
			if err := qt.Contains.Check(ykeys, []interface{}{xk}, note); err != nil {
				return err
			}
			yval, err := y.LookupByNode(xkey)
			if err != nil {
				panic(err)
			}
			if err := deepEqualCheck(xval, yval, note); err != nil {
				return err
			}
		}
		return nil
	case datamodel.Kind_List:
		if err := qt.Equals.Check(x.Length(), []interface{}{y.Length()}, note); err != nil {
			return err
		}
		xitr := x.ListIterator()
		yitr := y.ListIterator()
		for !xitr.Done() && !yitr.Done() {
			_, xval, err := xitr.Next()
			if err != nil {
				panic(err)
			}
			_, yval, err := yitr.Next()
			if err != nil {
				panic(err)
			}
			if err := deepEqualCheck(xval, yval, note); err != nil {
				return err
			}
		}
		return nil

	// As per the docs, other kinds such as Invalid are not deeply equal.
	default:
		panic("bad kind")
	}
}
