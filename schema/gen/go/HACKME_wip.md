misc notes during refresh
=========================

(This document will be deleted as soon as this update work cycle is complete;
the contents are a bit "top-of-the-head".)

- no major regrets about our stringoid approach
- use more "macros" for reusable regions
- don't bother making custom functions per symbol pattern
	- that is: conjoining "__Foo" as string in the template is fine; don't turn it into a shell game.
- consistently maintain separatation of **symbol** from **name**
	- **symbol** is what is used in type and function names in code.
	- **name** is the string that comes from the schema; it is never modified nor overridable.
- symbol processing all pipes through an adjunct configuration object.
	- we make this available in the templates via a funcmap so it's available context-free as a nice tidy pipe syntax.
- the typegen/nodegen/buildergen distinctions and how they turn into a 5-some still checks out.
- the embeddable starter kits per kind are still a perfectly reasonable idea.
	- ... though they also get to just shell out to the 'node/mixins' package a ton now, which is nice.
- we're gonna need to break down the methods for some types more than others.
	- namely, structs and unions... are gonna want an invok of template per field / member, i think.
	- there's no particularly useful way to expose that to the top level nodegen interface, afaict.  just do it concretely.
- generally do a pass on all templates to use consistent abbreviations for variable names.
- use a lot more pointers in method types, according to the understandings of the low cost of internal pointers.


corners needing mention in docs
-------------------------------

- iterating a type-level node with optional fields can yield the field and a maybe containing absent.
	- ...which is funny because if you feed that into a type-level builder, it doesn't like that absent.  it wants you to not feed that field.
	- alternatively, accepting explicit puts of absent: worse: we'd have to keep state track that it's been put, but to none, and reject future puts.
		- REVIEW: maybe this isn't as bad as first thought.  I think we end up with that state bit anyway.
		- REVIEW: maybe we should turn this on its head entirely: would it be clearer and more consistent if building a struct without explicitly assigning undef to any optional fields is actually *rejected* when using the type-level assemblers?
	- alternatively, not yielding them on iterate: worse: generic printer for structs would end up not reporting fields, and that would be both wrong and hard to hack your way out of without writing a metric ton more code that inspects the type info, which would ruin the point of the monomorphic methods in the first place to a much higher degree than this need to handle undefined/absent does.


underreviewed
-------------

### the maybe situation

- just one MaybeT per T still seems the 'least wrong' to me.
- I'm okay with that name.
	- we can support symbol customization for that too I guess?
	- feedback on other parts of the project seems to be that people respond negatively to seeing double-underscores.
	- but yeah, by default, if your schema has types named "Maybe*" in it, uh, you might be in for some bumps.
- do they implement the full damn Node interface?
	- some things seem fine about this
		- certainly makes sense to be able to 'IsNull' on it like any other Node.
		- if it's embeded, we can return an internal pointer to it just fine, so there's no obvious runtime perf reason not to.
	- why not?
		- it's another type with a ton of methods.  or two, or four.
			- may increase the binary size.  possibly by a significant constant multiplier.
		- would it have a 'Type' accessor on it?  if so, what does it say?
	- simply not sure how useful this is!
		- istm one will often either be passing the MaybeT to other speciated functions, or, fairly immediately de-maybing it.
			- if this is true, the number of times anyone wants to treat it as a Node are near zero.
		- i also don't see any reason to want to support giving a _MaybeT__Assembler (!) to something like unmarshalling.
			- if you have a null in the root, you can describe this with a kinded union, and probably would be better off for it.
			- if you have an undefined in the root... no.  no, you just don't.  that's called "EOF".
				- okay, this actually kinda gives me pause, because "EOF" handling does seem to be a spawner of darkness in practice.
			- does it show up usefully in the middle of a tree?
				- not sure
				- there's generally a _P_ValueAssembler type involved there anyway, which can exhibit the necessary traits.
- we should embed the MaybeT.
	- it's alternatively possible to bitpack them into one word at the top of a struct, but...
		- the complexity of this is high
		- it would be exposed to anyone who writes addntl code in-package, which is asking for errors
		- i don't think there are many antecedents for other langs doing this (e.g. does Rust?  i don't think so)
			- if they do, it's with pragmas.  (possibly because they worry about mem layout for cross-lang compat, though, which... i don't think we do.  uh.  much.  uff.)
		- the only thing this buys us is *slightly* less resident memory size
			- and long story short: no one in the world at large appears to care.
