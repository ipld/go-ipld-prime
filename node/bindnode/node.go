package bindnode

import (
	"fmt"
	"reflect"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/schema"
)

// Assert that we implement all the interfaces as expected.
// Grouped by the interfaces to implement, roughly.
var (
	_ datamodel.NodePrototype = (*_prototype)(nil)
	_ schema.TypedPrototype   = (*_prototype)(nil)
	_ datamodel.NodePrototype = (*_prototypeRepr)(nil)

	_ datamodel.Node   = (*_node)(nil)
	_ schema.TypedNode = (*_node)(nil)
	_ datamodel.Node   = (*_nodeRepr)(nil)

	_ datamodel.NodeBuilder   = (*_builder)(nil)
	_ datamodel.NodeAssembler = (*_assembler)(nil)
	_ datamodel.NodeBuilder   = (*_builderRepr)(nil)
	_ datamodel.NodeAssembler = (*_assemblerRepr)(nil)

	_ datamodel.MapAssembler = (*_structAssembler)(nil)
	_ datamodel.MapIterator  = (*_structIterator)(nil)
	_ datamodel.MapAssembler = (*_structAssemblerRepr)(nil)
	_ datamodel.MapIterator  = (*_structIteratorRepr)(nil)

	_ datamodel.ListAssembler = (*_listAssembler)(nil)
	_ datamodel.ListIterator  = (*_listIterator)(nil)
	_ datamodel.ListAssembler = (*_listAssemblerRepr)(nil)

	_ datamodel.MapAssembler = (*_unionAssembler)(nil)
	_ datamodel.MapIterator  = (*_unionIterator)(nil)
	_ datamodel.MapAssembler = (*_unionAssemblerRepr)(nil)
	_ datamodel.MapIterator  = (*_unionIteratorRepr)(nil)
)

type _prototype struct {
	schemaType schema.Type
	goType     reflect.Type // non-pointer
}

func (w *_prototype) NewBuilder() datamodel.NodeBuilder {
	return &_builder{_assembler{
		schemaType: w.schemaType,
		val:        reflect.New(w.goType).Elem(),
	}}
}

func (w *_prototype) Type() schema.Type {
	return w.schemaType
}

func (w *_prototype) Representation() datamodel.NodePrototype {
	return (*_prototypeRepr)(w)
}

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

func (w *_node) Representation() datamodel.Node {
	return (*_nodeRepr)(w)
}

func (w *_node) Kind() datamodel.Kind {
	return w.schemaType.TypeKind().ActsLike()
}

