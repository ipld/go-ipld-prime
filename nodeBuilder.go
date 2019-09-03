package ipld

// NodeBuilder is an interface that describes creating new Node instances.
//
// The Node interface is entirely read-only methods; a Node is immutable.
// Thus, we need a NodeBuilder system for creating new ones; the builder
// is mutable, and when we're done accumulating mutations, we take the
// accumulated data and produce an immutable Node out of it.
//
// Separating mutation into NodeBuilder and keeping Node immutable makes
// it possible to perform caching (or rather, memoization, since there's no
// such thing as cache invalidation for immutable systems) of computed
// properties of Node; use copy-on-write algorithms for memory efficiency;
// and to generally build pleasant APIs.
//
// Each package in `go-ipld-prime//impl/*` that implements ipld.Node also
// has a NodeBuilder implementation that produces new nodes of that same
// package's type.
//
// Most Node implementations also have a method which returns a NodeBuilder
// that produces more nodes of their same concrete implementation type.
// This is useful for algorithms that work on trees of nodes: this NodeBuilder
// getter will be used when an update deep in the tree causes a need to
// create several new nodes to propagate the change up through parent nodes.
//
// The NodeBuilder retrieved from a Node can also be used to do *updates*:
// consider the AmendMap and AmendList methods.  These methods are useful
// not just for programmer convenience, but also because they can reuse memory,
// sharing any common segments of memory with the earlier Node.
// (In the NodeBuilder exposed by the `go-ipld-prime//impl/*` packages, these
// methods are equivalent to their Create* counterparts.  As there's no
// "existing" node for them to refer to, it's treated the same as amending
// an empty node.)
//
// NodeBuilder instances obtained from Node.GetBuilder may carry some of the
// additional logic of their parent with them to the new Node they produce.
// For example, the NodeBuilder from typed.Node.GetBuilder may keep the type
// info and type constraints of their parent with them!
// (Continuing the typed.Node example: if you have a typed.Node that is
// constrained to be of some `type Foo = {Bar:Baz}` type, then any new Node
// produced from its NodeBuilder will still answer
// `n.(typed.Node).Type().Name()` as `Foo`; and if
// `n.GetBuilder().AmendMap().Insert(...)` is called with nodes of unmatching
// type given to the insertion, the builder will error!)
type NodeBuilder interface {
	CreateMap() (MapBuilder, error)
	AmendMap() (MapBuilder, error)
	CreateList() (ListBuilder, error)
	AmendList() (ListBuilder, error)
	CreateNull() (Node, error)
	CreateBool(bool) (Node, error)
	CreateInt(int) (Node, error)
	CreateFloat(float64) (Node, error)
	CreateString(string) (Node, error)
	CreateBytes([]byte) (Node, error)
	CreateLink(Link) (Node, error)
}

// MapBuilder is an interface for creating new Node instances of kind map.
//
// A MapBuilder is generally obtained by getting a NodeBuilder first,
// and then using CreateMap or AmendMap to begin.
//
// Methods mutate the builder's internal state; when done, call Build to
// produce a new immutable Node from the internal state.
// (After calling Build, future mutations may be rejected.)
//
// Insertion methods error if the key already exists.
//
// The BuilderForKeys and BuilderForValue functions return NodeBuilders
// that can be used to produce values for insertion.
// If you already have the data you're inserting, you can use those Nodes;
// if you don't, use these builders.
// (This is particularly relevant for typed nodes and bind nodes, since those
// have internal specializations, and not all NodeBuilders for them are equal.)
// Note that BuilderForValue requires a key as a parameter!
// This is because typed nodes which are structs may return different builders
// per field, specific to the field's type.
//
// You may be interested in the fluent package's fluent.MapBuilder equivalent
// for common usage with less error-handling boilerplate requirements.
type MapBuilder interface {
	Insert(k, v Node) error
	Delete(k Node) error
	Build() (Node, error)

	BuilderForKeys() NodeBuilder
	BuilderForValue(k string) NodeBuilder

	// FIXME we might need a way to reject invalid 'k' for BuilderForValue.
	//  We certainly can't wait until the insert, because we need the value builder before that (duh!).
	//  However, you probably should've applied the BuilderForKeys to the key value already,
	//   and that should've told you about most errors...?
	//   Hrm.  Not sure if we wanna rely on this.
	//  Panic?  or return an error, and be sad about breaking chaining of calls?  or return a curried-error thunk?
	//  Or can we shift all the responsibility to BuilderForKeys after all (with panic as 'unreachable' fallback)?
	//
	// - for maps with typed keys that have constraints, the rejection should've come from the key builder.  fine.
	//   - builderForValue also doesn't need to vary at all in this case; you could've given an empty 'k' and cached that one.
	//     - though note we haven't exposed a way to *detect* that yet (and it's questioned whether we should until more profiling/optimization info comes in).
	// - for structs, the rejection *could* come from the key builder, but we haven't decided if that's required or not.
	//   - requiring a builder for keys that whitelists the valid keys but still returns plain string nodes is viable...
	//     - but a little odd looking, since the returned thing is going to be a plain untyped string kind node.
	//       - in other words, the type wouldn't carry info about the constraints it has passed through; not wrong, but perhaps a design smell.
	//     - for codegen'd impls, this would be compiled into a string switch and be pretty cheap.  viable.
	//     - for runtime-wrapper typed.Node impls... still viable, but a little heavier.
}

// ListBuilder is an interface for creating new Node instances of kind list.
//
// A ListBuilder is generally obtained by getting a NodeBuilder first,
// and then using CreateList or AmendList to begin.
//
// Methods mutate the builder's internal state; when done, call Build to
// produce a new immutable Node from the internal state.
// (After calling Build, future mutations may be rejected.)
//
// Methods may error when handling typed lists if non-matching types are inserted.
//
// The BuilderForValue function returns a NodeBuilder
// that can be used to produce values for insertion.
// If you already have the data you're inserting, you can use those Nodes;
// if you don't, use these builders.
// (This is particularly relevant for typed nodes and bind nodes, since those
// have internal specializations, and not all NodeBuilders for them are equal.)
// Note that BuilderForValue requires an index as a parameter!
// In most cases, this is not relevant and the method returns a constant NodeBuilder;
// however, typed nodes which are structs and have list representations may
// return different builders per index, corresponding to the types of its fields.
//
// You may be interested in the fluent package's fluent.ListBuilder equivalent
// for common usage with less error-handling boilerplate requirements.
type ListBuilder interface {
	AppendAll([]Node) error
	Append(v Node) error
	Set(idx int, v Node) error
	Build() (Node, error)

	BuilderForValue(idx int) NodeBuilder

	// FIXME the question about rejection of invalid idx applies here as well,
	//  for all the same reasons it came up for BuilderForValue on maps:
	//   structs with tuple representation provoke all the exact same issues.
}

// future: add AppendIterator() methods (when we've implemented iterators!)

// future: add InsertConverting(map[string]interface{}) and similar methods.
//  (some open questions about how useful that is, given ipldbind should likely be more efficient, depending on use case.)

// future: define key ordering semantics during map insertion.
//  methods for re-ordering will probably be wanted someday.

// review: MapBuilder.Delete as an API is dangerously prone to usage which is accidentally quadratic.
//  https://accidentallyquadratic.tumblr.com/post/157496054437/ruby-reject describes a similar issue.
//   an internal implementation which accumulates oplogs is one fix, but not a joy either (not alloc free).
//   a totally different API -- say, `NodeBuilder.AmendMapWithout(...)` -- might be the best approach.
