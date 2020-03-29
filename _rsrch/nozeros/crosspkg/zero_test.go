package crosspkg

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime/_rsrch/nozeros"
)

func TestStuff(t *testing.T) {
	_ = nozeros.ExportedAlias{} // undesirable

	// _ = nozeros.ExportedPtr{} // undesirable
	// NOPE!  (Success!)
	// ./zero.go:14:25: invalid pointer type nozeros.ExportedPtr for composite literal (use &nozeros.internal instead)

	v := nozeros.NewExportedPtr("woo")
	var typecheckMe interface{} = v
	v2, ok := typecheckMe.(nozeros.ExportedPtr)

	Wish(t, ok, ShouldEqual, true) // well, duh.  Mostly we wanted to know if even asking was allowed.
	Wish(t, v2, ShouldEqual, v)    // well, duh.  But sanity check, I guess.

	v.Pow() // check that your IDE can autocomplete this, too.  it should.

	// this is all the semantics we wanted; awesome.
	//
	// still unfortunate: Pow won't show up in docs, since it's on an unexported type.

	// exporting an alias of the non-pointer makes a hole in your guarantees, of course:
	var hmm nozeros.ExportedPtr
	hmm = &nozeros.ExportedAlias{}
	_ = hmm
}

func TestStranger(t *testing.T) {
	var foo nozeros.Foo
	// foo = &nozeros.FooData{} // undesirable
	// NOPE!  (Success!)

	foo = nozeros.Foo(nil) // possible, yes.  but fairly irrelevant.
	_ = foo

	v := nozeros.NewFoo("woo")

	Wish(t, v.Read(), ShouldEqual, "woo")

	v.Pow()

	Wish(t, v.Read(), ShouldEqual, "waht")

	v.Pow()

	t.Logf("%#v", v) // this will log the internal type name, not the exported alias.  Arguably not ideal.
}
