package amend

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/traversal/patch"
)

func Eval(n datamodel.Node, ops []patch.Operation) (datamodel.Node, error) {
	var err error
	a := NewAmender(n) // One Amender To Patch Them All
	for _, op := range ops {
		err = EvalOne(a, op)
		if err != nil {
			return nil, err
		}
	}
	return a.Build(), nil
}

func EvalOne(a Amender, op patch.Operation) error {
	switch op.Op {
	case patch.Op_Add:
		return a.Add(op.Path, op.Value, true)
	case patch.Op_Remove:
		_, err := a.Remove(op.Path)
		return err
	case patch.Op_Replace:
		_, err := a.Replace(op.Path, op.Value)
		return err
	case patch.Op_Move:
		source, err := a.Remove(op.From)
		if err != nil {
			return err
		}
		// Similar to `replace` with the difference that the destination path might not exist and need to be created.
		return a.Add(op.Path, source, true)
	case patch.Op_Copy:
		source, err := a.Get(op.From)
		if err != nil {
			return err
		}
		return a.Add(op.Path, source, false)
	case patch.Op_Test:
		point, err := a.Get(op.Path)
		if err != nil {
			return err
		}
		if datamodel.DeepEqual(point, op.Value) {
			return nil
		}
		return fmt.Errorf("test failed") // TODO real error handling and a code
	default:
		return fmt.Errorf("misuse: invalid operation") // TODO real error handling and a code
	}
}
