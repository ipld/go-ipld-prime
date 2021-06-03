package bindnode

import (
	"fmt"
	"reflect"
	"strings"

	ipld "github.com/ipld/go-ipld-prime"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/schema"
)

// WrapNoSchema implements an ipld.Node given a pointer to a Go value.
//
// Same rules as PrototypeNoSchema apply.
func WrapNoSchema(ptr interface{}) ipld.Node {
	panic("TODO")
	// ptrVal := reflect.ValueOf(ptr)
	// if ptrVal.Kind() != reflect.Ptr {
	// 	panic("must be a pointer")
	// }
	// return &_node{val: ptrVal.Elem()}
}

// Unwrap takes an ipld.Node implemented by one of the Wrap* or Prototype* APIs,
// and returns a pointer to the inner Go value.
//
// Unwrap returns the input node if the node isn't implemented by this package.
func Unwrap(node ipld.Node) (ptr interface{}) {
	var val reflect.Value
	switch node := node.(type) {
	case *_node:
		val = node.val
	case *_nodeRepr:
		val = node.val
	default:
		return node
	}
	if val.Kind() == reflect.Ptr {
		panic("didn't expect val to be a pointer")
	}
	if !val.CanAddr() {
		// Not addressable? Just return the interface as-is.
		// TODO: This happens in some tests, figure out why.
		return val.Interface()
	}
	return val.Addr().Interface()
}

// PrototypeNoSchema implements an ipld.NodePrototype given a Go pointer type.
//
// In this form, no IPLD schema is used; it is entirely inferred from the Go
// type.
//
// Go types map to schema types in simple ways: Go string to schema String, Go
// []byte to schema Bytes, Go struct to schema Map, and so on.
//
// A Go struct field is optional when its type is a pointer. Nullable fields are
// not supported in this mode.
func PrototypeNoSchema(ptrType interface{}) ipld.NodePrototype {
	panic("TODO")
	// typ := reflect.TypeOf(ptrType)
	// if typ.Kind() != reflect.Ptr {
	// 	panic("must be a pointer")
	// }
	// return &_prototype{goType: typ.Elem()}
}

// PrototypeOnlySchema implements an ipld.NodePrototype given an IPLD schema type.
//
// In this form, Go values are constructed with types inferred from the IPLD
// schema, like a reverse of PrototypeNoSchema.
func PrototypeOnlySchema(schemaType schema.Type) ipld.NodePrototype {
	goType := inferGoType(schemaType)
	return prototype(goType, schemaType)
}

// from IPLD Schema field names like "foo" to Go field names like "Foo".
func fieldNameFromSchema(name string) string {
	return strings.Title(name)
}

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

// Prototype implements an ipld.NodePrototype given a Go pointer type and an
// IPLD schema type.
//
// In this form, it is assumed that the Go type and IPLD schema type are
// compatible. TODO: check upfront and panic otherwise
func Prototype(ptrType interface{}, schemaType schema.Type) ipld.NodePrototype {
	goPtrType := reflect.TypeOf(ptrType)
	if goPtrType.Kind() != reflect.Ptr {
		panic("ptrType must be a pointer")
	}
	return prototype(goPtrType.Elem(), schemaType)
}

func prototype(goType reflect.Type, schemaType schema.Type) ipld.NodePrototype {
	if goType.Kind() == reflect.Invalid {
		panic("goType must be valid")
	}
	if schemaType == nil {
		panic("schemaType must not be nil")
	}
	return &_prototype{schemaType: schemaType, goType: goType}
}

var (
	_ ipld.NodePrototype = (*_prototype)(nil)

	_ ipld.Node        = (*_node)(nil)
	_ schema.TypedNode = (*_node)(nil)

	_ ipld.NodeBuilder   = (*_builder)(nil)
	_ ipld.NodeAssembler = (*_assembler)(nil)

	_ ipld.MapAssembler = (*_structAssembler)(nil)
	_ ipld.MapIterator  = (*_structIterator)(nil)

	_ ipld.ListAssembler = (*_listAssembler)(nil)
	_ ipld.ListIterator  = (*_listIterator)(nil)
)

