package ipld

// NodeAssembler is the interface that describes all the ways we can set values
// in a node that's under construction.
//
// To create a Node, you should start with a NodeBuilder (which contains a
// superset of the NodeAssembler methods, and can return the finished Node
// from its `Build` method).
//
// Why do both this and the NodeBuilder interface exist?
// When creating trees of nodes, recursion works over the NodeAssembler interface.
// This is important to efficient library internals, because avoiding the
// requirement to be able to return a Node at any random point in the process
// relieves internals from needing to implement 'freeze' features.
// (This is useful in turn because implementing those 'freeze' features in a
// language without first-class/compile-time support for them (as golang is)
// would tend to push complexity and costs to execution time; we'd rather not.)
type NodeAssembler interface {
	BeginMap(sizeHint int) (MapNodeAssembler, error)
	BeginList(sizeHint int) (ListNodeAssembler, error)
	AssignNull() error
	AssignBool(bool) error
	AssignInt(int) error
	AssignFloat(float64) error
	AssignString(string) error
	AssignBytes([]byte) error

	Assign(Node) error // if you already have a completely constructed subtree, this method puts the whole thing in place at once.

	Style() NodeStyle // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
}

// MapNodeAssembler assembles a map node!  (You guessed it.)
//
// Methods on MapNodeAssembler must be called in a valid order:
// assemble a key, then assemble a value, then loop as long as desired;
// when finished, call 'Finish'.
//
// Incorrect order invocations will panic.
// Calling AssembleKey twice in a row will panic;
// calling AssembleValue before finishing using the NodeAssembler from AssembleKey will panic;
// calling AssembleValue twice in a row will panic;
// etc.
//
// Note that the NodeAssembler yielded from AssembleKey has additional behavior:
// if the node assembled there matches a key already present in the map,
// that assembler will emit the error!
type MapNodeAssembler interface {
	AssembleKey() NodeAssembler   // must be followed by call to AssembleValue.
	AssembleValue() NodeAssembler // must be called immediately after AssembleKey.

	AssembleDirectly(k string) (NodeAssembler, error) // shortcut combining AssembleKey and AssembleValue into one step; valid when the key is a string kind.

	Finish() error

	KeyStyle() NodeStyle   // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
	ValueStyle() NodeStyle // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
}

type ListNodeAssembler interface {
	AssembleValue() NodeAssembler

	Finish() error

	ValueStyle() NodeStyle // you probably don't need this (because you should be able to just feed data and check errors), but it's here.
}

type NodeBuilder interface {
	NodeAssembler

	// Build returns the new value after all other assembly has been completed.
	// A method on the NodeAssembler that finishes assembly of the data should
	// be called first (e.g., any of the "Assign*" methods, or "Finish" if
	// the assembly was for a map or a list); that method still has all
	// responsibility for validating the assembled data and returning
	// any errors from that process.
	//
	// REVIEW: can we just... get rid of the error return here?  Suspect yes.
	Build() (Node, error)

	// Resets the builder.  It can hereafter be used again.
	// Reusing a NodeBuilder can reduce allocations and improve performance.
	//
	// Only call this if you're going to reuse the builder.
	// (Otherwise, it's unnecessary, and may cause an unwanted allocation).
	Reset()
}

// Complex keys: What do they come from?  (What arrre they _good_ for? (Absolutely nothin, say it again))
// - in the Data Model, they don't exist.  They just... don't exist.
//   - json and javascript could never real deal reliably with numbers, so we just... nope.
//   - maps keyed by multiple kinds are also beyond the pale in many host languages, so, again... nope.
// - in the schema types layer, they can exist.
//   - a couple things can reach it:
//     - `type X struct {...} representation stringjoin`
//     - `type X struct {...} representation stringpairs`
//     - `type X map {...} representation stringpairs` // *maybe*.  go won't allow using this as a map key except in string form anyway.
//     - we don't know what the syntax would look like for a type-level int, but, haven't ruled it out either.
//   - but when feeding data in via representation: it's all strings, of course.
//   - if we have codegen and are app author, we can use native methods that pass-by-value.
//   - so it's ONLY when doing generic code, or typed but not using codegen, that we face these apis.
// - and it's ONLY -- this should go without saying, but let's say it -- relevant when thinking about map keys.
//   - structs can't have complex keys.  field names are strings.  done.
//   - lists can't have complex keys.  obvious category error.
//   - enums can't have complex keys.  because we said so; or, obvious category error; take your pick.
//   - why is this important?  well, it means complex keys only exist in a place where values are all going to use the same style/builder.
//     - which means `AssembleInsertion() (NodeAssembler, NodeAssembler)` is actually kinda back on the table.
//       - though since it would only work in *some* cases... still serious questions about whether that'd be a good API to show off.
//       - we still also need to be able to get a NodeAssembler for streaming in value content when handling structs, so, then we'd end up with two APIs for this?
//         - yeah, this design still does not go good places; abort.

// Hey what actually happens if you make it interally do `map[Node]Thingy` and all keys are concretely `K` (*not* `*K`)?
//  - You get the boxing to happen once, at least.
//  - But a boxing alloc does happen for *every single* key.
//  - Equality and equality-dependent behaviors are correct.
//  - But equality and lookups will be a tad (just a tad -- a few nanos) slower, because of all the interface checks.
// To say I'm disappointed by the idea of an alloc per key is an understatement.
// And it would end up with a bunch of `Node` iface flying around filled with non-pointer types, and that... is a source of unusualness that I suspect will bring major antijoy.
// But if you're going to iterate a thing literally twice, it would be preferable to the other status quo's.
//  (Except for the `[]Entry{K,V}` strategy or equivalently `[]K,[]V` strategy.  Those can still do batched-alloc keys.)
//   (Holy crap.  Actually... those can even return immutable keys during mid-construction.  Which fixes... oh my... one of the biggest issues that got us here.  Uh.)
//    (No, Falsey.  It could; but if you actually held onto those, and had to resize the array, you'd end up holding onto multiple arrays, and not be happy.  It would, in effect, be a feature you'd never want to use, unless you knew in a fractal of detail what you were doing, and even then, it's harder to imagine good rather than footgunny uses.)

// Something to consider for typed but non-codegen map implementations: `k struct{ body string; offsets [10]int }`.
//  This could make it possible to create the typed node relatively cheaply (no parsing).
//  But unclear if this is useful on the whole.  Doesn't avoid many allocs.  Doesn't make eq for map cheaper.
//  And obviously, it has a limit on how complex of a struct it works for (the array must be fixed size and not slice, in order to be comparable).
//  The less ornate option of just have a double map `{ keys map[string]Node; values map[string]Node }`, or `map[string]struct{k, v Node}` might remain superior to this idea.
