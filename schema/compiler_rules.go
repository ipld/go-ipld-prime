package schema

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
)

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

	// rule is the actual rule body.  If it encounters errors, it appends them to the list.
	//
	// The errors may do freetext details; the rule evaluator will wrap them with
	// another error which is based on the rule's Text.
	//
	// Same caveats about the Type's validity apply as they did for the predicate func.
	//
	// (Design note: considered taking an error slice as a param and accumulating it,
	// but this would've required pushing more error formatting down,
	// and while speed is nice, there's no alloc costs if you don't have errors,
	// so the speed really isn't a concern (e.g. when we're initializing a type system
	// that provides the selfdescription for codegen'd types, there's no error allocs).)
	rule func(*TypeSystem, Type) []error
}

type ErrInvalidTypeSpec struct {
	Rule   string
	Type   TypeReference
	Detail error
}

func (e ErrInvalidTypeSpec) Error() string {
	return fmt.Sprintf("type %s is invalid: %s: %s", e.Type, e.Rule, e.Detail)
}

func validate(ts *TypeSystem, typ Type, errs *[]error) {
	rules := rules[typ.TypeKind()]
	for _, rule := range rules {
		if rule.predicate(typ) {
			newErrors := rule.rule(ts, typ)
			if len(newErrors) > 0 {
				for _, err := range newErrors {
					*errs = append(*errs, ErrInvalidTypeSpec{rule.text, typ.Reference(), err})
				}
				break
			}
		}
	}
}