type _prototype struct {
	schemaType schema.Type
	goType     reflect.Type // non-pointer
}

func (w *_prototype) NewBuilder() ipld.NodeBuilder {
	return &_builder{_assembler{
		schemaType: w.schemaType,
		val:        reflect.New(w.goType).Elem(),
	}}
}

// TODO: consider these Typed interfaces for the schema package

type TypedPrototype interface {
	ipld.NodePrototype

	Representation() ipld.NodePrototype
}

type TypedAssembler interface {
	ipld.NodeAssembler

	Representation() ipld.NodeAssembler
}

func (w *_prototype) Representation() ipld.NodePrototype {
	return (*_prototypeRepr)(w)
}

var (
	goTypeBool   = reflect.TypeOf(false)
	goTypeInt    = reflect.TypeOf(int(0))
	goTypeFloat  = reflect.TypeOf(0.0)
	goTypeString = reflect.TypeOf("")
	goTypeBytes  = reflect.TypeOf([]byte{})
	goTypeLink   = reflect.TypeOf((*ipld.Link)(nil)).Elem()

	schemaTypeFieldName = schema.SpawnString("fieldNameString")
)

type _node struct {
	schemaType schema.Type

	val reflect.Value // non-pointer
}

// TODO: only expose TypedNode methods if the schema was explicit.
// type _typedNode struct {
// 	_node
// }

func (w *_node) Type() schema.Type {
	return w.schemaType
}

func (w *_node) Representation() ipld.Node {
	return (*_nodeRepr)(w)
}

func (w *_node) Kind() ipld.Kind {
	return w.schemaType.TypeKind().ActsLike()
}

func (w *_node) LookupByString(key string) (ipld.Node, error) {
	switch typ := w.schemaType.(type) {
	case *schema.TypeStruct:
		field := typ.Field(key)
		if field == nil {
			return nil, ipld.ErrInvalidKey{
				TypeName: typ.Name().String(),
				Key:      basicnode.NewString(key),
			}
		}
		fval := w.val.FieldByName(fieldNameFromSchema(key))
		if !fval.IsValid() {
			panic("TODO: go-schema mismatch")
		}
		if field.IsOptional() {
			if fval.IsNil() {
				return ipld.Absent, nil
			}
			fval = fval.Elem()
		}
		if field.IsNullable() {
			if fval.IsNil() {
				return ipld.Null, nil
			}
			fval = fval.Elem()
		}
		node := &_node{
			schemaType: field.Type(),
			val:        fval,
		}
		return node, nil
	case *schema.TypeMap:
		var kval reflect.Value
		valuesVal := w.val.FieldByName("Values")
		switch ktyp := typ.KeyType().(type) {
		case *schema.TypeString:
			kval = reflect.ValueOf(key)
		default:
			asm := &_assembler{
				schemaType: ktyp,
				val:        reflect.New(valuesVal.Type().Key()).Elem(),
			}
			if err := (*_assemblerRepr)(asm).AssignString(key); err != nil {
				return nil, err
			}
			kval = asm.val
		}
		fval := valuesVal.MapIndex(kval)
		if !fval.IsValid() { // not found
			return nil, ipld.ErrNotExists{Segment: ipld.PathSegmentOfString(key)}
		}
		// TODO: Error/panic if fval.IsNil() && !typ.ValueIsNullable()?
		// Otherwise we could have two non-equal Go values (nil map,
		// non-nil-but-empty map) which represent the exact same IPLD
		// node when the field is not nullable.
		if typ.ValueIsNullable() {
			if fval.IsNil() {
				return ipld.Null, nil
			}
			fval = fval.Elem()
		}
		node := &_node{
			schemaType: typ.ValueType(),
			val:        fval,
		}
		return node, nil
	case *schema.TypeUnion:
		var idx int
		var mtyp schema.Type
		for i, member := range typ.Members() {
			if member.Name().String() == key {
				idx = i
				mtyp = member
				break
			}
		}
		if mtyp == nil { // not found
			return nil, ipld.ErrNotExists{Segment: ipld.PathSegmentOfString(key)}
		}
		haveIdx := int(w.val.FieldByName("Index").Int())
		if haveIdx != idx { // mismatching type
			return nil, ipld.ErrNotExists{Segment: ipld.PathSegmentOfString(key)}
		}
		mval := w.val.FieldByName("Value").Elem()
		node := &_node{
			schemaType: mtyp,
			val:        mval,
		}
		return node, nil
	}
	return nil, ipld.ErrWrongKind{
		TypeName:   w.schemaType.Name().String(),
		MethodName: "LookupByString",
		// TODO
	}
}

