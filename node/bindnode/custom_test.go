package bindnode_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/multiformats/go-multihash"

	qt "github.com/frankban/quicktest"
)

type BoolSubst int

var errorDefault = errors.New("something went wrong")

const (
	BoolSubst_Yes = 100
	BoolSubst_No  = -100
)

func BoolSubstFromBool(b bool) (interface{}, error) {
	if b {
		return BoolSubst_Yes, nil
	}
	return BoolSubst_No, nil
}

func BoolSubstFromBoolError(b bool) (interface{}, error) {
	return BoolSubst_No, errorDefault
}

func BoolToBoolSubst(b interface{}) (bool, error) {
	bp, ok := b.(*BoolSubst)
	if !ok {
		return true, fmt.Errorf("expected *BoolSubst value")
	}
	switch *bp {
	case BoolSubst_Yes:
		return true, nil
	case BoolSubst_No:
		return false, nil
	default:
		return true, fmt.Errorf("bad BoolSubst")
	}
}
func BoolToBoolSubstError(b interface{}) (bool, error) {
	return false, errorDefault
}

type IntSubst string

func IntSubstFromInt(i int64) (interface{}, error) {
	if i == 1000 {
		return "one thousand", nil
	} else if i == 2000 {
		return "two thousand", nil
	}
	return nil, fmt.Errorf("unexpected value of IntSubst")
}

func IntSubstFromIntError(i int64) (interface{}, error) {
	return nil, errorDefault
}

func IntToIntSubst(i interface{}) (int64, error) {
	ip, ok := i.(*IntSubst)
	if !ok {
		return 0, fmt.Errorf("expected *IntSubst value")
	}
	switch *ip {
	case "one thousand":
		return 1000, nil
	case "two thousand":
		return 2000, nil
	default:
		return 0, fmt.Errorf("bad IntSubst")
	}
}

func IntToIntSubstError(i interface{}) (int64, error) {
	return 0, errorDefault
}

type BigFloat struct{ *big.Float }

func BigFloatFromFloat(f float64) (interface{}, error) {
	bf := big.NewFloat(f)
	return &BigFloat{bf}, nil
}

func BigFloatFromFloatError(f float64) (interface{}, error) {
	return nil, errorDefault
}

func FloatFromBigFloat(f interface{}) (float64, error) {
	fp, ok := f.(*BigFloat)
	if !ok {
		return 0, fmt.Errorf("expected *BigFloat value")
	}
	f64, _ := fp.Float64()
	return f64, nil
}

func FloatFromBigFloatError(f interface{}) (float64, error) {
	return 0, errorDefault
}

type ByteArray [][]byte

func ByteArrayFromString(s string) (interface{}, error) {
	sa := strings.Split(s, "|")
	ba := make([][]byte, 0)
	for _, a := range sa {
		ba = append(ba, []byte(a))
	}
	return ba, nil
}

func ByteArrayFromStringError(s string) (interface{}, error) {
	return nil, errorDefault
}

func StringFromByteArray(b interface{}) (string, error) {
	bap, ok := b.(*ByteArray)
	if !ok {
		return "", fmt.Errorf("expected *ByteArray value")
	}
	sb := strings.Builder{}
	for i, b := range *bap {
		sb.WriteString(string(b))
		if i != len(*bap)-1 {
			sb.WriteString("|")
		}
	}
	return sb.String(), nil
}

func StringFromByteArrayError(b interface{}) (string, error) {
	return "", errorDefault
}

// similar to cid/Cid, go-address/Address, go-graphsync/RequestID
type Boop struct{ str string }

func NewBoop(b []byte) *Boop {
	return &Boop{string(b)}
}

func (b Boop) Bytes() []byte {
	return []byte(b.str)
}

func (b Boop) String() string {
	return b.str
}

// similar to go-state-types/big/Int
type Frop struct{ *big.Int }

func NewFropFromString(str string) Frop {
	v, _ := big.NewInt(0).SetString(str, 10)
	return Frop{v}
}

func NewFropFromBytes(buf []byte) *Frop {
	var negative bool
	switch buf[0] {
	case 0:
		negative = false
	case 1:
		negative = true
	default:
		panic("can't handle this")
	}

	i := big.NewInt(0).SetBytes(buf[1:])
	if negative {
		i.Neg(i)
	}

	return &Frop{i}
}

func (b *Frop) Bytes() []byte {
	switch {
	case b.Sign() > 0:
		return append([]byte{0}, b.Int.Bytes()...)
	case b.Sign() < 0:
		return append([]byte{1}, b.Int.Bytes()...)
	default:
		return []byte{}
	}
}

