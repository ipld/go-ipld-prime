/*
Printer provides features for printing out IPLD nodes and their contained data in a human-readable diagnostic format.

Outputs should look like...

	map{
		"foo": "bar"
		"zot": struct<Foobar>{
			someFieldName: list{
				0: "this"
				1: "is untyped"
				2: int{400}
				3: bool{true}
			}
			otherField: list<ANamedListType>{
				0: "mind you: 'ANamedListType' is the name of the *list type*."
				1: "it is not the name of the value types.  those are not actually shown here."
				2: "you'd have to look at the schema for that information."
			}
			moreField: list<[nullable String]>{
				0: "this is a typed list"
				1: "but anonymous"
				2: null
			}
		}
		"frog": map<{String:String}>{
			"as you have probably imagined": "this is a typed (but anonymous type) map"
		}
		"numbers": int{1}
		"binary": bytes{ABCDEF0123456789}
		"typed numbers": int<MyNamedTypeInt>{9000}
		"typed string": string<MyNamedTypeString>{"okay, this one needed some marker prefixes."}
		"map with typed keys": map<{MyNamedTypeString:MyNamedTypeString}>{
			"despite being typed": string<MyNamedTypeString>{"map keys still never need prefixes; there's no ambiguity"}
			"but for the values": string<MyNamedTypeString>{"we still use them"}
			"well, maybe": string<MyNamedTypeString>{"this might actually be worth debating; it doesn't seem necessary (unless the map value is a variant type, but that's then already addressed by another case)"}
		}
		"structs": struct<FooBar>{
			foo: "do not need to have quoted field names"
			bar: "because (unlike map keys) their character range is already restricted"
		}
		"unit types": unit<TheTypeName>
		"notice unit types": "have no braces at all, because they have literally no further details.  they're all type info."
		"variants": variant<TheUnionName>{string<TheInhabitant>{
			"that was wild, wasn't it.  Check out these double closing braces, coming up, too!  also the string got forced to a new line, even though it usually would've clung closer to its type and kind marker."
		}}
		"enums": enum<TheEnumName>{"inhabitant name"}
		"typed bools": bool<TheBoolName>{true}
		"map with struct keys": map<{FooBar:String}>{
			struct<FooBar>{foo:"foo", bar:"bar"}: "that one probably surprised you, didn't it?"
			struct<FooBar>{foo:"hmmm", bar:"maybe"}: "we might be able to get away without the kind+type marker, actually.  but we need the one-liner struct content printing, at least, for sure."
		}
		"map with really wicked keys": map<{WickedNestedUnion:String}>{
			variant<WickedNestedUnion>{variant<AnotherUnion>{string<TheInhabitant>{"wow"}}}: "yeah, that happens sometimes"
		}
	}


Notice that strings can be emitted without a kind indicator.
This is optional, and configurable, but a default, because strings are so common
(and also, that the quotation marks already make them sufficiently clear)
that it's neither necessary nor desirable to burden every string with a kind indicator.
Everything else *does* have an explicit leading indicator that names the data's kind.
Maps and lists already need enough syntactic weight that adding another few characters isn't a significant weight.
For the various number kinds... well, it seems better to be clear, with those.  (Number parsers are otherwise often an annoying lookahead problem.)
For bytes, the need is obvious.  (Among other things, the hexidecimal up until the first letter could be confused with an integer, if we didn't label both of them.)
Anything that's typed also gets a leading indicator section again, even if its kind is something we'd otherwise elide, like string.

Notice that struct fields aren't quoted.  (It's not necessary, because field names are already constrained.)
But map keys are.  (They need quoting because they can be any string.)

Note that the output of printer is NOT INTENDED TO BE PARSABLE.
It is NOT an IPLD codec!
It is a diagnostic format only.
Much of the information included (especially about schema type information)
is _more_ information than the IPLD data model holds alone,
so trying to re-parse the printer output would be a strange choice.

The diagnostic format emitted by printer is not formally specified,
and is not necessarily language-agnostic.
It may not even remain stable across releases of this library.
It is intended to be used for diagnostics only.

*/
package printer

/*
How to print ADLs is not yet clear.

Perhaps something like `<!TheADLName>` will do;
this would also stack reasonably clearly with types as `<TheTypeName!TheADLName>`;
this style would have the downside of making ADLs look *very* different than other mere representation strategies,
which may be totally reasonable or mildly questionable depending on how purist you feel about that.
*/
