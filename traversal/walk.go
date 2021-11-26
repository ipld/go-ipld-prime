package traversal

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	"github.com/ipld/go-ipld-prime/traversal/selector"
)

// WalkLocal walks a tree of Nodes, visiting each of them,
// and calling the given VisitFn on all of them;
// it does not traverse any links.
//
// WalkLocal can skip subtrees if the VisitFn returns SkipMe,
// but lacks any other options for controlling or directing the visit;
// consider using some of the various Walk functions with Selector parameters if you want more control.
func WalkLocal(n datamodel.Node, fn VisitFn) error {
	return Progress{}.WalkLocal(n, fn)
}

// WalkMatching walks a graph of Nodes, deciding which to visit by applying a Selector,
// and calling the given VisitFn on those that the Selector deems a match.
//
// This function is a helper function which starts a new walk with default configuration.
// It cannot cross links automatically (since this requires configuration).
// Use the equivalent WalkMatching function on the Progress structure
// for more advanced and configurable walks.
func WalkMatching(n datamodel.Node, s selector.Selector, fn VisitFn) error {
	return Progress{}.WalkMatching(n, s, fn)
}

// WalkAdv is identical to WalkMatching, except it is called for *all* nodes
// visited (not just matching nodes), together with the reason for the visit.
// An AdvVisitFn is used instead of a VisitFn, so that the reason can be provided.
//
// This function is a helper function which starts a new walk with default configuration.
// It cannot cross links automatically (since this requires configuration).
// Use the equivalent WalkAdv function on the Progress structure
// for more advanced and configurable walks.
func WalkAdv(n datamodel.Node, s selector.Selector, fn AdvVisitFn) error {
	return Progress{}.WalkAdv(n, s, fn)
}

// WalkTransforming walks a graph of Nodes, deciding which to alter by applying a Selector,
// and calls the given TransformFn to decide what new node to replace the visited node with.
// A new Node tree will be returned (the original is unchanged).
//
// This function is a helper function which starts a new walk with default configuration.
// It cannot cross links automatically (since this requires configuration).
// Use the equivalent WalkTransforming function on the Progress structure
// for more advanced and configurable walks.
func WalkTransforming(n datamodel.Node, s selector.Selector, fn TransformFn) (datamodel.Node, error) {
	return Progress{}.WalkTransforming(n, s, fn)
}

// WalkMatching walks a graph of Nodes, deciding which to visit by applying a Selector,
// and calling the given VisitFn on those that the Selector deems a match.
//
// WalkMatching is a read-only traversal.
// See WalkTransforming if looking for a way to do "updates" to a tree of nodes.
//
// Provide configuration to this process using the Config field in the Progress object.
//
// This walk will automatically cross links, but requires some configuration
// with link loading functions to do so.
//
// Traversals are defined as visiting a (node,path) tuple.
// This is important to note because when walking DAGs with Links,
// it means you may visit the same node multiple times
// due to having reached it via a different path.
// (You can prevent this by using a LinkLoader function which memoizes a set of
// already-visited Links, and returns a SkipMe when encountering them again.)
//
// WalkMatching (and the other traversal functions) can be used again again inside the VisitFn!
// By using the traversal.Progress handed to the VisitFn,
// the Path recorded of the traversal so far will continue to be extended,
// and thus continued nested uses of Walk and Focus will see the fully contextualized Path.
//
func (prog Progress) WalkMatching(n datamodel.Node, s selector.Selector, fn VisitFn) error {
	prog.init()
	return prog.walkAdv(n, s, func(prog Progress, n datamodel.Node, tr VisitReason) error {
		if tr != VisitReason_SelectionMatch {
			return nil
		}
		return fn(prog, n)
	})
}

