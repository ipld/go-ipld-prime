package science

/*
A fairly substantial redraft of NodeBuilders -- in this one,
the recursion happens in such a way that you *don't* get a handle to a Node
immediately when completing every builder.

The idea is that this might make it easier to build large recursive values
(both for caller ergonomics, and for internal efficiency) since we don't
have to be concerned with returning an immutable Node at high granularity at
almost every step throughout the process of building the larger tree.
But does it?  Not sure yet.

This exploration is triggered by two things:

Firstly: in general, the NodeBuilder system at present is small-pieces-first...
and that's very bad for memory characteristics.  In trying to spec behaviors
that fix it, we've come up with possible 'memory visibility' models which
let us have the memory amortization we want while keeping the current
small-pieces-first interfaces... but it's... getting arcane.
It's unclear if the effort is worth it -- what if we just admit it and
fully abandon the small-pieces-first model?

Secondly: that increasingly-arcane pretends-to-be small-pieces-first model
is still failing in one critical place: handling keys for maps.
Making it possible to handle complex keys without heap allocs for
semiotically-unnecessary interface wrapping is turning out to be agnozing.
So maybe some different interfaces will make it easier.

Maybe.  Let's see.

*/

type NodeAssembler interface {
	BeginMap() MapNodeAssembler
	BeginList() ListNodeAssembler
	AssignNull()
	AssignBool(bool)
	AssignInt(int)
	AssignFloat(float64)
	AssignString(string)
	AssignBytes([]byte)

	Assign(Node)
	CheckError() error // where are these stored?  when each child is done, could check this.  and have it stored in each child.  means each allocs (nonzero struct)... but isn't that true anyway?  yeah.  well, except now you're apparently getting that for scalars, which is bad.
}

type MapNodeAssembler interface {
	AssembleKey() MapKeyAssembler

	Insert(string, Node)
	Insert2(Node, Node)

	Done()
}
type MapKeyAssembler interface {
	NodeAssembler
	AssembleValue() MapValueAssembler
}
type MapValueAssembler interface {
	NodeAssembler
}

type ListNodeAssembler interface {
	AssembleValue() NodeAssembler

	Append(Node)

	Done()
}

type NodeBuilder interface {
	NodeAssembler
	Build() (Node, error)
}

type Node interface {
	// all the familiar reader methods go here
}

func demo() {
	var nb NodeBuilder
	func(mb MapNodeAssembler) {
		mb.AssembleKey().AssignString("key")
		func(mb MapNodeAssembler) {
			mb.AssembleKey().AssignString("nested")
			mb.AssembleValue().AssignBool(true)
			mb.Done()
		}(mb.AssembleValue().BeginMap())
		mb.AssembleKey().AssignString("secondkey")
		mb.AssembleValue().AssignString("morevalue")
		mb.Done()
	}(nb.BeginMap())
	result, err := nb.Build()
	_, _ = result, err
}

/*

Pros:

- this nested pretty nicely.  not shuffling returned nodes (and errors!) everywhere is actually... prettier to use than what we've got presently.
- having the 'Build' at the end seems... fine.
- notice the lack of childbuilder getter methods.  it's just implied by normal operation.
	- also means it's intensely obvious that none of the assemblers can be reused freely -- whereas with the child builder getters at present just returning a NodeBuilder, that's nonobvious.

Cons:

- not very clear where error handling should go.  punting it *all* to the end is not great.
	- this is a super big deal!
	- when you have a typed map key, does the subsequent `AssembleValue` call trigger validating the key?
	- when you have a typed map value, does the next `AssembleKey` call trigger validating value?
	- what happens for the last value in a map?  now we're really in trouble.  the next action on the parent has to pick up the duty?  no, I really don't think that's even viable.
- is every builder gonna need to have room inside to curry an error?  that means a bunch of things might go from zero to nonzero struct size, which might be consequential.
- i really don't like that we have two different styles of usage appearing in the same interfaces (the 'Assemble' recursors vs 'Insert'/'Append').
	- ... maybe we can get rid of those now?  `mb.AssembleKey().Assign(n1).Assign(n2)`?  no, not quite a one-liner yet...
		- ... yeah, I can't get a one-liner out of this.  either the key or the value might require or design (respectively) a closure in order to flow well.  the former complicates the core interface, and the latter can still be done from the outside but then just looks weird.
- building single scalar values got more annoying.  it's at least two or three lines now: create builder, do assign, call build.
	- question is whether this is the thing to worry about -- we traded this for fewer lines when making recursive structures, which we presumably do more often than not.

Notes:

- yes, also ditched the "amend" methods.  those have been pretty consistently irritating; that said, I don't have a plan to replace their fastpath capabilities either.
- the types still don't particularly do much to force correct usage order.  e.g., your logic must alternate calling AssembleKey and AssembleValue, etc.
  - ...but as found in other explorations, there's a limit to what we can do at best anyway without linear types or rustacean 'move' semantics, etc.  so if we already have to give up a meter shy of the dream, what's another millimeter?
- i don't think there's any strong incentive _not_ to panic throughout this code, performance-wise at least.  something with this many vtables is already fractally uninlinable.
- interesting to note that if a user is creating map keys via a typed nodebuilder (presumably they're structs)... `Assign2(Node,Node)` **still** forces allocs for boxing.
	- while the AssembleKey style of usage avoids a boxing, it's awkward to use if the user wants to use a natively typed builder (`AssembleKey().Assign(nowthis)`?) and then still forces a temp copy-by-value in the middle for no good reason.
	- a natively-typed assign method that works by value is the only way to minimize all the costs at once.
	- we do still need it though: `Assign2(Node,Node)` is the only way to copy an existing typed map key generically (without the temp-copy lamented above for the `AssembleKey.Assign` chain).
- the `Append` and `Assign` method names might be too terse -- we didn't leave space for the natively typed methods, which we generally want to be the shortest ones.
- 'Assign' being a symbol on both NodeAssemler and MapNodeAssembler with different arity might be unfortunate -- would make it impossible to supply both with one concrete type.

*/
