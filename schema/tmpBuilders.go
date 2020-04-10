package schema

// Everything in this file is __a temporary hack__ and will be __removed__.
//
// These methods will only hang around until more of the "ast" packages are finished;
// thereafter, building schema.Type and schema.TypeSystem values will only be
// possible through first constructing a schema AST, and *then* using Reify(),
// which will validate things correctly, cycle-check, cross-link, etc.
//
// (Meanwhile, we're using these methods in the codegen prototypes.)

func SpawnString(name TypeName) TypeString {
	return TypeString{anyType{name, nil}}
}

func SpawnInt(name TypeName) TypeInt {
	return TypeInt{anyType{name, nil}}
}

func SpawnBytes(name TypeName) TypeBytes {
	return TypeBytes{anyType{name, nil}}
}

func SpawnLink(name TypeName) TypeLink {
	return TypeLink{anyType{name, nil}, nil, false}
}

func SpawnLinkReference(name TypeName, referenceType Type) TypeLink {
	return TypeLink{anyType{name, nil}, referenceType, true}
}
func SpawnList(name TypeName, typ Type, nullable bool) TypeList {
	return TypeList{anyType{name, nil}, false, typ, nullable}
}

func SpawnStruct(name TypeName, fields []StructField, repr StructRepresentation) TypeStruct {
	v := TypeStruct{
		anyType{name, nil},
		fields,
		make(map[string]StructField, len(fields)),
		repr,
	}
	for i := range fields {
		fields[i].parent = &v
		v.fieldsMap[fields[i].name] = fields[i]
	}
	return v
}
func SpawnStructField(name string, typ Type, optional bool, nullable bool) StructField {
	return StructField{nil /*populated later*/, name, typ, optional, nullable}
}
func SpawnStructRepresentationMap(renames map[string]string) StructRepresentation_Map {
	return StructRepresentation_Map{renames, nil}
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