func (w *_node) LookupByString(key string) (datamodel.Node, error) {
	switch typ := w.schemaType.(type) {
	case *schema.TypeStruct:
		field := typ.Field(key)
		if field == nil {
			return nil, schema.ErrInvalidKey{
				TypeName: typ.Name(),
				Key:      basicnode.NewString(key),
			}
		}
		fval := w.val.FieldByName(fieldNameFromSchema(key))
		if !fval.IsValid() {
			panic("TODO: go-schema mismatch")
		}
		if field.IsOptional() {
			if fval.IsNil() {
				return datamodel.Absent, nil
			}
			fval = fval.Elem()
		}
		if field.IsNullable() {
			if fval.IsNil() {
				return datamodel.Null, nil
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
			return nil, datamodel.ErrNotExists{Segment: datamodel.PathSegmentOfString(key)}
		}
		// TODO: Error/panic if fval.IsNil() && !typ.ValueIsNullable()?
		// Otherwise we could have two non-equal Go values (nil map,
		// non-nil-but-empty map) which represent the exact same IPLD
		// node when the field is not nullable.
		if typ.ValueIsNullable() {
			if fval.IsNil() {
				return datamodel.Null, nil
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
			if member.Name() == key {
				idx = i
				mtyp = member
				break
			}
		}
		if mtyp == nil { // not found
			return nil, datamodel.ErrNotExists{Segment: datamodel.PathSegmentOfString(key)}
		}
		// TODO: we could look up the right Go field straight away via idx.
		haveIdx, mval := unionMember(w.val)
		if haveIdx != idx { // mismatching type
			return nil, datamodel.ErrNotExists{Segment: datamodel.PathSegmentOfString(key)}
		}
		node := &_node{
			schemaType: mtyp,
			val:        mval,
		}
		return node, nil
	}
	return nil, datamodel.ErrWrongKind{
		TypeName:        w.schemaType.Name(),
		MethodName:      "LookupByString",
		AppropriateKind: datamodel.KindSet_JustMap,
		ActualKind:      w.Kind(),
	}
}

var invalidValue reflect.Value

func unionMember(val reflect.Value) (int, reflect.Value) {
	// The first non-nil field is a match.
	for i := 0; i < val.NumField(); i++ {
		elemVal := val.Field(i)
		if elemVal.Kind() != reflect.Ptr {
			panic("bindnode: found unexpected non-pointer in a union field")
		}
		if elemVal.IsNil() {
			continue
		}
		return i, elemVal.Elem()
	}
	return -1, invalidValue
}

func unionSetMember(val reflect.Value, memberIdx int, memberPtr reflect.Value) {
	// Reset the entire union struct to zero, to clear any non-nil pointers.
	val.Set(reflect.Zero(val.Type()))

	// Set the index pointer to the given value.
	val.Field(memberIdx).Set(memberPtr)
}

func (w *_node) LookupByIndex(idx int64) (datamodel.Node, error) {
	switch typ := w.schemaType.(type) {
	case *schema.TypeList:
		if idx < 0 || int(idx) >= w.val.Len() {
			return nil, datamodel.ErrNotExists{Segment: datamodel.PathSegmentOfInt(idx)}
		}
		val := w.val.Index(int(idx))
		if typ.ValueIsNullable() {
			if val.IsNil() {
				return datamodel.Null, nil
			}
			val = val.Elem()
		}
		return &_node{schemaType: typ.ValueType(), val: val}, nil
	}
	return nil, datamodel.ErrWrongKind{
		TypeName:        w.schemaType.Name(),
		MethodName:      "LookupByIndex",
		AppropriateKind: datamodel.KindSet_JustList,
		ActualKind:      w.Kind(),
	}
}

func (w *_node) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	switch w.Kind() {
	case datamodel.Kind_Map:
		return w.LookupByString(seg.String())
	case datamodel.Kind_List:
		idx, err := seg.Index()
		if err != nil {
			return nil, err
		}
		return w.LookupByIndex(idx)
	}
	return nil, datamodel.ErrWrongKind{
		TypeName:        w.schemaType.Name(),
		MethodName:      "LookupBySegment",
		AppropriateKind: datamodel.KindSet_Recursive,
		ActualKind:      w.Kind(),
	}
}

func (w *_node) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	switch w.Kind() {
	case datamodel.Kind_Map:
		s, err := key.AsString()
		if err != nil {
			return nil, err
		}
		return w.LookupByString(s)
	case datamodel.Kind_List:
		i, err := key.AsInt()
		if err != nil {
			return nil, err
		}
		return w.LookupByIndex(i)
	}
	return nil, datamodel.ErrWrongKind{
		TypeName:        w.schemaType.Name(),
		MethodName:      "LookupByNode",
		AppropriateKind: datamodel.KindSet_Recursive,
		ActualKind:      w.Kind(),
	}
}

func (w *_node) MapIterator() datamodel.MapIterator {
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

func (w *_node) ListIterator() datamodel.ListIterator {
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
	case datamodel.Kind_Map:
		switch typ := w.schemaType.(type) {
		case *schema.TypeStruct:
			return int64(len(typ.Fields()))
		case *schema.TypeUnion:
			return 1
		}
		return int64(w.val.FieldByName("Keys").Len())
	case datamodel.Kind_List:
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
	if w.Kind() != datamodel.Kind_Bool {
		return false, datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AsBool",
			AppropriateKind: datamodel.KindSet_JustBool,
			ActualKind:      w.Kind(),
		}
	}
	return w.val.Bool(), nil
}

func (w *_node) AsInt() (int64, error) {
	if w.Kind() != datamodel.Kind_Int {
		return 0, datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AsInt",
			AppropriateKind: datamodel.KindSet_JustInt,
			ActualKind:      w.Kind(),
		}
	}
	return w.val.Int(), nil
}

func (w *_node) AsFloat() (float64, error) {
	if w.Kind() != datamodel.Kind_Float {
		return 0, datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AsFloat",
			AppropriateKind: datamodel.KindSet_JustFloat,
			ActualKind:      w.Kind(),
		}
	}
	return w.val.Float(), nil
}

