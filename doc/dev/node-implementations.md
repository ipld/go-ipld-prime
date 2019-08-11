Dev Notes: on Node implementations
==================================

> (This document is aimed for developers and contributors to the library;
> if you're only using the library, it should not be necessary to read this.)

The concept of "Node" in IPLD is defined by the
[IPLD Data Model](https://github.com/ipld/specs/tree/master/data-model-layer),
and the interface of `ipld.Node` in this library is designed to make this
model manipulable in Golang.

`Node` is an interface rather than a concrete implementation because
there are many different performance tradeoffs which can be made,
and we have multiple implementations that make them!
Codegenerated types also necessitate having a `Node` interface they can conform to.



Designing a Node Implementation
-------------------------------

Concerns:

- 0. Correctness
- 1. Immutablity
- 2. Performance

A `Node` implementation must of course conform with the Data Model.
Some nodes (especially, those that also implement `typed.Node`) may also
have additional constraints.

A `Node` implementation must maintain immutablity, or it shatters abstractions
and breaks features that build upon the assumption of immutable nodes (e.g. caches).

A `Node` implementation should be as fast as possible.
(How fast this is, of course, varies -- and different implementations may make
different tradeoffs, e.g. often at the loss of generality.)

Note that "generality" is not on the list.  Some `Node` implementations are
concerned with being usable to store any shape of data; some are not.
(A `Node` implementation will usually document which camp it falls into.)


### Castability

Castability refers to whether the `Node` abstraction can be added or removed
(also referred to as "boxing" and "unboxing")
by use of a cast by user code outside the library.

Castability relates to all three of Correctness, Immutablity, and Performance.

- if something can be unboxed via cast, and thence become mutable, we have an Immutablity problem.
- if something mutable can be boxed via cast, staying mutable, we have an Immutablity problem.
- if something can be boxed via cast, and thence skip a validator, we have a Correctness problem.

(The relationship to performance is less black-and-white: though performance
considerations should always be backed up by benchmarks, casting can do well.)

If castability would run into one of these problems,
then a Node implementation must avoid it.
(A typical way to do this is to make a single-field struct.)

Whether a `Node` implementation will encounter these problems varies primarily on
the kind (literally, per `reflect.Kind`) of golang type is used in the implementation,
and whether the `Node` is "general" or can have an addition validators and constraints.

#### Castability cases by example

Castability for strings is safe when the `Node` is "general" (i.e. has no constraints).
With no constraints, there's no Correctness concern;
and since strings are immutable, there's no Immutablity concern.

Castability for strings is often *unsafe* when the `Node` is a `typed.Node`.
Typed nodes may have additional constraints, so we would have a Correctness problem.
(Note that the way we handle constraints in codegeneration means users can add
them *after* the code is generated, so the generation system can't presume
the absense of constraints.)

Castability for other scalar types (int, float, etc) are safe when the `Node` is "general"
for the same reasons it's safe for strings: all these things are pass-by-value
in golang, so they're effectively immutable, and thus we have no concerns.

Castability for bytes is a trickier topic.
See also [#when-performance-wins-over-immutablity].
(TODO: the recommended course of action here is not yet clear.
I'd default to suggesting it should use a `struct{[]byte}` wrapping,
but if the performance cost of that is high, the value may be dubious.)

#### Zero values

If a `Node` is implemented as a golang struct, zero values may be a concern.

If the struct type is unexported, the concern is absolved:
the zero value can't be initialized outside the package.

If the `Node` implementation has no other constraints
(e.g., it's not also a `typed.Node` in addition to just an `ipld.Node`),
the concern is (alomst certainly) absolved:
the zero value is simply a valid value.

For the remaining cases: it's possible that the zero value is not valid.
This is problematic, because in the system as a whole, we use the existence
of a value that's boxed into a `Node` as the statement that the value is valid,
rather than suffering the need for "defensive" checks cropping up everywhere.

(TODO: the recommended course of action here is not yet clear.
Making the type unexported and instead making an exported interface with a
single implementation may be one option, and it's possible it won't even be
noticably expensive if we already have to fit `Node`, but I'm not sure I've
reconnoitered all the other costs of that (e.g., godoc effects?).
It's possible this will be such a corner case in practice that we might
relegate the less-ergonomic mode to being an adjunct option for codegen.)



When Performance wins over Immutablity
--------------------------------------

Ideally?  Never.  In practice?  Unfortunately, sometimes.


### bytes

There is no way to have immutable byte slices in Go.
Defensive copying is also ridiculously costly.

Methods that return byte slices typically do so without defensive copy.

Methods doing this should consistently document that
"It is not lawful to mutate the slice returned by this function".
