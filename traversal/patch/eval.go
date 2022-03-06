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
		// The behavior of the 'add' op in jsonpatch varies based on if the parent of the target path is a list.
		// If the parent of the target path is a list, then 'add' is really more of an 'insert': it should slide the rest of the values down.
		// There's also a special case for "-", which means "append to the end of the list".
		// Otherwise, if the destination path exists, it's an error.  (No upserting.)
		// Handling this requires looking at the parent of the destination node, so we split this into *two* traversal.FocusedTransform calls.
		return traversal.FocusedTransform(n, op.Path.Pop(), func(prog traversal.Progress, parent datamodel.Node) (datamodel.Node, error) {
			if parent.Kind() == datamodel.Kind_List {
				seg := op.Path.Last()
				var idx int64
				if seg.String() == "-" {
					idx = -1
				}
				var err error
				idx, err = seg.Index()
				if err != nil {
					return nil, fmt.Errorf("patch-invalid-path-through-list: at %q", op.Path) // TODO error structuralization and review the code
				}

				nb := parent.Prototype().NewBuilder()
				la, err := nb.BeginList(parent.Length() + 1)
				if err != nil {
					return nil, err
				}
				for itr := n.ListIterator(); !itr.Done(); {
					i, v, err := itr.Next()
					if err != nil {
						return nil, err
					}
					if idx == i {
						la.AssembleValue().AssignNode(op.Value)
					}
					if err := la.AssembleValue().AssignNode(v); err != nil {
						return nil, err
					}
				}
				// TODO: is one-past-the-end supposed to be supported or supposed to be ruled out?
				if idx == -1 {
					la.AssembleValue().AssignNode(op.Value)
				}
				if err := la.Finish(); err != nil {
					return nil, err
				}
				return nb.Build(), nil
			}
			return prog.FocusedTransform(parent, datamodel.NewPath([]datamodel.PathSegment{op.Path.Last()}), func(prog traversal.Progress, point datamodel.Node) (datamodel.Node, error) {
				if point != nil && !point.IsAbsent() {
					return nil, fmt.Errorf("patch-target-exists: at %q", op.Path) // TODO error structuralization and review the code
				}
				return op.Value, nil
			}, false)
		}, false)
	case "remove":
		return traversal.FocusedTransform(n, op.Path, func(_ traversal.Progress, point datamodel.Node) (datamodel.Node, error) {
			return nil, nil // Returning a nil value here means "remove what's here".
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
