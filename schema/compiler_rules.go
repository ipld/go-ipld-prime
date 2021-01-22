package schema

import "fmt"

type rule struct {
	// text is the name of the rule and the start of the error message body if the rule is flunked.
	text string

	// predicate should return true if this rule will apply to this type at all.
	// Usually it just checks the representation strategy.
	// The typekind is already defacto checked by indexing into the rules map.
	//
	// The distinction between predicate and rule may be useful to build a more advanced diagnostic output.
	//
	// Note that since these functions are called on a type *during* its validation process,
	// some data in the type isn't known to be valid yet,
	// and as a result some helper functions may not be safe to use;
	// therefore it's often necessary to access the raw fields directly.
	predicate func(Type) bool

	// rule is the actual rule body.  If it's a non-nil return, the rule is flunked.
	// The error is for freetext detail formatting; it will be wrapped by another error
	// which is based on the rule's Text.
	//
	// Same caveats about the Type's validity apply as they did for the predicate func.
	rule func(*TypeSystem, Type) error
}

type ErrInvalidTypeSpec struct {
	Rule   string
	Type   TypeReference
	Detail error
}

// To validate a type:
//  - first get the slice fo rules that apply to its typekind
//  - then, for each rule:
//    - check if it applies (by virtue of the predicate); skip if not
//    - evaluate the rule and check if it errors
//    - errors accumulate and do not cause halts
var rules = map[TypeKind][]rule{
	TypeKind_Map: []rule{
		{"map declaration's key type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) error {
				tRef := TypeReference(t.(*TypeMap).keyTypeRef)
				if _, exists := ts.types[tRef]; !exists {
					return fmt.Errorf("missing type %q", tRef)
				}
				return nil
			},
		},
		{"map declaration's key type must be stringable",
			alwaysApplies,
			func(ts *TypeSystem, t Type) error {
				tRef := TypeReference(t.(*TypeMap).keyTypeRef)
				if hasStringRepresentation(ts.types[tRef]) {
					return fmt.Errorf("type %q is not a string typekind nor representation with string kind", tRef)
				}
				return nil
			},
		},
		{"map declaration's value type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) error {
				tRef := TypeReference(t.(*TypeMap).valueTypeRef)
				if _, exists := ts.types[tRef]; !exists {
					return fmt.Errorf("missing type %q", tRef)
				}
				return nil
			},
		},
	},
	TypeKind_List: []rule{
		{"list declaration's value type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) error {
				tRef := TypeReference(t.(*TypeList).valueTypeRef)
				if _, exists := ts.types[tRef]; !exists {
					return fmt.Errorf("missing type %q", tRef)
				}
				return nil
			},
		},
	},
	TypeKind_Link: []rule{
		{"link declaration's expected target type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) error {
				tRef := TypeReference(t.(*TypeLink).expectedTypeRef)
				if _, exists := ts.types[tRef]; !exists {
					return fmt.Errorf("missing type %q", tRef)
				}
				return nil
			},
		},
	},
	TypeKind_Struct: []rule{
		{"struct declaration's field type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) error {
				for _, field := range t.(*TypeStruct).fields {
					tRef := TypeReference(field.typeRef)
					if _, exists := ts.types[tRef]; !exists {
						return fmt.Errorf("missing type %q", tRef) // TODO: want this to return multiple errors.  time for another `*[]error` accumulator?  sigh.
					}
				}
				return nil
			},
		},
	},
	TypeKind_Union: []rule{
		{"union declaration's potential members must all be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) error {
				for _, member := range t.(*TypeUnion).members {
					tRef := TypeReference(member)
					if _, exists := ts.types[tRef]; !exists {
						return fmt.Errorf("missing type %q", tRef) // TODO: want this to return multiple errors.
					}
				}
				return nil
			},
		},
		// TODO continue with more union rules... but... they're starting to get conditional on the passage of prior rules.
		//   Unsure how much effort it's worth to represent this in detail.
		//     - Should we have flunk of one rule cause subsequent rules to be skipped on that type?
		//     - Should we just re-do all the prerequisite checks, but return nil if those fail (since another rule should've already reported those)?
		//     - Should we re-do all the prerequisite checks, and return a special 'inapplicable' error code if those fail?
		//     - Should we build a terribly complicated prerequisite tracking system?
		//       - (Okay, maybe it's not that complicated; a tree would probably suffice?)
		//   My original aim with this design was to get as close as possible to something table-driven,
		//    in the hope this would make it easier to port the semantics to other languages.
		//     As this code gets fancier, that goal fades fast, so a solution that's KISS is probably preferrable.
	},
}

var alwaysApplies = func(Type) bool { return true }

// hasStringRepresentation returns a bool for... well, what it says on the tin.
func hasStringRepresentation(t Type) bool {
	// Note that some of these cases might get more complicated in the future.
	//  We haven't defined or implemented features like "type Foo int representation base10str" yet, but we could.
	// This doesn't recursively check the sanity of types that claim to have string representation
	//  (e.g. that every field in a stringjoin struct is also stringable);
	//  the caller should do that (and the Compiler, which is the caller, does: on each type as it is looping over the whole set).
	switch t2 := t.(type) {
	case *TypeBool:
		return false
	case *TypeString:
		return true
	case *TypeBytes:
		return false
	case *TypeInt:
		return false
	case *TypeFloat:
		return false
	case *TypeMap:
		switch t2.Representation().(type) {
		case MapRepresentation_Map:
			return false
		case MapRepresentation_Stringpairs:
			return true
		case MapRepresentation_Listpairs:
			return false
		default:
			panic("unreachable")
		}
	case *TypeList:
		return false
	case *TypeLink:
		return false
	case *TypeStruct:
		switch t2.Representation().(type) {
		case StructRepresentation_Map:
			return false
		case StructRepresentation_Tuple:
			return false
		case StructRepresentation_Stringpairs:
			return true
		case StructRepresentation_Stringjoin:
			return true
		case StructRepresentation_Listpairs:
			return false
		default:
			panic("unreachable")
		}
	case *TypeUnion:
		switch t2.Representation().(type) {
		case UnionRepresentation_Keyed:
			return false
		case UnionRepresentation_Kinded:
			return false
		case UnionRepresentation_Envelope:
			return false
		case UnionRepresentation_Inline:
			return false
		case UnionRepresentation_Stringprefix:
			return true
		case UnionRepresentation_Byteprefix:
			return false
		default:
			panic("unreachable")
		}
	case *TypeEnum:
		switch t2.Representation().(type) {
		case EnumRepresentation_String:
			return true
		case EnumRepresentation_Int:
			return false
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
}
