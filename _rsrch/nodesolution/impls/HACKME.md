hackme
======

Design rational are documented here.

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

### about nodestyles

`ipld.NodeStyle` is a type that opaquely represents some information about how
a node was constructed and is implemented.  The general contract for what
should happen when asking a node for its style
(via the `ipld.Node.Style() NodeStyle` interface) is that style should be
effective instructions for how one could build a copy of that node, using
the same implementation details.

By example, if something was made as a `plainString` -- i.e.,
either via `String()` or via `Style__String{}.NewBuilder()` --
then its `Style()` will be `Style__String`.

Note there are also limits to this: if a node was built in a flexible way,
the style it reports later may only report what it is now, and not return
that same flexibility again.
By example, if something was made as an "any" -- i.e.,
via `Style__Any{}.NewBuilder()`, and then *happened* to be assigned a string value --
the resulting string node will carry a `Style()` property that returns
`Style__String` -- **not** `Style__Any`.

#### nodestyles meet generic transformation

One of the core purposes of the `NodeStyle` interface (and all the different
ways you can get it from existing data) is to enable the `traversal` package
(or other user-written packages like it) to do transformations on data.

// work-in-progress warning: generic transformations are not fully implemented.

When implementating a transformation that works over unknown data,
the signiture of function a user provides is roughly:
`func(oldValue Node, acceptableValues NodeStyle) (Node, error)`.
(This signiture may vary by the strategy taken by the transformation -- this
signiture is useful because it's capable of no-op'ing; an alternative signiture
might give the user a `NodeAssembler` instead of the `NodeStyle`.)

In this situation, the transformation system determines the `NodeStyle`
(or `NodeAssembler`) to use by asking the parent value of the one we're visiting.
This is because we want to give the update function the ability to create
any kind of value that would be accepted in this position -- not just create a
value of the same style as the one currently there!  It is for this reason
the `oldValue.Style()` property can't be used directly.

At the root of such a transformation, we use the `node.Style()` property to
determine how to get started building a new value.

#### nodestyles meet recursive assemblers

Asking for a NodeStyle in a recursive assembly process tells you about what
kind of node would be accepted in an `AssignNode(Node)` call.
It does *not* make any remark on the fact it's a key assembler or value assembler
and might be wrapped with additional rules (such as map key uniqueness, field
name expectations, etc).

(Note that it's also not an exclusive statement about what `AssignNode(Node)` will
accept; e.g. in many situations, while a `Style__MyStringType` might be the style
returned, any string kinded node can be used in `AssignNode(Node)` and will be
appropriately converted.)

Any of these paths counts as "recursive assembly process":

- `MapNodeAssembler.KeyStyle()`
- `MapNodeAssembler.ValueStyle()`
- `ListNodeAssembler.ValueStyle()`
- `MapNodeAssembler.AssembleKey().Style()`
- `MapNodeAssembler.AssembleValue().Style()`
- `ListNodeAssembler.AssembleValue().Style()`

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