// The rules table contains all the logical validations that apply to a schema during compilation.
//  Some forms of validation of the data are already done by nature of the schema-schema; others will require more work here.
//   In general: rules which stretch across multiple types (especially, if they're graph properties) can't be implemented in schemas alone, and so end up here.
//  The most common example is that any type that has some kind of recursion (maps, lists, structs, unions, links with target type info)
//   will need to do a lookup to see if the referenced types were defined elsewhere in the schema document.
//  Some kinds of types involve other more specific checks,
//   such as maps verifying that their keys are stringable (which is a rule we enforce for reasons relating to pathing),
//   and unions verifying that all their discriminant tables are complete (which is a rule that's necessary for sanity!),
//   and etc.
//
// To validate a type:
//  - first get the slice of rules that apply to its typekind
//  - then, for each rule:
//    - check if it applies (by virtue of the predicate); skip if not
//    - evaluate the rule; doing so may accumulate errors
//    - if any errors accumulated, skip all further rules for this type; it already flunked.
//
// The short circuiting logic between subsequent rules means that
// later rules are allowed to make presumptions of things checked by earlier rules.
//
// The table-like design here hopefully will make the semantics defined within
// easier to port to other implementations in other languages.
var rules = map[TypeKind][]rule{
	// FUTURE: after adding unit types, we'll need most recursives to additionally do some checks for correct composition of nullability.
	TypeKind_Map: []rule{
		{"map declaration's key type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) []error {
				tRef := TypeReference(t.(*TypeMap).keyTypeRef)
				if _, exists := ts.types[tRef]; !exists {
					return []error{fmt.Errorf("missing type %q", tRef)}
				}
				return nil
			},
		},
		{"map declaration's key type must be stringable",
			alwaysApplies,
			func(ts *TypeSystem, t Type) []error {
				tRef := TypeReference(t.(*TypeMap).keyTypeRef)
				if hasStringRepresentation(ts.types[tRef]) {
					return []error{fmt.Errorf("type %q is not a string typekind nor representation with string kind", tRef)}
				}
				return nil
			},
		},
		{"map declaration's value type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) []error {
				tRef := t.(*TypeMap).valueTypeRef
				if _, exists := ts.types[tRef]; !exists {
					return []error{fmt.Errorf("missing type %q", tRef)}
				}
				return nil
			},
		},
	},
	TypeKind_List: []rule{
		{"list declaration's value type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) []error {
				tRef := t.(*TypeList).valueTypeRef
				if _, exists := ts.types[tRef]; !exists {
					return []error{fmt.Errorf("missing type %q", tRef)}
				}
				return nil
			},
		},
	},
	TypeKind_Link: []rule{
		{"link declaration's expected target type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) []error {
				tRef := TypeReference(t.(*TypeLink).expectedTypeRef)
				if tRef == "" {
					return nil
				}
				if _, exists := ts.types[tRef]; !exists {
					return []error{fmt.Errorf("missing type %q", tRef)}
				}
				return nil
			},
		},
	},
	TypeKind_Struct: []rule{
		{"struct declaration's field type must be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) (errs []error) {
				for _, field := range t.(*TypeStruct).fields {
					tRef := field.typeRef
					if _, exists := ts.types[tRef]; !exists {
						errs = append(errs, fmt.Errorf("missing type %q", tRef))
					}
				}
				return
			},
		},
	},
	TypeKind_Union: []rule{
		{"union declaration's potential members must all be defined",
			alwaysApplies,
			func(ts *TypeSystem, t Type) (errs []error) {
				for _, member := range t.(*TypeUnion).members {
					tRef := TypeReference(member)
					if _, exists := ts.types[tRef]; !exists {
						errs = append(errs, fmt.Errorf("missing type %q", tRef))
					}
				}
				return nil
			},
		},
		{"union's representation must specify exactly one discriminant for each member",
			alwaysApplies,
			func(ts *TypeSystem, t Type) (errs []error) {
				t2 := t.(*TypeUnion)
				// All of these are very similar, but they store the info in technically distinct places, so we have to destructure to get at it.
				switch r := t2.rstrat.(type) {
				case UnionRepresentation_Keyed:
					checkUnionDiscriminantInfo(t2.members, r.discriminantTable, &errs)
				case UnionRepresentation_Kinded:
					checkUnionDiscriminantInfo2(t2.members, r.discriminantTable, &errs)
				case UnionRepresentation_Envelope:
					checkUnionDiscriminantInfo(t2.members, r.discriminantTable, &errs)
				case UnionRepresentation_Inline:
					checkUnionDiscriminantInfo(t2.members, r.discriminantTable, &errs)
				case UnionRepresentation_Stringprefix:
					checkUnionDiscriminantInfo(t2.members, r.discriminantTable, &errs)
				case UnionRepresentation_Byteprefix:
					checkUnionDiscriminantInfo(t2.members, r.discriminantTable, &errs)
				}
				return nil
			},
		},
		{"kinded union's discriminants must match the member's kinds",
			func(t Type) bool { _, ok := t.(*TypeUnion).rstrat.(UnionRepresentation_Kinded); return ok },
			func(ts *TypeSystem, t Type) (errs []error) {
				r := t.(*TypeUnion).rstrat.(UnionRepresentation_Kinded)
				for k, v := range r.discriminantTable {
					vrb := ts.types[TypeReference(v)].RepresentationBehavior()
					if vrb == ipld.Kind_Invalid { // this indicates a kinded union (the only thing that can't statically state its representation behavior), which deserves a special error message.
						errs = append(errs, fmt.Errorf("%s is not a valid member: kinded unions cannot be nested and %s is also a kinded union", v, v))
					} else if vrb != k {
						errs = append(errs, fmt.Errorf("kind mismatch: %s is declared to be received as type %s, but that type's representation's kind is %s", k, v, vrb))
					}
				}
				return
			},
		},
		{"envelope union's magic keys must be distinct",
			func(t Type) bool { _, ok := t.(*TypeUnion).rstrat.(UnionRepresentation_Envelope); return ok },
			func(ts *TypeSystem, t Type) []error {
				r := t.(*TypeUnion).rstrat.(UnionRepresentation_Envelope)
				if r.discriminantKey == r.contentKey {
					return []error{fmt.Errorf("content key and discriminant key are the same")}
				}
				return nil
			},
		},
		{"inline union's members must all have map representations and not collide with the union's discriminant key",
			func(t Type) bool { _, ok := t.(*TypeUnion).rstrat.(UnionRepresentation_Inline); return ok },
			func(ts *TypeSystem, t Type) (errs []error) {
				r := t.(*TypeUnion).rstrat.(UnionRepresentation_Inline)
				// This is one of the more complicated rules.
				// - many typekinds can be rejected as members outright, because they don't have any map representations.
				// - maps themselves are acceptable, if they still have a map representation (although this is a janky thing to do in a protocol design, because it means there's at least one key that's now illeagal in that map, and we can't help you with that).
				// - structs are acceptable if they have map representation... but we also validate in advance that the discriminant key doesn't collide with any field names in any of the structs.
				// - other unions aren't ever valid members, even if their representation has a map kind, because the logical rules don't fit together.  So we give distinct error messages for this.
				for _, v := range r.discriminantTable {
					switch vt := ts.types[TypeReference(v)].(type) {
					case *TypeBool, *TypeString, *TypeBytes, *TypeInt, *TypeFloat, *TypeList, *TypeLink, *TypeEnum:
						errs = append(errs, fmt.Errorf("%s is not a valid member: has representation kind %s", v, vt.RepresentationBehavior()))
					case *TypeMap:
						if vt.RepresentationBehavior() != ipld.Kind_Map {
							errs = append(errs, fmt.Errorf("%s is not a valid member: has representation kind %s", v, vt.RepresentationBehavior()))
						}
					case *TypeUnion:
						errs = append(errs, fmt.Errorf("%s is not a valid member: inline unions cannot directly contain other unions, because their representation rules would conflict", v))
					case *TypeStruct:
						if vt.RepresentationBehavior() != ipld.Kind_Map {
							errs = append(errs, fmt.Errorf("%s is not a valid member: has representation kind %s", v, vt.RepresentationBehavior()))
						}
						for _, f := range vt.Fields() {
							if r.DiscriminantKey() == vt.Representation().(StructRepresentation_Map).GetFieldKey(f) {
								errs = append(errs, fmt.Errorf("%s is not a valid member: key collision: %s has a field that collides with %s's discriminant key when represented", v, v, t.Name()))
							}
						}
					}
				}
				return
			},
		},
		// FUTURE: UnionRepresentation_Stringprefix will probably have additional rules too
		// FUTURE: UnionRepresentation_Bytesprefix will probably have additional rules too
	},
	TypeKind_Enum: []rule{
		{"enums's representation must specify exactly one discriminant for each member",
			alwaysApplies,
			func(ts *TypeSystem, t Type) (errs []error) {
				t2 := t.(*TypeEnum)
				covered := make([]bool, len(t2.members))
				switch r := t2.RepresentationStrategy().(type) {
				case EnumRepresentation_String:
					for k, v := range r.labels {
						found := false
						for i, m := range t2.members {
							if k == m {
								if found {
									errs = append(errs, fmt.Errorf("more than one discriminant pointing to member %q", m))
								}
								found = true
								covered[i] = true
							}
						}
						if !found {
							errs = append(errs, fmt.Errorf("discriminant %q refers to a non-member %q", v, k))
						}
					}
				case EnumRepresentation_Int:
					for k, v := range r.labels {
						found := false
						for i, m := range t2.members {
							if k == m {
								if found {
									errs = append(errs, fmt.Errorf("more than one discriminant pointing to member %q", m))
								}
								found = true
								covered[i] = true
							}
						}
						if !found {
							errs = append(errs, fmt.Errorf("discriminant \"%d\" refers to a non-member %q", v, k))
						}
					}
				}
				for i, m := range t2.members {
					if !covered[i] {
						errs = append(errs, fmt.Errorf("missing discriminant for member %q", m))
					}
				}
				return
			},
		},
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

func checkUnionDiscriminantInfo(members []TypeName, discriminantsMap map[string]TypeName, ee *[]error) {
	covered := make([]bool, len(members))
	for _, v := range discriminantsMap {
		found := false
		for i, v2 := range members {
			if v == v2 {
				if found {
					*ee = append(*ee, fmt.Errorf("more than one discriminant pointing to member type %s", v2))
				}
				found = true
				covered[i] = true
			}
		}
		if !found {
			*ee = append(*ee, fmt.Errorf("discriminant refers to a non-member type %s", v))
		}
	}
	for i, m := range members {
		if !covered[i] {
			*ee = append(*ee, fmt.Errorf("missing discriminant info for member type %s", m))
		}
	}
}

func checkUnionDiscriminantInfo2(members []TypeName, discriminantsMap map[ipld.Kind]TypeName, ee *[]error) {
	covered := make([]bool, len(members))
	for _, v := range discriminantsMap {
		found := false
		for i, v2 := range members {
			if v == v2 {
				if found {
					*ee = append(*ee, fmt.Errorf("more than one discriminant pointing to member type %s", v2))
				}
				found = true
				covered[i] = true
			}
		}
		if !found {
			*ee = append(*ee, fmt.Errorf("discriminant refers to a non-member type %s", v))
		}
	}
	for i, m := range members {
		if !covered[i] {
			*ee = append(*ee, fmt.Errorf("missing discriminant info for member type %s", m))
		}
	}
}
