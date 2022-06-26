// Package patch provides an implementation of the IPLD Patch specification.
// IPLD Patch is a system for declaratively specifying patches to a document,
// which can then be applied to produce a new, modified document.
//
//
// This package is EXPERIMENTAL; its behavior and API might change as it's still
// in development.
package patch

import (
	"fmt"
	"github.com/ipld/go-ipld-prime/datamodel"
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
	a := NewAmender(n) // One Amender To Patch Them All
	for _, op := range ops {
		_, err := EvalOne(a.(datamodel.Node), op)
		if err != nil {
			return nil, err
		}
	}
	return a.Build(), nil
}

func EvalOne(n datamodel.Node, op Operation) (datamodel.Node, error) {
	// If the `Node` being modified is already an `Amender` reuse it, otherwise create a fresh one.
	var a Amender
	if amd, castOk := n.(Amender); castOk {
		a = amd
	} else {
		a = NewAmender(n)
	}
	an := a.(datamodel.Node)
	switch op.Op {
	case Op_Add:
		// The behavior of the 'add' op in jsonpatch varies based on if the parent of the target path is a list.
		// If the parent of the target path is a list, then 'add' is really more of an 'insert': it should slide the rest of the values down.
		// There's also a special case for "-", which means "append to the end of the list". (TODO: implement this)
		// Otherwise, if the destination path exists, it's an error.  (No upserting.)
		return an, a.Add(op.Path, op.Value, true)
	case Op_Remove:
		_, err := a.Remove(op.Path)
		return an, err
	case Op_Replace:
		_, err := a.Replace(op.Path, op.Value)
		return an, err
	case Op_Move:
		source, err := a.Remove(op.From)
		if err != nil {
			return nil, err
		}
		// Similar to `replace` with the difference that the destination path might not exist and need to be created.
		return an, a.Add(op.Path, source, true)
	case Op_Copy:
		source, err := a.Get(op.From)
		if err != nil {
			return nil, err
		}
		return an, a.Add(op.Path, source, false)
	case Op_Test:
		point, err := a.Get(op.Path)
		if err != nil {
			return nil, err
		}
		if datamodel.DeepEqual(point, op.Value) {
			return an, nil
		}
		return nil, fmt.Errorf("test failed") // TODO real error handling and a code
	default:
		return nil, fmt.Errorf("misuse: invalid operation") // TODO real error handling and a code
	}
}
