package ipld

// Traversal is an applicative function which takes one node and returns another,
// while also returning a Path describing a way to repeat the traversal, and
// an error if any part of the traversal failed.
//
// The most common type of Traversal is ipld.Path.Traversal, but it's possible
// to implement other kinds of Traversal function: for example, one could
// implement a traversal algorithm which performs some sort of search to
// select a target node (rather than knowing where it's going before it
// starts, as Path.Traversal does).
//
// The Traversal interface is specified to return a Path primarily for
// logging/debugability/comprehensibility reasons.
//
// In the event of an error, the Traversal may return a nil finish node
// and empty path; or, it may return partial progress values for node
// and path in addition to the error.  It is thus not generally correct to
// use the finish node reference until having checked for an error.
type Traversal func(start Node) (finish Node, path Path, err error)
