package schema

// Everything in this file is __a temporary hack__ and will be __removed__.
//
// These methods will only hang around until more of the "ast" packages are finished;
// thereafter, building schema.Type and schema.TypeSystem values will only be
// possible through first constructing a schema AST, and *then* using Reify(),
// which will validate things correctly, cycle-check, cross-link, etc.
//
// (Meanwhile, we're using these methods in the codegen prototypes.)

// These methods use Type objects as parameters when pointing to other things,
//  but this is... turning out consistently problematic.
//   Even when we're doing this hacky direct-call doesn't-need-to-be-serializable temp stuff,
//    as written, this doesn't actually let us express cyclic things viably!
//   The same initialization questions are also going to come up again when we try to make
//    concrete values in the output of codegen.
// Maybe it's actually just a bad idea to have our reified Type types use Type pointers at all.
//  (I will never get tired of the tongue twisters, evidently.)
//  I'm not actually using that much, and it's always avoidable (it's trivial to replace with a map lookup bouncing through a 'ts' variable somewhere).
//  And having the AST gen'd types be... just... the thing... sounds nice.  It could save a lot of work.
//   (It would mean the golang types don't tell you whether the values have been checked for global properties or not, but, eh.)
//   (It's not really compatible with "Prototype and Type are the same thing for codegen'd stuff", either (or, we need more interfaces, and to *really* lean into them), but maybe that's okay.)

func SpawnString(name TypeName) TypeString {
	return TypeString{typeBase{name, nil}}
}

func SpawnInt(name TypeName) TypeInt {
	return TypeInt{typeBase{name, nil}}
}

func SpawnBytes(name TypeName) TypeBytes {
	return TypeBytes{typeBase{name, nil}}
}

func SpawnLink(name TypeName) TypeLink {
	return TypeLink{typeBase{name, nil}, nil, false}
}

func SpawnLinkReference(name TypeName, referenceType Type) TypeLink {
	return TypeLink{typeBase{name, nil}, referenceType, true}
}

func SpawnList(name TypeName, typ Type, nullable bool) TypeList {
	return TypeList{typeBase{name, nil}, false, typ, nullable}
}

func SpawnMap(name TypeName, keyType Type, valueType Type, nullable bool) TypeMap {
	return TypeMap{typeBase{name, nil}, false, keyType, valueType, nullable}
}

func SpawnStruct(name TypeName, fields []StructField, repr StructRepresentation) TypeStruct {
	v := TypeStruct{
		typeBase{name, nil},
		fields,
		make(map[string]StructField, len(fields)),
		repr,
	}
	for i := range fields {
		fields[i].parent = &v
		v.fieldsMap[fields[i].name] = fields[i]
	}
	switch repr.(type) {
	case StructRepresentation_Stringjoin:
		for _, f := range fields {
			if f.IsMaybe() {
				panic("neither nullable nor optional is supported on struct stringjoin representation")
			}
		}
	}
	return v
}
func SpawnStructField(name string, typ Type, optional bool, nullable bool) StructField {
	return StructField{nil /*populated later*/, name, typ, optional, nullable}
}
func SpawnStructRepresentationMap(renames map[string]string) StructRepresentation_Map {
	return StructRepresentation_Map{renames, nil}
}
func SpawnStructRepresentationStringjoin(delim string) StructRepresentation_Stringjoin {
	return StructRepresentation_Stringjoin{delim}
}

func SpawnUnion(name TypeName, members []Type, repr UnionRepresentation) TypeUnion {
	return TypeUnion{typeBase{name, nil}, members, repr}
}
func SpawnUnionRepresentationKeyed(table map[string]Type) UnionRepresentation_Keyed {
	return UnionRepresentation_Keyed{table}
}

// The methods relating to TypeSystem are also mutation-heavy and placeholdery.

func (ts *TypeSystem) Init() {
	ts.namedTypes = make(map[TypeName]Type)
}
func (ts *TypeSystem) Accumulate(typ Type) {
	ts.namedTypes[typ.Name()] = typ
}
func (ts TypeSystem) GetTypes() map[TypeName]Type {
	return ts.namedTypes
}
func (ts TypeSystem) TypeByName(n string) Type {
	return ts.namedTypes[TypeName(n)]
}
