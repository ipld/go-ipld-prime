package bindnode

import (
	"fmt"
	"reflect"
	"strings"

	ipld "github.com/ipld/go-ipld-prime"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/schema"
)

func reprNode(node ipld.Node) ipld.Node {
	if node, ok := node.(schema.TypedNode); ok {
		return node.Representation()
	}
	// ipld.Absent and ipld.Null are not typed.
	// TODO: is this a problem? surely a typed struct's fields are always
	// typed, even when absent or null.
	return node
}

func reprStrategy(typ schema.Type) interface{} {
	switch typ := typ.(type) {
	case *schema.TypeStruct:
		return typ.RepresentationStrategy()
	case *schema.TypeUnion:
		return typ.RepresentationStrategy()
	}
	return nil
}

type _prototypeRepr _prototype

func (w *_prototypeRepr) NewBuilder() ipld.NodeBuilder {
	return &_builderRepr{_assemblerRepr{
		schemaType: w.schemaType,
		val:        reflect.New(w.goType).Elem(),
	}}
}

type _nodeRepr _node

func (w *_nodeRepr) Kind() ipld.Kind {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Stringjoin:
		return ipld.Kind_String
	case schema.StructRepresentation_Map:
		return ipld.Kind_Map
	case schema.StructRepresentation_Tuple:
		return ipld.Kind_List
	case schema.UnionRepresentation_Keyed:
		return ipld.Kind_Map
	case schema.UnionRepresentation_Kinded:
		haveIdx := int(w.val.FieldByName("Index").Int())
		mtyp := w.schemaType.(*schema.TypeUnion).Members()[haveIdx]
		return mtyp.TypeKind().ActsLike()
	case schema.UnionRepresentation_Stringprefix:
		return ipld.Kind_String
	case nil:
		return (*_node)(w).Kind()
	default:
		panic(fmt.Sprintf("TODO Kind: %T", stg))
	}
}

func outboundMappedKey(stg schema.StructRepresentation_Map, key string) string {
	// TODO: why doesn't stg just allow us to "get" by the key string?
	field := schema.SpawnStructField(key, "", false, false)
	mappedKey := stg.GetFieldKey(field)
	return mappedKey
}

func inboundMappedKey(typ *schema.TypeStruct, stg schema.StructRepresentation_Map, key string) string {
	// TODO: can't do a "reverse" lookup... needs better API probably.
	fields := typ.Fields()
	for _, field := range fields {
		mappedKey := stg.GetFieldKey(field)
		if key == mappedKey {
			// println(key, "rev-mapped to", field.Name())
			return field.Name()
		}
	}
	// println(key, "had no mapping")
	return key // fallback to the same key
}

func outboundMappedType(stg schema.UnionRepresentation_Keyed, key string) string {
	// TODO: why doesn't stg just allow us to "get" by the key string?
	typ := schema.SpawnBool(schema.TypeName(key))
	mappedKey := stg.GetDiscriminant(typ)
	return mappedKey
}

func inboundMappedType(typ *schema.TypeUnion, stg schema.UnionRepresentation_Keyed, key string) string {
	// TODO: can't do a "reverse" lookup... needs better API probably.
	for _, member := range typ.Members() {
		mappedKey := stg.GetDiscriminant(member)
		if key == mappedKey {
			// println(key, "rev-mapped to", field.Name())
			return member.Name().String()
		}
	}
	// println(key, "had no mapping")
	return key // fallback to the same key
}

func (w *_nodeRepr) asKinded(stg schema.UnionRepresentation_Kinded, kind ipld.Kind) *_nodeRepr {
	name := stg.GetMember(kind)
	members := w.schemaType.(*schema.TypeUnion).Members()
	for _, member := range members {
		if member.Name() != name {
			continue
		}
		w2 := *w
		w2.val = w.val.FieldByName("Value").Elem()
		w2.schemaType = member
		return &w2
	}
	return nil
}