func (w *_node) LookupByNode(key ipld.Node) (ipld.Node, error) {
	panic("TODO: LookupByNode")
}

func (w *_node) LookupByIndex(idx int64) (ipld.Node, error) {
	switch typ := w.schemaType.(type) {
	case *schema.TypeList:
		if idx < 0 || int(idx) >= w.val.Len() {
			return nil, ipld.ErrNotExists{Segment: ipld.PathSegmentOfInt(idx)}
		}
		val := w.val.Index(int(idx))
		if typ.ValueIsNullable() {
			if val.IsNil() {
				return ipld.Null, nil
			}
			val = val.Elem()
		}
		return &_node{schemaType: typ.ValueType(), val: val}, nil
	}
	return nil, ipld.ErrWrongKind{
		TypeName:   w.schemaType.Name().String(),
		MethodName: "LookupByIndex",
		// TODO
	}
}

func (w *_node) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	panic("TODO: LookupBySegment")
}

func (w *_node) MapIterator() ipld.MapIterator {
	switch typ := w.schemaType.(type) {
	case *schema.TypeStruct:
		return &_structIterator{
			schemaType: typ,
			fields:     typ.Fields(),
			val:        w.val,
		}
	case *schema.TypeUnion:
		return &_unionIterator{
			schemaType: typ,
			members:    typ.Members(),
			val:        w.val,
		}
	case *schema.TypeMap:
		return &_mapIterator{
			schemaType: typ,
			keysVal:    w.val.FieldByName("Keys"),
			valuesVal:  w.val.FieldByName("Values"),
		}
	}
	return nil
}

func (w *_node) ListIterator() ipld.ListIterator {
	val := w.val
	if val.Type().Kind() == reflect.Ptr {
		if !val.IsNil() {
			val = val.Elem()
		}
	}
	switch typ := w.schemaType.(type) {
	case *schema.TypeList:
		return &_listIterator{schemaType: typ, val: val}
	}
	return nil
}

func (w *_node) Length() int64 {
	switch w.Kind() {
	case ipld.Kind_Map:
		switch typ := w.schemaType.(type) {
		case *schema.TypeStruct:
			return int64(len(typ.Fields()))
		case *schema.TypeUnion:
			return 1
		}
		return int64(w.val.FieldByName("Keys").Len())
	case ipld.Kind_List:
		return int64(w.val.Len())
	}
	return -1
}

// TODO: better story around pointers and absent/null

func (w *_node) IsAbsent() bool {
	return false
}

func (w *_node) IsNull() bool {
	return false
}

func (w *_node) AsBool() (bool, error) {
	if w.Kind() != ipld.Kind_Bool {
		return false, ipld.ErrWrongKind{
			TypeName:   w.schemaType.Name().String(),
			MethodName: "AsBool",
			// TODO
		}
	}
	return w.val.Bool(), nil
}

func (w *_node) AsInt() (int64, error) {
	if w.Kind() != ipld.Kind_Int {
		return 0, ipld.ErrWrongKind{
			TypeName:   w.schemaType.Name().String(),
			MethodName: "AsInt",
			// TODO
		}
	}
	return w.val.Int(), nil
}

