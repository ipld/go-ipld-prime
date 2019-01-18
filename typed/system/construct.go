package typesystem

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
	typedeclaration "github.com/ipld/go-ipld-prime/typed/declaration"
)

func Construct(decls ...typedeclaration.Type) (_ *Universe, err error) {
	defer recover()
	var tsu = Universe{
		namedTypes: make(map[TypeName]Type),
	}
	var checkName = func(name TypeName) {
		if _, exists := tsu.namedTypes[name]; exists {
			err = fmt.Errorf("repeated name (%q)", name)
			panic(nil)
		}
	}

	// Inject the base types.
	tsu.namedTypes["Bool"] = TypeBool{anyType{name: "Bool", universe: &tsu}}
	tsu.namedTypes["String"] = TypeString{anyType{name: "String", universe: &tsu}}
	tsu.namedTypes["Bytes"] = TypeBytes{anyType{name: "Bytes", universe: &tsu}}
	tsu.namedTypes["Int"] = TypeInt{anyType{name: "Int", universe: &tsu}}
	tsu.namedTypes["Float"] = TypeFloat{anyType{name: "Float", universe: &tsu}}

	// Do everything in two giant passes:
	//  - Pass 1: get everybody's name in the map (and do any other basic checks);
	//  - Pass 2: fill in connections for recursive types (since now we have something to point to from Pass 1).
	for _, d := range decls {
		switch decl := d.(type) {
		case typedeclaration.TypeBool:
			checkName(decl.Name)
			tsu.namedTypes[decl.Name] = TypeBool{
				anyType{decl.Name, &tsu},
			}
		case typedeclaration.TypeString:
			checkName(decl.Name)
			tsu.namedTypes[decl.Name] = TypeString{
				anyType{decl.Name, &tsu},
			}
		case typedeclaration.TypeBytes:
			checkName(decl.Name)
			tsu.namedTypes[decl.Name] = TypeBytes{
				anyType{decl.Name, &tsu},
			}
		case typedeclaration.TypeInt:
			checkName(decl.Name)
			tsu.namedTypes[decl.Name] = TypeInt{
				anyType{decl.Name, &tsu},
			}
		case typedeclaration.TypeFloat:
			checkName(decl.Name)
			tsu.namedTypes[decl.Name] = TypeFloat{
				anyType{decl.Name, &tsu},
			}
		case typedeclaration.TypeMap:
			checkName(decl.Name)
			if decl.Anon {
				return nil, fmt.Errorf("anonymous declarations are nonsense except as child nodes of other recursive types (%q should not be anonymous)", decl.Name)
			}
			tsu.namedTypes[decl.Name] = TypeMap{
				anyType: anyType{decl.Name, &tsu},
			}
		case typedeclaration.TypeList:
			checkName(decl.Name)
			if decl.Anon {
				return nil, fmt.Errorf("anonymous declarations are nonsense except as child nodes of other recursive types (%q should not be anonymous)", decl.Name)
			}
			tsu.namedTypes[decl.Name] = TypeList{
				anyType: anyType{decl.Name, &tsu},
			}
		case typedeclaration.TypeLink:
			checkName(decl.Name)
			tsu.namedTypes[decl.Name] = TypeLink{
				anyType{decl.Name, &tsu},
			}
		case typedeclaration.TypeUnion:
			checkName(decl.Name)
			t := TypeUnion{
				anyType: anyType{decl.Name, &tsu},
			}
			switch decl.Style {
			case typedeclaration.UnionStyle_Kinded:
				t.style = UnionStyle_Kinded
			case typedeclaration.UnionStyle_Keyed:
				t.style = UnionStyle_Keyed
			case typedeclaration.UnionStyle_Envelope:
				t.style = UnionStyle_Envelope
				t.typeHintKey = decl.TypeHintKey
				t.contentKey = decl.ContentKey
			case typedeclaration.UnionStyle_Inline:
				t.style = UnionStyle_Inline
				t.typeHintKey = decl.TypeHintKey
			}
			tsu.namedTypes[decl.Name] = t
		case typedeclaration.TypeStruct:
			checkName(decl.Name)
			tsu.namedTypes[decl.Name] = TypeStruct{
				anyType:    anyType{decl.Name, &tsu},
				tupleStyle: decl.TupleStyle,
			}
		case typedeclaration.TypeEnum:
			checkName(decl.Name)
			t := TypeEnum{
				anyType: anyType{decl.Name, &tsu},
				members: make([]string, len(decl.Members)),
			}
			copy(t.members, decl.Members)
			tsu.namedTypes[decl.Name] = t
		}
	}

	// Now that 2nd pass for filling in recursives.
	for _, d := range decls {
		switch decl := d.(type) {
		case typedeclaration.TypeMap:
			t := tsu.namedTypes[decl.Name].(TypeMap)
			keyType := tsu.namedTypes[decl.KeyType]
			if keyType == nil {
				fmt.Errorf("type %q references an undeclared type %q for its key type", decl.Name, decl.KeyType)
			}
			if keyType.ReprKind() != ipld.ReprKind_String {
				fmt.Errorf("type %q has an invalid key type: map key types must be of string representation kind (%q is %s)", decl.Name, decl.KeyType, keyType.ReprKind())
			}
			t.keyType = keyType
			t.valueType = tsu.namedTypes[decl.ValueType] // FIXME will need updating to handle recursive TypeTerm, see matching comment in tdecl types.
			if t.valueType == nil {
				fmt.Errorf("type %q references an undeclared type %q for its value type", decl.Name, decl.ValueType)
			}
			t.valueNullable = decl.ValueNullable
		case typedeclaration.TypeList:
			t := tsu.namedTypes[decl.Name].(TypeList)
			t.valueType = tsu.namedTypes[decl.ValueType] // FIXME will need updating to handle recursive TypeTerm, see matching comment in tdecl types.
			if t.valueType == nil {
				fmt.Errorf("type %q references an undeclared type %q for its value type", decl.Name, decl.ValueType)
			}
			t.valueNullable = decl.ValueNullable
		case typedeclaration.TypeUnion:
			t := tsu.namedTypes[decl.Name].(TypeUnion)
			switch t.style {
			case UnionStyle_Kinded:
				// TODO
			case UnionStyle_Keyed:
				fallthrough
			case UnionStyle_Envelope:
				t.values = make(map[string]Type, len(decl.Values))
				for k, v := range decl.Values {
					t.values[k] = tsu.namedTypes[v]
					if t.values[k] == nil {
						fmt.Errorf("type %q references an undeclared type %q for one of its members", decl.Name, v)
					}
				}
			case UnionStyle_Inline:
				t.values = make(map[string]Type, len(decl.Values))
				for k, v := range decl.Values {
					t.values[k] = tsu.namedTypes[v]
					if t.values[k] == nil {
						fmt.Errorf("type %q references an undeclared type %q for one of its members", decl.Name, v)
					}
					// TODO need to check that all of the members are reprKind==map!
					// TODO need to check that none of the members have conflicting keys!
					//   ... so that might also mean all the members need to be kind==struct, even more specifically than reprKind==map.
				}
			}
		case typedeclaration.TypeStruct:
			panic("TODO") // TODO
		}
	}

	// Particularly fun edge cases: might want to check for cyclic
	//  references that don't have any nullables or optionals that allow
	//   the cycle to be broken in concrete values.  But we'll leave that
	//    for a later fun PR :)

	return &tsu, nil
}
