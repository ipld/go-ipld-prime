/*
Printer provides features for printing out IPLD nodes and their contained data in a human-readable diagnostic format.

Outputs should look like...

	map{
		string{"foo"}: string{"bar"}
		string{"zot"}: struct<Foobar>{
			someFieldName: list{
				0: string{"this list is untyped"}
				1: string{"and contains a mixture of kinds of values"}
				2: int{400}
				3: bool{true}
			}
			otherField: list<ANamedListType>{
				0: string<String>{"mind you: 'ANamedListType' is the name of the *list type*."}
				1: string<String>{"it is not the name of the types of the value."}
				2: string<String>{"you'd have to look at the schema for that information."}
				3: string<String>{"or, of course, you can see it at the start of each of these entries, since they are also each annotated."}
			}
			moreField: list<[nullable String]>{
				0: string<String>{"this is a typed list"}
				1: string<String>{"but anonymous (meaning you see the value type in the 'name' of it)"}
				2: null
			}
		}
		string{"frog"}: map<{String:String}>{
			string<String>{"as you have probably imagined"}: string<String>{"this is a typed (but anonymous type) map"}
		}
		string{"numbers"}: int{1}
		string{"binary"}: bytes{ABCDEF0123456789}
		string{"typed numbers"}: int<MyNamedTypeInt>{9000}
		string{"typed string"}: string<MyNamedTypeString>{"okay, this one needed some marker prefixes."}
		string{"map with typed keys"}: map<{MyNamedTypeString:MyNamedTypeString}>{
			string<MyNamedStringType>{"work just fine"}: string<MyNamedTypeString>{"there's no ambiguity"}
			string<MyNamedStringType>{"you could elide key type info"}: string<MyNamedTypeString>{"as long as its a string kind, anyway"}
			string<MyNamedStringType>{"but we don't by default"}: string<MyNamedTypeString>{"explicit is good, especially in a debug tool!"}
		}
		string{"structs"}: struct<FooBar>{
			foo: string<String>{"do not need to have quoted field names"}
			bar: string<String>{"because (unlike map keys) their character range is already restricted"}
		}
		string{"unit types"}: unit<TheTypeName>
		string{"notice unit types"}: string{"have no braces at all, because they have literally no further details.  they're all type info."}
		string{"unions"}: union<TheUnionName>{string<TheInhabitant>{
			"that was wild, wasn't it.  Check out these double closing braces, coming up, too!  also the string got forced to a new line, even though it usually would've clung closer to its type and kind marker."
		}}
		string{"enums"}: enum<TheEnumName>{"inhabitant name"}
		string{"typed bools"}: bool<TheBoolName>{true}
		string{"map with struct keys"{: map<{FooBar:String}>{
			struct<FooBar>{foo:"foo", bar:"bar"}: string<String>{"that one probably surprised you, didn't it?"}
			struct<FooBar>{foo:"hmmm", bar:"maybe"}: string<String>{"we might be able to get away without the kind+type marker, actually.  but we need the one-liner struct content printing, at least, for sure."}
		}
		string{"map with really wicked keys"}: map<{WickedNestedUnion:String}>{
			union<WickedNestedUnion>{union<AnotherUnion>{string<TheInhabitant>{"wow"}}}: "yeah, that happens sometimes"
		}
	}

The pattern is a preamble saying what kind the value is (and what type, if applicable), followed by the actual value content, in braces.
For untyped nodes, this means `kindname{"value"}` (so: `string{"foo"}` and `int{12}` and `bool{true}` etc),
or for typed nodes, we get `typekindname<TheTypeName>{"value"}`.
In addition to the example above, you can check out the tests for a few more examples of how it looks.

Some configuration options are available to elide some information.
For example, some configuration can reduce the amount of annotational weight around strings
(which is possible to do without getting completely vague because the quotation markings for strings already are syntatically distinctive).
Not all things can be configured for elision, however.
For example, for the various number kinds, the kind preambles are always required.  (Number parsers are otherwise often an annoying lookahead problem.)
Similarly, for bytes, the kind preamble is always required.  (Among other things, the hexidecimal up until the first letter could be confused with an integer, if we didn't label both of them.)
Anything that's typed also gets a preamble with the type and kind information, even if its kind is something we'd otherwise elide, like string.

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
