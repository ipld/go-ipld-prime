package schema

import (
	"fmt"
	"path"

	"github.com/ipld/go-ipld-prime"
)

// FUTURE: we also want something *almost* identical to this Validate method,
//  but returning a `typed.Node` in the case of no error.
//  (Such a method would go in the same package as `typed.Node`, presumably.)
// How do we avoid writing this method twice?
//  Maybe both a Validate and Reify method belong in `typed` package,
//   and Validate just returns less?
//  No... Reify should probably short-circuit sooner?
// Unclear.  Guess first step is that we need to decide the intended UX!

func Validate(ts TypeSystem, t Type, node ipld.Node) []error {
	return validate(ts, t, node, "/")
}

// review: 'ts' param might not actually be necessary; everything relevant can be reached from t so far.
func validate(ts TypeSystem, t Type, node ipld.Node, pth string) []error {
	switch t2 := t.(type) {
	case TypeBool:
		if node.ReprKind() != ipld.ReprKind_Bool {
			return []error{fmt.Errorf("Schema match failed: expected type %q (which is kind %v) at path %q, but found kind %v", t2.Name(), t.ReprKind(), pth, node.ReprKind())}
		}
		return nil
	case TypeString:
		if node.ReprKind() != ipld.ReprKind_String {
			return []error{fmt.Errorf("Schema match failed: expected type %q (which is kind %v) at path %q, but found kind %v", t2.Name(), t.ReprKind(), pth, node.ReprKind())}
		}
		return nil
	case TypeBytes:
		if node.ReprKind() != ipld.ReprKind_Bytes {
			return []error{fmt.Errorf("Schema match failed: expected type %q (which is kind %v) at path %q, but found kind %v", t2.Name(), t.ReprKind(), pth, node.ReprKind())}
		}
		return nil
	case TypeInt:
		if node.ReprKind() != ipld.ReprKind_Int {
			return []error{fmt.Errorf("Schema match failed: expected type %q (which is kind %v) at path %q, but found kind %v", t2.Name(), t.ReprKind(), pth, node.ReprKind())}
		}
		return nil
	case TypeFloat:
		if node.ReprKind() != ipld.ReprKind_Float {
			return []error{fmt.Errorf("Schema match failed: expected type %q (which is kind %v) at path %q, but found kind %v", t2.Name(), t.ReprKind(), pth, node.ReprKind())}
		}
		return nil
	case TypeMap:
		if node.ReprKind() != ipld.ReprKind_Map {
			return []error{fmt.Errorf("Schema match failed: expected type %q (which is kind %v) at path %q, but found kind %v", t2.Name(), t.ReprKind(), pth, node.ReprKind())}
		}
		errs := []error(nil)
		for itr := node.MapIterator(); !itr.Done(); {
			k, v, err := itr.Next()
			if err != nil {
				return []error{err}
			}
			// FUTURE: if KeyType is an enum rather than string, do membership check.
			ks, _ := k.AsString()
			if v.IsNull() {
				if !t2.ValueIsNullable() {
					errs = append(errs, fmt.Errorf("Schema match failed: map at path %q contains unpermitted null in key %q", pth, ks))
				}
			} else {
				errs = append(errs, validate(ts, t2.ValueType(), v, path.Join(pth, ks))...)
			}
		}
		return errs
	case TypeList:
	case TypeLink:
		// TODO interesting case: would need resolver to keep checking.
	case TypeUnion:
		// TODO *several* interesting errors
	case TypeStruct:
		switch t2.tupleStyle {
		case false: // as map!
			if node.ReprKind() != ipld.ReprKind_Map {
				return []error{fmt.Errorf("Schema match failed: expected type %q (which is kind %v) at path %q, but found kind %v", t2.Name(), t.ReprKind(), pth, node.ReprKind())}
			}
			// TODO loop over em
			// TODO REVIEW order strictness questions?
		case true: // as array!

		}
	case TypeEnum:
		// TODO another interesting error
	}
	return nil
}