func BoopFromBytes(b []byte) (interface{}, error) {
	return NewBoop(b), nil
}

func BoopFromBytesError(b []byte) (interface{}, error) {
	return nil, errorDefault
}

func BoopToBytes(iface interface{}) ([]byte, error) {
	if boop, ok := iface.(*Boop); ok {
		return boop.Bytes(), nil
	}
	return nil, fmt.Errorf("did not get expected type")
}

func BoopToBytesError(iface interface{}) ([]byte, error) {
	return nil, errorDefault
}

func FropFromBytes(b []byte) (interface{}, error) {
	return NewFropFromBytes(b), nil
}

func FropFromBytesError(b []byte) (interface{}, error) {
	return nil, errorDefault
}

func FropToBytes(iface interface{}) ([]byte, error) {
	if frop, ok := iface.(*Frop); ok {
		return frop.Bytes(), nil
	}
	return nil, fmt.Errorf("did not get expected type")
}

func FropToBytesError(iface interface{}) ([]byte, error) {
	return nil, errorDefault
}

// Bitcoin's version of "links" is a hex form of the dbl-sha2-256 digest reversed
type BtcId string

func FromCidToBtcId(c cid.Cid) (interface{}, error) {
	if c.Prefix().Codec != cid.BitcoinBlock { // should be able to do BitcoinTx too .. but ..
		return nil, fmt.Errorf("can only convert IDs for BitcoinBlock codecs")
	}
	// and multihash must be dbl-sha2-256
	dig, err := multihash.Decode(c.Hash())
	if err != nil {
		return nil, err
	}
	hid := make([]byte, 0)
	for i := len(dig.Digest) - 1; i >= 0; i-- {
		hid = append(hid, dig.Digest[i])
	}
	return BtcId(hex.EncodeToString(hid)), nil
}

func FromCidToBtcIdError(c cid.Cid) (interface{}, error) {
	return BtcId(""), errorDefault
}

func FromBtcIdToCid(iface interface{}) (cid.Cid, error) {
	bid, ok := iface.(*BtcId)
	if !ok {
		return cid.Undef, fmt.Errorf("expected *BtcId value")
	}
	dig := make([]byte, 0)
	hid, err := hex.DecodeString(string(*bid))
	if err != nil {
		return cid.Undef, err
	}
	for i := len(hid) - 1; i >= 0; i-- {
		dig = append(dig, hid[i])
	}
	mh, err := multihash.Encode(dig, multihash.DBL_SHA2_256)
	if err != nil {
		return cid.Undef, err
	}
	return cid.NewCidV1(cid.BitcoinBlock, mh), nil
}

func FromBtcIdToCidError(iface interface{}) (cid.Cid, error) {
	return cid.Undef, errorDefault
}

type Boom struct {
	S    string
	St   ByteArray
	B    Boop
	Bo   BoolSubst
	Bptr *Boop
	F    Frop
	Fl   BigFloat
	I    int
	In   IntSubst
	L    BtcId
}

const boomSchema = `
type ByteArray string
type Boop bytes
type BoolSubst bool
type Frop bytes
type BigFloat float
type IntSubst int
type BtcId &Any

type Boom struct {
	S String
	St   ByteArray
	B    Boop
	Bo   BoolSubst
	Bptr nullable Boop
	F    Frop
	Fl   BigFloat
	I    Int
	In   IntSubst
	L    BtcId
} representation map
`

const boomFixtureDagJson = `{"B":{"/":{"bytes":"dGhlc2UgYXJlIGJ5dGVz"}},"Bo":false,"Bptr":{"/":{"bytes":"dGhlc2UgYXJlIHBvaW50ZXIgYnl0ZXM"}},"F":{"/":{"bytes":"AAH3fubjrGlwOMpClAkh/ro13L5Uls4/CtI"}},"Fl":1.12,"I":10101,"In":2000,"L":{"/":"bagyacvra2e6qt2fohajauxceox55t3gedsyqap2phmv7q2qaaaaaaaaaaaaa"},"S":"a string here","St":"a|byte|array"}`

