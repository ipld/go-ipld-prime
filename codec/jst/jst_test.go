package jst

import (
	"bytes"
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func TestSimple(t *testing.T) {
	fixture := Dedent(`
		[
		  {"path": "./foo",  "moduleName": "whiz.org/teamBar/foo", "status": "changed"},
		  {"path": "./baz",  "moduleName": "whiz.org/teamBar/baz", "status": "green"},
		  {"path": "./quxx", "moduleName": "example.net/quxx",     "status": "lit"}
		]
	`)
	nb := basicnode.Style.Any.NewBuilder()
	Require(t, dagjson.Decoder(nb, bytes.NewBufferString(fixture)), ShouldEqual, nil)
	n := nb.Build()

	st := state{}
	Wish(t, stride(&st, n), ShouldEqual, nil)
	Wish(t, st.tables, ShouldEqual, map[tableGroupID]*table{
		"path": &table{
			entryStyles: map[columnName]entryStyle{"path": entryStyle_column, "moduleName": entryStyle_column, "status": entryStyle_column},
			keySize:     map[columnName]int{}, // TODO not yet supported but this shouldn't be empty when it is
			cols:        []columnName{"path", "moduleName", "status"},
			colSize:     map[columnName]int{"path": 8, "moduleName": 22, "status": 9},
			ownLine:     nil,
		},
	})

	var buf bytes.Buffer
	Wish(t, Marshal(n, &buf), ShouldEqual, nil)
	Wish(t, buf.String(), ShouldEqual, fixture)
}

func TestTrailing(t *testing.T) {
	// FUTURE
}

func TestSubTables(t *testing.T) {
	fixture := Dedent(`
		[
		  {"path": "./foo",  "moduleName": "whiz.org/teamBar/foo", "status": "changed"},
		  {"path": "./baz",  "moduleName": "whiz.org/teamBar/baz", "status": "green",
		    "subtable": [
		      {"frob": "zozzle", "zim": "boink"},
		      {"frob": "narf",   "zim": "zamf"},
		    ]},
		  {"path": "./quxx", "moduleName": "example.net/quxx",     "status": "lit"}
		]
	`)
	nb := basicnode.Style.Any.NewBuilder()
	Require(t, dagjson.Decoder(nb, bytes.NewBufferString(fixture)), ShouldEqual, nil)
	n := nb.Build()

	var buf bytes.Buffer
	Wish(t, Marshal(n, &buf), ShouldEqual, nil)
	Wish(t, buf.String(), ShouldEqual, fixture)
}
