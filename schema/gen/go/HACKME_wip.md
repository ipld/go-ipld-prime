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
	- we don't _currently_ support symbol customization separate from name, but likely will need to someday.
- the typegen/nodegen/buildergen distinctions and how they turn into a 5-some still checks out.
- the embeddable starter kits per kind are still a perfectly reasonable idea.
	- ... though they also get to just shell out to the 'node/mixins' package a ton now, which is nice.
- we're gonna need to break down the methods for some types more than others.
	- namely, structs and unions... are gonna want an invok of template per field / member, i think.
	- there's no particularly useful way to expose that to the top level nodegen interface, afaict.  just do it concretely.
- for symbols, have "{Situation}Symbol[A/B]"
	- *don't* do a "titlize" function; that works but is lazy and will obstruct rather than ease customizabilty later.
	- ... what are A and B?  Upper/Lower?  Exported/Unexported?
	- not all things have this: type symbols effectively never appear in lowercase.
	- things that do: fields, for example.
		- interesting because `Field{{ .FieldSymbolUpper }}` is something we'll probably see, and it's not *exactly* the case that the capital letter effects export in this case.
- generally do a pass on all templates to use consistent abbreviations for variable names.
- use a lot more pointers in method types, according to the understandings of the low cost of internal pointers.


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
		- can't think of significant arguments against this.  returning pointers to maybes embedded in larger structures does seem to be the way to go anyway.
- it's an adjunct parameter whether the value field in the MaybeT is a pointer; may vary for each T.
	- need to sometimes use pointers: for cycle-breaking!
	- want customizability of optimism (oversize-allocation versus alloc-count-amortization).
		- is it possible to want different Maybe implementation strategies in different areas?  Perhaps, but hopefully vanishingly unlikely in practice.  Do not want to support this; complexity add high.
	- thinking this may want to default to useptr=yes.  less likely to generate end user surprise.  can opt into noptr fastness when/where you know it works.

### where does any name override info live

- we've currently hacked the idea of type name overrides into place but mostly because they're at the top and don't fractilate.
	- this doesn't generalize to fields.
- putting it in the schema.Type area makes the easiest to access, but seems weird because it sure doesn't have anything to do with not-codegen inspection.
- is it the case that any/all "adjunct" info shouldn't survive into the schema.Type interface?
	- probably
- seems like we're gonna have to build a tree of adjunct info and pipe type and field references into functions on that after all.
	- whee.
	- although note this still doesn't mean a single thing will appear in the funcmap.  this isn't a repeat of the 'titlize' mistake.
	- FIXME yeah this is needed ALL over.  when dealing with fields in a struct: you also need to get their symbol override, if any.
		- this is probably going to require schema.StructField to sprout a parent ptr ...wellll maybe, alternatives exist
		- `{{ .Adjutant.FieldSymbol($type, $field) }}` ?
		- `{{ $field | $adjutant.SymbolForField }}` ?
