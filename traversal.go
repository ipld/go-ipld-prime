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

func ContinuedTraversal(tfn Traversal, alreadyReachedPath Path) Traversal {
	return func(start Node) (finish Node, path Path, err error) {
		n, p, e := tfn(start)
		combinedSegments := make([]string, len(alreadyReachedPath.segments)+len(p.segments))
		copy(combinedSegments, alreadyReachedPath.segments)
		copy(combinedSegments[len(alreadyReachedPath.segments):], p.segments)
		return n, Path{combinedSegments}, e
		// REVIEW: this composition is imperfect!  the error might embed a copy of of a shorter less global path in it.
		//  We have two options for dealing with this:
		//   1) Errors become strongly typed and we make sure we can reach in and fix it up;
		//   2) We replace Traversal with something that already can do the "continued" stuff;
		//   3) (we ignore it.)
		//  Option 2 is probably pretty defensible; we don't need Traversal to be simple to implement since
		//   the vast, vast majority of users will simply use Path.
		//  (We should also pursue Option 1 *regardless*, of course; it's just an implementation dependency/order question.)
	}
}