func (w *_node) AsFloat() (float64, error) {
	if w.Kind() != ipld.Kind_Float {
		return 0, ipld.ErrWrongKind{
			TypeName:   w.schemaType.Name().String(),
			MethodName: "AsFloat",
			// TODO
		}
	}
	return w.val.Float(), nil
}

func (w *_node) AsString() (string, error) {
	if w.Kind() != ipld.Kind_String {
		return "", ipld.ErrWrongKind{
			TypeName:   w.schemaType.Name().String(),
			MethodName: "AsString",
			// TODO
		}
	}
	return w.val.String(), nil
}

func (w *_node) AsBytes() ([]byte, error) {
	if w.Kind() != ipld.Kind_Bytes {
		return nil, ipld.ErrWrongKind{
			TypeName:   w.schemaType.Name().String(),
			MethodName: "AsBytes",
			// TODO
		}
	}
	return w.val.Bytes(), nil
}

func (w *_node) AsLink() (ipld.Link, error) {
	if w.Kind() != ipld.Kind_Link {
		return nil, ipld.ErrWrongKind{
			TypeName:   w.schemaType.Name().String(),
			MethodName: "AsLink",
			// TODO
		}
	}
	link, _ := w.val.Interface().(ipld.Link)
	return link, nil
}

func (w *_node) Prototype() ipld.NodePrototype {
	panic("TODO: Prototype")
}

type _builder struct {
	_assembler
}

func (w *_builder) Build() ipld.Node {
	// TODO: should we panic if no Assign call was made, just like codegen?
	return &_node{schemaType: w.schemaType, val: w.val}
}

func (w *_builder) Reset() {
	panic("TODO: Reset")
}

type _assembler struct {
	schemaType schema.Type
	val        reflect.Value // non-pointer
	finish     func() error

	// kinded   bool // true if val is interface{} for a kinded union
	nullable bool // true if field or map value is nullable
}

func (w *_assembler) nonPtrVal() reflect.Value {
	val := w.val
	if w.nullable {
		val.Set(reflect.New(val.Type().Elem()))
		val = val.Elem()
	}
	return val
}

func (w *_assembler) kind() ipld.Kind {
	return w.schemaType.TypeKind().ActsLike()
}

func (w *_assembler) Representation() ipld.NodeAssembler {
	return (*_assemblerRepr)(w)
}

func (w *_assembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	switch typ := w.schemaType.(type) {
	case *schema.TypeStruct:
		val := w.nonPtrVal()
		doneFields := make([]bool, val.NumField())
		return &_structAssembler{
			schemaType: typ,
			val:        val,
			doneFields: doneFields,
			finish:     w.finish,
		}, nil
	case *schema.TypeMap:
		val := w.nonPtrVal()
		keysVal := val.FieldByName("Keys")
		valuesVal := val.FieldByName("Values")
		if valuesVal.IsNil() {
			valuesVal.Set(reflect.MakeMap(valuesVal.Type()))
		}
		return &_mapAssembler{
			schemaType: typ,
			keysVal:    keysVal,
			valuesVal:  valuesVal,
			finish:     w.finish,
		}, nil
	case *schema.TypeUnion:
		val := w.nonPtrVal()
		return &_unionAssembler{
			schemaType: typ,
			val:        val,
			finish:     w.finish,
		}, nil
	}
	return nil, ipld.ErrWrongKind{
		TypeName:   w.schemaType.Name().String(),
		MethodName: "BeginMap",
		// TODO
	}
}

func (w *_assembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	switch typ := w.schemaType.(type) {
	case *schema.TypeList:
		val := w.nonPtrVal()
		return &_listAssembler{
			schemaType: typ,
			val:        val,
			finish:     w.finish,
		}, nil
	}
	return nil, ipld.ErrWrongKind{
		TypeName:   w.schemaType.Name().String(),
		MethodName: "BeginList",
		// TODO
	}
}

