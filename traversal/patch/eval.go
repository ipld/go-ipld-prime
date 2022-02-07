package patch

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/traversal"
)

type Op string

const (
	Op_Add     = "add"
	Op_Remove  = "remove"
	Op_Replace = "replace"
	Op_Move    = "move"
	Op_Copy    = "copy"
	Op_Test    = "test"
)

type Operation struct {
	Op    Op             // Always required.
	Path  datamodel.Path // Always required.
	Value datamodel.Node // Present on 'add', 'replace', 'test'.
	From  datamodel.Path // Present on 'move', 'copy'.
}

func Eval(n datamodel.Node, ops []Operation) (datamodel.Node, error) {
	var err error
	for _, op := range ops {
		n, err = EvalOne(n, op)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}

func EvalOne(n datamodel.Node, op Operation) (datamodel.Node, error) {
	switch op.Op {
	case "add":
		return traversal.FocusedTransform(n, op.Path, func(_ traversal.Progress, point datamodel.Node) (datamodel.Node, error) {
			return op.Value, nil // is this right?  what does FocusedTransform do re upsert?
		}, false)
	case "remove":
		return traversal.FocusedTransform(n, op.Path, func(_ traversal.Progress, point datamodel.Node) (datamodel.Node, error) {
			return nil, nil
		}, false)
	case "replace":
		// TODO i think you need a check that it's not landing under itself here
		return traversal.FocusedTransform(n, op.Path, func(_ traversal.Progress, point datamodel.Node) (datamodel.Node, error) {
			return op.Value, nil // is this right?  what does FocusedTransform do re upsert?
		}, false)
	case "move":
		// TODO i think you need a check that it's not landing under itself here
		source, err := traversal.Get(n, op.From)
		if err != nil {
			return nil, err
		}
		return traversal.FocusedTransform(n, op.Path, func(_ traversal.Progress, point datamodel.Node) (datamodel.Node, error) {
			return source, nil // is this right?  what does FocusedTransform do re upsert?
		}, false)
	case "copy":
		// TODO i think you need a check that it's not landing under itself here
		source, err := traversal.Get(n, op.From)
		if err != nil {
			return nil, err
		}
		return traversal.FocusedTransform(n, op.Path, func(_ traversal.Progress, point datamodel.Node) (datamodel.Node, error) {
			return source, nil // is this right?  what does FocusedTransform do re upsert?
		}, false)
	case "test":
		point, err := traversal.Get(n, op.Path)
		if err != nil {
			return nil, err
		}
		if datamodel.DeepEqual(point, op.Value) {
			return n, nil
		}
		return n, fmt.Errorf("test failed") // TODO real error handling and a code
	default:
		return nil, fmt.Errorf("misuse: invalid operation") // TODO real error handling and a code
	}
}
