package traversal_test

import (
	"errors"
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/traversal"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	selectorparse "github.com/ipld/go-ipld-prime/traversal/selector/parse"
)

func TestWalkResume(t *testing.T) {
	seen := 0
	count := func(p traversal.Progress, n datamodel.Node, _ traversal.VisitReason) error {
		seen++
		return nil
	}

	lsys := cidlink.DefaultLinkSystem()
	lsys.SetReadStorage(&store)
	p := traversal.Progress{
		Cfg: &traversal.Config{
			LinkSystem:                     lsys,
			LinkTargetNodePrototypeChooser: basicnode.Chooser,
		},
	}
	resumer, err := traversal.WithTraversingLinksystem(&p)
	if err != nil {
		t.Fatal(err)
	}
	sd := selectorparse.CommonSelector_ExploreAllRecursively
	s, _ := selector.CompileSelector(sd)
	if err := p.WalkAdv(rootNode, s, count); err != nil {
		t.Fatal(err)
	}
	if seen != 14 {
		t.Fatalf("expected total traversal to visit 14 nodes, got %d", seen)
	}

	// resume from beginning.
	resumer(datamodel.NewPath(nil))
	seen = 0
	if err := p.WalkAdv(rootNode, s, count); err != nil {
		t.Fatal(err)
	}
	if seen != 14 {
		t.Fatalf("expected resumed traversal to visit 14 nodes, got %d", seen)
	}

	// resume from middle.
	resumer(datamodel.NewPath([]datamodel.PathSegment{datamodel.PathSegmentOfString("linkedMap")}))
	seen = 0
	if err := p.WalkAdv(rootNode, s, count); err != nil {
		t.Fatal(err)
	}
	// one less: will not visit 'linkedString' before linked map.
	if seen != 13 {
		t.Fatalf("expected resumed traversal to visit 13 nodes, got %d", seen)
	}

	// resume from middle.
	resumer(datamodel.NewPath([]datamodel.PathSegment{datamodel.PathSegmentOfString("linkedList")}))
	seen = 0
	if err := p.WalkAdv(rootNode, s, count); err != nil {
		t.Fatal(err)
	}
	// will not visit 'linkedString' or 'linkedMap' before linked list.
	if seen != 7 {
		t.Fatalf("expected resumed traversal to visit 7 nodes, got %d", seen)
	}
}

func TestWalkResumePartialWalk(t *testing.T) {
	seen := 0
	limit := 0
	countUntil := func(p traversal.Progress, n datamodel.Node, _ traversal.VisitReason) error {
		seen++
		if seen >= limit {
			return traversal.SkipMe{}
		}
		return nil
	}

	lsys := cidlink.DefaultLinkSystem()
	lsys.SetReadStorage(&store)
	p := traversal.Progress{
		Cfg: &traversal.Config{
			LinkSystem:                     lsys,
			LinkTargetNodePrototypeChooser: basicnode.Chooser,
		},
	}
	resumer, err := traversal.WithTraversingLinksystem(&p)
	if err != nil {
		t.Fatal(err)
	}
	sd := selectorparse.CommonSelector_ExploreAllRecursively
	s, _ := selector.CompileSelector(sd)
	limit = 9
	if err := p.WalkAdv(rootNode, s, countUntil); !errors.Is(err, traversal.SkipMe{}) {
		t.Fatal(err)
	}
	if seen != limit {
		t.Fatalf("expected partial traversal, got %d", seen)
	}

	// resume.
	resumer(datamodel.NewPath([]datamodel.PathSegment{datamodel.PathSegmentOfString("linkedMap")}))
	seen = 0
	limit = 14
	if err := p.WalkAdv(rootNode, s, countUntil); err != nil {
		t.Fatal(err)
	}
	if seen != 13 {
		t.Fatalf("expected resumed traversal to visit 13 nodes, got %d", seen)
	}
}
