package ipld

import "context"

// Traversal is an applicative function which takes one node and returns another,
// while also returning a Path describing a way to repeat the traversal, and
// an error if any part of the traversal failed.
//
// Traversal requires a TraversalProgress argument (which may be zero-valued),
// and returns a new TraversalProgress containing an updated Path.
//
// The most common type of Traversal is ipld.Path.Traversal, but it's possible
// to implement other kinds of Traversal function: for example, one could
// implement a traversal algorithm which performs some sort of search to
// select a target node (rather than knowing where it's going before it
// starts, as Path.Traversal does).
//
// In the case of error, the returned TraversalProgress may be zero and the
// Node may be nil.  (The particular Path at which the error was encountered
// may be encoded in the error type.)
type Traversal func(tp TraversalProgress, start Node) (tp2 TraversalProgress, finish Node, err error)

type TraversalProgress struct {
	Ctx  context.Context
	Path Path
}
