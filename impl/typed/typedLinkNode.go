package typed

import "github.com/ipld/go-ipld-prime"

// typed.LinkNode is a superset of the schema.TypedNode interface, and has one additional behavior.
//
// A typed.LinkNode contains a hint for the appropriate node builder to use for loading data
// on the other side of the link contained within the node, so that it can be assembled
// into a node representation and validated against the schema as quickly as possible
//
// So, for example, if you wanted to support loading the other side of a link
// with a code-gen'd node builder while utilizing the automatic loading facilities
// of the traversal package, you could write a LinkNodeBuilderChooser as follows:
//
// func LinkNodeBuilderChooser(lnk ipld.Link, lnkCtx ipld.LinkContext) ipld.NodeBuilder {
//  if tlnkNd, ok := lnkCtx.LinkNode.(typed.LinkNode); ok {
//    return tlnkNd.SuggestedNodeBuilder()
//  }
//	return ipldfree.NodeBuilder()
// }
//
type LinkNode interface {
	ReferencedNodeBuilder() ipld.NodeBuilder
}
