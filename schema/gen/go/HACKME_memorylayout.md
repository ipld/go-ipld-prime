about memory layout
===================

Memory layout is important when designing a system for going fast.
It also shows up in exported types (whether or not they're pointers, etc).

For the most part, we try to hide these details;
or, failing that, at least make them appear consistent.
There's some deeper logic required to *pick* which way we do things, though.



Prerequisite understandings
---------------------------

The following headings contain brief summaries of information that's important
to know in order to understand how we designed the IPLD data structure
memory layouts (and how to tune them).

Most of these concepts are common to many programming languages, so you can
likely skim those sections if you know them.  Others are fairly golang-specific.

### heap vs stack

The concept of heap vs stack in Golang is pretty similar to the concept
in most other languages with garbage collection, so we won't cover it
in great detail here.

The key concept to know: the *count* of allocations which are made on
the heap significantly affects performance.  Allocations on the heap
consume CPU time both when made, and later, as part of GC.

The *size* of the allocations affects the total memory needed, but
does *not* significantly affect the speed of execution.

Allocations which are made on the stack are (familiarly) effectively free.

### escape analysis

"Escape Analysis" refers to the efforts the compiler makes to figure out if some
piece of memory can be kept on the stack or if it must "escape" to the heap.
If escape analysis finds that some memory can be kept on the stack,
it will prefer to do so (and this is faster/preferable because it both means
allocation is simple and that no 'garbage' is generated to collect later).

Since whether things are allocated on the stack or the heap affects performance,
the concept of escape analysis is important.  The details (fortunately) are not:
For the purposes of what we need to do in in our IPLD data structures,
our goal with our code is to A) flunk out and escape to heap
as soon as possible, but B) do that in one big chunk of memory at once
(because we'll be able to use [internal pointers](#internal-pointers)
thereafter).

One implication of escape analysis that's both useful and easy to note is that
whether or not you use a struct literal (`Foo{}`) or a pointer (`&Foo{}`)
*does not determine* whether that memory gets allocated on the heap or stack.
If you use a pointer, the escape analysis can still prove that the pointer
never escapes, it will still end up allocated on the stack.

Another way to thing about this is: use pointers freely!  By using pointers,
you're in effect giving the compiler *more* freedom to decide where memory resides;
in contrast, avoiding the use of pointers in method signitures, etc, will
give the compiler *less* choice about where the memory should reside,
and typically forces copying.  Giving the compiler more freedom generally
has better results.

**pro-tip**: you can compile a program with the arguments `-gcflags "-m -m"` to
get lots of information about the escape analysis the compiler performs.

### embed vs pointer

Structs can be embeded -- e.g. `type Foo struct { field Otherstruct }` --
or referenced by a pointer -- e.g. `type Foo struct { field *Otherstruct }`.

The difference is substantial.

