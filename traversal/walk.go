package traversal

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

func WalkMatching(n ipld.Node, s selector.Selector, fn VisitFn) error {
	return Progress{}.WalkMatching(n, s, fn)
}

func WalkAdv(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	return Progress{}.WalkAdv(n, s, fn)
}

func WalkTransforming(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	return Progress{}.WalkTransforming(n, s, fn)
}

func (prog Progress) WalkMatching(n ipld.Node, s selector.Selector, fn VisitFn) error {
	prog.init()
	return prog.walkAdv(n, s, func(prog Progress, n ipld.Node, tr VisitReason) error {
		if tr != VisitReason_SelectionMatch {
			return nil
		}
		return fn(prog, n)
	})
}

func (prog Progress) WalkAdv(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	prog.init()
	return prog.walkAdv(n, s, fn)
}

func (prog Progress) walkAdv(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	if s.Decide(n) {
		if err := fn(prog, n, VisitReason_SelectionMatch); err != nil {
			return err
		}
	} else {
		if err := fn(prog, n, VisitReason_SelectionCandidate); err != nil {
			return err
		}
	}
	nk := n.ReprKind()
	switch nk {
	case ipld.ReprKind_Map, ipld.ReprKind_List: // continue
	default:
		return nil
	}
	attn := s.Interests()
	if attn == nil {
		return prog.walkAdv_iterateAll(n, s, fn)
	}
	return prog.walkAdv_iterateSelective(n, attn, s, fn)

}

func (prog Progress) walkAdv_iterateAll(n ipld.Node, s selector.Selector, fn AdvVisitFn) error {
	for itr := selector.NewSegmentIterator(n); !itr.Done(); {
		ps, v, err := itr.Next()
		if err != nil {
			return err
		}
		sNext := s.Explore(n, ps)
		if sNext != nil {
			progNext := prog
			progNext.Path = prog.Path.AppendSegment(ps)
			if v.ReprKind() == ipld.ReprKind_Link {
				lnk, _ := v.AsLink()
				progNext.LastBlock.Path = progNext.Path
				progNext.LastBlock.Link = lnk
				v, err = progNext.loadLink(v, n)
				if err != nil {
					return err
				}
			}

			err = progNext.walkAdv(v, sNext, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (prog Progress) walkAdv_iterateSelective(n ipld.Node, attn []ipld.PathSegment, s selector.Selector, fn AdvVisitFn) error {
	for _, ps := range attn {
		v, err := n.LookupSegment(ps)
		if err != nil {
			continue
		}
		sNext := s.Explore(n, ps)
		if sNext != nil {
			progNext := prog
			progNext.Path = prog.Path.AppendSegment(ps)
			if v.ReprKind() == ipld.ReprKind_Link {
				lnk, _ := v.AsLink()
				progNext.LastBlock.Path = progNext.Path
				progNext.LastBlock.Link = lnk
				v, err = progNext.loadLink(v, n)
				if err != nil {
					return err
				}
			}

			err = progNext.walkAdv(v, sNext, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (prog Progress) loadLink(v ipld.Node, parent ipld.Node) (ipld.Node, error) {
	lnk, err := v.AsLink()
	if err != nil {
		return nil, err
	}
	// Assemble the LinkContext in case the Loader or NBChooser want it.
	lnkCtx := ipld.LinkContext{
		LinkPath:   prog.Path,
		LinkNode:   v,
		ParentNode: parent,
	}
	// Load link!
	v, err = lnk.Load(
		prog.Cfg.Ctx,
		lnkCtx,
		prog.Cfg.LinkNodeBuilderChooser(lnk, lnkCtx),
		prog.Cfg.LinkLoader,
	)
	if err != nil {
		return nil, fmt.Errorf("error traversing node at %q: could not load link %q: %s", prog.Path, lnk, err)
	}
	return v, nil
}

func (prog Progress) WalkTransforming(n ipld.Node, s selector.Selector, fn TransformFn) (ipld.Node, error) {
	panic("TODO")
}
