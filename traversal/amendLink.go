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
	amendCfg

	nLink   datamodel.Link // newest link: can be `nil` if a transformation occurred, but it hasn't been recomputed
	pLink   datamodel.Link // previous link: will always have a valid value, even if not the latest
	child   Amender
	linkCtx linking.LinkContext
	linkSys linking.LinkSystem
}

func (opts AmendOptions) newLinkAmender(base datamodel.Node, parent Amender, create bool) Amender {
	// If the base node is already a link-amender, reuse the mutation state but reset `parent` and `created`.
	if amd, castOk := base.(*linkAmender); castOk {
		la := &linkAmender{amendCfg{&opts, amd.base, parent, create}, amd.nLink, amd.pLink, amd.child, amd.linkCtx, amd.linkSys}
		// Make a copy of the child amender so that it has its own mutation state
		if la.child != nil {
			child := la.child.Build()
			la.child = opts.newAmender(child, la, child.Kind(), false)
		}
		return la
	} else {
		// Start with fresh state because existing metadata could not be reused.
		link, err := base.AsLink()
		if err != nil {
			panic("misuse")
		}
		// `linkCtx` and `linkSys` can be defaulted since they're only needed for recomputing the link after a
		// transformation occurs, and such a transformation would have populated them correctly.
		return &linkAmender{amendCfg{&opts, base, parent, create}, link, link, nil, linking.LinkContext{}, linking.LinkSystem{}}
	}
}

// -- Node -->

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
	return a.computeLink()
}

func (a *linkAmender) Prototype() datamodel.NodePrototype {
	return basicnode.Prototype.Link
}

// -- Amender -->

func (a *linkAmender) Fetch(prog *Progress, path datamodel.Path, trackProgress bool) (datamodel.Node, error) {
	// Check the budget
	if prog.Budget != nil {
		if prog.Budget.LinkBudget <= 0 {
			return nil, &ErrBudgetExceeded{BudgetKind: "link", Path: prog.Path, Link: a.validLink()}
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
	return a.child.Fetch(prog, path, trackProgress)
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
			*a = *a.opts.newLinkAmender(newNode, a.parent, a.created).(*linkAmender)
			return prevNode, nil
		}
	}
	err := a.loadLink(prog, true)
	if err != nil {
		return nil, err
	}
	childVal, err := a.child.Transform(prog, path, fn, createParents)
	if err != nil {
		return nil, err
	}
	if a.opts.LazyLinkUpdate {
		// Reset the link and lazily compute it when it is needed instead of on every transformation.
		a.nLink = nil
	} else {
		newLink, err := a.linkSys.Store(a.linkCtx, a.nLink.Prototype(), a.child.Build())
		if err != nil {
			return nil, fmt.Errorf("transform: error storing transformed node at %q: %w", prog.Path, err)
		}
		a.nLink = newLink
	}
	return childVal, nil
}

func (a *linkAmender) Build() datamodel.Node {
	// `linkAmender` is also a `Node`.
	return (datamodel.Node)(a)
}

func (a *linkAmender) isCreated() bool {
	return a.created
}

// validLink will return a valid `Link`, whether the base value, an intermediate recomputed value, or the latest value.
func (a *linkAmender) validLink() datamodel.Link {
	if a.nLink == nil {
		return a.pLink
	}
	return a.nLink
}

func (a *linkAmender) computeLink() (datamodel.Link, error) {
	// `nLink` can be `nil` if lazy link computation is enabled and the child node has been transformed, but the updated
	// link has not yet been requested (and thus not recomputed).
	if a.nLink == nil {
		// We've already validated that `base` is a valid `Link` and so we don't care about a conversion error here.
		baseLink, _ := a.base.AsLink()
		lp := baseLink.Prototype()
		// `nLink` will only be `nil` if a transformation made it "dirty", indicating that it needs to be recomputed. In
		// this case, `child` will always have a valid value since it would have already been loaded/updated, so we
		// don't need to check.
		newLink, err := a.linkSys.ComputeLink(lp, a.child.Build())
		if err != nil {
			return nil, err
		}
		a.nLink = newLink
	}
	return a.nLink, nil
}

func (a *linkAmender) loadLink(prog *Progress, trackProgress bool) error {
	if a.child == nil {
		// Check the budget
		if prog.Budget != nil {
			if prog.Budget.LinkBudget <= 0 {
				return &ErrBudgetExceeded{BudgetKind: "link", Path: prog.Path, Link: a.validLink()}
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
		// `child` will only be `nil` if it was never loaded. In this case, `nLink` will always be valid, so we don't
		// need to check.
		// Pick what in-memory format we will build.
		np, err := prog.Cfg.LinkTargetNodePrototypeChooser(a.nLink, a.linkCtx)
		if err != nil {
			return fmt.Errorf("error traversing node at %q: could not load link %q: %w", prog.Path, a.nLink, err)
		}
		// Load link
		child, err := a.linkSys.Load(a.linkCtx, a.nLink, np)
		if err != nil {
			return fmt.Errorf("error traversing node at %q: could not load link %q: %w", prog.Path, a.nLink, err)
		}
		a.child = a.opts.newAmender(child, a, child.Kind(), false)
		if trackProgress {
			prog.LastBlock.Path = prog.Path
			prog.LastBlock.Link = a.nLink
		}
	}
	return nil
}