func (w *_node) AsString() (string, error) {
	if w.Kind() != datamodel.Kind_String {
		return "", datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AsString",
			AppropriateKind: datamodel.KindSet_JustString,
			ActualKind:      w.Kind(),
		}
	}
	return w.val.String(), nil
}

func (w *_node) AsBytes() ([]byte, error) {
	if w.Kind() != datamodel.Kind_Bytes {
		return nil, datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AsBytes",
			AppropriateKind: datamodel.KindSet_JustBytes,
			ActualKind:      w.Kind(),
		}
	}
	return w.val.Bytes(), nil
}

func (w *_node) AsLink() (datamodel.Link, error) {
	if w.Kind() != datamodel.Kind_Link {
		return nil, datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AsLink",
			AppropriateKind: datamodel.KindSet_JustLink,
			ActualKind:      w.Kind(),
		}
	}
	switch val := w.val.Interface().(type) {
	case datamodel.Link:
		return val, nil
	case cid.Cid:
		return cidlink.Link{Cid: val}, nil
	default:
		panic(fmt.Sprintf("bindnode: unexpected link type %T", val))
	}
}

func (w *_node) Prototype() datamodel.NodePrototype {
	return &_prototype{schemaType: w.schemaType, goType: w.val.Type()}
}

type _builder struct {
	_assembler
}

