package typed

import "github.com/ipld/go-ipld-prime"

// typed.Link is a superset of the ipld.Link interface, and has one additional behavior.
//
// A typed.Link contains a hint for the appropriate node builder to use for loading data
// on the other side of a link, so that it can be assembled into a node representation
// and validated against the schema as quickly as possible
//
// So, for example, if you wanted to support loading the other side of a link
// with a code-gen'd node builder while utilizing the automatic loading facilities
// of the traversal package, you could write a LinkNodeBuilderChooser as follows:
//
// func LinkNodeBuilderChooser(lnk ipld.Link, lnkCtx ipld.LinkContext) ipld.NodeBuilder {
//  if tlnk, ok := lnk.(typed.Link); ok {
//    return tlnk.SuggestedNodeBuilder()
//  }
//	return ipldfree.NodeBuilder()
// }
//
type Link interface {
	SuggestedNodeBuilder() ipld.NodeBuilder
}
