Schema Syntax
=============

Kinds
-----

Kinds in schemas are a superset of the kinds already known at the Data Model layer:

- "representation kinds" are map, array, string, bytes, boolean, integer, etc -- all of these are clear at the Data Model layer.
- "perceived kinds" include struct, union, and enum -- these are introduced at the Schema layer!

Perceived kinds are not inherently representable; rather, they're a way we look
at data in the other representation kinds and constrain their behaviors to be
easier to work with.  Schemas can declare types which are the perceived kinds,
but to do so, the schema simultaneously has to declare how to map that type onto
representation kinds.

Maps and lists at the schema level are more constrained than their equivalents
at the Data Model layer.  At the Data Model layer, heterogenous contents are
always allowed; in the Schema system, key and value types must be declared!
(Structs still typically describe heterogenous-content maps (or lists).
Wildcard types can also be used for value types explicitly when necessary.)


Types
-----

Types are declared with a name, which of the kinds it is, any details of the
type (e.g., fields in a struct; value types in a list; etc), and a representation.

The type declaration syntax varies by kind (e.g. unions and maps use visually
distinctive syntaxes) and follow a variety of rules:

- Types of recursive kinds (maps and lists) have terse declaration syntaxes:
  `[valueType]` defines arrays, and `{keyType:valueType}` defines maps.
- Struct fields may be declared to contain any named type, or alternatively may
  be declared to contain an anonymous type for which the definition is inlined
  (e.g. arrays can be defined inline for a struct field's type: `fieldName [Type]`).
- Recursive kind types may also use inline definitions
  (e.g. `type Foo {Bar:[Baz]}` is a map containing lists containing `Baz` elements).
- Struct fields, map values, and array values can be defined as "nullable".
- Struct fields can also be defined as "optional" (distinct from nullable: the
  key may be absent, but if present, the value must be non-null).

By example:

```ipldsch
type MyString string // "type" is a keyword; "MyString" is the name; "string" is the kind.
type MyInt int       // "type" is a keyword; "MyInt" is the name; "int" is the kind.
```

Recursive types have more details:

```ipldsch
type MyList [String] // "MyList" is the name; brackets indicate list kind;
                     //  and "String" is the contained value type.
type MyMap {String:String} // Curly-braces indicate map kind.
```

Note that in the above examples, `String` is capitalized, not lowercase:
this is because these are references to the string *type*, not a bare *kind*.
(`String` is a built-in/default type: `type String string`.)

We can also have structs, which are composed of fields:

```ipldsch
type MyStruct struct {
	AnInt Int      // "AnInt" is the field name; "Int" is the type.
	AString String // "AString" is the field name; "String" is the type.
}
```

TODO introduce examples of anonymous recursive types

TODO introduce examples of non-default representation

TODO introduce examples of nullable

TODO introduce examples of optional

TODO introduce examples of enums

TODO introduce examples of unions (in all their representations)


Fully Worked Examples
---------------------

See the self-representing Schema schema:

- [Schema schema in Schema DSL format](../typed/declaration/schema-schema.ipldsch)
- [Schema schema in IPLD-over-JSON format](../typed/declaration/schema-schema.ipldsch.json)

See also some other example schema:

- [Example schema in Schema DSL format](../typed/declaration/examples.ipldsch)
- [Example schema in IPLD-over-JSON format](../typed/declaration/examples.ipldsch.json)