func (w *_nodeRepr) LookupByString(key string) (ipld.Node, error) {
	if stg, ok := reprStrategy(w.schemaType).(schema.UnionRepresentation_Kinded); ok {
		w = w.asKinded(stg, ipld.Kind_Map)
	}
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Map:
		revKey := inboundMappedKey(w.schemaType.(*schema.TypeStruct), stg, key)
		v, err := (*_node)(w).LookupByString(revKey)
		if err != nil {
			return nil, err
		}
		return reprNode(v), nil
	case schema.UnionRepresentation_Keyed:
		revKey := inboundMappedType(w.schemaType.(*schema.TypeUnion), stg, key)
		v, err := (*_node)(w).LookupByString(revKey)
		if err != nil {
			return nil, err
		}
		return reprNode(v), nil
	case nil:
		v, err := (*_node)(w).LookupByString(key)
		if err != nil {
			return nil, err
		}
		return reprNode(v), nil
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_nodeRepr) LookupByNode(key ipld.Node) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{
		TypeName:   w.schemaType.Name().String(),
		MethodName: "LookupByNode", AppropriateKind: ipld.KindSet_JustList, ActualKind: ipld.Kind_Map,
	}
}

func (w *_nodeRepr) LookupByIndex(idx int64) (ipld.Node, error) {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Tuple:
		fields := w.schemaType.(*schema.TypeStruct).Fields()
		field := fields[idx]
		v, err := (*_node)(w).LookupByString(field.Name())
		if err != nil {
			return nil, err
		}
		return reprNode(v), nil
	case nil:
		v, err := (*_node)(w).LookupByIndex(idx)
		if err != nil {
			return nil, err
		}
		return reprNode(v), nil
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_nodeRepr) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{
		TypeName:   w.schemaType.Name().String(),
		MethodName: "LookupBySegment",
		// TODO
	}
}

func (w *_nodeRepr) MapIterator() ipld.MapIterator {
	if stg, ok := reprStrategy(w.schemaType).(schema.UnionRepresentation_Kinded); ok {
		w = w.asKinded(stg, ipld.Kind_Map)
	}
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Map:
		itr := (*_node)(w).MapIterator().(*_structIterator)
		itr.reprEnd = int(w.lengthMinusTrailingAbsents())
		return (*_structIteratorRepr)(itr)
	case schema.UnionRepresentation_Keyed:
		itr := (*_node)(w).MapIterator().(*_unionIterator)
		return (*_unionIteratorRepr)(itr)
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_nodeRepr) ListIterator() ipld.ListIterator {
	return nil
}

func (w *_nodeRepr) lengthMinusTrailingAbsents() int64 {
	fields := w.schemaType.(*schema.TypeStruct).Fields()
	for i := len(fields) - 1; i >= 0; i-- {
		field := fields[i]
		if !field.IsOptional() || !w.val.Field(i).IsNil() {
			return int64(i + 1)
		}
	}
	return 0
}

func (w *_nodeRepr) Length() int64 {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Stringjoin:
		return -1
	case schema.StructRepresentation_Map:
		return w.lengthMinusTrailingAbsents()
	case schema.StructRepresentation_Tuple:
		return w.lengthMinusTrailingAbsents()
	case schema.UnionRepresentation_Keyed:
		return (*_node)(w).Length()
	case schema.UnionRepresentation_Kinded:
		w = w.asKinded(stg, w.Kind())
		// continues below
	case nil:
		// continues below
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
	return (*_node)(w).Length()
}

func (w *_nodeRepr) IsAbsent() bool {
	return false
}

func (w *_nodeRepr) IsNull() bool {
	return false
}

func (w *_nodeRepr) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{
		TypeName:   w.schemaType.Name().String(),
		MethodName: "AsBool", AppropriateKind: ipld.KindSet_JustBool, ActualKind: ipld.Kind_Map,
	}
}

