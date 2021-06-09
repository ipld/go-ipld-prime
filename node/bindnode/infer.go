package bindnode

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ipld/go-ipld-prime/schema"
)

// Consider exposing these APIs later, if they might be useful.

func inferGoType(typ schema.Type) reflect.Type {
	switch typ := typ.(type) {
	case *schema.TypeBool:
		return goTypeBool
	case *schema.TypeInt:
		return goTypeInt
	case *schema.TypeFloat:
		return goTypeFloat
	case *schema.TypeString:
		return goTypeString
	case *schema.TypeBytes:
		return goTypeBytes
	case *schema.TypeStruct:
		fields := typ.Fields()
		goFields := make([]reflect.StructField, len(fields))
		for i, field := range fields {
			ftyp := inferGoType(field.Type())
			if field.IsNullable() {
				ftyp = reflect.PtrTo(ftyp)
			}
			if field.IsOptional() {
				ftyp = reflect.PtrTo(ftyp)
			}
			goFields[i] = reflect.StructField{
				Name: fieldNameFromSchema(field.Name()),
				Type: ftyp,
			}
		}
		return reflect.StructOf(goFields)
	case *schema.TypeMap:
		ktyp := inferGoType(typ.KeyType())
		vtyp := inferGoType(typ.ValueType())
		if typ.ValueIsNullable() {
			vtyp = reflect.PtrTo(vtyp)
		}
		// We need an extra field to keep the map ordered,
		// since IPLD maps must have stable iteration order.
		// We could sort when iterating, but that's expensive.
		// Keeping the insertion order is easy and intuitive.
		//
		//	struct {
		//		Keys   []K
		//		Values map[K]V
		//	}
		goFields := []reflect.StructField{
			{
				Name: "Keys",
				Type: reflect.SliceOf(ktyp),
			},
			{
				Name: "Values",
				Type: reflect.MapOf(ktyp, vtyp),
			},
		}
		return reflect.StructOf(goFields)
	case *schema.TypeList:
		etyp := inferGoType(typ.ValueType())
		if typ.ValueIsNullable() {
			etyp = reflect.PtrTo(etyp)
		}
		return reflect.SliceOf(etyp)
	case *schema.TypeUnion:
		// We need an extra field to record what member we stored.
		type goUnion struct {
			Index int // 0..len(typ.Members)-1
			Value interface{}
		}
		return reflect.TypeOf(goUnion{})
	}
	panic(fmt.Sprintf("%T\n", typ))
}

// from IPLD Schema field names like "foo" to Go field names like "Foo".
func fieldNameFromSchema(name string) string {
	return strings.Title(name)
}

func inferSchema(typ reflect.Type) schema.Type {
	panic("TODO")
}