func (w *_assembler) AssignNull() error {
	if !w.nullable {
		return ipld.ErrWrongKind{
			TypeName:   w.schemaType.Name().String(),
			MethodName: "AssignNull",
			// TODO
		}
	}
	w.val.Set(reflect.Zero(w.val.Type()))
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_assembler) AssignBool(b bool) error {
	if w.kind() != ipld.Kind_Bool {
		return ipld.ErrWrongKind{
			TypeName:        w.schemaType.Name().String(),
			MethodName:      "AssignBool",
			AppropriateKind: ipld.KindSet{ipld.Kind_Bool},
			ActualKind:      w.kind(),
		}
	}
	w.nonPtrVal().SetBool(b)
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_assembler) AssignInt(i int64) error {
	if w.kind() != ipld.Kind_Int {
		return ipld.ErrWrongKind{
			TypeName:        w.schemaType.Name().String(),
			MethodName:      "AssignInt",
			AppropriateKind: ipld.KindSet{ipld.Kind_Int},
			ActualKind:      w.kind(),
		}
	}
	w.nonPtrVal().SetInt(i)
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_assembler) AssignFloat(f float64) error {
	if w.kind() != ipld.Kind_Float {
		return ipld.ErrWrongKind{
			TypeName:        w.schemaType.Name().String(),
			MethodName:      "AssignFloat",
			AppropriateKind: ipld.KindSet{ipld.Kind_Float},
			ActualKind:      w.kind(),
		}
	}
	w.nonPtrVal().SetFloat(f)
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_assembler) AssignString(s string) error {
	if w.kind() != ipld.Kind_String {
		return ipld.ErrWrongKind{
			TypeName:        w.schemaType.Name().String(),
			MethodName:      "AssignString",
			AppropriateKind: ipld.KindSet{ipld.Kind_String},
			ActualKind:      w.kind(),
		}
	}
	w.nonPtrVal().SetString(s)
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_assembler) AssignBytes(p []byte) error {
	if w.kind() != ipld.Kind_Bytes {
		return ipld.ErrWrongKind{
			TypeName:        w.schemaType.Name().String(),
			MethodName:      "AssignBytes",
			AppropriateKind: ipld.KindSet{ipld.Kind_Bytes},
			ActualKind:      w.kind(),
		}
	}
	w.nonPtrVal().SetBytes(p)
	return nil
}

func (w *_assembler) AssignLink(link ipld.Link) error {
	newVal := reflect.ValueOf(link)
	if !newVal.Type().AssignableTo(w.val.Type()) {
		return ipld.ErrWrongKind{
			TypeName:   w.schemaType.Name().String(),
			MethodName: "AssignLink",
			// TODO
		}
	}
	w.nonPtrVal().Set(newVal)
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_assembler) AssignNode(node ipld.Node) error {
	// TODO: does this ever trigger?
	// newVal := reflect.ValueOf(node)
	// if newVal.Type().AssignableTo(w.val.Type()) {
	// 	w.val.Set(newVal)
	// 	return nil
	// }
	switch node.Kind() {
	case ipld.Kind_Map:
		itr := node.MapIterator()
		// TODO: consider reusing this code from elsewhere,
		// via something like ipld.BlindCopyMap.
		am, err := w.BeginMap(-1) // TODO: length?
		if err != nil {
			return err
		}
		for !itr.Done() {
			k, v, err := itr.Next()
			if err != nil {
				return err
			}
			if err := am.AssembleKey().AssignNode(k); err != nil {
				return err
			}
			if err := am.AssembleValue().AssignNode(v); err != nil {
				return err
			}
		}
		return am.Finish()
	case ipld.Kind_List:
		itr := node.ListIterator()
		am, err := w.BeginList(-1) // TODO: length?
		if err != nil {
			return err
		}
		for !itr.Done() {
			_, v, err := itr.Next()
			if err != nil {
				return err
			}
			if err := am.AssembleValue().AssignNode(v); err != nil {
				return err
			}
		}
		return am.Finish()

	case ipld.Kind_Bool:
		b, err := node.AsBool()
		if err != nil {
			return err
		}
		return w.AssignBool(b)
	case ipld.Kind_Int:
		i, err := node.AsInt()
		if err != nil {
			return err
		}
		return w.AssignInt(i)
	case ipld.Kind_Float:
		f, err := node.AsFloat()
		if err != nil {
			return err
		}
		return w.AssignFloat(f)
	case ipld.Kind_String:
		s, err := node.AsString()
		if err != nil {
			return err
		}
		return w.AssignString(s)
	case ipld.Kind_Bytes:
		p, err := node.AsBytes()
		if err != nil {
			return err
		}
		return w.AssignBytes(p)
	case ipld.Kind_Link:
		l, err := node.AsLink()
		if err != nil {
			return err
		}
		return w.AssignLink(l)
	case ipld.Kind_Null:
		return w.AssignNull()
	}
	// fmt.Println(w.val.Type(), reflect.TypeOf(node))
	panic(fmt.Sprintf("TODO: %v %v", w.val.Type(), node.Kind()))
}

