hacking gengo
=============

What the heck?
--------------

We're doing code generation.

The name of the game is "keep it simple".
Most of this is implemented as string templating.
No, we didn't use the Go AST system.  We could have; we didn't.
Implementing this as string templating seemed easier to mentally model,
and the additional value provided by use of AST libraries seems minimal
since we feed the outputs into a compiler for verification immediately anyway.

Some things seem significantly redundant.
That's probably because they are.
In general, if there's a choice between apparent redundancy in the generator itself
versus almost any other tradeoff which affects the outputs, we prioritize the outputs.
(This may be especially noticable when it comes to error messages: we emit a lot
of them... while making sure they contain very specific references.  This leads
to some seemingly redundant code, but good error messages are worth it.)


Entrypoints
-----------

The `gen_test.go` file is the effective "main" method right now.
It contains substantial amounts of hardcoded testcases.

Run the tests in the `./_test` subpackage explicitly to make sure the
generated code passes its own interface contracts and tests.

The eventual plan is be able to drive this whole apparatus around via a CLI
which consumes IPLD Schema files.
Implementing this can come after more of the core is done.
(Seealso the `schema/tmpBuilders.go` file a couple directories up for why
this is currently filed as nontrivial/do-later.)


Organization
------------

### How many things are generated, anyway?

There are roughly *seven* categories of API to generate per type:

- 1: the readonly thing a native caller uses
- 2: the builder thing a native caller uses
- 3: the readonly typed node
- 4: the builder for typed node
- 5: the readonly representation node
- 6: the builder via representation
- 7: and a maybe wrapper

(And these are just the ones nominally visible in the exported API surface!
There are several more concrete types than this implied by some parts of that list,
such as iterators for the nodes, internal parts of builders, and so forth.)

These numbers will be used to describe some further organization.

### How are the generator components grouped?

There are three noteworthy types of generator internals:

- `typedNodeGenerator`
- `nodeGenerator`
- `nodebuilderGenerator`

The first one is where you start; the latter two do double duty for each type.

Exported types for purpose 1, 2, 3, and 7 are emitted from `typedNodeGenerator` (3 from the embedded `nodeGenerator`).

The exported type for purpose 5 is emitted from another `nodeGenerator` instance.

The exported types for purposes 4 and 6 are emitted from two distinct `nodebuilderGenerator` instances.

For kinds that have more than one known representation strategy,
there may be more than two implementations of `nodeGenerator` and `nodebuilderGenerator`!
(There's always one for the type-semantics node+builder,
and then one more *for each* representation strategy.)

### How are files and their contents grouped?

Most of the files in this package are following a pattern:

- for each kind:
	- `genKind{Kind}.go` -- has emitters for the native type parts (1, 2, 7).
	- `genKind{Kind}Node.go` -- has emitters for the typed node parts (3, 4), and the entrypoint to (5).
	- for each representation that kind can have:
		- `genKind{Kind}Repr{ReprStrat}.go` -- has emitters for (5, 6).