func (w *_builder) Build() datamodel.Node {
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

func (w *_assembler) kind() datamodel.Kind {
	return w.schemaType.TypeKind().ActsLike()
}

func (w *_assembler) Representation() datamodel.NodeAssembler {
	return (*_assemblerRepr)(w)
}

func (w *_assembler) BeginMap(sizeHint int64) (datamodel.MapAssembler, error) {
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
	return nil, datamodel.ErrWrongKind{
		TypeName:        w.schemaType.Name(),
		MethodName:      "BeginMap",
		AppropriateKind: datamodel.KindSet_JustMap,
		ActualKind:      w.kind(),
	}
}

func (w *_assembler) BeginList(sizeHint int64) (datamodel.ListAssembler, error) {
	switch typ := w.schemaType.(type) {
	case *schema.TypeList:
		val := w.nonPtrVal()
		return &_listAssembler{
			schemaType: typ,
			val:        val,
			finish:     w.finish,
		}, nil
	}
	return nil, datamodel.ErrWrongKind{
		TypeName:        w.schemaType.Name(),
		MethodName:      "BeginList",
		AppropriateKind: datamodel.KindSet_JustList,
		ActualKind:      w.kind(),
	}
}

func (w *_assembler) AssignNull() error {
	if !w.nullable {
		return datamodel.ErrWrongKind{
			TypeName:   w.schemaType.Name(),
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
	if w.kind() != datamodel.Kind_Bool {
		return datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AssignBool",
			AppropriateKind: datamodel.KindSet{datamodel.Kind_Bool},
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
	if w.kind() != datamodel.Kind_Int {
		return datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AssignInt",
			AppropriateKind: datamodel.KindSet{datamodel.Kind_Int},
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
	if w.kind() != datamodel.Kind_Float {
		return datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AssignFloat",
			AppropriateKind: datamodel.KindSet{datamodel.Kind_Float},
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
	if w.kind() != datamodel.Kind_String {
		return datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AssignString",
			AppropriateKind: datamodel.KindSet{datamodel.Kind_String},
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
	if w.kind() != datamodel.Kind_Bytes {
		return datamodel.ErrWrongKind{
			TypeName:        w.schemaType.Name(),
			MethodName:      "AssignBytes",
			AppropriateKind: datamodel.KindSet{datamodel.Kind_Bytes},
			ActualKind:      w.kind(),
		}
	}
	w.nonPtrVal().SetBytes(p)
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_assembler) AssignLink(link datamodel.Link) error {
	newVal := reflect.ValueOf(link)
	if !newVal.Type().AssignableTo(w.val.Type()) {
		if newVal.Type() == goTypeCidLink && goTypeCid.AssignableTo(w.val.Type()) {
			// Unbox a cidlink.Link to assign to a go-cid.Cid value.
			newVal = newVal.FieldByName("Cid")
		} else {
			// The target value cannot be assigned a datamodel.Link or go-cid.Cid.
			return datamodel.ErrWrongKind{
				TypeName:        w.schemaType.Name(),
				MethodName:      "AssignLink",
				AppropriateKind: datamodel.KindSet_JustLink,
				ActualKind:      w.kind(),
			}
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

func (w *_assembler) AssignNode(node datamodel.Node) error {
	// TODO: does this ever trigger?
	// newVal := reflect.ValueOf(node)
	// if newVal.Type().AssignableTo(w.val.Type()) {
	// 	w.val.Set(newVal)
	// 	return nil
	// }
	switch node.Kind() {
	case datamodel.Kind_Map:
		itr := node.MapIterator()
		// TODO: consider reusing this code from elsewhere,
		// via something like datamodel.BlindCopyMap.
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
	case datamodel.Kind_List:
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

	case datamodel.Kind_Bool:
		b, err := node.AsBool()
		if err != nil {
			return err
		}
		return w.AssignBool(b)
	case datamodel.Kind_Int:
		i, err := node.AsInt()
		if err != nil {
			return err
		}
		return w.AssignInt(i)
	case datamodel.Kind_Float:
		f, err := node.AsFloat()
		if err != nil {
			return err
		}
		return w.AssignFloat(f)
	case datamodel.Kind_String:
		s, err := node.AsString()
		if err != nil {
			return err
		}
		return w.AssignString(s)
	case datamodel.Kind_Bytes:
		p, err := node.AsBytes()
		if err != nil {
			return err
		}
		return w.AssignBytes(p)
	case datamodel.Kind_Link:
		l, err := node.AsLink()
		if err != nil {
			return err
		}
		return w.AssignLink(l)
	case datamodel.Kind_Null:
		return w.AssignNull()
	}
	return fmt.Errorf("AssignNode TODO: %v %v", w.val.Type(), node.Kind())
}

func (w *_assembler) Prototype() datamodel.NodePrototype {
	return &_prototype{schemaType: w.schemaType, goType: w.val.Type()}
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

func (w *_structAssembler) AssembleKey() datamodel.NodeAssembler {
	w.curKey = _assembler{
		schemaType: schemaTypeString,
		val:        reflect.New(goTypeString).Elem(),
	}
	return &w.curKey
}

func (w *_structAssembler) AssembleValue() datamodel.NodeAssembler {
	// TODO: optimize this to do one lookup by name
	name := w.curKey.val.String()
	field := w.schemaType.Field(name)
	if field == nil {
		// TODO: should've been raised when the key was submitted (we have room to return errors there, but can only panic at this point in the game).
		// TODO: should make well-typed errors for this.
		panic(fmt.Sprintf("TODO: invalid key: %q is not a field in type %s", name, w.schemaType.Name()))
		// panic(schema.ErrInvalidKey{
		// 	TypeName: w.schemaType.Name(),
		// 	Key:      basicnode.NewString(name),
		// })
	}
	ftyp, ok := w.val.Type().FieldByName(fieldNameFromSchema(name))
	if !ok {
		// It is unfortunate this is not detected proactively earlier during bind.
		panic(fmt.Sprintf("schema type %q has field %q, we expect go struct to have field %q", w.schemaType.Name(), field.Name(), fieldNameFromSchema(name)))
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

func (w *_structAssembler) AssembleEntry(k string) (datamodel.NodeAssembler, error) {
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
		return schema.ErrMissingRequiredField{Missing: missing}
	}
	if w.finish != nil {
		if err := w.finish(); err != nil {
			return err
		}
	}
	return nil
}

func (w *_structAssembler) KeyPrototype() datamodel.NodePrototype {
	// TODO: if the user provided their own schema with their own typesystem,
	// the schemaTypeString here may be using the wrong typesystem.
	return &_prototype{schemaType: schemaTypeString, goType: goTypeString}
}

func (w *_structAssembler) ValuePrototype(k string) datamodel.NodePrototype {
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

func (w *_mapAssembler) AssembleKey() datamodel.NodeAssembler {
	w.curKey = _assembler{
		schemaType: w.schemaType.KeyType(),
		val:        reflect.New(w.valuesVal.Type().Key()).Elem(),
	}
	return &w.curKey
}

func (w *_mapAssembler) AssembleValue() datamodel.NodeAssembler {
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

func (w *_mapAssembler) AssembleEntry(k string) (datamodel.NodeAssembler, error) {
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

func (w *_mapAssembler) KeyPrototype() datamodel.NodePrototype {
	return &_prototype{schemaType: w.schemaType.KeyType(), goType: w.valuesVal.Type().Key()}
}

func (w *_mapAssembler) ValuePrototype(k string) datamodel.NodePrototype {
	return &_prototype{schemaType: w.schemaType.ValueType(), goType: w.valuesVal.Type().Elem()}
}

type _listAssembler struct {
	schemaType *schema.TypeList
	val        reflect.Value // non-pointer
	finish     func() error
}

func (w *_listAssembler) AssembleValue() datamodel.NodeAssembler {
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

func (w *_listAssembler) ValuePrototype(idx int64) datamodel.NodePrototype {
	return &_prototype{schemaType: w.schemaType.ValueType(), goType: w.val.Type().Elem()}
}

type _unionAssembler struct {
	schemaType *schema.TypeUnion
	val        reflect.Value // non-pointer
	finish     func() error

	// TODO: more state checks

	curKey _assembler

	nextIndex int // only used by repr.go
}

func (w *_unionAssembler) AssembleKey() datamodel.NodeAssembler {
	w.curKey = _assembler{
		schemaType: schemaTypeString,
		val:        reflect.New(goTypeString).Elem(),
	}
	return &w.curKey
}

func (w *_unionAssembler) AssembleValue() datamodel.NodeAssembler {
	name := w.curKey.val.String()
	var idx int
	var mtyp schema.Type
	for i, member := range w.schemaType.Members() {
		if member.Name() == name {
			idx = i
			mtyp = member
			break
		}
	}
	if mtyp == nil {
		panic(fmt.Sprintf("TODO: missing member %s in %s", name, w.schemaType.Name()))
		// return nil, datamodel.ErrInvalidKey{
		// 	TypeName: w.schemaType.Name(),
		// 	Key:      basicnode.NewString(name),
		// }
	}

	goType := w.val.Field(idx).Type().Elem()
	valPtr := reflect.New(goType)
	finish := func() error {
		// fmt.Println(kval.Interface(), val.Interface())
		unionSetMember(w.val, idx, valPtr)
		return nil
	}
	return &_assembler{
		schemaType: mtyp,
		val:        valPtr.Elem(),
		finish:     finish,
	}
}

func (w *_unionAssembler) AssembleEntry(k string) (datamodel.NodeAssembler, error) {
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

func (w *_unionAssembler) KeyPrototype() datamodel.NodePrototype {
	return &_prototype{schemaType: schemaTypeString, goType: goTypeString}
}

func (w *_unionAssembler) ValuePrototype(k string) datamodel.NodePrototype {
	panic("TODO: union ValuePrototype")
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

func (w *_structIterator) Next() (key, value datamodel.Node, _ error) {
	if w.Done() {
		return nil, nil, datamodel.ErrIteratorOverread{}
	}
	field := w.fields[w.nextIndex]
	val := w.val.Field(w.nextIndex)
	w.nextIndex++
	key = basicnode.NewString(field.Name())
	if field.IsOptional() {
		if val.IsNil() {
			return key, datamodel.Absent, nil
		}
		val = val.Elem()
	}
	if field.IsNullable() {
		if val.IsNil() {
			return key, datamodel.Null, nil
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
}

func (w *_mapIterator) Next() (key, value datamodel.Node, _ error) {
	if w.Done() {
		return nil, nil, datamodel.ErrIteratorOverread{}
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
			return key, datamodel.Null, nil
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

func (w *_listIterator) Next() (index int64, value datamodel.Node, _ error) {
	if w.Done() {
		return 0, nil, datamodel.ErrIteratorOverread{}
	}
	idx := int64(w.nextIndex)
	val := w.val.Index(w.nextIndex)
	w.nextIndex++
	if w.schemaType.ValueIsNullable() {
		if val.IsNil() {
			return idx, datamodel.Null, nil
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

func (w *_unionIterator) Next() (key, value datamodel.Node, _ error) {
	if w.Done() {
		return nil, nil, datamodel.ErrIteratorOverread{}
	}
	w.done = true

	haveIdx, mval := unionMember(w.val)
	if haveIdx < 0 {
		return nil, nil, fmt.Errorf("bindnode: union %s has no member", w.val.Type())
	}
	mtyp := w.members[haveIdx]

	node := &_node{
		schemaType: mtyp,
		val:        mval,
	}
	key = basicnode.NewString(mtyp.Name())
	return key, node, nil
}

func (w *_unionIterator) Done() bool {
	return w.done
}