func (w *_assembler) Prototype() ipld.NodePrototype {
	panic("TODO: Assembler.Prototype")
}

type _structAssembler struct {
	// TODO: embed _assembler?

	schemaType *schema.TypeStruct
	val        reflect.Value // non-pointer
	finish     func() error

	// TODO: more state checks

	// TODO: Consider if we could do this in a cheaper way,
	// such as looking at the reflect.Value directly.
	// If not, at least avoid an extra alloc.
	doneFields []bool

	// TODO: optimize for structs

	curKey _assembler

	nextIndex int // only used by repr.go
}

func (w *_structAssembler) AssembleKey() ipld.NodeAssembler {
	w.curKey = _assembler{
		schemaType: schemaTypeFieldName,
		val:        reflect.New(goTypeString).Elem(),
	}
	return &w.curKey
}

func (w *_structAssembler) AssembleValue() ipld.NodeAssembler {
	// TODO: optimize this to do one lookup by name
	name := w.curKey.val.String()
	field := w.schemaType.Field(name)
	if field == nil {
		panic(name)
		// return nil, ipld.ErrInvalidKey{
		// 	TypeName: w.schemaType.Name().String(),
		// 	Key:      basicnode.NewString(name),
		// }
	}
	ftyp, ok := w.val.Type().FieldByName(fieldNameFromSchema(name))
	if !ok {
		panic("TODO: go-schema mismatch")
	}
	if len(ftyp.Index) > 1 {
		panic("TODO: embedded fields")
	}
	w.doneFields[ftyp.Index[0]] = true
	fval := w.val.FieldByIndex(ftyp.Index)
	if field.IsOptional() {
		fval.Set(reflect.New(fval.Type().Elem()))
		fval = fval.Elem()
	}
	// TODO: reuse same assembler for perf?
	return &_assembler{
		schemaType: field.Type(),
		val:        fval,
		nullable:   field.IsNullable(),
	}
}

func (w *_structAssembler) AssembleEntry(k string) (ipld.NodeAssembler, error) {
	if err := w.AssembleKey().AssignString(k); err != nil {
		return nil, err
	}
	am := w.AssembleValue()
	return am, nil
}

func (w *_structAssembler) Finish() error {
	fields := w.schemaType.Fields()
	var missing []string
	for i, field := range fields {
		if !field.IsOptional() && !w.doneFields[i] {
			missing = append(missing, field.Name())
		}
	}
	if len(missing) > 0 {
		return ipld.ErrMissingRequiredField{Missing: missing}
	}
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_structAssembler) KeyPrototype() ipld.NodePrototype {
	return &_prototype{schemaType: schemaTypeFieldName, goType: goTypeString}
}

func (w *_structAssembler) ValuePrototype(k string) ipld.NodePrototype {
	panic("TODO: struct ValuePrototype")
}

type _mapAssembler struct {
	schemaType *schema.TypeMap
	keysVal    reflect.Value // non-pointer
	valuesVal  reflect.Value // non-pointer
	finish     func() error

	// TODO: more state checks

	curKey _assembler

	nextIndex int // only used by repr.go
}

