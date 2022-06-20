package registry_test

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/node/bindnode/registry"

	qt "github.com/frankban/quicktest"
)

type HexString string
type Foo struct {
	Int  int
	Bool bool
}

func TestRegistry(t *testing.T) {
	reg := registry.NewRegistry()
	qt.Assert(t, reg.IsRegistered((*Foo)(nil)), qt.IsFalse)
	qt.Assert(t, reg.IsRegistered((*HexString)(nil)), qt.IsFalse)

	err := reg.RegisterType((*Foo)(nil),
		`type Foo struct {
			Int Int
			Bool Bool
		}`, "Foo")
	qt.Assert(t, err, qt.IsNil)

	err = reg.RegisterType((*HexString)(nil), "type HS bytes", "HS", bindnode.TypedBytesConverter(
		(*HexString)(nil),
		func(b []byte) (interface{}, error) {
			return HexString(hex.EncodeToString(b)), nil
		},
		func(i interface{}) ([]byte, error) {
			s, _ := i.(*HexString)
			return hex.DecodeString(string(*s))
		}))
	qt.Assert(t, err, qt.IsNil)

	qt.Assert(t, reg.IsRegistered((*Foo)(nil)), qt.IsTrue)
	qt.Assert(t, reg.IsRegistered((*HexString)(nil)), qt.IsTrue)

	hsi, err := reg.TypeFromNode(basicnode.NewBytes([]byte{0, 1, 2, 3, 4}), (*HexString)(nil))
	qt.Assert(t, err, qt.IsNil)
	hs, ok := hsi.(*HexString)
	qt.Assert(t, ok, qt.IsTrue)
	qt.Assert(t, string(*hs), qt.Equals, "0001020304")

	byts, _ := hex.DecodeString("a263496e74386364426f6f6cf4")
	fooi, err := reg.TypeFromBytes(byts, (*Foo)(nil), dagcbor.Decode)
	qt.Assert(t, err, qt.IsNil)
	foo, ok := fooi.(*Foo)
	qt.Assert(t, ok, qt.IsTrue)
	qt.Assert(t, *foo, qt.Equals, Foo{Int: -100, Bool: false})

	byts, err = reg.TypeToBytes(&Foo{Int: -100, Bool: false}, dagjson.Encode)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, string(byts), qt.Equals, `{"Bool":false,"Int":-100}`)

	byts, _ = hex.DecodeString("a263496e741a7fffffff64426f6f6cf5")
	fooi, err = reg.TypeFromReader(bytes.NewReader(byts), (*Foo)(nil), dagcbor.Decode)
	qt.Assert(t, err, qt.IsNil)
	foo, ok = fooi.(*Foo)
	qt.Assert(t, ok, qt.IsTrue)
	qt.Assert(t, *foo, qt.Equals, Foo{Int: math.MaxInt32, Bool: true})

	w := bytes.Buffer{}
	err = reg.TypeToWriter(&Foo{Int: math.MaxInt32, Bool: true}, &w, dagjson.Encode)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, w.String(), qt.Equals, `{"Bool":true,"Int":2147483647}`)
}

func TestRegistryErrors(t *testing.T) {
	reg := registry.NewRegistry()
	err := reg.RegisterType((*Foo)(nil), `type Nope nope {}`, "Foo")
	qt.Assert(t, err, qt.ErrorMatches, `.*unknown type keyword: "nope".*`)

	err = reg.RegisterType((*HexString)(nil), "type HS string", "HS")
	qt.Assert(t, err, qt.IsNil)

	err = reg.RegisterType((*HexString)(nil), "type HS2 string", "HS2")
	qt.Assert(t, err, qt.ErrorMatches, `.*type already registered: HexString`)

	err = reg.RegisterType((*Foo)(nil), "type NotFoo string", "Foo")
	qt.Assert(t, err, qt.ErrorMatches, `.*does not contain that named type.*`)

	err = reg.RegisterType((*Foo)(nil),
		`type Foo struct {
			NotInt String
			NotBool Float
		}`, "Foo")
	qt.Assert(t, err, qt.ErrorMatches, `.*kind mismatch.*`)
}