var boomFixtureInstance = Boom{
	B:    *NewBoop([]byte("these are bytes")),
	Bo:   BoolSubst_No,
	Bptr: NewBoop([]byte("these are pointer bytes")),
	F:    NewFropFromString("12345678901234567891234567890123456789012345678901234567890"),
	Fl:   BigFloat{big.NewFloat(1.12)},
	I:    10101,
	In:   IntSubst("two thousand"),
	S:    "a string here",
	St:   ByteArray([][]byte{[]byte("a"), []byte("byte"), []byte("array")}),
	L:    BtcId("00000000000000006af82b3b4f3f00b11cc4ecd9fb75445c0a1238aee8093dd1"),
}

func TestCustom(t *testing.T) {
	opts := []bindnode.Option{
		bindnode.TypedBytesConverter(&Boop{}, BoopFromBytes, BoopToBytes),
		bindnode.TypedBytesConverter(&Frop{}, FropFromBytes, FropToBytes),
		bindnode.TypedBoolConverter(BoolSubst(0), BoolSubstFromBool, BoolToBoolSubst),
		bindnode.TypedIntConverter(IntSubst(""), IntSubstFromInt, IntToIntSubst),
		bindnode.TypedFloatConverter(&BigFloat{}, BigFloatFromFloat, FloatFromBigFloat),
		bindnode.TypedStringConverter(&ByteArray{}, ByteArrayFromString, StringFromByteArray),
		bindnode.TypedLinkConverter(BtcId(""), FromCidToBtcId, FromBtcIdToCid),
	}

	typeSystem, err := ipld.LoadSchemaBytes([]byte(boomSchema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("Boom")
	proto := bindnode.Prototype(&Boom{}, schemaType, opts...)

	builder := proto.Representation().NewBuilder()
	err = dagjson.Decode(builder, bytes.NewReader([]byte(boomFixtureDagJson)))
	qt.Assert(t, err, qt.IsNil)

	typ := bindnode.Unwrap(builder.Build())
	inst, ok := typ.(*Boom)
	qt.Assert(t, ok, qt.IsTrue)

	cmpr := qt.CmpEquals(
		cmp.Comparer(func(x, y Boop) bool { return x.String() == y.String() }),
		cmp.Comparer(func(x, y Frop) bool { return x.String() == y.String() }),
		cmp.Comparer(func(x, y BigFloat) bool { return x.String() == y.String() }),
	)
	qt.Assert(t, *inst, cmpr, boomFixtureInstance)

	tn := bindnode.Wrap(inst, schemaType, opts...)
	var buf bytes.Buffer
	err = dagjson.Encode(tn.Representation(), &buf)
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, buf.String(), qt.Equals, boomFixtureDagJson)
}

func TestCustomNamed(t *testing.T) {
	opts := []bindnode.Option{
		bindnode.NamedBytesConverter("Boop", BoopFromBytes, BoopToBytes),
		bindnode.NamedBytesConverter("Frop", FropFromBytes, FropToBytes),
		bindnode.NamedBoolConverter("BoolSubst", BoolSubstFromBool, BoolToBoolSubst),
		bindnode.NamedIntConverter("IntSubst", IntSubstFromInt, IntToIntSubst),
		bindnode.NamedFloatConverter("BigFloat", BigFloatFromFloat, FloatFromBigFloat),
		bindnode.NamedStringConverter("ByteArray", ByteArrayFromString, StringFromByteArray),
		bindnode.NamedLinkConverter("BtcId", FromCidToBtcId, FromBtcIdToCid),
		// these will error, but shouldn't get called cause the named converters take precedence
		bindnode.TypedBytesConverter(&Boop{}, BoopFromBytesError, BoopToBytesError),
		bindnode.TypedBytesConverter(&Frop{}, FropFromBytesError, FropToBytesError),
		bindnode.TypedBoolConverter(BoolSubst(0), BoolSubstFromBoolError, BoolToBoolSubstError),
		bindnode.TypedIntConverter(IntSubst(""), IntSubstFromIntError, IntToIntSubstError),
		bindnode.TypedFloatConverter(&BigFloat{}, BigFloatFromFloatError, FloatFromBigFloatError),
		bindnode.TypedStringConverter(&ByteArray{}, ByteArrayFromStringError, StringFromByteArrayError),
		bindnode.TypedLinkConverter(BtcId(""), FromCidToBtcIdError, FromBtcIdToCidError),
	}

	typeSystem, err := ipld.LoadSchemaBytes([]byte(boomSchema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("Boom")
	proto := bindnode.Prototype(&Boom{}, schemaType, opts...)

	builder := proto.Representation().NewBuilder()
	err = dagjson.Decode(builder, bytes.NewReader([]byte(boomFixtureDagJson)))
	qt.Assert(t, err, qt.IsNil)

	typ := bindnode.Unwrap(builder.Build())
	inst, ok := typ.(*Boom)
	qt.Assert(t, ok, qt.IsTrue)

	cmpr := qt.CmpEquals(
		cmp.Comparer(func(x, y Boop) bool { return x.String() == y.String() }),
		cmp.Comparer(func(x, y Frop) bool { return x.String() == y.String() }),
		cmp.Comparer(func(x, y BigFloat) bool { return x.String() == y.String() }),
	)
	qt.Assert(t, *inst, cmpr, boomFixtureInstance)

	tn := bindnode.Wrap(inst, schemaType, opts...)
	var buf bytes.Buffer
	err = dagjson.Encode(tn.Representation(), &buf)
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, buf.String(), qt.Equals, boomFixtureDagJson)
}

type AnyExtend struct {
	Name         string
	Blob         AnyExtendBlob
	Count        int
	Null         AnyCborEncoded
	NullPtr      *AnyCborEncoded
	NullableWith *AnyCborEncoded
	Bool         AnyCborEncoded
	Int          AnyCborEncoded
	Float        AnyCborEncoded
	String       AnyCborEncoded
	Bytes        AnyCborEncoded
	Link         AnyCborEncoded
	Map          AnyCborEncoded
	List         AnyCborEncoded
	BoolPtr      *BoolSubst // included to test that a null entry won't call a non-Any converter
	XListAny     []AnyCborEncoded
	XMapAny      anyMap
}

type anyMap struct {
	Keys   []string
	Values map[string]*AnyCborEncoded
}

const anyExtendSchema = `
type AnyExtend struct {
	Name String
	Blob Any
	Count Int
	Null nullable Any
	NullPtr nullable Any
	NullableWith nullable Any
	Bool Any
	Int Any
	Float Any
	String Any
	Bytes Any
	Link Any
	Map Any
	List Any
	BoolPtr nullable Bool
	XListAny [Any]
	XMapAny {String:Any}
}
`

type AnyExtendBlob struct {
	f string
	x int64
	y int64
	z int64
}

func AnyExtendBlobFromNode(node datamodel.Node) (interface{}, error) {
	foo, err := node.LookupByString("foo")
	if err != nil {
		return nil, err
	}
	fooStr, err := foo.AsString()
	if err != nil {
		return nil, err
	}
	baz, err := node.LookupByString("baz")
	if err != nil {
		return nil, err
	}
	x, err := baz.LookupByIndex(0)
	if err != nil {
		return nil, err
	}
	xi, err := x.AsInt()
	if err != nil {
		return nil, err
	}
	y, err := baz.LookupByIndex(1)
	if err != nil {
		return nil, err
	}
	yi, err := y.AsInt()
	if err != nil {
		return nil, err
	}
	z, err := baz.LookupByIndex(2)
	if err != nil {
		return nil, err
	}
	zi, err := z.AsInt()
	if err != nil {
		return nil, err
	}
	return &AnyExtendBlob{f: fooStr, x: xi, y: yi, z: zi}, nil
}

func (aeb AnyExtendBlob) ToNode() (datamodel.Node, error) {
	return qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "foo", qp.String(aeb.f))
		qp.MapEntry(ma, "baz", qp.List(-1, func(la datamodel.ListAssembler) {
			qp.ListEntry(la, qp.Int(aeb.x))
			qp.ListEntry(la, qp.Int(aeb.y))
			qp.ListEntry(la, qp.Int(aeb.z))
		}))
	})
}

