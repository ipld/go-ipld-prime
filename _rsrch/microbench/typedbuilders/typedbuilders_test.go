/*
	Designing how the user should interact with, build, and mutate data
	when using natively-typed/generated code for IPLD-Schema-based structures...
	is nontrivial.

	We have to account for the distinction between nullable and optional/absent fields.
	This is a celebrated feature and critical to our ability to round-trip binary data losslessly;
	it's also tricky to provide users with a smooth access to that distinction.
	There's very little prior art to follow, either -- we're one of first projects (if not 'the' first project)
	to try to directly handle the "+2 cardinality" problem (most others ignore it, shifting blame
	for inevitable consequent bugs on users, which is a 'solution' of sorts, but certainly not ours);
	in the long run, this is valuable, but at present it means we need some creativity.

	(Fasinatingly, by engaging directly with the "+1"'ing, we might also have stumbled upon
	one of the better solutions in history to the "how do I denote patches to data" problem:
	it's just another "+1" in the 'Maybe' enumeration seen below.  This makes golang's "zero"
	values work for us in a useful way: fail to mention a field?  It won't be patched.  Neat.)

	Here are some experiments, which may or may not outline a solution,
	but definitely should at least concretely explain the problem.
*/
package typedbuilders

import "testing"

var sink interface{}

// Stroct is an example of the kind of struct we're already generating from schemas.
type Stroct struct {
	f1 String  // `ipldsch:"String"`
	f2 *String // `ipldsch:"optional String"`
	f3 *String // `ipldsch:"optional nullable String"`
	f4 *String // `ipldsch:"nullable String"`

	f3__exists bool
}

// Stroct__Content is a proposed way we might have typed builders.
// This is the subject of discussion.
type Stroct__Content struct {
	F1 MaybeString
	F2 MaybeString
	F3 MaybeString
	F4 MaybeString
}

// String is a another type we're already generating from schemas.
// (It has to conform to the ipld 'Node' interfaces, so it's more than a plain golang string.)
type String string

// MaybeString is a proposed part of the way we might handle typed builders.
// A "Maybe{Foo}" type would be generated for all "{Foo}" types.
type MaybeString struct {
	Maybe  Maybe
	String String
}

// Maybe is a const used for an enum.  It would be defined in the ipld packages.
// (It's all in the same file here -- but the other types would be in codegen output; this one would not.)
type Maybe uint8

const (
	Maybe_Ignore = Maybe(0)
	Maybe_Value  = Maybe(1)
	Maybe_Null   = Maybe(2)
	Maybe_Absent = Maybe(3)
)

/*
	Next two funcs are 'Build' and 'Apply'.

	'Apply' "patches" an existing struct to have new values.

	Both are look complex -- but semantically simple, emit helpful errors, and are the product of codegen.
	Imagine doing these semantics *without* having these helpers -- the odds of user error would be *high*.

	These functions also mean the user can write code without knowing the details
	of how the internal implementation uses pointers or bools (or both) to represent optionality and nulls.
	This is important, because note how in this example, a pointer is created for any nullable or optional field?
	Right now, we've hardcoded that as the only option.  But it mean a heap alloc!
	What if we implement adjunct config that allows different implementation choices (e.g., more bools for less heap pointers)?
	This abstraction will make that possible without disrup.
*/

func (p Stroct__Content) Build() Stroct {
	v := Stroct{}
	switch p.F1.Maybe {
	case Maybe_Ignore:
		panic("must explicitly state all fields when building")
	case Maybe_Value:
		v.f1 = p.F1.String
	case Maybe_Null:
		panic("cannot assign null to nonnullable field")
	case Maybe_Absent:
		panic("cannot have absent nonoptional field")
	default:
		goto no
	}
	switch p.F2.Maybe {
	case Maybe_Ignore:
		panic("must explicitly state all fields when building")
	case Maybe_Value:
		v.f2 = &p.F2.String
	case Maybe_Null:
		panic("cannot assign null to nonnullable field")
	case Maybe_Absent:
		v.f2 = nil
	default:
		goto no
	}
	switch p.F3.Maybe {
	case Maybe_Ignore:
		panic("must explicitly state all fields when building")
	case Maybe_Value:
		v.f3 = &p.F3.String
		v.f3__exists = true
	case Maybe_Null:
		v.f3 = nil
		v.f3__exists = true
	case Maybe_Absent:
		v.f3 = nil
		v.f3__exists = false
	default:
		goto no
	}
	switch p.F4.Maybe {
	case Maybe_Ignore:
		panic("must explicitly state all fields when building")
	case Maybe_Value:
		v.f4 = &p.F4.String
	case Maybe_Null:
		v.f4 = nil
	case Maybe_Absent:
		panic("cannot have absent nonoptional field")
	default:
		goto no
	}
	return v
no:
	panic("invalid maybe")
}

