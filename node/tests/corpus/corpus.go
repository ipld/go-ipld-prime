/*
	The corpus package exports some values useful for building tests and benchmarks.

	Values come as JSON strings.  It is assumed you can unmarshal those.
	The reason we do this is so that this corpus package doesn't import
	any particular concrete implementation of ipld.Node... since that would
	make it ironically incapable of being used for that Node's tests.

	The naming convention is roughly as follows:

		- {Kind}{{Count}|N}{KeyKind}{ValueKind}
		- 'Kind' is usually 'Map' or 'List'.
		  It can also be a scalar like 'Int', in which case that's it.
		- If a specific int is given for 'Count', that's the size of the thing;
		- if 'N' is present, it's a scalable corpus and you can decide the size.
		- 'KeyKind' is present for maps (it will be string...).
		- 'ValueKind' is present for maps and lists.  It can recurse.

	Of course, this naming convention is not perfectly specific,
	but it's usually enough for our needs, or at least enough to get started.
	Some corpuses designed for probing (for example) tuple-represented structs
	will end up with interesting designations for various reasons:

		- some corpuses are meant to test struct semantics.
		  This is usually what it means when you see fixed size maps.
		  "List5Various" can also be this reason (it's for tuple-represented structs).
		- some corpuses are meant to test nullable or optional semantics.
		  These might have name suffixes like "WithNull" to indicate this.

	Everything is exported as a function, for consistency.
	Many functions need no args.
	Some functions need an argument for "N".

	If you're using these corpuses in a benchmark, don't forget to call
	`b.ResetTimer()` after getting the corpus.
*/
package corpus

import (
	"fmt"
)

func Map3StrInt() string {
	return `{"whee":1,"woot":2,"waga":3}`
}

func MapNStrInt(n int) string {
	return `{` + ents(n, func(i int) string {
		return fmt.Sprintf(`"k%d":%d`, i, i)
	}) + `}`
}

func MapNStrMap3StrInt(n int) string {
	return `{` + ents(n, func(i int) string {
		return fmt.Sprintf(`"k%d":`, i) +
			fmt.Sprintf(`{"whee":%d,"woot":%d,"waga":%d}`, i*3+1, i*3+2, i*3+3)
	}) + `}`
}
