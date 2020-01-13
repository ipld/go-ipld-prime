package impls

// justMap and this file is how ipldfree would do it: no concrete types, just interfaces.

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

type justMap struct {
	m map[string]ipld.Node // string key -- even if a runtime schema wrapper is using us for storage, we must have a comparable type here, and string is all we know.
	t []mapEntry
}

type mapEntry struct {
	k string    // address of this used when we return keys as nodes, such as in iterators.  Need in one place to amortize shifts to heap when ptr'ing for iface.
	v ipld.Node // ... actually, might not be needed.  is, for codegen'd, to slab alloc; but here?  already have em in the map values.
}

type justMapAssembler struct {
	w *justMap

	midappend bool // if true, next call must be 'AssembleValue'.
}
type justMapKeyAssembler struct{ justMapAssembler }
type justMapValueAssembler struct{ justMapAssembler }

// NOT yet an interface function... but wanting benchmarking on it.
// (how much it might show up for structs is another question.)
//  (there's... enough differences vs things with concrete type knowledge we might wanna do this all with one of those, too.)
func (ma *justMapAssembler) AssembleDirectly(k string) (ipld.NodeAssembler, error) {
	if ma.midappend == true {
		panic("misuse")
	}
	_, exists := ma.w.m[k]
	if exists {
		return nil, ipld.ErrRepeatedMapKey{String(k)}
	}
	//l := len(ma.w.t)
	ma.w.t = append(ma.w.t, mapEntry{k: k})
	// configure and return an anyAssembler, similar to below in prepareAssigner
	panic("todo")
}

func (ma *justMapAssembler) AssembleKey() ipld.NodeAssembler {
	return &justMapKeyAssembler{*ma} // todo: yeah, pretty sure need to embed these.  sigh.
}
func (ma *justMapAssembler) AssembleValue() ipld.NodeAssembler {
	return &justMapValueAssembler{*ma} // todo: yeah, pretty sure need to embed these.  sigh.
}

func (justMapKeyAssembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error)   { panic("no") }
func (justMapKeyAssembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) { panic("no") }
func (justMapKeyAssembler) AssignNull() error                                      { panic("no") }
func (justMapKeyAssembler) AssignBool(bool) error                                  { panic("no") }
func (justMapKeyAssembler) AssignInt(int) error                                    { panic("no") }
func (justMapKeyAssembler) AssignFloat(float64) error                              { panic("no") }
func (mka *justMapKeyAssembler) AssignString(v string) error {
	if mka.midappend == true {
		panic("misuse")
	}
	_, exists := mka.w.m[v]
	if exists {
		return ipld.ErrRepeatedMapKey{String(v)}
	}
	mka.w.t = append(mka.w.t, mapEntry{k: v})
	mka.midappend = true
	return nil
}
func (justMapKeyAssembler) AssignBytes([]byte) error { panic("no") }
func (mka *justMapKeyAssembler) Assign(v ipld.Node) error {
	vs, err := v.AsString()
	if err != nil {
		return fmt.Errorf("cannot assign non-string node into map key assembler") // FIXME:errors: this doesn't quite fit in ErrWrongKind cleanly; new error type?
	}
	return mka.AssignString(vs)
}
func (justMapKeyAssembler) Style() ipld.NodeStyle { panic("later") } // probably should give the style of justString, which could say "only stores string kind" (though we haven't made such a feature part of the interface yet).

func (mva *justMapValueAssembler) prepareAssigner() ipld.NodeAssembler {
	//l := len(mva.w.t) - 1
	//_ = anyAssembler{&mva.w.t[l].v} // is this even helpful?!  I don't think it is...
	panic("todo")
}
func (mva *justMapValueAssembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error) {
	return mva.prepareAssigner().BeginMap(sizeHint)
}
func (mva *justMapValueAssembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) {
	return mva.prepareAssigner().BeginList(sizeHint)
}
func (mva *justMapValueAssembler) AssignNull() error           { panic("yes") }
func (mva *justMapValueAssembler) AssignBool(bool) error       { panic("yes") }
func (mva *justMapValueAssembler) AssignInt(int) error         { panic("yes") }
func (mva *justMapValueAssembler) AssignFloat(float64) error   { panic("yes") }
func (mva *justMapValueAssembler) AssignString(v string) error { panic("yes") }
func (mva *justMapValueAssembler) AssignBytes([]byte) error    { panic("yes") }
func (mva *justMapValueAssembler) Assign(v ipld.Node) error {
	if mva.midappend == false {
		panic("misuse")
	}
	l := len(mva.w.t) - 1
	mva.w.t[l].v = v
	mva.w.m[mva.w.t[l].k] = v
	return nil
}
func (justMapValueAssembler) Style() ipld.NodeStyle { panic("later") }
