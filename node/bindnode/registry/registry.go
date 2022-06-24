package registry

import (
	"fmt"
	"io"
	"reflect"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
)

type prototypeData struct {
	proto   schema.TypedPrototype
	options []bindnode.Option
}

// BindnodeRegistry holds TypedPrototype and bindnode options for Go types and
// will use that data for conversion operations.
type BindnodeRegistry map[reflect.Type]prototypeData

// NewRegistry creates a new BindnodeRegistry
func NewRegistry() BindnodeRegistry {
	return make(BindnodeRegistry)
}

func typeOf(ptrValue interface{}) reflect.Type {
	val := reflect.ValueOf(ptrValue).Type()
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

func (br BindnodeRegistry) prototypeDataFor(ptrType interface{}) prototypeData {
	typ := typeOf(ptrType)
	proto, ok := br[typ]
	if !ok {
		panic(fmt.Sprintf("bindnode utils: type has not been registered: %s", typ.Name()))
	}
	return proto
}

// RegisterType registers ptrType with schema such that it can be wrapped and
// unwrapped without needing the schema, Type, or TypedPrototype.
// Typically the typeName will match the Go type name, but it can be whatever
// is defined in the schema for the type being registered.
// Registering the same type twice on this registry will cause an error.
// This call may also error if the schema is invalid or the type doesn't match
// the schema. Additionally, panics from within bindnode's initial prototype
// checks will be captured and returned as errors from this function.
func (br BindnodeRegistry) RegisterType(ptrType interface{}, schema string, typeName string, options ...bindnode.Option) (err error) {
	typ := typeOf(ptrType)
	if _, ok := br[typ]; ok {
		return fmt.Errorf("bindnode utils: type already registered: %s", typ.Name())
	}
	typeSystem, err := ipld.LoadSchemaBytes([]byte(schema))
	if err != nil {
		return fmt.Errorf("bindnode utils: failed to load schema: %s", err.Error())
	}
	schemaType := typeSystem.TypeByName(typeName)
	if schemaType == nil {
		return fmt.Errorf("bindnode utils: schema for [%T] does not contain that named type [%s]", ptrType, typ.Name())
	}

	// focusing on bindnode setup panics
	defer func() {
		if rec := recover(); rec != nil {
			switch v := rec.(type) {
			case string:
				err = fmt.Errorf(v)
			case error:
				err = v
			default:
				panic(rec)
			}
		}
	}()

	proto := bindnode.Prototype(ptrType, schemaType, options...)
	br[typ] = prototypeData{
		proto,
		options,
	}

	return err
}

// IsRegistered can be used to determine if the type has already been registered
// within this registry.
// Using RegisterType on an already registered type will cause a panic, so where
// this may be the case, IsRegistered can be used to check.
func (br BindnodeRegistry) IsRegistered(ptrType interface{}) bool {
	_, ok := br[typeOf(ptrType)]
	return ok
}

// TypeFromReader deserializes bytes using the given codec from a Reader and
// instantiates the Go type that's provided as a pointer via the ptrValue
// argument.
func (br BindnodeRegistry) TypeFromReader(r io.Reader, ptrValue interface{}, decoder codec.Decoder) (interface{}, error) {
	protoData := br.prototypeDataFor(ptrValue)
	node, err := ipld.DecodeStreamingUsingPrototype(r, decoder, protoData.proto)
	if err != nil {
		return nil, err
	}
	typ := bindnode.Unwrap(node)
	return typ, nil
}

// TypeFromBytes deserializes bytes using the given codec from its bytes and
// instantiates the Go type that's provided as a pointer via the ptrValue
// argument.
func (br BindnodeRegistry) TypeFromBytes(byts []byte, ptrValue interface{}, decoder codec.Decoder) (interface{}, error) {
	protoData := br.prototypeDataFor(ptrValue)
	node, err := ipld.DecodeUsingPrototype(byts, decoder, protoData.proto)
	if err != nil {
		return nil, err
	}
	typ := bindnode.Unwrap(node)
	return typ, nil
}

// TypeFromNode converts an datamodel.Node into an appropriate Go type that's
// provided as a pointer via the ptrValue argument.
func (br BindnodeRegistry) TypeFromNode(node datamodel.Node, ptrValue interface{}) (interface{}, error) {
	protoData := br.prototypeDataFor(ptrValue)
	if tn, ok := node.(schema.TypedNode); ok {
		node = tn.Representation()
	}
	builder := protoData.proto.Representation().NewBuilder()
	err := builder.AssignNode(node)
	if err != nil {
		return nil, err
	}
	typ := bindnode.Unwrap(builder.Build())
	return typ, nil
}

// TypeToNode converts a Go type that's provided as a pointer via the ptrValue
// argument to an schema.TypedNode.
func (br BindnodeRegistry) TypeToNode(ptrValue interface{}) schema.TypedNode {
	protoData := br.prototypeDataFor(ptrValue)
	return bindnode.Wrap(ptrValue, protoData.proto.Type(), protoData.options...)
}

// TypeToWriter is a utility method that serializes a Go type that's provided as
// a pointer via the ptrValue argument through the given codec to a Writer.
func (br BindnodeRegistry) TypeToWriter(ptrValue interface{}, w io.Writer, encoder codec.Encoder) error {
	return ipld.EncodeStreaming(w, br.TypeToNode(ptrValue), encoder)
}

// TypeToBytes is a utility method that serializes a Go type that's provided as
// a pointer via the ptrValue argument through the given codec and returns the
// bytes.
func (br BindnodeRegistry) TypeToBytes(ptrValue interface{}, encoder codec.Encoder) ([]byte, error) {
	return ipld.Encode(br.TypeToNode(ptrValue), encoder)
}