func (w *_nodeRepr) AsInt() (int64, error) {
	return 0, ipld.ErrWrongKind{
		TypeName: w.schemaType.Name().String(), MethodName: "AsInt",
		AppropriateKind: ipld.KindSet_JustInt, ActualKind: ipld.Kind_Map,
	}
}

func (w *_nodeRepr) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{
		TypeName: w.schemaType.Name().String(), MethodName: "AsFloat",
		AppropriateKind: ipld.KindSet_JustFloat, ActualKind: ipld.Kind_Map,
	}
}

func (w *_nodeRepr) AsString() (string, error) {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Stringjoin:
		var b strings.Builder
		itr := (*_node)(w).MapIterator()
		for !itr.Done() {
			_, v, err := itr.Next()
			if err != nil {
				return "", err
			}
			s, err := reprNode(v).AsString()
			if err != nil {
				return "", err
			}
			if b.Len() > 0 {
				b.WriteString(stg.GetDelim())
			}
			b.WriteString(s)
		}
		return b.String(), nil
	case schema.UnionRepresentation_Stringprefix:
		haveIdx := int(w.val.FieldByName("Index").Int())
		mtyp := w.schemaType.(*schema.TypeUnion).Members()[haveIdx]

		w2 := *w
		w2.val = w.val.FieldByName("Value").Elem()
		w2.schemaType = mtyp
		s, err := w2.AsString()
		if err != nil {
			return "", err
		}

		name := stg.GetDiscriminant(mtyp)
		return name + stg.GetDelim() + s, nil
	case schema.UnionRepresentation_Kinded:
		w = w.asKinded(stg, ipld.Kind_String)
		// continues below
	case nil:
		// continues below
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
	return (*_node)(w).AsString()
}

func (w *_nodeRepr) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{
		TypeName: w.schemaType.Name().String(), MethodName: "AsBytes",
		AppropriateKind: ipld.KindSet_JustBytes, ActualKind: ipld.Kind_Map,
	}
}

func (w *_nodeRepr) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{
		TypeName: w.schemaType.Name().String(), MethodName: "AsLink",
		AppropriateKind: ipld.KindSet_JustLink, ActualKind: ipld.Kind_Map,
	}
}

func (w *_nodeRepr) Prototype() ipld.NodePrototype {
	panic("TODO: Prototype")
}

type _builderRepr struct {
	_assemblerRepr
}

// TODO: returning a repr node here is probably good, but there's a gotcha: one
// can go from a typed node to a repr node via the Representation method, but
// not the other way. That's probably why codegen returns a typed node here.
// The solution might be to add a way to go from the repr node to its parent
// typed node.

func (w *_builderRepr) Build() ipld.Node {
	// TODO: see the notes above.
	// return &_nodeRepr{schemaType: w.schemaType, val: w.val}
	return &_node{schemaType: w.schemaType, val: w.val}
}

func (w *_builderRepr) Reset() {
	panic("TODO: Reset")
}

type _assemblerRepr struct {
	schemaType schema.Type
	val        reflect.Value // non-pointer
	finish     func() error

	nullable bool
}

func (w *_assemblerRepr) asKinded(stg schema.UnionRepresentation_Kinded, kind ipld.Kind) *_assemblerRepr {
	name := stg.GetMember(kind)
	members := w.schemaType.(*schema.TypeUnion).Members()
	for idx, member := range members {
		if member.Name() != name {
			continue
		}
		w2 := *w
		goType := inferGoType(member) // TODO: do this upfront
		w2.val = reflect.New(goType).Elem()
		w2.schemaType = member

		// Layer a new finish func on top, to set Index/Value.
		w2.finish = func() error {
			if w.finish != nil {
				if err := w.finish(); err != nil {
					return err
				}
			}
			w.val.FieldByName("Index").SetInt(int64(idx))
			w.val.FieldByName("Value").Set(w2.val)
			return nil
		}
		return &w2
	}
	return nil
}

