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
	AssignLink(Link) error

	AssignNode(Node) error // if you already have a completely constructed subtree, this method puts the whole thing in place at once.

	// Style returns a NodeStyle describing what kind of value we're assembling.
	//
	// You often don't need this (because you should be able to
	// just feed data and check errors), but it's here.
	//
	// Using `this.Style().NewBuilder()` to produce a new `Node`,
	// then giving that node to `this.AssignNode(n)` should always work.
	// (Note that this is not necessarily an _exclusive_ statement on what
	// sort of values will be accepted by `this.AssignNode(n)`.)
	Style() NodeStyle
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

	// KeyStyle returns a NodeStyle that knows how to build keys of a type this map uses.
	//
	// You often don't need this (because you should be able to
	// just feed data and check errors), but it's here.
	//
	// For all Data Model maps, this will answer with a basic concept of "string".
	// For Schema typed maps, this may answer with a more complex type (potentially even a struct type).
	KeyStyle() NodeStyle

	// ValueStyle returns a NodeStyle that knows how to build values this map can contain.
	//
	// You often don't need this (because you should be able to
	// just feed data and check errors), but it's here.
	//
	// ValueStyle requires a parameter describing the key in order to say what
	// NodeStyle will be acceptable as a value for that key, because when using
	// struct types (or union types) from the Schemas system, they behave as maps
	// but have different acceptable types for each field (or member, for unions).
	// For plain maps (that is, not structs or unions masquerading as maps),
	// the empty string can be used as a parameter, and the returned NodeStyle
	// can be assumed applicable for all values.
	// Using an empty string for a struct or union will return a nil NodeStyle.
	// (Design note: a string is sufficient for the parameter here rather than
	// a full Node, because the only cases where the value types vary are also
	// cases where the keys may not be complex.)
	ValueStyle(k string) NodeStyle
}

type ListNodeAssembler interface {
	AssembleValue() NodeAssembler

	Finish() error

	// ValueStyle returns a NodeStyle that knows how to build values this map can contain.
	//
	// You often don't need this (because you should be able to
	// just feed data and check errors), but it's here.
	//
	// In contrast to the `MapNodeAssembler.ValueStyle(key)` function,
	// to determine the ValueStyle for lists we need no parameters;
	// lists always contain one value type (even if it's "any").
	ValueStyle() NodeStyle
}

type NodeBuilder interface {
	NodeAssembler

	// Build returns the new value after all other assembly has been completed.
	//
	// A method on the NodeAssembler that finishes assembly of the data must
	// be called first (e.g., any of the "Assign*" methods, or "Finish" if
	// the assembly was for a map or a list); that finishing method still has
	// all responsibility for validating the assembled data and returning
	// any errors from that process.
	// (Correspondingly, there is no error return from this method.)
	Build() Node

	// Resets the builder.  It can hereafter be used again.
	// Reusing a NodeBuilder can reduce allocations and improve performance.
	//
	// Only call this if you're going to reuse the builder.
	// (Otherwise, it's unnecessary, and may cause an unwanted allocation).
	Reset()
}

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
