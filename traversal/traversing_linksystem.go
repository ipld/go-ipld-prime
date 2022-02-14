package traversal

import (
	"io"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
)

type pathNode struct {
	link     datamodel.Link
	children map[datamodel.PathSegment]*pathNode
}

func newPath(link datamodel.Link) *pathNode {
	return &pathNode{
		link:     link,
		children: make(map[datamodel.PathSegment]*pathNode),
	}
}

func (pn pathNode) addPath(p []datamodel.PathSegment, link datamodel.Link) {
	if len(p) == 0 {
		return
	}
	if _, ok := pn.children[p[0]]; !ok {
var child *pathNode
if len(p) == 1 {
  child = newPath(link)
} else {
   child = newPath(nil)
}
		pn.children[p[0]] = child
	}
	pn.children[p[0]].addPath(p[1:], link)
}

func (pn pathNode) allLinks() []datamodel.Link {
	if len(pn.children) == 0 {
		return []datamodel.Link{pn.link}
	}
	links := make([]datamodel.Link, 0)
	if pn.link != nil {
		links = append(links, pn.link)
	}
	for _, v := range pn.children {
		links = append(links, v.allLinks()...)
	}
	return links
}

// getPaths returns reconstructed paths in the tree rooted at 'root'
func (pn pathNode) getLinks(root datamodel.Path) []datamodel.Link {
	segs := root.Segments()
	switch len(segs) {
	case 0:
		if pn.link != nil {
			return []datamodel.Link{pn.link}
		}
		return []datamodel.Link{}
	case 1:
		// base case 1: get all paths below this child.
		next := segs[0]
		if child, ok := pn.children[next]; ok {
			return child.allLinks()
		}
		return []datamodel.Link{}
	default:
	}

	next := segs[0]
	if _, ok := pn.children[next]; !ok {
		// base case 2: not registered sub-path.
		return []datamodel.Link{}
	}
	return pn.children[next].getLinks(datamodel.NewPathNocopy(segs[1:]))
}

// TraverseResumer allows resuming a progress from a previously encountered path in the selector.
type TraverseResumer func(from datamodel.Path) error

type traversalState struct {
	underlyingOpener linking.BlockReadOpener
	position         int
	pathOrder        map[int]datamodel.Path
	pathTree         *pathNode
	target           *datamodel.Path
	progress         *Progress
}

func (ts *traversalState) resume(from datamodel.Path) error {
	if ts.progress == nil {
		return nil
	}
	// reset progress and traverse until target.
	ts.progress.SeenLinks = make(map[datamodel.Link]struct{})
	ts.position = 0
	ts.target = &from
	return nil
}

func (ts *traversalState) traverse(lc linking.LinkContext, l ipld.Link) (io.Reader, error) {
	// when not in replay mode, we track metadata
	if ts.target == nil {
		ts.pathOrder[ts.position] = lc.LinkPath
		ts.pathTree.addPath(lc.LinkPath.Segments(), l)
		ts.position++
		return ts.underlyingOpener(lc, l)
	}

	// if we reach the target, we exit replay mode (by removing target)
	if lc.LinkPath.String() == ts.target.String() {
		ts.target = nil
		return ts.underlyingOpener(lc, l)
	}

	// when replaying, we skip links not of our direct ancestor,
	// and add all links on the path under them as 'seen'
	targetSegments := ts.target.Segments()
	seg := lc.LinkPath.Segments()
	for i, s := range seg {
		if i >= len(targetSegments) {
			break
		}
		if targetSegments[i].String() != s.String() {
			links := ts.pathTree.getLinks(datamodel.NewPathNocopy(seg[0:i]))
			for _, l := range links {
				ts.progress.SeenLinks[l] = struct{}{}
			}
			return nil, SkipMe{}
		}
	}

	// descend.
	return ts.underlyingOpener(lc, l)
}

// WithTraversingLinksystem extends a progress for traversal such that it can
// subsequently resume and perform subsets of the walk efficiently from
// an arbitrary position within the selector traversal.
func WithTraversingLinksystem(p *Progress) (TraverseResumer, error) {
	ts := &traversalState{
		underlyingOpener: p.Cfg.LinkSystem.StorageReadOpener,
		pathOrder:        make(map[int]datamodel.Path),
		pathTree:         newPath(nil),
		progress:         p,
	}
	p.Cfg.LinkSystem.StorageReadOpener = ts.traverse
	return ts.resume, nil
}