func (p Stroct__Content) Apply(old Stroct) Stroct {
	v := Stroct{}
	switch p.F1.Maybe {
	case Maybe_Ignore: // pass
	case Maybe_Value:
		v.f1 = p.F1.String
	case Maybe_Null:
		panic("cannot assign null to nonnullable field")
	case Maybe_Absent:
		panic("cannot have absent nonoptional field")
	default:
		goto no
	}
	switch p.F2.Maybe {
	case Maybe_Ignore: // pass
	case Maybe_Value:
		v.f2 = &p.F2.String
	case Maybe_Null:
		panic("cannot assign null to nonnullable field")
	case Maybe_Absent:
		v.f2 = nil
	default:
		goto no
	}
	switch p.F3.Maybe {
	case Maybe_Ignore: // pass
	case Maybe_Value:
		v.f3 = &p.F3.String
		v.f3__exists = true
	case Maybe_Null:
		v.f3 = nil
		v.f3__exists = true
	case Maybe_Absent:
		v.f3 = nil
		v.f3__exists = false
	default:
		goto no
	}
	switch p.F4.Maybe {
	case Maybe_Ignore: // pass
	case Maybe_Value:
		v.f4 = &p.F4.String
	case Maybe_Null:
		v.f4 = nil
	case Maybe_Absent:
		panic("cannot have absent nonoptional field")
	default:
		goto no
	}
	return v
no:
	panic("invalid maybe")
}

/*
	Misc considerations:

	- Could make Stroct__Content use functions instead of exported fields.
	   Purely syntactic concern.
	   Probably makes no difference to perf: setters would be inlined.
	   Would let us call B.S. on invalid Maybe states closer to incident site.

	- It's wild, but what if we made everything include its own Maybe states all the time?
	   Would make it possible to have zero-values scream -- desirable.
	   Would raise questions of where to put funcs for making them -- undesirable if grows the reserved words lists.
	   Would make dense lists of ints twice the size -- wildly undesirable.
	   Might be interesting to roll around this thought and see if there are ways to deburr it of the undesirable parts.

	- For internal sanity, we should consider renaming the '__exists' fields to '__absent'.
	   It's six of one and a half dozen of the other, but look at the struct initializers.  The former makes you say two things at once.

	- Syntax: what if we generated a whole 'maybe' package?
	   Each type would exist in your output package, and its 'maybe' variant in the maybe subpackage.
	   Alternative: maybe the `Maybe_*` constants deserve a small package their own?
	   Worried about the sheer char count around using these things.

	- Obviously, we've given up on compile-time checking against bugs like trying to assign null to a nonnullable field.
	   Ideally, I'd love to have those checks.
	   The only way I can see to get them involves generating a combinatoric set of MaybeFoo types (undesirable),
	   or, making a lot more methods on the builder (worth consideration, but syntactically raises eyebrows).
*/

func BenchmarkStroctConstructionDirect(b *testing.B) {
	var v Stroct
	for i := 0; i < b.N; i++ {
		v = Stroct{
			f1: "foo",
			f2: strptr("bar"),
			f3: strptr("baz"), f3__exists: true,
			f4: strptr("woz"),
		}
	}
	sink = v
}

func strptr(s String) *String {
	return &s
}

func BenchmarkStroctConstructionBuilder(b *testing.B) {
	var v Stroct
	for i := 0; i < b.N; i++ {
		v = Stroct__Content{
			F1: MaybeString{Maybe_Value, "foo"},
			F2: MaybeString{Maybe_Value, "bar"},
			F3: MaybeString{Maybe_Value, "baz"},
			F4: MaybeString{Maybe_Value, "woz"},
		}.Build()
	}
	sink = v
}

// n.b. both of these are a little untrue, syntactically, because we didn't box the strings *enough*.
// the penalty visually is probably similar for both, but might actually change perf of number one(...?).
//
// worth considering: we might want to special-case any builders for the prelude types (such as 'String')
// such that they need less boilerplate to construct.
// it's one of those things that's technically a constact factor improvement, but on like 50% of cases, so aggregate impactful.
// and the example here seems to suggest it's actually very conveniently possible to do so,
// making it overall uncontentious that we then should.

