package traversal

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func Traverse(n ipld.Node, s selector.Selector, fn VisitFn) error {
	return TraversalProgress{}.Traverse(n, s, fn)
}

func TraverseInformatively(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	return TraversalProgress{}.TraverseInformatively(n, s, fn)
}

func TraverseTransform(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	return TraversalProgress{}.TraverseTransform(n, s, fn)
}

func (tp TraversalProgress) Traverse(n ipld.Node, s selector.Selector, fn VisitFn) error {
	tp.init()
	return tp.TraverseInformatively(n, s, func(tp TraversalProgress, n ipld.Node, tr TraversalReason) error {
		if tr != 1 {
			return nil
		}
		return fn(tp, n)
	})
}

func (tp TraversalProgress) TraverseInformatively(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	panic("TODO")
}

func (tp TraversalProgress) TraverseTransform(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	panic("TODO")
}