// WalkLocal is the same as the package-scope function of the same name,
// but considers an existing Progress state (and any config it might reference).
func (prog Progress) WalkLocal(n datamodel.Node, fn VisitFn) error {
	// Check the budget!
	if prog.Budget != nil {
		if prog.Budget.NodeBudget <= 0 {
			return &ErrBudgetExceeded{BudgetKind: "node", Path: prog.Path}
		}
		prog.Budget.NodeBudget--
	}
	// Visit the current node.
	if err := fn(prog, n); err != nil {
		if _, ok := err.(SkipMe); ok {
			return nil
		}
		return err
	}
	// Recurse on nodes with a recursive kind; otherwise just return.
	switch n.Kind() {
	case datamodel.Kind_Map:
		for itr := n.MapIterator(); !itr.Done(); {
			k, v, err := itr.Next()
			if err != nil {
				return err
			}
			ks, _ := k.AsString()
			progNext := prog
			progNext.Path = prog.Path.AppendSegmentString(ks)
			if err := progNext.WalkLocal(v, fn); err != nil {
				return err
			}
		}
		return nil
	case datamodel.Kind_List:
		for itr := n.ListIterator(); !itr.Done(); {
			idx, v, err := itr.Next()
			if err != nil {
				return err
			}
			progNext := prog
			progNext.Path = prog.Path.AppendSegmentInt(idx)
			if err := progNext.WalkLocal(v, fn); err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}

// WalkAdv is identical to WalkMatching, except it is called for *all* nodes
// visited (not just matching nodes), together with the reason for the visit.
// An AdvVisitFn is used instead of a VisitFn, so that the reason can be provided.
//
func (prog Progress) WalkAdv(n datamodel.Node, s selector.Selector, fn AdvVisitFn) error {
	prog.init()
	return prog.walkAdv(n, s, fn)
}

func (prog Progress) walkAdv(n datamodel.Node, s selector.Selector, fn AdvVisitFn) error {
	// Check the budget!
	if prog.Budget != nil {
		if prog.Budget.NodeBudget <= 0 {
			return &ErrBudgetExceeded{BudgetKind: "node", Path: prog.Path}
		}
		prog.Budget.NodeBudget--
	}
	// Decide if this node is matched -- do callbacks as appropriate.
	if s.Decide(n) {
		if err := fn(prog, n, VisitReason_SelectionMatch); err != nil {
			return err
		}
	} else {
		if err := fn(prog, n, VisitReason_SelectionCandidate); err != nil {
			return err
		}
	}
	// If we're handling scalars (e.g. not maps and lists) we can return now.
	nk := n.Kind()
	switch nk {
	case datamodel.Kind_Map, datamodel.Kind_List: // continue
	default:
		return nil
	}
	// For maps and lists: recurse (in one of two ways, depending on if the selector also states specific interests).
	attn := s.Interests()
	if attn == nil {
		return prog.walkAdv_iterateAll(n, s, fn)
	}
	return prog.walkAdv_iterateSelective(n, attn, s, fn)

}

func (prog Progress) walkAdv_iterateAll(n datamodel.Node, s selector.Selector, fn AdvVisitFn) error {
	for itr := selector.NewSegmentIterator(n); !itr.Done(); {
		ps, v, err := itr.Next()
		if err != nil {
			return err
		}
		sNext, err := s.Explore(n, ps)
		if err != nil {
			return err
		}
		if sNext != nil {
			progNext := prog
			progNext.Path = prog.Path.AppendSegment(ps)
			if v.Kind() == datamodel.Kind_Link {
				lnk, _ := v.AsLink()
				if prog.Cfg.LinkVisitOnlyOnce {
					if _, seen := prog.SeenLinks[lnk]; seen {
						continue
					}
					prog.SeenLinks[lnk] = struct{}{}
				}
				progNext.LastBlock.Path = progNext.Path
				progNext.LastBlock.Link = lnk
				v, err = progNext.loadLink(v, n)
				if err != nil {
					if _, ok := err.(SkipMe); ok {
						continue
					}
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

func (prog Progress) walkAdv_iterateSelective(n datamodel.Node, attn []datamodel.PathSegment, s selector.Selector, fn AdvVisitFn) error {
	for _, ps := range attn {
		v, err := n.LookupBySegment(ps)
		if err != nil {
			continue
		}
		sNext, err := s.Explore(n, ps)
		if err != nil {
			return err
		}
		if sNext != nil {
			progNext := prog
			progNext.Path = prog.Path.AppendSegment(ps)
			if v.Kind() == datamodel.Kind_Link {
				lnk, _ := v.AsLink()
				if prog.Cfg.LinkVisitOnlyOnce {
					if _, seen := prog.SeenLinks[lnk]; seen {
						continue
					}
					prog.SeenLinks[lnk] = struct{}{}
				}
				progNext.LastBlock.Path = progNext.Path
				progNext.LastBlock.Link = lnk
				v, err = progNext.loadLink(v, n)
				if err != nil {
					if _, ok := err.(SkipMe); ok {
						continue
					}
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

func (prog Progress) loadLink(v datamodel.Node, parent datamodel.Node) (datamodel.Node, error) {
	lnk, err := v.AsLink()
	if err != nil {
		return nil, err
	}
	// Check the budget!
	if prog.Budget != nil {
		if prog.Budget.LinkBudget <= 0 {
			return nil, &ErrBudgetExceeded{BudgetKind: "link", Path: prog.Path, Link: lnk}
		}
		prog.Budget.LinkBudget--
	}
	// Put together the context info we'll offer to the loader and prototypeChooser.
	lnkCtx := linking.LinkContext{
		Ctx:        prog.Cfg.Ctx,
		LinkPath:   prog.Path,
		LinkNode:   v,
		ParentNode: parent,
	}
	// Pick what in-memory format we will build.
	np, err := prog.Cfg.LinkTargetNodePrototypeChooser(lnk, lnkCtx)
	if err != nil {
		return nil, fmt.Errorf("error traversing node at %q: could not load link %q: %w", prog.Path, lnk, err)
	}
	// Load link!
	n, err := prog.Cfg.LinkSystem.Load(lnkCtx, lnk, np)
	if err != nil {
		if _, ok := err.(SkipMe); ok {
			return nil, err
		}
		return nil, fmt.Errorf("error traversing node at %q: could not load link %q: %w", prog.Path, lnk, err)
	}
	return n, nil
}

// WalkTransforming walks a graph of Nodes, deciding which to alter by applying a Selector,
// and calls the given TransformFn to decide what new node to replace the visited node with.
// A new Node tree will be returned (the original is unchanged).
//
// If the TransformFn returns the same Node which it was called with,
// then the transform is a no-op; if every visited node is a no-op,
// then the root node returned from the walk as a whole will also be
// the same as its starting Node (no new memory will be used).
//
// When a Node is replaced, no further recursion of this walk will occur on its contents.
// (You can certainly do a additional traversals, including transforms,
// from inside the TransformFn while building the replacement node.)
//
// The prototype (that is, implementation) of Node returned will be the same as the
// prototype of the Nodes at the same positions in the existing tree
// (literally, builders used to construct any new needed intermediate nodes
// are chosen by asking the existing nodes about their prototype).
//
// This feature is not yet implemented.
func (prog Progress) WalkTransforming(n datamodel.Node, s selector.Selector, fn TransformFn) (datamodel.Node, error) {
	panic("TODO")
}