/*
	Okay, let's now investigate the idea of having a bunch of builder methods.
	This might give us better brevity (we'll see),
	and also can give us better compile-time checks (and autocomplete) for which special modes are valid on a field.

	This could be offered in combination with the exported fields.
	Or, we could unexport the fields (and possibly even the maybeFoo types) by using this approach.

	Big question: can chaining syntax like this be optimized by the compiler enough to even be viable?
*/

func (p *Stroct__Content) Set_f1(v String) *Stroct__Content {
	p.F1 = MaybeString{Maybe_Value, v}
	return p
}
func (p *Stroct__Content) Set_f2(v String) *Stroct__Content {
	p.F2 = MaybeString{Maybe_Value, v}
	return p
}
func (p *Stroct__Content) Set_f3(v String) *Stroct__Content {
	p.F3 = MaybeString{Maybe_Value, v}
	return p
}
func (p *Stroct__Content) Set_f4(v String) *Stroct__Content {
	p.F4 = MaybeString{Maybe_Value, v}
	return p
}

// TODO/elided: more `Set_fN_Null()` and `Set_fN_Absent()` methods,
// as appropriate for the fields that should have them.

func BenchmarkStroctConstructionBuilderMethodic(b *testing.B) {
	var v Stroct
	for i := 0; i < b.N; i++ {
		b := Stroct__Content{}
		v = b.
			Set_f1("foo").
			Set_f2("bar").
			Set_f3("baz").
			Set_f4("woz").
			Build()
	}
	sink = v
}

/*
Okay, disassembly time.  What have we learned?

- The vast majority of our time spent in *any* of these is the heap allocation for the string pointer.
   This isn't a super huge surprise, honestly.
   This is why the "hide whether pointers are involved from users" thing is so dang important.

- Oh wow.  I think if we *do* keep pointers in play... I see why the builder is advantaged now.
   It looks like *the whole builder struct* is getting moved to the heap by the Build func,
    and so then getting an address for each string doesn't generate subsequent allocs.
	 This moves the alloc count from 3 to 1 in total.
	 In exchange, it's probably moving a larger amount of things to the heap.
   Yes, benchmem confirms (that probably would've been easier than reading asm...):

	BenchmarkStroctConstructionDirect-8             10000000    71.6 ns/op    48 B/op   3 allocs/op
	BenchmarkStroctConstructionBuilder-8            10000000    51.0 ns/op    96 B/op   1 allocs/op
	BenchmarkStroctConstructionBuilderMethodic-8    10000000    62.5 ns/op    96 B/op   1 allocs/op

   This has... interesting implications.
    Generally speaking, it's true that the *count* of allocs is the main cost to worry about;
     alloc *size* doesn't matter... in terms of time that's later going to cost via gc.
    But these pointers are going to stay around for the life of the object.
     That means for a struct with just one pointer field and lots of non-ptrs,
      the total memory *for its lifetime, not just to the builder's lifetime* will rise sharply.
     The words for the maybe-enum also will stick around!  Ouch!!
      (Scalars wouldn't encounter this themselves, so it's not death to the packed list-of-int case, but still, ouch.)
      Wow... this would actually mean... using the pointers in the struct but enums in the builder...
       is a doublespend in 100% of cases.
        If that doesn't prompt a rethink, I don't know what does.
         Jeez.

*/

/*
Other syntax ideas, just for fun:

	```
	Stroct__Content{
		f1: "foo",
		f2: "bar",
	}.Uncommonly(
		Stroct_f3_Absent,
		Stroct_f3_Null,
	).Build()
	```

	```
	Stroct__Content{
		f1: "foo",
		f2: "bar",
	}.Uncommonly(
		Absent("f3"),
		Null("f4"),
	).Build()
	```

Neither of these seems very pleasing, though.

The former adds tonnes of exported symbols.
The latter makes it possible to misspell fields entirely.
Neither provides any useful compile-time checks.

A third variant could address the misspelling problem, but is comically verbose:


	```
	Stroct__Content{
		f1: "foo",
		f2: "bar",
	}.Uncommonly(
		Absent(TypeSystem.Stroct.Fields("f3")),
		Null(TypeSystem.Stroct.Fields("f4")),
	).Build()
	```

... and also, actually didn't even solve the problem (unless we also add
generatation of field constants where we didn't previously have it), so, uh.

*/