- fields of the maybe types definitely need to be unexported.
	- now that we have a way to syntactically prevent zero-value creation and thus prevent dodging validation, want to maintain that.
		- another type from the same package which you *can* create zero-values of and which embeds the shielded one... breaks the shield.
	- may also just want to shield it entirely.
		- might be easier than having its enum component include an 'invalid' state which can kick you in the shins at runtime.
			- roughly the same decision we previously considered for other types (aside from the word for such an enum already existing here).
			- if we make the zero state be 'absent' rather than introducing a new 'invalid' state, that makes this less bad.
		- can't think of significant arguments against this.  returning pointers to maybes embedded in larger structures does seem to be the way to go anyway.
- it's an adjunct parameter whether the value field in the MaybeT is a pointer; may vary for each T.
	- need to sometimes use pointers: for cycle-breaking!
	- want customizability of optimism (oversize-allocation versus alloc-count-amortization).
		- is it possible to want different Maybe implementation strategies in different areas?  Perhaps, but hopefully vanishingly unlikely in practice.  Do not want to support this; complexity add high.
	- thinking this may want to default to useptr=yes.  less likely to generate end user surprise.  can opt into noptr fastness when/where you know it works.
	- there are a couple very predictable places where we *do* want useptr=no, also.
		- strings!  and other scalars.  essentially *never* useful to have a ptr to those inside their MaybeT.
- MaybeT should be an alias of `*_T_Maybe`.
	- Sure, sometimes we could get along fine by passing them around by value.
		- But sometimes -- when they contain a large value and don't use a ptr internally -- we don't want to.
			- And we don't want that to then embroil into *this* MaybeT requiring you to *handle* it externally as a ptr at the same time when some MaybeT2 requires the opposite.
	- Coincidentally, removes a fair number of conditionalizations of '&' when returning fields.
		- Not a factor in the decision, but a nice bonus.
		- A bunch of remaining conditionals in templates also all consistently move to the end of their area, which is... sort of pleasing.
	- This means we get to "shield" it entirely.
		- Not a factor in the decision (perhaps surprisingly).  But doesn't hurt.
			- The zero values of the maybe type already didn't threaten to reveal zeros of the T, as long as the field was private and the zero state for maybes isn't 'Maybe_Value'.
		- We could also just use an exported `*MaybeT`.  But I see no reason to favor this.
	- Since we've already resolved to embed MaybeT values (rather than bitpack their states together), there's no cost to this.
		- When speciated methods return a `MaybeT`, it's an internal pointer to existing memory.
		- When creating a populated maybe...
			- If it's a useptr=yes, this might mean you get two allocs... except I think the outer one should optimize out:
				- the conversion into maybe should be inlinable, and that should make the whole thing escape-analyzable to get rid of the maybe alloc?
			- If it's a useptr=no, construction is probably going to involve a big ol' memcopy.
				- This is unfortunate, but I can't think of a way around it short of doing the full dance of a builder.
				- This probably won't be an issue in practice, if we're imagining useptr=no is only likely to be used on fairly small values.
			- We can also just try to avoid the freestanding creation of maybes entirely.
				- The ideal look and field of speciated builders has not been determined.  SetAbsent, etc, methods are certainly in bounds.
	- Structs and lists can certainly amortize these.
		- Maps would turn tricky, except we're already happy with having an internal slice in those too.  So, it's covered.
			- Interestingly, maps also have the ability to easily store all of (null,absent,value{...}) in them, without extra space.
				- But considering the thing about slices-in-maps, that may not be relevant.  We want quick linear iterator-friendly reads of Maybe state too.
- We discarded the idea of sometimes collapsing a maybe with one mode (e.g., 'optional' or 'nullable' but not 'optional nullable') down to one pointer and no struct.
	- mostly because it just fell off the plate of considerations.  we could still try to do this, and save a word of memory in some cases.
	- if we did this, complexity increases.
		- take the already existing total complexity of maybes with-and-without-useptr and for the three combos... and double it again.  ow.
	- if we did this, the typedef of `type MaybeT = *_T__Maybe` wouldn't work anymore.
		- satisfying that type would mean any case that uses a pointer without embedded maybe struct would now need allocations.
		  - this would be unacceptably high performance cost, so, we'd then end up needed to pursue more type signatures to avoid those costs...
		- this might actually be a pretty hard-stop reason not to pursue this possibilty.  exposing *more* user-facing complexity in this area, in the form of more types in the golang code the user has to reason about despite having no equivalent in the IPLD Schema, is strongly undesirable.
