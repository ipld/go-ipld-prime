package wow

type String string
type MaybeString struct {
	exists, null bool
	value        String
}

func (x MaybeString) Must() String {
	if !x.exists || x.null {
		panic("unbox of a maybe rejected")
	}
	return x.value
}

// Which bits are which?
// We probably precompute a table based on the struct as a whole.
// One might suppose a bit should be on to indicate 'present' so that in inaction the default doesn't become zero values?
// Ah, but this conflicts with another goal: it would make you say two things at once to create a value; that's silly.
// Honestly.  This whole set of questions sucks and is stupid.
type _Foo__bitfield uint32

type Foo struct {
	d Foo__Content
}

type Foo__Content struct {
	Modifiers _Foo__bitfield // shoot -- this symbol has to be exported too, then?  or can we do enough with builders?
	// even if you kick the can down the road with builder funcs... your choices are put all symbols in munges of func name,
	//  or else you'll still just end up with more constants needing to be exported somewhere.

	F1 String // `ipldsch:"String"`
	F2 String // `ipldsch:"optional String"`
	F3 String // `ipldsch:"optional nullable String"`
	F4 String // `ipldsch:"nullable String"`
}

func (d Foo__Content) Build() Foo {
	return Foo{d}
}

// and we still need something to handle all the accessors.
// maybe we should just... go back to the MaybeFoo stuff.
// use bitfields for longlived stuff, but maybestuff, fat though it is, during both build and yield.
// the main reason that approach got unpalatable was the pointer fuffery, right?  which we're deciding to avoid regardless.

func (x Foo) F1() String {
	return x.d.F1
}
func (x Foo) F2() (exists bool, value String) {
	return true, x.d.F1 // we can do both fo these methods
}
func (x Foo) F2_Maybe() MaybeString {
	return MaybeString{value: x.d.F1}
}

func frob(x Foo) {

}
