package impls

// Map_K2_T2 and this file is how a codegen'd map type would work.  it's allowed to use concrete key and value types.
// In constrast with Map_K_T, this one has both complex keys and a struct for the value.

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

// --- we need some types to use for keys and values: --->
/*	ipldsch:
	type K2 struct { u string, i string } representation stringjoin (":")
	type T2 struct { a int, b int, c int, d int }
*/

// Note how we're not able to use `int` in the structs, but instead `plainInt`: this is so we can take address of those fields directly and return them as nodes.
//  We don't currently have concrete exported types that allow us to do this.  Maybe we should?

type K2 struct{ u, i plainString }
type T2 struct{ a, b, c, d plainInt }

func (K2) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (n *K2) LookupString(key string) (ipld.Node, error) {
	switch key {
	case "u":
		return &n.u, nil
	case "i":
		return &n.i, nil
	default:
		return nil, fmt.Errorf("no such field")
	}
}
func (n *K2) Lookup(key ipld.Node) (ipld.Node, error) {
	ks, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return n.LookupString(ks)
}
func (K2) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: "K2", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (n *K2) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return n.LookupString(seg.String())
}
func (n *K2) MapIterator() ipld.MapIterator {
	return &_K2_MapIterator{n, 0}
}
func (K2) ListIterator() ipld.ListIterator {
	return nil
}
func (K2) Length() int {
	return -1
}
func (K2) IsUndefined() bool {
	return false
}
func (K2) IsNull() bool {
	return false
}
func (K2) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsString() (string, error) {
	return "", ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Map}
}
func (K2) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{TypeName: "K2", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Map}
}
func (K2) Style() ipld.NodeStyle {
	panic("todo")
}

type _K2_MapIterator struct {
	n   *K2
	idx int
}

