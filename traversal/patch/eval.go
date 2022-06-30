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
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
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
	a := traversal.NewAmender(n) // One Amender To Patch Them All
	prog := traversal.Progress{}
	for _, op := range ops {
		_, err := evalOne(&prog, a.Build(), op)
		if err != nil {
			return nil, err
		}
	}
	return a.Build(), nil
}

func EvalOne(n datamodel.Node, op Operation) (datamodel.Node, error) {
	return evalOne(&traversal.Progress{}, n, op)
}

func evalOne(prog *traversal.Progress, n datamodel.Node, op Operation) (datamodel.Node, error) {
	// If the node being modified is already an `Amender` reuse it, otherwise create a fresh one.
	var a traversal.Amender
	if amd, castOk := n.(traversal.Amender); castOk {
		a = amd
	} else {
		a = traversal.NewAmender(n)
	}
	switch op.Op {
	case Op_Add:
		// The behavior of the 'add' op in jsonpatch varies based on if the parent of the target path is a list.
		// If the parent of the target path is a list, then 'add' is really more of an 'insert': it should slide the
		// rest of the values down. There's also a special case for "-", which means "append to the end of the list".
		// Otherwise, if the destination path exists, it's an error.  (No upserting.)
		if _, err := a.Transform(prog, op.Path, func(progress traversal.Progress, prev datamodel.Node) (datamodel.Node, error) {
			if n.Kind() == datamodel.Kind_List {
				// Since jsonpatch expects list "add" operations to insert the element, return the transformed list
				// "[previous node, new node]" so that the transformation can expand this list at the specified index in
				// the original list. This allows jsonpatch to continue inserting elements for "add" operations, while
				// also allowing transformations that update list elements in place (default behavior), currently used
				// by `FocusedTransform`.
				return qp.BuildList(basicnode.Prototype.Any, 2, func(la datamodel.ListAssembler) {
					qp.ListEntry(la, qp.Node(op.Value))
					qp.ListEntry(la, qp.Node(prev))
				})
			}
			return op.Value, nil
		}, true); err != nil {
			return nil, err
		}
	case Op_Remove:
		if _, err := a.Transform(prog, op.Path, func(progress traversal.Progress, node datamodel.Node) (datamodel.Node, error) {
			return nil, nil
		}, false); err != nil {
			return nil, err
		}
	case Op_Replace:
		if _, err := a.Transform(prog, op.Path, func(progress traversal.Progress, node datamodel.Node) (datamodel.Node, error) {
			return op.Value, nil
		}, false); err != nil {
			return nil, err
		}
	case Op_Move:
		if source, err := a.Transform(prog, op.From, func(progress traversal.Progress, node datamodel.Node) (datamodel.Node, error) {
			// Returning `nil` will cause the target node to be deleted.
			return nil, nil
		}, false); err != nil {
			return nil, err
		} else if _, err := a.Transform(prog, op.Path, func(progress traversal.Progress, node datamodel.Node) (datamodel.Node, error) {
			return source, nil
		}, true); err != nil {
			return nil, err
		}
	case Op_Copy:
		if source, err := a.Get(prog, op.From, true); err != nil {
			return nil, err
		} else if _, err := a.Transform(prog, op.Path, func(progress traversal.Progress, node datamodel.Node) (datamodel.Node, error) {
			return source, nil
		}, true); err != nil {
			return nil, err
		}
	case Op_Test:
		if point, err := a.Get(prog, op.Path, true); err != nil {
			return nil, err
		} else if !datamodel.DeepEqual(point, op.Value) {
			return nil, fmt.Errorf("test failed") // TODO real error handling and a code
		}
	default:
		return nil, fmt.Errorf("misuse: invalid operation") // TODO real error handling and a code
	}
	return a.Build(), nil
}
