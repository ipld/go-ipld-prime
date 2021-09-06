package bindnode

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

var (
	goTypeBool   = reflect.TypeOf(false)
	goTypeInt    = reflect.TypeOf(int(0))
	goTypeFloat  = reflect.TypeOf(0.0)
	goTypeString = reflect.TypeOf("")
	goTypeBytes  = reflect.TypeOf([]byte{})
	goTypeLink   = reflect.TypeOf((*datamodel.Link)(nil)).Elem()

	schemaTypeBool   = schema.SpawnBool("Bool")
	schemaTypeInt    = schema.SpawnInt("Int")
	schemaTypeFloat  = schema.SpawnFloat("Float")
	schemaTypeString = schema.SpawnString("String")
	schemaTypeBytes  = schema.SpawnBytes("Bytes")
	schemaTypeLink   = schema.SpawnLink("Link")
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
		fieldsGo := make([]reflect.StructField, len(fields))
		for i, field := range fields {
			ftypGo := inferGoType(field.Type())
			if field.IsNullable() {
				ftypGo = reflect.PtrTo(ftypGo)
			}
			if field.IsOptional() {
				ftypGo = reflect.PtrTo(ftypGo)
			}
			fieldsGo[i] = reflect.StructField{
				Name: fieldNameFromSchema(field.Name()),
				Type: ftypGo,
			}
		}
		return reflect.StructOf(fieldsGo)
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
		fieldsGo := []reflect.StructField{
			{
				Name: "Keys",
				Type: reflect.SliceOf(ktyp),
			},
			{
				Name: "Values",
				Type: reflect.MapOf(ktyp, vtyp),
			},
		}
		return reflect.StructOf(fieldsGo)
	case *schema.TypeList:
		etyp := inferGoType(typ.ValueType())
		if typ.ValueIsNullable() {
			etyp = reflect.PtrTo(etyp)
		}
		return reflect.SliceOf(etyp)
	case *schema.TypeUnion:
		// type goUnion struct {
		// 	Type1 *Type1
		// 	Type2 *Type2
		// 	...
		// }
		members := typ.Members()
		fieldsGo := make([]reflect.StructField, len(members))
		for i, ftyp := range members {
			ftypGo := inferGoType(ftyp)
			fieldsGo[i] = reflect.StructField{
				Name: fieldNameFromSchema(string(ftyp.Name())),
				Type: reflect.PtrTo(ftypGo),
			}
		}
		return reflect.StructOf(fieldsGo)
	}
	panic(fmt.Sprintf("%T\n", typ))
}

// from IPLD Schema field names like "foo" to Go field names like "Foo".
func fieldNameFromSchema(name string) string {
	return strings.Title(name)
}

var defaultTypeSystem schema.TypeSystem

func init() {
	defaultTypeSystem.Init()

	defaultTypeSystem.Accumulate(schemaTypeBool)
	defaultTypeSystem.Accumulate(schemaTypeInt)
	defaultTypeSystem.Accumulate(schemaTypeFloat)
	defaultTypeSystem.Accumulate(schemaTypeString)
	defaultTypeSystem.Accumulate(schemaTypeBytes)
	defaultTypeSystem.Accumulate(schemaTypeLink)
}

// TODO: support IPLD maps and unions in inferSchema

// TODO: support bringing your own TypeSystem?

// TODO: we should probably avoid re-spawning the same types if the TypeSystem
// has them, and test that that works as expected

func inferSchema(typ reflect.Type) schema.Type {
	switch typ.Kind() {
	case reflect.Bool:
		return schemaTypeBool
	case reflect.Int64:
		return schemaTypeInt
	case reflect.Float64:
		return schemaTypeFloat
	case reflect.String:
		return schemaTypeString
	case reflect.Struct:
		fieldsSchema := make([]schema.StructField, typ.NumField())
		for i := range fieldsSchema {
			field := typ.Field(i)
			ftyp := field.Type
			ftypSchema := inferSchema(ftyp)
			fieldsSchema[i] = schema.SpawnStructField(
				field.Name, // TODO: allow configuring the name with tags
				ftypSchema.Name(),

				// TODO: support nullable/optional with tags
				false,
				false,
			)
		}
		name := schema.TypeName(typ.Name())
		if name == "" {
			panic("TODO: anonymous composite types")
		}
		typSchema := schema.SpawnStruct(name, fieldsSchema, nil)
		defaultTypeSystem.Accumulate(typSchema)
		return typSchema
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			// Special case for []byte.
			return schemaTypeBytes
		}

		etyp := typ.Elem()
		nullable := false
		if etyp.Kind() == reflect.Ptr {
			etyp = etyp.Elem()
			nullable = true
		}
		etypSchema := inferSchema(typ.Elem())
		name := schema.TypeName(typ.Name())
		if name == "" {
			name = "List_" + etypSchema.Name()
		}
		typSchema := schema.SpawnList(name, etypSchema.Name(), nullable)
		defaultTypeSystem.Accumulate(typSchema)
		return typSchema
	}
	panic(fmt.Sprintf("%s\n", typ.Kind()))
}
