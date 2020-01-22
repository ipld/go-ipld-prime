/*
	This system is going to write a *lot* of code,
	and then defacto dictate style and constraints on the shape of *even more*
	code that *other people* subsequently write (and invest a lot of time in --
	both writing, and then later yet debugging and maintaining).

	We basically have to imagine every possible way to write code in this language,
	and then think about which ones we want to encourage,
	and then figure out how to encourage that.

	The impact is HIGH.

	This requires a lot (...LOT...) of thinking.
*/
package typedbuilders

import "testing"

func BenchmarkRound2(b *testing.B) {
	var v Stroct
	for i := 0; i < b.N; i++ {
		v = Stroct__Cntnt{
			F1: "foo",
			F2: "bar",
			F3: "baz",
			F4: "woz",
		}.Build()
		// The 'Patch' function is the only one who can tolerate 'skip this' bools.  Would be nice to confine that.
		//  Adding a third^W fourth type tho?

		// We've now got proposal for:
		// - the readonly type
		// - the mutable builder ('content') which disregards
		// - the null/absent bitmask type
		// - the patch/ignore bits type.
		// each of these would need exported symbols if used as a struct;
		//  so at least the first two do (maybe the latter two can be builder.func gated).
		// four is two many.  two is acceptable.  three is stretching.

		// if the mutable builder also avoids nulls by default: that may *also* actua...... no, don't care.  yes shit do damnit.
		//  if you want to hoist one value out of the mutable builder by pointer, it's gonna move the whole thing,
		//   and this actually becomes *more* of a concern if you do recursively-flat objects, because the constants get bigger.
	}
	sink = v
}
