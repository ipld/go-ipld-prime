package traversal

import "github.com/ipld/go-ipld-prime/datamodel"

type Amender interface {
	// Get returns the node at the specified path. It will not create any intermediate nodes because this is just a
	// retrieval and not a modification operation.
	Get(prog *Progress, path datamodel.Path, trackProgress bool) (datamodel.Node, error)

	// Transform will do an in-place transformation of the node at the specified path and return its previous value.
	// If `createParents = true`, any missing parents will be created, otherwise this function will return an error.
	Transform(prog *Progress, path datamodel.Path, fn TransformFn, createParents bool) (datamodel.Node, error)

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
	} else if kind == datamodel.Kind_Link {
		return newLinkAmender(base, parent, create)
	} else {
		return newAnyAmender(base, parent, create)
	}
}

func isCreated(n datamodel.Node) bool {
	if a, isAmd := n.(Amender); isAmd {
		if ma, isMa := a.(*mapAmender); isMa {
			return ma.created
		} else if la, isLa := a.(*listAmender); isLa {
			return la.created
		} else if aa, isAa := a.(*anyAmender); isAa {
			return aa.created
		}
	}
	return false
}