func (w *_mapAssembler) AssembleKey() ipld.NodeAssembler {
	w.curKey = _assembler{
		schemaType: w.schemaType.KeyType(),
		val:        reflect.New(w.valuesVal.Type().Key()).Elem(),
	}
	return &w.curKey
}

func (w *_mapAssembler) AssembleValue() ipld.NodeAssembler {
	kval := w.curKey.val
	val := reflect.New(w.valuesVal.Type().Elem()).Elem()
	finish := func() error {
		// fmt.Println(kval.Interface(), val.Interface())

		// TODO: check for duplicates in keysVal
		w.keysVal.Set(reflect.Append(w.keysVal, kval))

		w.valuesVal.SetMapIndex(kval, val)
		return nil
	}
	return &_assembler{
		schemaType: w.schemaType.ValueType(),
		val:        val,
		nullable:   w.schemaType.ValueIsNullable(),
		finish:     finish,
	}
}

func (w *_mapAssembler) AssembleEntry(k string) (ipld.NodeAssembler, error) {
	if err := w.AssembleKey().AssignString(k); err != nil {
		return nil, err
	}
	am := w.AssembleValue()
	return am, nil
}

func (w *_mapAssembler) Finish() error {
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_mapAssembler) KeyPrototype() ipld.NodePrototype {
	return &_prototype{schemaType: w.schemaType.KeyType(), goType: w.valuesVal.Type().Key()}
}

func (w *_mapAssembler) ValuePrototype(k string) ipld.NodePrototype {
	panic("TODO: struct ValuePrototype")
}

type _listAssembler struct {
	schemaType *schema.TypeList
	val        reflect.Value // non-pointer
	finish     func() error
}

func (w *_listAssembler) AssembleValue() ipld.NodeAssembler {
	goType := w.val.Type().Elem()
	// TODO: use a finish func to append
	w.val.Set(reflect.Append(w.val, reflect.New(goType).Elem()))
	return &_assembler{
		schemaType: w.schemaType.ValueType(),
		val:        w.val.Index(w.val.Len() - 1),
		nullable:   w.schemaType.ValueIsNullable(),
	}
}

func (w *_listAssembler) Finish() error {
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_listAssembler) ValuePrototype(idx int64) ipld.NodePrototype {
	panic("TODO: list ValuePrototype")
}

type _unionAssembler struct {
	schemaType *schema.TypeUnion
	val        reflect.Value // non-pointer
	finish     func() error

	// TODO: more state checks

	curKey _assembler

	nextIndex int // only used by repr.go
}

func (w *_unionAssembler) AssembleKey() ipld.NodeAssembler {
	w.curKey = _assembler{
		schemaType: schemaTypeFieldName,
		val:        reflect.New(goTypeString).Elem(),
	}
	return &w.curKey
}

func (w *_unionAssembler) AssembleValue() ipld.NodeAssembler {
	name := w.curKey.val.String()
	var idx int
	var mtyp schema.Type
	for i, member := range w.schemaType.Members() {
		if member.Name().String() == name {
			idx = i
			mtyp = member
			break
		}
	}
	if mtyp == nil {
		panic("TODO: missing member")
		// return nil, ipld.ErrInvalidKey{
		// 	TypeName: w.schemaType.Name().String(),
		// 	Key:      basicnode.NewString(name),
		// }
	}
	goType := inferGoType(mtyp) // TODO: do this upfront
	val := reflect.New(goType).Elem()
	finish := func() error {
		// fmt.Println(kval.Interface(), val.Interface())
		w.val.FieldByName("Index").SetInt(int64(idx))
		w.val.FieldByName("Value").Set(val)
		return nil
	}
	return &_assembler{
		schemaType: mtyp,
		val:        val,
		finish:     finish,
	}
}

func (w *_unionAssembler) AssembleEntry(k string) (ipld.NodeAssembler, error) {
	if err := w.AssembleKey().AssignString(k); err != nil {
		return nil, err
	}
	am := w.AssembleValue()
	return am, nil
}

