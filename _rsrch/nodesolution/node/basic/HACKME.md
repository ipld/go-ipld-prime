hackme
======

Design rationale are documented here.

This doc is not necessary reading for users of this package,
but if you're considering submitting patches -- or just trying to understand
why it was written this way, and check for reasoning that might be dated --
then it might be useful reading.

### scalars are just typedefs

This is noteworthy because in codegen, this is typically *not* the case:
in codegen, even scalar types are boxed in a struct, such that it prevents
casting values into those types.

This casting is not a concern for ipldfree types, because
A) we don't have any kind of validation rules to make such casting worrying; and
B) since our types are unexported, casting is still blocked by this anyway.

### pointer-vs-value inhabitant consistency

Builders always return interfaces inhabited by pointers.
Some constructor functions (e.g. `String()`) return an interface inhabited by
the bare typedef and no pointer.

This means you can get `*plainString` and `plainString` inhabitants of `Node`.

(In contrast, codegen systems usually make a stance of consistently using
only pointer inhabitants -- so you can do type cast checks more easily.)

This hasn't been a problem yet, but it's something to keep an eye on.
If this causes the slightest irritation, we'll standardize to pointer inhabitants.

### about builders for scalars

The assembler types for scalars (string, int, etc) are pretty funny-looking.
You might wish to make them work without any state at all!

The reason this doesn't fly is that we have to keep the "wip" value in hand
just long enough to return it from the `NodeBuilder.Build` method -- the
`NodeAssembler` contract for `Assign*` methods doesn't permit just returning
their results immediately.

(Another possible reason is if we expected to use these assemblers on
slab-style allocations (say, `[]plainString`)...
however, this is inapplicable at present, because
A) we don't (except places that have special-case internal paths anyway); and
B) the types aren't exported, so users can't either.)

Does this mean that using `NodeBuilder` for scalars has a completely
unnecessary second allocation, which is laughably inefficient?  Yes.
It's unfortunate the interfaces constrain us to this.
**But**: one typically doesn't actually use builders for scalars much;
they're just here for completeness.
So this is less of a problem in practice than it might at first seem.

More often, one will use the "any" builder (which is has a whole different set
of design constraints and tradeoffs);
or, if one is writing code and knows which scalar they need, the exported
direct constructor function for that kind
(e.g., `String("foo")` instead of `Style__String{}.NewBuilder().AssignString("foo")`)
will do the right thing and do it in one allocation (and it's less to type, too).

### maps and list keyAssembler and valueAssemblers have custom scalar handling

Related to the above heading.

Maps and lists in this package do their own internal handling of scalars,
using unexported features inside the package, because they can more efficient.

### when to invalidate the 'w' pointers

The 'w' pointer -- short for 'wip' node pointer -- has an interesting lifecycle.

In a NodeAssembler, the 'w' pointer should be intialized before the assembler is used.
This means either the matching NodeBuilder type does so; or,
if we're inside recursive structure, the parent assembler did so.

The 'w' pointer is used throughout the life of the assembler.

When assembly becomes "finished", the 'w' pointer should be set to nil,
in order to make it impossible to continue to mutate that node.
However, this doesn't *quite* work... because in the case of builders (at the root),
we need to continue to hold onto the node between when it becomes "finished"
and when Build is called, allowing us to return it (and then finally nil 'w').

This has some kinda wild implications.  In recursive structures, it means the
child assembler wrapper type is the one who takes reponsibility for nilling out
the 'w' pointer at the same time as it updates the parent's state machine to
proceed with the next entry.  In the case of scalars at the root of a build,
it means *you can actually use the assign method more than once*.

We could normalize the case with scalars at the root of a tree by adding another
piece of memory to the scalar builders... but currently we haven't bothered.
We're not in trouble on compositional correctness because nothing's visible
until Build is called (whereas recursives have to be a lot stricter around this,
because if something gets validated and finished, it's already essential that
it now be unable to mutate out of the validated state, even if it's also not
yet 'visible').

Note that these remarks are for the `basicnode` package, but may also
apply to other implementations too (e.g., our codegen output follows similar
overall logic).

### nodestyles are implemented as exported concrete types

This is sorta arbitrary.  We could equally easily make them exported as
package scope vars (which happen to be singletons), or as functions
(which happen to have fixed returns).

These choices affect syntax, but not really semantics.

In any of these cases, we'd still end up with concrete types per style
(so that there's a place to hang the methods)... so, it seemed simplest
to just export them.