func AnyExtendBlobToNode(ptr interface{}) (datamodel.Node, error) {
	aeb, ok := ptr.(*AnyExtendBlob)
	if !ok {
		return nil, fmt.Errorf("expected *AnyExtendBlob type")
	}
	return aeb.ToNode()
}

// take a datamodel.Node, dag-cbor encode it and store it here, do the reverse
// to get the datamodel.Node back
type AnyCborEncoded struct{ str []byte }

func AnyCborEncodedFromNode(node datamodel.Node) (interface{}, error) {
	if tn, ok := node.(schema.TypedNode); ok {
		node = tn.Representation()
	}
	var buf bytes.Buffer
	err := dagcbor.Encode(node, &buf)
	if err != nil {
		return nil, err
	}
	acb := AnyCborEncoded{str: buf.Bytes()}
	return &acb, nil
}

func AnyCborEncodedToNode(ptr interface{}) (datamodel.Node, error) {
	acb, ok := ptr.(*AnyCborEncoded)
	if !ok {
		return nil, fmt.Errorf("expected *AnyCborEncoded type")
	}
	na := basicnode.Prototype.Any.NewBuilder()
	err := dagcbor.Decode(na, bytes.NewReader(acb.str))
	if err != nil {
		return nil, err
	}
	return na.Build(), nil
}