func (w *_assemblerRepr) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	if stg, ok := reprStrategy(w.schemaType).(schema.UnionRepresentation_Kinded); ok {
		w = w.asKinded(stg, ipld.Kind_Map)
	}
	asm, err := (*_assembler)(w).BeginMap(sizeHint)
	if err != nil {
		return nil, err
	}
	switch asm := asm.(type) {
	case *_structAssembler:
		return (*_structAssemblerRepr)(asm), nil
	case *_mapAssembler:
		return (*_mapAssemblerRepr)(asm), nil
	case *_unionAssembler:
		return (*_unionAssemblerRepr)(asm), nil
	default:
		panic(fmt.Sprintf("%T", asm))
	}
}

func (w *_assemblerRepr) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Tuple:
		asm, err := (*_assembler)(w).BeginMap(sizeHint)
		if err != nil {
			return nil, err
		}
		return (*_listStructAssemblerRepr)(asm.(*_structAssembler)), nil
	case nil:
		asm, err := (*_assembler)(w).BeginList(sizeHint)
		if err != nil {
			return nil, err
		}
		return (*_listAssemblerRepr)(asm.(*_listAssembler)), nil
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_assemblerRepr) AssignNull() error {
	switch stg := reprStrategy(w.schemaType).(type) {
	case nil:
		return (*_assembler)(w).AssignNull()
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_assemblerRepr) AssignBool(bool) error {
	panic("TODO: AssignBool")
}

func (w *_assemblerRepr) AssignInt(i int64) error {
	panic("TODO")
}

func (w *_assemblerRepr) AssignFloat(float64) error {
	panic("TODO: AssignFloat")
}

func (w *_assemblerRepr) AssignString(s string) error {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Stringjoin:
		fields := w.schemaType.(*schema.TypeStruct).Fields()
		parts := strings.Split(s, stg.GetDelim())
		if len(parts) != len(fields) {
			panic("TODO: len mismatch")
		}
		mapAsm, err := (*_assembler)(w).BeginMap(-1)
		if err != nil {
			return err
		}
		for i, field := range fields {
			entryAsm, err := mapAsm.AssembleEntry(field.Name())
			if err != nil {
				return err
			}
			entryAsm = entryAsm.(TypedAssembler).Representation()
			if err := entryAsm.AssignString(parts[i]); err != nil {
				return err
			}
		}
		return mapAsm.Finish()
	case schema.UnionRepresentation_Kinded:
		name := stg.GetMember(ipld.Kind_String)
		members := w.schemaType.(*schema.TypeUnion).Members()
		for idx, member := range members {
			if member.Name() != name {
				continue
			}
			w.val.FieldByName("Index").SetInt(int64(idx))
			w.val.FieldByName("Value").Set(reflect.ValueOf(s))
			return nil
		}
		panic("TODO: GetMember result is missing?")
	case schema.UnionRepresentation_Stringprefix:
		parts := strings.SplitN(s, stg.GetDelim(), 2)
		if len(parts) != 2 {
			panic("TODO: bad format")
		}
		name, value := parts[0], parts[1]
		members := w.schemaType.(*schema.TypeUnion).Members()
		for idx, member := range members {
			if stg.GetDiscriminant(member) != name {
				continue
			}

			w2 := *w
			goType := inferGoType(member) // TODO: do this upfront
			w2.val = reflect.New(goType).Elem()
			w2.schemaType = member
			w2.finish = func() error {
				if w.finish != nil {
					if err := w.finish(); err != nil {
						return err
					}
				}
				w.val.FieldByName("Index").SetInt(int64(idx))
				w.val.FieldByName("Value").Set(w2.val)
				return nil
			}

			return w2.AssignString(value)
		}
		panic("TODO: GetMember result is missing?")
	case nil:
		return (*_assembler)(w).AssignString(s)
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_assemblerRepr) AssignBytes(p []byte) error {
	panic("TODO")
}

func (w *_assemblerRepr) AssignLink(link ipld.Link) error {
	panic("TODO")
}

func (w *_assemblerRepr) AssignNode(node ipld.Node) error {
	panic("TODO")
}

func (w *_assemblerRepr) Prototype() ipld.NodePrototype {
	panic("TODO: Assembler.Prototype")
}

type _structAssemblerRepr _structAssembler

func (w *_structAssemblerRepr) AssembleKey() ipld.NodeAssembler {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Map:
		return (*_structAssembler)(w).AssembleKey()
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_structAssemblerRepr) AssembleValue() ipld.NodeAssembler {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Map:
		key := w.curKey.val.String()
		revKey := inboundMappedKey(w.schemaType, stg, key)
		w.curKey.val.SetString(revKey)

		valAsm := (*_structAssembler)(w).AssembleValue()
		valAsm = valAsm.(TypedAssembler).Representation()
		return valAsm
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_structAssemblerRepr) AssembleEntry(k string) (ipld.NodeAssembler, error) {
	if err := w.AssembleKey().AssignString(k); err != nil {
		return nil, err
	}
	am := w.AssembleValue()
	return am, nil
}

func (w *_structAssemblerRepr) Finish() error {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Map:
		err := (*_structAssembler)(w).Finish()
		if err, ok := err.(ipld.ErrMissingRequiredField); ok {
			for i, name := range err.Missing {
				serial := outboundMappedKey(stg, name)
				if serial != name {
					err.Missing[i] += fmt.Sprintf(" (serial:%q)", serial)
				}
			}
		}
		return err
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_structAssemblerRepr) KeyPrototype() ipld.NodePrototype {
	panic("TODO")
}

func (w *_structAssemblerRepr) ValuePrototype(k string) ipld.NodePrototype {
	panic("TODO: struct ValuePrototype")
}

type _mapAssemblerRepr _mapAssembler

func (w *_mapAssemblerRepr) AssembleKey() ipld.NodeAssembler {
	switch stg := reprStrategy(w.schemaType).(type) {
	case nil:
		asm := (*_mapAssembler)(w).AssembleKey()
		return (*_assemblerRepr)(asm.(*_assembler))
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_mapAssemblerRepr) AssembleValue() ipld.NodeAssembler {
	switch stg := reprStrategy(w.schemaType).(type) {
	case nil:
		asm := (*_mapAssembler)(w).AssembleValue()
		return (*_assemblerRepr)(asm.(*_assembler))
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_mapAssemblerRepr) AssembleEntry(k string) (ipld.NodeAssembler, error) {
	if err := w.AssembleKey().AssignString(k); err != nil {
		return nil, err
	}
	am := w.AssembleValue()
	return am, nil
}

func (w *_mapAssemblerRepr) Finish() error {
	switch stg := reprStrategy(w.schemaType).(type) {
	case nil:
		return (*_mapAssembler)(w).Finish()
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_mapAssemblerRepr) KeyPrototype() ipld.NodePrototype {
	panic("TODO")
}

func (w *_mapAssemblerRepr) ValuePrototype(k string) ipld.NodePrototype {
	panic("TODO: struct ValuePrototype")
}

type _listStructAssemblerRepr _structAssembler

func (w *_listStructAssemblerRepr) AssembleValue() ipld.NodeAssembler {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Tuple:
		fields := w.schemaType.Fields()
		field := fields[w.nextIndex]
		w.nextIndex++

		entryAsm, err := (*_structAssembler)(w).AssembleEntry(field.Name())
		if err != nil {
			panic(err) // TODO: probably return an assembler that always errors?
		}
		entryAsm = entryAsm.(TypedAssembler).Representation()
		return entryAsm
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_listStructAssemblerRepr) Finish() error {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Tuple:
		return (*_structAssembler)(w).Finish()
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_listStructAssemblerRepr) ValuePrototype(idx int64) ipld.NodePrototype {
	panic("TODO: list ValuePrototype")
}

// Note that lists do not have any representation strategy right now.
type _listAssemblerRepr _listAssembler

func (w *_listAssemblerRepr) AssembleValue() ipld.NodeAssembler {
	asm := (*_listAssembler)(w).AssembleValue()
	return (*_assemblerRepr)(asm.(*_assembler))
}

func (w *_listAssemblerRepr) Finish() error {
	return (*_listAssembler)(w).Finish()
}

func (w *_listAssemblerRepr) ValuePrototype(idx int64) ipld.NodePrototype {
	panic("TODO: list ValuePrototype")
}

type _unionAssemblerRepr _unionAssembler

func (w *_unionAssemblerRepr) AssembleKey() ipld.NodeAssembler {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.UnionRepresentation_Keyed:
		return (*_unionAssembler)(w).AssembleKey()
	case nil:
		asm := (*_unionAssembler)(w).AssembleKey()
		return (*_assemblerRepr)(asm.(*_assembler))
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_unionAssemblerRepr) AssembleValue() ipld.NodeAssembler {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.UnionRepresentation_Keyed:
		key := w.curKey.val.String()
		revKey := inboundMappedType(w.schemaType, stg, key)
		w.curKey.val.SetString(revKey)

		valAsm := (*_unionAssembler)(w).AssembleValue()
		valAsm = valAsm.(TypedAssembler).Representation()
		return valAsm
	case nil:
		asm := (*_unionAssembler)(w).AssembleValue()
		return (*_assemblerRepr)(asm.(*_assembler))
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_unionAssemblerRepr) AssembleEntry(k string) (ipld.NodeAssembler, error) {
	if err := w.AssembleKey().AssignString(k); err != nil {
		return nil, err
	}
	am := w.AssembleValue()
	return am, nil
}

func (w *_unionAssemblerRepr) Finish() error {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.UnionRepresentation_Keyed:
		return (*_unionAssembler)(w).Finish()
	case nil:
		return (*_unionAssembler)(w).Finish()
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_unionAssemblerRepr) KeyPrototype() ipld.NodePrototype {
	panic("TODO")
}

func (w *_unionAssemblerRepr) ValuePrototype(k string) ipld.NodePrototype {
	panic("TODO: struct ValuePrototype")
}

type _structIteratorRepr _structIterator

func (w *_structIteratorRepr) Next() (key, value ipld.Node, _ error) {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Map:
	_skipAbsent:
		key, value, err := (*_structIterator)(w).Next()
		if err != nil {
			return nil, nil, err
		}
		if value.IsAbsent() {
			goto _skipAbsent
		}
		keyStr, _ := key.AsString()
		mappedKey := outboundMappedKey(stg, keyStr)
		if mappedKey != keyStr {
			key = basicnode.NewString(mappedKey)
		}
		return key, reprNode(value), nil
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_structIteratorRepr) Done() bool {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.StructRepresentation_Map:
		// TODO: the fact that repr map iterators skip absents should be
		// documented somewhere
		return w.nextIndex >= w.reprEnd
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

type _unionIteratorRepr _unionIterator

func (w *_unionIteratorRepr) Next() (key, value ipld.Node, _ error) {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.UnionRepresentation_Keyed:
		key, value, err := (*_unionIterator)(w).Next()
		if err != nil {
			return nil, nil, err
		}
		keyStr, _ := key.AsString()
		mappedKey := outboundMappedType(stg, keyStr)
		if mappedKey != keyStr {
			key = basicnode.NewString(mappedKey)
		}
		return key, reprNode(value), nil
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}

func (w *_unionIteratorRepr) Done() bool {
	switch stg := reprStrategy(w.schemaType).(type) {
	case schema.UnionRepresentation_Keyed:
		return (*_unionIterator)(w).Done()
	default:
		panic(fmt.Sprintf("TODO: %T", stg))
	}
}
