hackme: the schema compiler system
==================================


why is this like this
---------------------

These packages could've been arranged a lot of other ways.

We have a lot of constraints to solve for:

- we definitely want a `schema` package which provides (read-only) access to type system concepts
- we definitely want a `schema/dmt` package which knows how to interact with the IPLD Data Model tree forms and serializations of Schemas
- we want any instances of the `schema.Type*` information to be validated and correct by the time they're accessable
- validation of `schema.Type*` information is a graph problem -- and that includes cycles (e.g. imagine the definition of linked lists)
	- means it cannot trivally be validated with tree structures alone (meaning Schemas alone are *not* powerful enough to validate all conditions).
- we want `schema.Type*` information to have accessor methods which returns other `schema.Type*` handles (not just type names)
	- again means we have a graph problem, and need to do pointer constructions (which may include cycles) which cannot be trees alone.
- the existence of cycles means we cannot use `schema/dmt` types alone (they can't express cycles very well).
- it's important to remember that codegen'd nodes implement `schema.TypedNode`
	- and we want to use codegen to produce the `schema/dmt` package
- it's important to remember that `schema.TypedNode` includes a requirement for a function `Type() schema.Type`.
	- this scenario (an interface with a method refering to another interface in the same package) means we can't declare an equivalent interface in another package (which is otherwise sometimes a valid technique to simplify golang package dependency trees).
- it's desireable for the `schema.*` information to be immutable, simply because this would be reasonable.
- it would be ideal to use the `schema/dmt` types to do as much of the work as possible, since we need them regardless.
- it would be nice to have efficient programmatic constructors for `schema.Type*` info, since codegen outputs would like to use them.

(When we refer to `schema.Type*`, this means types like `schema.TypeStruct`, `schema.TypeMap`, `schema.TypeString`, `schema.TypeEnum`, etc.)

Conceptually, it seems like there's approximately three components in play:

- the `schema` types, being read-only information
- the `dmt` types, being what handles IPLD Data Model tree forms of the data
- the `compiler` logic, being what does graph-property checks and deeper validations on data (presumably usually fed from the `dmt` system) and produces `schema` data.

However, exactly how to relate these three concepts in golang code turns out a little nontrivial.

Here's a bunch of things we considered and/or tried,
one of which describes the code you're actually going to find around this hackme file now:

### two packages, schema wraps dmt and adds pointers

If we could have all the `schema.Type*` types simply be wrappers of the `schema/dmt.*` types,
and just decorate them with pointers (so they can have helpful accessor methods to related types)
and a single constructor function that just takes a `schema/dmt` tree and does all validation at once,
that'd produce a lot of great results:

- we would trivially have the whole immutability thing sorted out, because codegen would give us that in the `dmt.*` types.
- programmatic construction would always be "just use the dmt" -- fine!
- all of our validation code could be written against the `dmt`, which is just quite direct and elegant and guarantees no barrier to usability on schema documents
- generally, minimal redundancy -- no need to write any near-duplicate golang types that have the same logical content as the `dmt.*` codegen'd types.

There'd be one small downside to all this, which is that in practice,
the need for those wrapper types would mean another allocation per type when compiling.
This is very unlikely to be of any real consequence, though.

And, okay, we'd still be writing a lot of helpful accessor methods that wrap and prettify access to the `dmt.*` types just a bit.
Still, that's okay.  That's something we probably want to do regardless.

So this seems overall pretty good!

But... no.  It can't be done, because this would generate an import cycle, which is an instant show-stopper.
The `schema/dmt`, being codegen output, and typed nodes, contains many types that implement `schema.TypedNode`, and thus refer to `schema.Type`.
That means `schema` has no chance of being able to depend on structures from `schema/dmt` without an import cycle problem.
Shoot.

### three packages, schema and compiler and dmt

This seemed like it would've been one of the most legible options.
Compared to the "schema wraps dmt" approach, we'd probably be writing a lot of near-duplicate types,
but at least they'd be well-organized into task-specific packages.

Unfortunately, it doesn't work out well.

Golang immutability is what torpedoes this, more or less.
We would have to have all the `Type*` implementations in the `compiler` package,
because that's where they're written to from,
but at the same time, have all the reading behaviors specified by interfaces in the `schema` package.

Mostly, this would be workable, if fairly highly redundant
(because we'd have all these interfaces wrapping exactly one concrete type).

However, there's one thing significantly incorrect in this outcome:
we wouldn't be able to make the `schema.Type` interface "closed".
We don't actually want an interface for that type at all -- we really want a sum type.
Since golang doesn't have sum types, we hack it, with an interface.
But it's really not desirable nor correct at all to allow other implementations of this interface,
because we regularly write code that explicitly assumes it can enumerate all the concrete types that might inhabit that interface,
and we consider such code to be correct.

### two packages, schema contains compiler

**This is what we actually did**.

This approach threatened to make the godocs for the schema package fairly illegible,
because it might've resulted in a mishmash of mutable types used for compilation,
and then read-only types that would be used when not working with compilation --
only one of which end-users would be expected to interact with.
This was addressed by doing a fairly large amount of work to attach all the compilation-relevant
methods to the `Compiler` type (including constructors for any intermediate values, which otherwise simply could've been exported struct types),
thus keeping them in one place.

This also comes out hugely, hugely redundant in code needed.
There are tons of near-duplicate types which have the same logical content as the `dmt.*` types,
just slightly flipped about.  (See `schema.TypeStruct` and `dmt.TypeStruct`, for example.)
Re-building the immutability guarantees that we get for free from the codegen'd `dmt.*` types also
took a rather staggering amount of highly redundant code.

Since we were unable to used the codegen'd types, and we use golang native maps instead in some places,
this tends to make our logic less deterministic in its evaluation order.
We've mostly ignored this, but you may notice it in situations such as validation of a schema
that has more than one thing wrong within the same union representation specification, for example,
since those happen to use golang maps in the (mostly redundant) golang structures.

Our validation logic ended up written against the `schema.Type*` types rather than the `dmt.*` types.
It was necessary to do it this way because now we have the possibility of compiling schema types without going through `dmt.*` at all.
Since there's still a path from the `dmt.*` types to here, we didn't _lose_ any important features
like the ability to run full validations on a schema document, but it's become a bit indirected.
And you can see some scar tissue in the `schema.Compiler` methods, particularly in how it handles errors,
where the behavior is still _functionally_ dictated by the need to be useful for `dmt` document validation.

And we can't *get* `dmt.*` information trivially from `schema.Type*` values.
(You could write a function for it, I'm sure.  But: "can't _trivally_".)
Maybe this isn't all that important.  But it seems unfortunate.

But for all those drawbacks: it works.

(This whole section has a double duty: it also serves as a nice list
of cool features you get for free when using our codegen.)

### two packages, compiler is with dmt

Doesn't really fly for the same reasons as three packages.
(Separating the read-only interfaces from the writable implementations doesn't really come out well.)

### three packages, schema aliasing from compiler heavily

This could've worked... it just seemed like it might be reaching into slightly-too-clever territory in its use of type aliases.

More practically, it would've resulted in user-facing messages with weird content
in the event users got compiler errors relating to something they're doing with the readable `schema.*` types --
golang's error messages use the canonical type name, not its aliases,
which means it would be talking about `compiler.*` types, when users would expect it to be talking about `schema.*` types.
This could be potentially confusing to users.

### ... could we have done this yet differently, by breaking some of the constraints entirely?

Yeah, maybe.

For example, if we removed the requirement for `schema.TypedNode` to have `Type() schema.Type` accessor,
that would've removed a big import cycle problem, and made different solutions viable.
(This would seem to have pretty huge knock-on effects, though, reaching well beyond these packages.)