const anyExtendDagJson = `{"Blob":{"baz":[2,3,4],"foo":"bar"},"Bool":false,"BoolPtr":null,"Bytes":{"/":{"bytes":"AgMEBQYHCA"}},"Count":101,"Float":2.34,"Int":123456789,"Link":{"/":"bagyacvra2e6qt2fohajauxceox55t3gedsyqap2phmv7q2qaaaaaaaaaaaaa"},"List":[null,"one","two","three",1,2,3,true],"Map":{"foo":"bar","one":1,"three":3,"two":2},"Name":"Any extend test","Null":null,"NullPtr":null,"NullableWith":123456789,"String":"this is a string","XListAny":[1,2,true,null,"bop"],"XMapAny":{"a":1,"b":2,"c":true,"d":null,"e":"bop"}}`

var anyExtendFixtureInstance = AnyExtend{
	Name:         "Any extend test",
	Count:        101,
	Blob:         AnyExtendBlob{f: "bar", x: 2, y: 3, z: 4},
	Null:         AnyCborEncoded{mustFromHex("f6")}, // normally these two fields would be `nil`, but we now get to decide whether it should be something concrete
	NullPtr:      &AnyCborEncoded{mustFromHex("f6")},
	NullableWith: &AnyCborEncoded{mustFromHex("1a075bcd15")},
	Bool:         AnyCborEncoded{mustFromHex("f4")},
	Int:          AnyCborEncoded{mustFromHex("1a075bcd15")},                                                                           // 123456789
	Float:        AnyCborEncoded{mustFromHex("fb4002b851eb851eb8")},                                                                   // 2.34
	String:       AnyCborEncoded{mustFromHex("7074686973206973206120737472696e67")},                                                   // "this is a string"
	Bytes:        AnyCborEncoded{mustFromHex("4702030405060708")},                                                                     // [2,3,4,5,6,7,8]
	Link:         AnyCborEncoded{mustFromHex("d82a58260001b0015620d13d09e8ae38120a5c4475fbd9ecc41cb1003f4f3b2bf86a0000000000000000")}, // bagyacvra2e6qt2fohajauxceox55t3gedsyqap2phmv7q2qaaaaaaaaaaaaa
	Map:          AnyCborEncoded{mustFromHex("a463666f6f63626172636f6e65016374776f0265746872656503")},                                 // {"one":1,"two":2,"three":3,"foo":"bar"}
	List:         AnyCborEncoded{mustFromHex("88f6636f6e656374776f657468726565010203f5")},                                             // [null,'one','two','three',1,2,3,true]
	BoolPtr:      nil,
	XListAny:     []AnyCborEncoded{{mustFromHex("01")}, {mustFromHex("02")}, {mustFromHex("f5")}, {mustFromHex("f6")}, {mustFromHex("63626f70")}}, // [1,2,true,null,"bop"]
	XMapAny: anyMap{
		Keys: []string{"a", "b", "c", "d", "e"},
		Values: map[string]*AnyCborEncoded{
			"a": {mustFromHex("01")},
			"b": {mustFromHex("02")},
			"c": {mustFromHex("f5")},
			"d": {mustFromHex("f6")},
			"e": {mustFromHex("63626f70")}}}, // {"a":1,"b":2,"c":true,"d":null,"e":"bop"}
}

