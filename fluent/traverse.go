package fluent

import (
	ipld "github.com/ipld/go-ipld-prime"
)

func Traverse(
	node ipld.Node,
	path ipld.Path,
) (reachedNode ipld.Node, reachedPath ipld.Path) {
	return TraverseUsingTraversal(node, path.Traverse)
}

func TraverseUsingTraversal(
	node ipld.Node,
	traversal ipld.Traversal,
) (reachedNode ipld.Node, reachedPath ipld.Path) {
	return nil, ipld.Path{} // TODO this method doesn't have much to do except the error bouncing
}
