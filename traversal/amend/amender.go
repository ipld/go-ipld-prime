package amend

import "github.com/ipld/go-ipld-prime/datamodel"

type Amender interface {
	// Get returns the node at the specified path. It will not create any intermediate nodes because this is just a
	// retrieval and not a modification operation.
	Get(path datamodel.Path) (datamodel.Node, error)

	// Add will add the specified Node at the specified path. If `createParents = true`, any missing parents will be
	// created, otherwise this function will return an error.
	Add(path datamodel.Path, value datamodel.Node, createParents bool) error

	// Remove will remove the node at the specified path and return its value. This is useful for implementing a "move"
	// operation, where a node can be "removed" and then "added" at a different path.
	Remove(path datamodel.Path) (datamodel.Node, error)

	// Replace will do an in-place replacement of the node at the specified path and return its previous value.
	Replace(path datamodel.Path, value datamodel.Node) (datamodel.Node, error)

	// Build returns a traversable node that can be used with existing codec implementations. An `Amender` does not
	// *have* to be a `Node` although currently, all `Amender` implementations are also `Node`s.
	Build() datamodel.Node
}

// NewAmender returns a new amender of the right "type" (i.e. map, list, any) using the specified base node.
func NewAmender(base datamodel.Node) Amender {
	// Do not allow externally creating a new amender without a base node to refer to. Amendment assumes that there is
	// something to amend.
	if base == nil {
		panic("misuse")
	}
	return newAmender(base, nil, base.Kind(), false)
}

func newAmender(base datamodel.Node, parent Amender, kind datamodel.Kind, create bool) Amender {
	if kind == datamodel.Kind_Map {
		return newMapAmender(base, parent, create)
	} else if kind == datamodel.Kind_List {
		return newListAmender(base, parent, create)
	} else {
		return newAnyAmender(base, parent, create)
	}
}

func isCreated(a Amender) bool {
	if ma, castOk := a.(*mapAmender); castOk {
		return ma.created
	} else if la, castOk := a.(*listAmender); castOk {
		return la.created
	} else if aa, castOk := a.(*anyAmender); castOk {
		return aa.created
	}
	panic("misuse")
}