func (w *_unionAssembler) Finish() error {
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_unionAssembler) KeyPrototype() ipld.NodePrototype {
	return &_prototype{schemaType: schemaTypeFieldName, goType: goTypeString}
}

func (w *_unionAssembler) ValuePrototype(k string) ipld.NodePrototype {
	panic("TODO: struct ValuePrototype")
}

type _structIterator struct {
	// TODO: support embedded fields?
	schemaType *schema.TypeStruct
	fields     []schema.StructField
	val        reflect.Value // non-pointer
	nextIndex  int

	// these are only used in repr.go
	reprEnd int
}

func (w *_structIterator) Next() (key, value ipld.Node, _ error) {
	if w.Done() {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	field := w.fields[w.nextIndex]
	val := w.val.Field(w.nextIndex)
	w.nextIndex++
	key = basicnode.NewString(field.Name())
	if field.IsOptional() {
		if val.IsNil() {
			return key, ipld.Absent, nil
		}
		val = val.Elem()
	}
	if field.IsNullable() {
		if val.IsNil() {
			return key, ipld.Null, nil
		}
		val = val.Elem()
	}
	node := &_node{
		schemaType: field.Type(),
		val:        val,
	}
	return key, node, nil
}

func (w *_structIterator) Done() bool {
	return w.nextIndex >= len(w.fields)
}

type _mapIterator struct {
	schemaType *schema.TypeMap
	keysVal    reflect.Value // non-pointer
	valuesVal  reflect.Value // non-pointer
	nextIndex  int

	// these are only used in repr.go
	reprEnd int
}

func (w *_mapIterator) Next() (key, value ipld.Node, _ error) {
	if w.Done() {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	goKey := w.keysVal.Index(w.nextIndex)
	val := w.valuesVal.MapIndex(goKey)
	w.nextIndex++

	key = &_node{
		schemaType: w.schemaType.KeyType(),
		val:        goKey,
	}
	if w.schemaType.ValueIsNullable() {
		if val.IsNil() {
			return key, ipld.Null, nil
		}
		val = val.Elem()
	}
	node := &_node{
		schemaType: w.schemaType.ValueType(),
		val:        val,
	}
	return key, node, nil
}

func (w *_mapIterator) Done() bool {
	return w.nextIndex >= w.keysVal.Len()
}

type _listIterator struct {
	schemaType *schema.TypeList
	val        reflect.Value // non-pointer
	nextIndex  int
}

func (w *_listIterator) Next() (index int64, value ipld.Node, _ error) {
	if w.Done() {
		return 0, nil, ipld.ErrIteratorOverread{}
	}
	idx := int64(w.nextIndex)
	val := w.val.Index(w.nextIndex)
	w.nextIndex++
	if w.schemaType.ValueIsNullable() {
		if val.IsNil() {
			return idx, ipld.Null, nil
		}
		val = val.Elem()
	}
	return idx, &_node{schemaType: w.schemaType.ValueType(), val: val}, nil
}

func (w *_listIterator) Done() bool {
	return w.nextIndex >= w.val.Len()
}

type _unionIterator struct {
	// TODO: support embedded fields?
	schemaType *schema.TypeUnion
	members    []schema.Type
	val        reflect.Value // non-pointer

	done bool
}

func (w *_unionIterator) Next() (key, value ipld.Node, _ error) {
	if w.Done() {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	w.done = true

	haveIdx := int(w.val.FieldByName("Index").Int())
	mtyp := w.members[haveIdx]
	mval := w.val.FieldByName("Value").Elem()

	node := &_node{
		schemaType: mtyp,
		val:        mval,
	}
	key = basicnode.NewString(mtyp.Name().String())
	return key, node, nil
}

func (w *_unionIterator) Done() bool {
	return w.done
}

// TODO: consider making our own Node interface, like:
//
// type WrappedNode interface {
//     ipld.Node
//     Unwrap() (ptr interface)
// }
//
// Pros: API is easier to understand, harder to mix up with other ipld.Nodes.
// Cons: One usually only has an ipld.Node, and type assertions can be weird.