When structs are embedded, the layout in memory of the larger struct is simply
a concatenation of the embedded structs.  This means the amount of memory
that structure takes is the sum of the size of all of the embedded things;
and by the other side of the same coint, the *count* of allocations needed
(remember! the *count* affects performance more than the *size*, as we briefly
discussed in the [heap-vs-stack](#heap-vs-stack) section) is exactly *one*.

When pointers are used instead of embedding, the parent struct is typically
smaller (pointers are one word of memory, whereas the embedded thing can often
be larger), and null values can be used... but if fields are assigned to some
other value than null, there's a very high likelihood that heap allocations
will start cropping up in the process of creating values to take pointers
to before then assigning the pointer field!  (This can be subverted by
either [escape analysis](#escape-analysis) (though it's fairly uncommon),
or by [internal pointers](#internal-pointers) (which are going to turn out
very important, and will be discussed later)... but it's wise to default
to worrying about it until you can prove that one of the two will save you.)

When setting fields, another difference appears: a pointer field takes one
instruction (assuming the value already exists, and we're not invoking heap
allocation to get the pointer!) to assign,
whereas an embedded field generally signifies a memcopy, which
may take several instructions if the embedded value is large.

You can see how the choice between use of pointers and embeds results
in significantly different memory usage and performance characteristics!

(Quick mention in passing: "cache lines", etc, are also potential concerns that
can be addressed by embedding choices.  However, it's probably wise to attend
to GC first.  While cache alignment *can* be important, it's almost always going
to be a winning bet that GC will be a much higher impact concern.)

It is an unfortunate truth that whether or not a field can be null in Golang
and whether or not it's a pointer are two properties that are conflated --
you can't choose one independently of the other.  (The reasoning for this is
based on intuitions around mechanical sympathy -- but it's worth mentioning that
a sufficiently smart compiler *could* address both the logical separation
and simultaneously have the compiler solve for the mechanical sympathy concerns
in order to reach good performance in many cases; Golang just doesn't do so.)

### interfaces are two words and may cause implicit allocation

Interfaces in Golang are always two words in size.  The first word is a pointer
to the type information for what the interface contains.  The second word is
a pointer to the data itself.

This means if some data is assigned into an interface value, it *must* become
a pointer -- the compiler will do this implicitly; and this is the case even if
the type info in the first word retains a claim that the data is not a pointer.
In practice, this also almost guarantees in practice that the data in question
will escape to the heap.

(This applies even to primitives that are one word in size!  At least, as of
golang version 1.13 -- keep an eye on on the `runtime.convT32` functions
if you want to look into this further; the `mallocgc` call is clear to see.
There's a special case inside `malloc` which causes zero values to get a
free pass (!), but in all other cases, allocation will occur.)

Knowing this, you probably can conclude a general rule of thumb: if your
application is going to put a value in an interface, and *especially* if it's
going to do that more than once, you're probably best off explicitly handling
it as a pointer rather than a value.  Any other approach wil be very likely to
provoke unnecessary copy behavior and/or multiple unnecessary heap allocations
as the value moves in and out of pointer form.

(Fun note: if attempting to probe this by microbenchmarking experiments, be
careful to avoid using zero values!  Zero values get special treatment and avoid
allocations in ways that aren't general.)

### internal pointers

"Internal pointers" refer to any pointer taken to some position in a piece
of memory that was already allocated somewhere.

For example, given some `type Foo struct { a, b, c Otherstruct }`, the
value of `f := &Foo{}` and `b := &f.b` will be very related: they will
differ by the size of `Otherstruct`!

The main consequence of this is: using internal pointers can allow you to
construct large structure containing many pointers... *without* using a
correspondingly large *count of allocations*.  This unlocks a lot of potential
choices for how to build data structures in memory while minimizing allocs!

Internal pointers are not without their tradeoffs, however: in particular,
internal pointers have an interesting relationship with garbage collection.
When there's an internal pointer to some field in a large struct, that pointer
will cause the *entire* containing struct to be still considered to be
referenced for garbage collection purposes -- that is, *it won't be collected*.
So, in our example above, keeping a reference to `&f.b` will in fact cause
memory of the size of *three* `Otherstruct`s to be uncollectable, not one.

You can find more information about internal pointers in this talk:
https://blog.golang.org/ismmkeynote

### inlining functions

Function inlining is an important compiler optimization.

Inlining optimizes in two regards: one, can remove some of the overhead of
function calls; and two, it can enable *other* optimizations by getting the
relevant instruction blocks to be located together and thus rearrangable.
(Inlining does increase the compiled binary size, so it's not all upside.)

Calling a function has some fixed overhead -- shuffling arguments from registers
into calling convention order on the stack; potentially growing the stack; etc.
While these overheads are small in practice... if the function is called many
(many) times, this overhead can still add up.  Inlining can remove these costs!

More interestingly, function inlining can also enable *other* optimizations.
For example, a function that *would* have caused escape analysis to flunk
something out to the heap *if* that function as called was alone... can
potentially be inlined in such a way that in its contextual usage,
the escape analysis flunking can actually disappear entirely.
Many other kinds of optimizations can similarly be enabled by inlining.
This makes designing library code to be inline-friendly a potentially
high-impact concern -- sometimes even more so than can be easily seen.

The exact mechanisms used by the compiler to determine what can (and should)
be inlined may vary significantly from version to version of the Go compiler,
which means one should be cautious of spending too much time in the details.
However, we *can* make useful choices around things that will predictably
obstruct inlining -- such as [virtual function calls](#virtual-function-calls).
Occasionally there are positive stories in teasing the inliner to do well,
such as https://blog.filippo.io/efficient-go-apis-with-the-inliner/ (but these
seem to generally require a lot of thought and probably aren't the first stop
on most optimization quests).

### virtual function calls

Function calls which are intermediated by interfaces are called "virtual"
function calls.  (You may also encounter the term "v-table" in compiler
and runtime design literature -- this 'v' stands for "virtual".)

Virtual function calls generally can't be inlined.  This can have significant
effects, as described in the [inlining functions](#inlining-functions) section --
it both means function call overhead can't be removed, and it can have cascading
consequences by making other potential optimizations unreachable.



Resultant Design Features
-------------------------

### concrete implementations

We generate a concrete type for each type in the schema.

Using a concrete type means methods on it are possible to inline.
This is interesting because most of the methods are "accessors" -- that is,
a style of function that has a small body and does little work -- and these
are precisely the sort of function where inlining can add up.

There is one one downside to using an exported concrete type (rather than
keeping it unexported and hidden behind and exported interface):
it means any code external to the package can produce Golang's natural "zero"
for the type.  This is problematic because it's true even if the Golang "zero"
value for the type doesn't correspond to a valid value.
This is an unfortunate but overall practical tradeoff.

### embed by default

Embedded structs amortizes the count of memory allocations.
This addresses what is typically our biggest concern.

The increase in size is generally not consequential.  We expect most fields
end up filled anyway, so reserving that memory up front is reasonable.
(Indeed, unfilled fields are only possible for nullable or optional fields
which are implemented as embedded.)

Assignment into embedded fields may have the cost of a memcopy.
(By contrast, if fields were pointers, assigning them would be cheap...
though at the same time, we would've had to pay the allocation cost, elsewhere.)
However, combined with (other tricks)[#child-nodebuilders-point-into-embedded-fields],
a shortcut becomes possible: if we at some point used shared memory as the
scratch space for the child nodebuilder... and it's since been finalized...
and that very same pointer (into ourselves!) is now being assigned to us...
we can cheaply detect that and fastpath it.  (This sounds contrived, but it's
actually the common case.)

### nullable and optional struct fields embed too

TODO intro

There is some chance of over-allocation in the event of nullable or optional
fields.  We support tuning that via adjunct configuration to the code generator
which allows you to opt in to using pointers for fields; choosing to do this
will of course cause you to lose out on alloc amortization features in exchange.

TODO also resolve the loops note, at bottom

### nodebuilders point to the concrete type

We generate NodeBuilder types which contain a pointer to the type they build.

This means a single NodeBuilder and its produced Node will require
**two** allocations -- one for the NodeBuilder, and a separate one for the Node.

An alternative would be to embed the concrete Node value in the NodeBuilder,
and return a pointer to when finalizing the creation of the Node;
however, because due to the garbage collection semantics around
[internal pointers](#internal-pointers), such a design would cause the entirety
of the memory needed in the NodeBuilder to remain uncollectable as long as
completed Node is reachable!  This would be an unfortunate trade --
we can do better, and will... via [racking builders](#racking-builders).

### child nodebuilders point into embedded fields

TODO this is the critical secret sauce

### racking builders

(This where things start to get decidedly less-than-obvious.)

After generating the NodeBuilder for each type, we **additionally** generate
a "racker" type.  This "racker" is a struct which embeds the NodeBuilder...
and the racker (and thus NodeBuilder) for each of the fields within a struct.
This lets us amortize the allocations for all the *builders* in the same way
as embedding in the actual value structs let us amortized allocations there.

With racking builders, we can amortize all the allocations of working memory
needed for a whole family of NodeBuilders... **and** amortize all the
allocations for the value structures into a second allocation...
and that's it, it's just those two.  Further more, the separation means that
once the construction of the Node is done, we can drop all the NodeBuilder
memory and expect it to be immediately garbage collectable.  Win!

The code for this gets a little complex, and the result also carries several
additional limitations to the usage, but it does keep the allocations finite,
and thus makes the overall performance fast.

### visibility rules

It's perfectly fine to let builders accumulate mutations... right up until
the moment where a Node is returned.

(While it's less than ideal that different nodebuilders might interact with
each other... it's technically not a violation of terms: the one final
concern is whether or not Node immutablity is violated.  Experiencing
spooky-action-at-a-distance between NodeBuilder instances is irrelevant.)

So, we reach the following rules:

- when a NodeBuilder.Build method returns a Node, that memory must be frozen:
	- that NodeBuilder of course sets its target pointer to nil, jamming itself;
	- no other set methods on the *parent* NodeBuilder may assign to that field;
	- and the *parent* NodeBuilder may never return another child NodeBuilder for that field.

This set of rules around visibility lets us do amortized allocations
of a big hunk of working memory, and still comply with the familiar
small-pieces-first creation model of the NodeBuilder interface
by returning piecemeal read-only pointers into that big amortized memory hunk.

In order to satisfy these rules (namely, ensuring we never return a NodeBuilder
that addresses memory that's already been frozen) -- and do so without
consuming linearly more memory to track it! -- maps and lists end up with some
notable limitations:

- Lists can only be appended linearly, not populated in free order.
  (This means we can condense the 'isFrozen' attribute to an int offset.)
- Maps can only build one new value at a time.
  (Even though we usually handle map contents as pointers, being able to build
  more than one value at a time would require unknown amounts of memory for
  the any NodeBuilder state after the first, which is undesirable.)
- Structs need no special handling -- they can still be regarded in any order.
  (We know how much memory we need at compile time, so we can swallow that.)

### natively-typed methods in addition to the general interface

We generate two sets of methods: **both** the general interface methods to
comply with Node and NodeBuilder interfaces, **and** also natively-typed
variants of the same methods (e.g. a `Lookup` method for maps that takes
the concrete type key and returns the concrete type value, rather than
taking and returning `Node` interfaces).

While both sets of methods can accomplish the same end goals, both are needed.
There are two distinct advantages to natively-typed methods;
and at the same time, the need for the general methods is system critical.

Firstly, to programmers writing code that can use the concrete types, the
natively-typed methods provide more value in the form of compile-time type
checking, autocompletion and other tooling assist opportunities, and
less verbosity.

Secondly, natively-typed funtions on concrete types can be higher performance:
since they're not [virtual function calls](#virtual-function-calls), we
can expect [inlining](#inlining-functions) to work.  We might expect this to
be particularly consequential in builders and in accessor methods, since these
involve numerous calls to methods with small bodies -- precisely the sort of
situation that often substantially benefits from inlining.

At the same time, it goes without saying that we need the general Node and
NodeBuilder interfaces to be satisfied, so that we can write generic library
code such as reusable traversals, etc.  It is not possible to satisfy both
needs with a single set of methods with the Golang typesystem.



Amusing Details and Edge Cases
------------------------------

### looped references

// who's job is it to detect this?
// the schema validator should check it...
// but something that breaks the cycle *there* doesn't necessarily do so for the emitted code!  aggh!
//  ... unless we go back to optional and nullable both making ptrs unconditionally.



Learning more (the hard way)
----------------------------

If this document doesn't provide enough information for you,
you've probably graduated to the point where doing experiments is next.  :)

Prototypes and research examples can be found in the
`go-ipld-prime/_rsrch/` directories.
In particular, the "multihoisting" and "nodeassembler" packages are relevant,
containing research that lead to the drafting of this doc,
as well as some partially-worked alternative interface drafts.

Tests there include some benchmarks (self-explanitory);
some tests based on runtime memory stats inspection;
and some tests which are simply meant to be disassembled and read thusly.

Compiler flags can provide useful insights:

- `-gcflags '-S'` -- gives you assembler dump.
	- read this to see for sure what's inlined and not.
	- easy to quickly skim for calls like `runtime.newObject`, etc.
	- often critically useful to ensure a benchmark hasn't optimized out the question you meant to ask it!
	- generally gives a ground truth which puts an end to guessing.
- `-gcflags '-m -m'` -- reports escape analysis and other decisions.
   - note the two m's -- not a typo: this gives you info in stack form,
	  which is radically more informative than the single-m output.
- `-gcflags '-l'` -- disables inlining!
	- useful on benchmarks to quickly detect whether inlining is a major part of performance.

These flags can apply to any command like `go install`... as well as `go test`.

Profiling information collected from live systems in use is of course always
intensely useful... if you have any on hand.  When handling this, be aware of
how data-dependent performance can be when handling serialization systems:
different workload content can very much lead to different hot spots.

Happy hunting.