func TestCustomAny(t *testing.T) {
	opts := []bindnode.Option{
		bindnode.TypedAnyConverter(&AnyExtendBlob{}, AnyExtendBlobFromNode, AnyExtendBlobToNode),
		bindnode.TypedAnyConverter(&AnyCborEncoded{}, AnyCborEncodedFromNode, AnyCborEncodedToNode),
		bindnode.TypedBoolConverter(BoolSubst(0), BoolSubstFromBool, BoolToBoolSubst),
	}

	typeSystem, err := ipld.LoadSchemaBytes([]byte(anyExtendSchema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("AnyExtend")
	proto := bindnode.Prototype(&AnyExtend{}, schemaType, opts...)

	builder := proto.Representation().NewBuilder()
	err = dagjson.Decode(builder, bytes.NewReader([]byte(anyExtendDagJson)))
	qt.Assert(t, err, qt.IsNil)

	typ := bindnode.Unwrap(builder.Build())
	inst, ok := typ.(*AnyExtend)
	qt.Assert(t, ok, qt.IsTrue)

	cmpr := qt.CmpEquals(
		cmp.Comparer(func(x, y AnyExtendBlob) bool {
			return x.f == y.f && x.x == y.x && x.y == y.y && x.z == y.z
		}),
		cmp.Comparer(func(x, y AnyCborEncoded) bool {
			return bytes.Equal(x.str, y.str)
		}),
	)
	qt.Assert(t, *inst, cmpr, anyExtendFixtureInstance)

	tn := bindnode.Wrap(inst, schemaType, opts...)
	var buf bytes.Buffer
	err = dagjson.Encode(tn.Representation(), &buf)
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, buf.String(), qt.Equals, anyExtendDagJson)
}

func mustFromHex(hexStr string) []byte {
	byt, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return byt
}

type ClosedUnion interface {
	isClosedUnion()
}

type intUnion uint64

func (intUnion) isClosedUnion() {}

type stringUnion string

func (stringUnion) isClosedUnion() {}

func ClosedUnionFromNode(node datamodel.Node) (interface{}, error) {
	asInt, err := node.AsInt()
	if err == nil {
		return intUnion(asInt), nil
	}
	asString, err := node.AsString()
	if err == nil {
		return stringUnion(asString), nil
	}
	return nil, errors.New("unrecognized type")
}

func ClosedUnionToNode(val interface{}) (datamodel.Node, error) {
	cu, ok := val.(*ClosedUnion)
	if !ok {
		return nil, errors.New("should be a ClosedUnion")
	}
	switch concrete := (*cu).(type) {
	case intUnion:
		return basicnode.NewInt(int64(concrete)), nil
	case stringUnion:
		return basicnode.NewString(string(concrete)), nil
	default:
		return nil, errors.New("unexpected union type")
	}
}

type StructWithUnion struct {
	Cu ClosedUnion
}

const closedUnionSchema = `
type ClosedUnion any

type StructWithUnion struct {
  cu ClosedUnion
}
`

const closedUnionFixtureIntDagJson = `{"cu":8}`
const closedUnionFixtureStringDagJson = `{"cu":"happy"}`

var closedUnionIntInst = StructWithUnion{
	Cu: intUnion(8),
}
var closedUnionStringInst = StructWithUnion{
	Cu: stringUnion("happy"),
}

func TestCustomAnyWithInterface(t *testing.T) {
	opts := []bindnode.Option{
		bindnode.NamedAnyConverter("ClosedUnion", ClosedUnionFromNode, ClosedUnionToNode),
	}

	typeSystem, err := ipld.LoadSchemaBytes([]byte(closedUnionSchema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("StructWithUnion")
	proto := bindnode.Prototype(&StructWithUnion{}, schemaType, opts...)

	// test one union variant
	builder := proto.Representation().NewBuilder()
	err = dagjson.Decode(builder, bytes.NewReader([]byte(closedUnionFixtureIntDagJson)))
	qt.Assert(t, err, qt.IsNil)

	typ := bindnode.Unwrap(builder.Build())
	inst, ok := typ.(*StructWithUnion)
	qt.Assert(t, ok, qt.IsTrue)

	cmpr := qt.CmpEquals()
	qt.Assert(t, *inst, cmpr, closedUnionIntInst)

	tn := bindnode.Wrap(inst, schemaType, opts...)
	var buf bytes.Buffer
	err = dagjson.Encode(tn.Representation(), &buf)
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, buf.String(), qt.Equals, closedUnionFixtureIntDagJson)

	// test other union variant
	builder = proto.Representation().NewBuilder()
	err = dagjson.Decode(builder, bytes.NewReader([]byte(closedUnionFixtureStringDagJson)))
	qt.Assert(t, err, qt.IsNil)

	typ = bindnode.Unwrap(builder.Build())
	inst, ok = typ.(*StructWithUnion)
	qt.Assert(t, ok, qt.IsTrue)

	cmpr = qt.CmpEquals()
	qt.Assert(t, *inst, cmpr, closedUnionStringInst)

	tn = bindnode.Wrap(inst, schemaType, opts...)
	buf = bytes.Buffer{}
	err = dagjson.Encode(tn.Representation(), &buf)
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, buf.String(), qt.Equals, closedUnionFixtureStringDagJson)
}
