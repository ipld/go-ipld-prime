package nodeworkshop

/*
	Experiment closed: we didn't reach a clearly better error handling
	strategy that's ergonomic and performant, and what we did reach is
	probably stil better implemented as a wrapper package like 'fluent'.
*/
/*
	The hypothesis I want to explore here is if having a "workshop" that
	stores some global state and behaviors (and is very visibly passed down)
	is a good idea.

	So far we don't have such a thing; even with the assembler design,
	the interface doesn't really *say* much about state coming down with it.

	One particularly useful thing this would enable is storing any
	error -- _and any error handling strategy_ -- in one clear place.

	This might also give us a place to store the "extras" tree in case of
	tolerant unmarshal modes.

	Maybe we can store a pointer to a 'prior' node and do the "amend" features
	this way, instead of having them appear at fractal depth?  Would that help?

	Naming bikeshed: this isn't really about building *a* node anymore;
	it's about building a whole tree of them at once.  Maybe another name
	could reflect that better.  (We don't really have a term for "the tree of
	nodes you get from a single block", yet.  Maybe we would benefit from one.)
*/

type NodeWorkshop struct {
	Err     error       // used to return errors
	OnError func(error) // error will always be stored in `nw.Err`, but this also allows to react instantly (or, cancel it by nilling it).

	// making it non lexically obvious if error checking is required or not... no, we're not doing that;
	// no builder and assembler methods should return errors anymore; this is purely an immediacy thing.
	// (questionable whether the callback is needed, really.  panicking is the only useful thing it can do.  and is that useful, really?  might be better to just add a 'MustBeUnerrored' checkpoint func.)
}

type ErrInvalidKindForNodeStyle struct{}
type ErrValidationFailed struct{ nested error } // almost all typed errors go in here.  even ErrInvalidStructKey, possibly?  not sure on bind.Node with atlas miss.
type ErrInvalidStructKey struct{}
type ErrListOverrun struct{}          // only really possible with struct types with list reprs, so... again, categorization questions.
type ErrMissingRequiredField struct{} // definitely valid to fit inside ErrValidationFailed.
type ErrInvalidKindForKindedUnion struct{}
type ErrInvalidUnionDiscriminant struct{} // could break this down further into like three different kinds of error (missing descrim key, missing content key, wrong descrim) and/or merge with above.
type ErrRepeatedMapKey struct{}
type ErrCannotBeNull struct{} // arguably either ErrInvalidKindForNodeStyle or ErrValidationFailed.

// We had speculated previously about "unwinding" construction errors to provide context about where in the tree they occurred.
// We don't currently have a way to do that.

/*
	About ErrValidationFailed:

	- might be better to make this an interface?
		- dubious: hard to imagine what code would branch on that.
	- can we really even argue ErrInvalidKindForNodeStyle is not a validation error, if it comes up deep in a tree?
		- for typed nodes... no, I don't think so.  maybe typed nodes just shouldn't use that error (even though it exists, for bind nodes or other applications like free.justString).
		- n.b. ErrInvalidKindForNodeStyle might be particularly likely to arise for (untyped) map keys.
*/

type Node interface {
	// ELIDED: all the familiar reader methods
}

type NodeAssembler interface {
	// comes implicitly bound to some memory already
	// is not terribly reusable because of the above
}

type NodeBuilder interface {
	// comes with knowledge of type but not already having bound memory
	// is reusable because it starts the process
}

func demo() {
	//nsA.
	nsB.Construct(field1nb, arg0scalar, arg1scalar, arg2scalar)
	nsC.Construct(field2nb, arg0scalar, arg1scalar)
	// i still can't get it to nest for a dang.
	// we can make a flip like this, but that doesn't help nest.
	// i literally don't know how to do this without some currying or context syntax.

	// oh.  unless it returns a function pointer -- `func(Assembler) (sideeffecting)` -- in a reasonably cheap way.
	//  (i don't know if that's even possible: I doubt it is: it's a closure no matter how you do it; it _will_ shove an alloc somewhere.)

	// okay, you *can* also do it combined with deep-dotting syntax:
	nsC.Construct(ws.MutA().MutB(), arg0scalar, arg1scalar)
	// ... but it's not entirely clear how much this really buys you.
	//  I guess it means you can do one struct (at maxdepth-1) per line rather than roughly one field per line, and that's something.
	//  But this is still an entirely unnatural syntax that has nothing to do with how people would normally create structs in golang, so we're sorta barking after cars that are quickly vanishing down entirely the wrong highway.
}