func (itr *_K2_MapIterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.idx >= 2 {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	switch itr.idx {
	case 0:
		k = plainString("u") // TODO: I guess we should generate const pools for struct field names?
		v = &itr.n.u
	case 1:
		k = plainString("i")
		v = &itr.n.i
	default:
		panic("unreachable")
	}
	itr.idx++
	return
}
func (itr *_K2_MapIterator) Done() bool {
	return itr.idx >= 2
}

type _K2__Assembler struct {
	w *K2

	state maState

	// ca_u // this field is primitive kind and we know it, so we can just do a bunch of stuff directly.
	// ca_i // this field is primitive kind and we know it, so we can just do a bunch of stuff directly.
	// ^ not at all sure this is true, we need SOMETHING to yield, and it's not a uniform thing in general like it is for maps.
	// ^ and it's not clear we're saving much by trying to avoid it or specialize it.  it'll get largely inlined anyway.  ... wait
	// holy kjott i continue to need finish checkers in all places and
	// i do not want
	// to generate a wrapper PER FIELD TYPE.
	// that was possible for maps because it's just one, but here?!  oh my god.  no.
	// we need interfaces for this.  we do.
	//...
	// `typed.wrapper` can use just one monomorphic type, i think.
	// datamodel stuff alone just doesn't have this issue.
	// codegen'd stuff is the only one who has an issue here.
	// not sure if that observation offers any usefully clearer path to solution.
	//...
	// if would seem to imply that if we do use an interface here, it can actually (bizarrely enough) be an unexported one.
	// the idea of pointers in every assembler to a callback, and the parent just sets that like it does the 'w' node, is also still on the table.
	//  (as long as we can verify that we won't accidentally spawn a new closure alloc to do that -- but this should be easy enough.)
	//  what to do if this is nil bothers me, but... actually, i think a fixed noop is best; there's already constructors fencing all assemblers, so we can do this.
	// or we can indeed generate a type per field, and that's honestly gonna be at least one more thing inlined rather than funcptr call,
	//  if at fairly horrendous cost in GLOC as well as binary size.
	//  it's the same cost in memory, surprisingly: either a ptr to the callback, or the generated type has a pointer to the parent 'w'.
	//...
	// okay, so i'm coming towards prefering the callback, given the realizations about {moot memory size + unexportability + no nil checks needed}.
	// an interesting consequence of this, though, is that yep, every codegen'd package emitted is going to have its own typedefs of primitives.
	//  but this is just a "sigh" thing, not any kind of problem, since we already defacto expect most packages providing Node to have their own on these.
	//...
	// we could also make a single va{*ca,*maState} per type and have it do a bunch of delegating calls.
	//  the result of this would be a vtable jump for the delegation (otherwise we have unreusablity)... but that's similar to the update hook idea;
	//  and at the same time, it's adding two words to the parent assembler, instead of a word to every field's assembler, which is likely better.
	//  ... hang on, no, this is nontrivial for recursives: we don't want to have to decorate every intermediate key and value assemble call in a map builder just to keep control on the stack so we can intercept the 'finish'; that's *way* worse than just one ptrfunc jump at the end.
	//...
	// a single clear correct choice is not yet clear here.  some further analysis may be needed.

	// - untyped maps can have two more types: plainMap__ValueAssemblerMap and plainMap__ValueAssemblerList...
	//   - and all they do is override the Finish method.
	//   - they have to be allocated as pointers, but this ain't news.
	// - typed maps can be basically the same, in fact.
	//   - the Map__K__V__ValueAssembler will also vary scalar values, but this ain't news.
	// - structs...
	//   - honestly, is a type per field so bad?  do types have a meaningful cost, when you're not banging them out at the keyboard?
	//     - i don't think they do but actually come to think of it i'd love to measure that.
	//       - okayyyyyy.  ballpark it at 42kb per type, if you believe `go build -o`.  wow; that's a lot more than I expected.
	//       - or 11kb per type if you believe the change in size of 'impls.test' binaries produced as a sideeffect of '-memprofile'.
	//         - why are these different?  no idea.  the thing from `go build -o` is an 'ar' file rather than binary, though.
	//         - pretty sure this is a more accurate number than 42.
	//       - this is for the addition of 'plainMap__ValueAssemblerMap', to be specific.
	//       - now, just for context, this comment block here has gotte up to 4.4kb, so, uh.  lol.  kilobytes go fast.
	//       - it depends a LOT on how many methods the type has.  embedding the node interface can actually be a lot.
	//         - there's a measurable difference between {every method on NodeAssembler} and {just MapNodeAssembler} methods.
	//         - also depends on how exactly the methods either inline or autogenerate; some forms of autogen are much longer than others (some just JMP; others do a ton of setup and then CALL).
	//           - and I don't at present have a predictive understanding of why this is.
	//       - okay, the 'plainMap__ValueAssemblerMap' type is down to adding 9717 bytes after bothering to get more unneeded autogen methods out.  So we also have some control over this.
	//     - approaching ballparks from the other side: schema-schema has about 37 fields in it.
	//       - so we if take '11kb' as a random number, that means generating a type per field might add about .4mb to final binaries in this practical example.
	//         - is that acceptable?  it's not great, but I think it's probably fine, also.  I'd usually rather pay that than lose run speed.

	// in review: what are we balancing?
	// - AS: asm size
	// - BM: builder mem
	// - SP: execution speed
	// - AC: allocation count -- but only because of its outsized impact on SP.
	//
	// I care about SP>BM>AS.
	//  (Unless BM gets > 2x, or AS gets really out of hand.)
	//  (BM also has a special cond where if it increases on recursive kinds, but not on scalars,
	//   we regard that as half price, because generally most of a tree is leaves.)

	// So in this case, the question is "did AS get 'really out of hand'" if we choose the type-per-field route.
	//  And the answer seems to be leaning towards "good question... but no, it's not (quite) 'really out of hand'".

	// Go figure: in this situation, SP>BM>AS does seem to resolve to generating more code and more types.

	isset_u bool
	isset_i bool
}
type _K2__ReprAssembler struct {
	w *K2

	// note how this is totally different than the type-level assembler -- that's map-like, this is string.
}

func (ta *_K2__Assembler) BeginMap(_ int) (ipld.MapNodeAssembler, error) { panic("no") }
func (_K2__Assembler) BeginList(_ int) (ipld.ListNodeAssembler, error)   { panic("no") }
func (_K2__Assembler) AssignNull() error                                 { panic("no") }
func (_K2__Assembler) AssignBool(bool) error                             { panic("no") }
func (_K2__Assembler) AssignInt(v int) error                             { panic("no") }
func (_K2__Assembler) AssignFloat(float64) error                         { panic("no") }
func (_K2__Assembler) AssignString(v string) error                       { panic("no") }
func (_K2__Assembler) AssignBytes([]byte) error                          { panic("no") }
func (ta *_K2__Assembler) AssignNode(v ipld.Node) error {
	if v2, ok := v.(*K2); ok {
		*ta.w = *v2
		return nil
	}
	panic("todo implement generic copy and use it here")
}
func (_K2__Assembler) Style() ipld.NodeStyle { panic("later") }

func (ma *_K2__Assembler) AssembleDirectly(k string) (ipld.NodeAssembler, error) {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_midValue
	// Figure out which field we're addressing,
	//  check if it's already been assigned (error if so),
	//   grab a pointer to it and init its value assembler with that,
	//    and yield that value assembler.
	//  (Note that `isset_foo` bools may be inside the 'ma.w' node if
	//   that field is optional; if it's required, they stay in 'ma'.)
	switch k {
	case "u":
		if ma.isset_u {
			return nil, ipld.ErrRepeatedMapKey{plainString("u")} // REVIEW: interesting to note this is a place we *keep* needing a basic string node impl, *everywhere*.
		}
		// TODO initialize the field child assembler 'w' *and* 'finish' callback to us; return it.
		panic("todo")
	case "i":
		// TODO same as above
		panic("todo")
	default:
		panic("invalid field key")
	}
}

func (ma *_K2__Assembler) AssembleKey() ipld.NodeAssembler {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_midKey
	// TODO return a fairly dummy assembler which just contains a string switch (probably sharing code with AssembleDirectly).
	panic("todo")
}
func (ma *_K2__Assembler) AssembleValue() ipld.NodeAssembler {
	// Sanity check, then update, assembler state.
	if ma.state != maState_expectValue {
		panic("misuse")
	}
	ma.state = maState_midValue
	// TODO initialize the field child assembler 'w' *and* 'finish' callback to us; return it.
	panic("todo")
}
func (ma *_K2__Assembler) Finish() error {
	// Sanity check assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_finished
	// validators could run and report errors promptly, if this type had any.
	return nil
}
func (_K2__Assembler) KeyStyle() ipld.NodeStyle           { panic("later") }
func (_K2__Assembler) ValueStyle(k string) ipld.NodeStyle { panic("later") }

func (T2) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (n *T2) LookupString(key string) (ipld.Node, error) {
	switch key {
	case "a":
		return &n.a, nil
	case "b":
		return &n.b, nil
	case "c":
		return &n.c, nil
	case "d":
		return &n.d, nil
	default:
		return nil, fmt.Errorf("no such field")
	}
}
func (n *T2) Lookup(key ipld.Node) (ipld.Node, error) {
	ks, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return n.LookupString(ks)
}
func (T2) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: "T2", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (n *T2) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return n.LookupString(seg.String())
}
func (n *T2) MapIterator() ipld.MapIterator {
	return &_T2_MapIterator{n, 0}
}
func (T2) ListIterator() ipld.ListIterator {
	return nil
}
func (T2) Length() int {
	return -1
}
func (T2) IsUndefined() bool {
	return false
}
func (T2) IsNull() bool {
	return false
}
func (T2) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsString() (string, error) {
	return "", ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Map}
}
func (T2) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{TypeName: "T2", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Map}
}
func (T2) Style() ipld.NodeStyle {
	panic("todo")
}

