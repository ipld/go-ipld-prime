package traversal

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

var (
	_ datamodel.Node = &linkAmender{}
	_ Amender        = &linkAmender{}
)

type linkAmender struct {
	base    datamodel.Node
	parent  Amender
	created bool
	link    datamodel.Link
	linkCtx linking.LinkContext
	linkSys linking.LinkSystem
	child   Amender
}

func newLinkAmender(base datamodel.Node, parent Amender, create bool) Amender {
	// If the base node is already a link-amender, reuse the metadata that encapsulates all accumulated modifications
	// but reset `parent` and `created`.
	if amd, castOk := base.(*linkAmender); castOk {
		return &linkAmender{amd.base, parent, create, amd.link, amd.linkCtx, amd.linkSys, amd.child}
	} else {
		// Start with fresh state because existing metadata could not be reused.
		link, _ := base.AsLink()
		return &linkAmender{base, parent, create, link, linking.LinkContext{}, linking.LinkSystem{}, nil}
	}
}

func (a *linkAmender) Build() datamodel.Node {
	// `linkAmender` is also a `Node`.
	return (datamodel.Node)(a)
}

func (a *linkAmender) Kind() datamodel.Kind {
	return datamodel.Kind_Link
}

func (a *linkAmender) LookupByString(key string) (datamodel.Node, error) {
	return mixins.Link{TypeName: "linkAmender"}.LookupByString(key)
}

func (a *linkAmender) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	return mixins.Link{TypeName: "linkAmender"}.LookupByNode(key)
}

func (a *linkAmender) LookupByIndex(idx int64) (datamodel.Node, error) {
	return mixins.Link{TypeName: "linkAmender"}.LookupByIndex(idx)
}

func (a *linkAmender) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	return mixins.Link{TypeName: "link"}.LookupBySegment(seg)
}

func (a *linkAmender) MapIterator() datamodel.MapIterator {
	return nil
}

func (a *linkAmender) ListIterator() datamodel.ListIterator {
	return nil
}

func (a *linkAmender) Length() int64 {
	return -1
}

func (a *linkAmender) IsAbsent() bool {
	return false
}

func (a *linkAmender) IsNull() bool {
	return false
}

func (a *linkAmender) AsBool() (bool, error) {
	return mixins.Link{TypeName: "linkAmender"}.AsBool()
}

func (a *linkAmender) AsInt() (int64, error) {
	return mixins.Link{TypeName: "linkAmender"}.AsInt()
}

func (a *linkAmender) AsFloat() (float64, error) {
	return mixins.Link{TypeName: "linkAmender"}.AsFloat()
}

func (a *linkAmender) AsString() (string, error) {
	return mixins.Link{TypeName: "linkAmender"}.AsString()
}

func (a *linkAmender) AsBytes() ([]byte, error) {
	return mixins.Link{TypeName: "linkAmender"}.AsBytes()
}

func (a *linkAmender) AsLink() (datamodel.Link, error) {
	return a.link, nil
}

func (a *linkAmender) Prototype() datamodel.NodePrototype {
	return basicnode.Prototype.Link
}

func (a *linkAmender) Get(prog *Progress, path datamodel.Path, trackProgress bool) (datamodel.Node, error) {
	// Check the budget
	if prog.Budget != nil {
		if prog.Budget.LinkBudget <= 0 {
			return nil, &ErrBudgetExceeded{BudgetKind: "link", Path: prog.Path, Link: a.link}
		}
		prog.Budget.LinkBudget--
	}
	err := a.loadLink(prog, trackProgress)
	if err != nil {
		return nil, err
	}
	if path.Len() == 0 {
		return a.Build(), nil
	}
	return a.child.Get(prog, path, trackProgress)
}

func (a *linkAmender) Transform(prog *Progress, path datamodel.Path, fn TransformFn, createParents bool) (datamodel.Node, error) {
	// Allow the base node to be replaced.
	if path.Len() == 0 {
		prevNode := a.Build()
		if newNode, err := fn(*prog, prevNode); err != nil {
			return nil, err
		} else if newNode.Kind() != datamodel.Kind_Link {
			return nil, fmt.Errorf("transform: cannot transform root into incompatible type: %q", newNode.Kind())
		} else {
			// Go through `newLinkAmender` in case `newNode` is already a link-amender.
			newAmd := newLinkAmender(newNode, a.parent, a.created).(*linkAmender)
			// Reset the current amender to use the transformed node.
			a.base = newAmd.base
			a.link = newAmd.link
			a.linkCtx = newAmd.linkCtx
			a.linkSys = newAmd.linkSys
			a.child = newAmd.child
			return prevNode, nil
		}
	}
	err := a.loadLink(prog, true)
	if err != nil {
		return nil, err
	}
	if _, err = a.child.Transform(prog, path, fn, createParents); err != nil {
		return nil, err
	}
	prevNode := a.child.Build()
	// TODO: Is it possible to store lazily, and not everytime a child is modified? Perhaps this can be configured?
	newLink, err := a.linkSys.Store(a.linkCtx, a.link.Prototype(), prevNode)
	if err != nil {
		return nil, fmt.Errorf("transform: error storing transformed node at %q: %w", prog.Path, err)
	}
	a.link = newLink
	return prevNode, nil
}

func (a *linkAmender) loadLink(prog *Progress, trackProgress bool) error {
	if a.child == nil {
		// Check the budget
		if prog.Budget != nil {
			if prog.Budget.LinkBudget <= 0 {
				return &ErrBudgetExceeded{BudgetKind: "link", Path: prog.Path, Link: a.link}
			}
			prog.Budget.LinkBudget--
		}
		// Put together the context info we'll offer to the loader and prototypeChooser.
		a.linkCtx = linking.LinkContext{
			Ctx:        prog.Cfg.Ctx,
			LinkPath:   prog.Path,
			LinkNode:   a.Build(),
			ParentNode: a.parent.Build(),
		}
		a.linkSys = prog.Cfg.LinkSystem
		// Pick what in-memory format we will build.
		np, err := prog.Cfg.LinkTargetNodePrototypeChooser(a.link, a.linkCtx)
		if err != nil {
			return fmt.Errorf("error traversing node at %q: could not load link %q: %w", prog.Path, a.link, err)
		}
		// Load link
		child, err := a.linkSys.Load(a.linkCtx, a.link, np)
		if err != nil {
			return fmt.Errorf("error traversing node at %q: could not load link %q: %w", prog.Path, a.link, err)
		}
		// A node loaded from a link can never be considered "created".
		a.child = newAmender(child, a, child.Kind(), false)
		if trackProgress {
			prog.LastBlock.Path = prog.Path
			prog.LastBlock.Link = a.link
		}
	}
	return nil
}
