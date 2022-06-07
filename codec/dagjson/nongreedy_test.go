package dagjson

import (
	"bytes"
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestNonGreedy(t *testing.T) {
	buf := bytes.NewBufferString(`{"a": 1}{"b": 2}`)
	opts := DecodeOptions{
		ParseLinks:         false,
		ParseBytes:         false,
		DontParseBeyondEnd: true,
	}
	nb1 := basicnode.Prototype.Map.NewBuilder()
	err := opts.Decode(nb1, buf)
	if err != nil {
		t.Fatalf("first decode (%v)", err)
	}
	n1 := nb1.Build()
	if n1.Kind() != datamodel.Kind_Map {
		t.Errorf("expecting a map")
	}
	nb2 := basicnode.Prototype.Map.NewBuilder()
	err = opts.Decode(nb2, buf)
	if err != nil {
		t.Fatalf("second decode (%v)", err)
	}
	n2 := nb2.Build()
	if n2.Kind() != datamodel.Kind_Map {
		t.Errorf("expecting a map")
	}
}
