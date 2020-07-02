
### absent values

The handling of absent values is still not consistent.

Currently:

- reading (via accessors or iterators) yields `ipld.Absent` values for absent fields
- putting those ipld.Absent values via NodeAssembler.AssignNode will result in `ErrWrongKind`.
- *the recursive copies embedded in AssignNode methods don't handle absents either*.

The first two are defensible and consistent (if not necessarily ergonomic).
The third is downright a bug, and needs to be fixed.

How we fix it is not entirely settled.

- Option 1: keep the hostility to absent assignment
- Option 2: *require* explicit absent assignment
- Option 3: become indifferent to absent assignment when it's valid
- Option 4: don't yield values that are absent during iteration at all

Option 3 seems the most preferrable (and least user-hostile).
(Options 1 and 2 create work for end users;
Option 4 has questionable logical consistency.)

Updating the codegen to do Option 3 needs some work, though.

It's likely that the way to go about this would involve adding two more valid
bit states to the extended schema.Maybe values: one for allowAbsent (similar to
the existing allowNull), and another for both (for "nullable optional" fields).
Every NodeAssembler would then have to support that, just as they each support allowNull now.

I think the above design is valid, but it's not implemented nor tested yet.


### AssignNode optimality

The AssignNode methods we generate currently do pretty blithe things with large structures:
they iterate over the given node, and hurl entries into the assembler's AssignKey and AssignValue methods.

This isn't always optimal.
For any structure that is more efficient when fed info in an ideal order, we might want to take account of that.

For example, unions with representation mode "inline" are a stellar example of this:
if the discriminant key comes first, they can work *much, much* more efficiently.
By contrast, if the discriminant key shows up late in the object, it is necessary to
have buffered *all the other* data, then backtrack to handle it once the discriminant is found and parsed.

At best, this probably means iterating once, plucking out the discriminant entry,
and then *getting a new iterator* that starts from the beginning (which shifts
the buffer problem to the Node we're consuming data from).

Even more irritatingly: since NodeAssembler has to accept entries in any order
if it is to accept information streamingly from codecs, the NodeAssembler
*also* has to be ready to do the buffering work...
TODO ouch what are the ValueAssembler going to yield for dealing with children?
TODO we have to hand out dummy ValueAssembler types that buffer... a crazy amount of stuff.  (Reinvent refmt.Tok??  argh.)  cannot avoid???
TODO this means where errors arise from will be nuts: you cant say if anything is wrong until you figure out the discriminant.  then we replay everything?  your errors for deeper stuff will appear... uh... midway, from a random AssembleValue finishing that happens to be for the discriminant.  that is not pleasant.

... let's leave that thought aside: suffice to say, some assemblers are *really*
not happy or performant if they have to accept things in unpleasant orderings.

So.

We should flip all this on its head.  The AssignNode methods should lean in
on the knowledge they have about the structure they're building, and assume
that the Node we're copying content from supports random access:
pluck the fields that we care most about out first with direct lookups,
and only use iteration to cover the remaining data that the new structure
doesn't care about the ordering of.

Perhaps this only matters for certain styles of unions.


### sidenote about codec interfaces

Perhaps we should get used to the idea of codec packages offering two styles of methods:

- `UnmarshalIntoAssembler(io.Reader, ipld.NodeAssembler) error`
	- this is for when you have opinions about what kind of in-memory format should be used
- `Unmarshal(io.Reader) (ipld.Node, error)`
	- this is for when you want to let the codec pick.

We might actually end up preferring the latter in a fair number of cases.

Looking at this inline union ordering situation described above:
the best path through that (other than saying "don't fking use inline unions,
and if you do, put the discriminant in the first fking entry or gtfo") would probably be
to do a cbor (or whatever) unmarshal that produces the half-deserialized skip-list nodes
(which are specialized to the cbor format rather than general purpose, but we want that in this story)...
and those can then claim to do random access, thereby letting them take on the "buffering".
This approach would let the serialization-specialized nodes take on the work,
rather than forcing the union's NodeAssembler to do buffer at a higher level...
which is good because doing that buffering in a structured way at a higher level
is actually more work and causes more memory fragmentation and allocations.

Whew.

I have not worked out what this implies for multicodecs or other muxes that do compositions of codecs.


### enums of union keys

It's extremely common to have an enum that is the discrimant values of a union.

We should make a schema syntax for that.

We tend to generate such an enum in codegen anyway, for various purposes.
Might as well let people name it outright too, if they have the slightest desire to do so.

(Doesn't apply to kinded unions.)