type _T2_MapIterator struct {
	n   *T2
	idx int
}

func (itr *_T2_MapIterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.idx >= 4 {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	switch itr.idx {
	case 0:
		k = plainString("a") // TODO: I guess we should generate const pools for struct field names?
		v = &itr.n.a
	case 1:
		k = plainString("b")
		v = &itr.n.b
	case 2:
		k = plainString("c")
		v = &itr.n.c
	case 3:
		k = plainString("d")
		v = &itr.n.d
	default:
		panic("unreachable")
	}
	itr.idx++
	return
}
func (itr *_T2_MapIterator) Done() bool {
	return itr.idx >= 4
}

type _T2__Assembler struct {
	w *T2
}
type _T2__ReprAssembler struct {
	w *T2
}

func (ta *_T2__Assembler) BeginMap(_ int) (ipld.MapNodeAssembler, error) {
	return ta, nil
}
func (_T2__Assembler) BeginList(_ int) (ipld.ListNodeAssembler, error) { panic("no") }
func (_T2__Assembler) AssignNull() error                               { panic("no") }
func (_T2__Assembler) AssignBool(bool) error                           { panic("no") }
func (_T2__Assembler) AssignInt(int) error                             { panic("no") }
func (_T2__Assembler) AssignFloat(float64) error                       { panic("no") }
func (_T2__Assembler) AssignString(v string) error                     { panic("no") }
func (_T2__Assembler) AssignBytes([]byte) error                        { panic("no") }
func (ta *_T2__Assembler) AssignNode(v ipld.Node) error {
	if v2, ok := v.(*T2); ok {
		*ta.w = *v2
		return nil
	}
	// todo: apply a generic 'copy' function.
	panic("later")
}
func (_T2__Assembler) Style() ipld.NodeStyle { panic("later") }

func (ta *_T2__Assembler) AssembleDirectly(string) (ipld.NodeAssembler, error) {
	// this'll be fun
	panic("soon")
}
func (ta *_T2__Assembler) AssembleKey() ipld.NodeAssembler {
	// this'll be fun
	panic("soon")
}
func (ta *_T2__Assembler) AssembleValue() ipld.NodeAssembler {
	// also fun
	panic("soon")
}
func (ta *_T2__Assembler) Finish() error {
	panic("soon")
}
func (_T2__Assembler) KeyStyle() ipld.NodeStyle           { panic("later") }
func (_T2__Assembler) ValueStyle(k string) ipld.NodeStyle { panic("later") }

// --- okay, now the type of interest: the map. --->
/*	ipldsch:
	type Root struct { mp {K2:T2} } # nevermind the root part, the anonymous map is the point.
*/

type Map_K2_T2 struct {
	m map[K2]*T2          // used for quick lookup.
	t []_Map_K2_T2__entry // used both for order maintainence, and for allocation amortization for both keys and values.
}

type _Map_K2_T2__entry struct {
	k K2 // address of this used when we return keys as nodes, such as in iterators.  Need in one place to amortize shifts to heap when ptr'ing for iface.
	v T2 // address of this is used in map values and to return.
}

func (n *Map_K2_T2) LookupString(key string) (ipld.Node, error) {
	panic("decision") // FIXME: What's this supposed to do?  does this error for maps with complex keys?
}

type _Map_K2_T2__Assembler struct {
	w  *Map_K2_T2
	ka _K2__Assembler
	va _T2__Assembler
}
type _Map_K2_T2__ReprAssembler struct {
	w  *Map_K2_T2
	ka _K2__ReprAssembler
	va _T2__ReprAssembler
}
