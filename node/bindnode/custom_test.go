package bindnode_test

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"

	qt "github.com/frankban/quicktest"
)

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
type Blop struct{ *big.Int }

func NewBlopFromString(str string) Blop {
	v, _ := big.NewInt(0).SetString(str, 10)
	return Blop{v}
}

func NewBlopFromBytes(buf []byte) Blop {
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

	return Blop{i}
}

func (b *Blop) Bytes() []byte {
	switch {
	case b.Sign() > 0:
		return append([]byte{0}, b.Int.Bytes()...)
	case b.Sign() < 0:
		return append([]byte{1}, b.Int.Bytes()...)
	default:
		return []byte{}
	}
}

type Boom struct {
	S    string
	B    Boop
	Bptr *Boop
	BI   Blop
	I    int
}

const boomSchema = `
type Boom struct {
	S String
	B Bytes
	Bptr nullable Bytes
	BI Bytes
	I Int
} representation map
`

const boomFixtureDagJson = `{"B":{"/":{"bytes":"dGhlc2UgYXJlIGJ5dGVz"}},"BI":{"/":{"bytes":"AAH3fubjrGlwOMpClAkh/ro13L5Uls4/CtI"}},"Bptr":{"/":{"bytes":"dGhlc2UgYXJlIHBvaW50ZXIgYnl0ZXM"}},"I":10101,"S":"a string here"}`

var boomFixtureInstance = Boom{
	S:    "a string here",
	B:    *NewBoop([]byte("these are bytes")),
	BI:   NewBlopFromString("12345678901234567891234567890123456789012345678901234567890"),
	Bptr: NewBoop([]byte("these are pointer bytes")),
	I:    10101,
}

type BoopConverter struct {
}

func (bc BoopConverter) FromBytes(b []byte) (interface{}, error) {
	return NewBoop(b), nil
}

func (bc BoopConverter) ToBytes(typ interface{}) ([]byte, error) {
	if boop, ok := typ.(*Boop); ok {
		return boop.Bytes(), nil
	}
	return nil, fmt.Errorf("did not get a Boop type")
}

type BlopConverter struct {
}

func (bc BlopConverter) FromBytes(b []byte) (interface{}, error) {
	return NewBlopFromBytes(b), nil
}

func (bc BlopConverter) ToBytes(typ interface{}) ([]byte, error) {
	if blop, ok := typ.(*Blop); ok {
		return blop.Bytes(), nil
	}
	return nil, fmt.Errorf("did not get a Blop type")
}

var (
	_ bindnode.CustomTypeBytesConverter = (*BoopConverter)(nil)
	_ bindnode.CustomTypeBytesConverter = (*BlopConverter)(nil)
)

func TestCustom(t *testing.T) {
	opts := []bindnode.Option{
		bindnode.AddCustomTypeBytesConverter(Boop{}, BoopConverter{}),
		bindnode.AddCustomTypeBytesConverter(Blop{}, BlopConverter{}),
	}

	nb := basicnode.Prototype.Any.NewBuilder()
	err := dagjson.Decode(nb, bytes.NewReader([]byte(boomFixtureDagJson)))
	qt.Assert(t, err, qt.IsNil)

	typeSystem, err := ipld.LoadSchemaBytes([]byte(boomSchema))
	qt.Assert(t, err, qt.IsNil)
	schemaType := typeSystem.TypeByName("Boom")
	proto := bindnode.Prototype(&Boom{}, schemaType, opts...)

	node := nb.Build()
	builder := proto.Representation().NewBuilder()
	err = builder.AssignNode(node)
	qt.Assert(t, err, qt.IsNil)

	typ := bindnode.Unwrap(builder.Build())
	inst, ok := typ.(*Boom)
	qt.Assert(t, ok, qt.IsTrue)

	cmpr := qt.CmpEquals(
		cmp.Comparer(func(x, y Boop) bool { return x.String() == y.String() }),
		cmp.Comparer(func(x, y Blop) bool { return x.String() == y.String() }),
	)
	qt.Assert(t, *inst, cmpr, boomFixtureInstance)

	tn := bindnode.Wrap(&boomFixtureInstance, schemaType, opts...)
	var buf bytes.Buffer
	err = dagjson.Encode(tn.Representation(), &buf)
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, buf.String(), qt.Equals, boomFixtureDagJson)
}
