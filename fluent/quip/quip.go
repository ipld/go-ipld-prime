// quip is a package of quick ipld patterns.
//
// Most quip functions take a pointer to an error as their first argument.
// This has two purposes: if there's an error there, the quip function will do nothing;
// and if the quip function does something and creates an error, it puts it there.
// The effect of this is that most logic can be written very linearly.
//
// quip functions can be used to increase brevity without worrying about performance costs.
// None of the quip functions cause additional allocations in the course of their work.
// Benchmarks indicate no measurable speed penalties versus longhand manual error checking.
//
// Several functions perform comparable operations but with different arguments,
// and so these function names follow a pattern:
//
//   - `Build*` functions take a NodePrototype and return a Node.
//   - `Assemble*` functions take a NodeAssembler and feed data into it.
//   - There is no analog of `NodeAssembler.Begin*` functions
//     (we simply always use callbacks for structuring, because this is reasonably optimal).
//   - `Assign*` functions handle values of the scalar kinds
//     (these of course also never need callbacks, since there's no possible recursion).
//
// The `Assemble*` functions are used recursively.
// The `Build*` functions can be used instead of `Assemble*` at the top of a tree
// in order to save on writing a few additional lines of NodeBuilder setup and usage.
// (The `Assemble*` functions can also be used at the top of a tree if you
// wish to control the NodeBuilder yourself.  This may be desirable for being
// able to reset and reuse the NodeBuilder when performance is critical, for example.)
//
// The usual IPLD NodeAssembler, MapAssembler, and ListAssembler interfaces are still
// available while using quip functions, should you wish to interact with them directly,
// or compose the use of quip functions with other styles of data manipulation.
//
package quip

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
)

// - removed the "Begin*" functions; no need to expose that kind of raw operation when the callback forms are zero-cost.
// - renamed functions for consistent "Build" vs "Assemble".
// - added "Assign*" functions for all scalar kinds, which reduces the usage of "AbsorbError" (but also, left AbsorbError in).
// - renamed the ListEntry/MapEntry functions to also have "Assemble*" forms (still callback style).
// - while also adding Assign{Map|List}Entry{Kind} functions (lets you get rid of another callback whenever the value is a scalar).
// - added Assign{|MapEntry|ListEntry} functions, which shell out to fluent.Reflect for even more convenience (at the cost of performance).
// - moved higher level functions like CopyRange to a separate file.
//
// Varations on map key arguments (which could be PathSegment or even Node, in addition to string) still aren't made available this.
// Perhaps that's just okay.  If you're really up some sort of creek where you need that, you can still just use the MapAssembler.AssembleKey system directly.

func AbsorbError(e *error, err error) {
	if *e != nil {
		return
	}
	if err != nil {
		*e = err
	}
}

func BuildMap(e *error, np ipld.NodePrototype, sizeHint int64, fn func(ma ipld.MapAssembler)) ipld.Node {
	if *e != nil {
		return nil
	}
	nb := np.NewBuilder()
	ma, err := nb.BeginMap(sizeHint)
	if err != nil {
		*e = err
		return nil
	}
	fn(ma)
	if *e != nil {
		return nil
	}
	*e = ma.Finish()
	if *e != nil {
		return nil
	}
	return nb.Build()
}

func AssembleMap(e *error, na ipld.NodeAssembler, sizeHint int64, fn func(ma ipld.MapAssembler)) {
	if *e != nil {
		return
	}
	ma, err := na.BeginMap(sizeHint)
	if err != nil {
		*e = err
		return
	}
	fn(ma)
	if *e != nil {
		return
	}
	*e = ma.Finish()
}

func AssembleMapEntry(e *error, ma ipld.MapAssembler, k string, fn func(va ipld.NodeAssembler)) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	fn(va)
}

func BuildList(e *error, np ipld.NodePrototype, sizeHint int64, fn func(la ipld.ListAssembler)) ipld.Node {
	if *e != nil {
		return nil
	}
	nb := np.NewBuilder()
	la, err := nb.BeginList(sizeHint)
	if err != nil {
		*e = err
		return nil
	}
	fn(la)
	if *e != nil {
		return nil
	}
	*e = la.Finish()
	if *e != nil {
		return nil
	}
	return nb.Build()
}

func AssembleList(e *error, na ipld.NodeAssembler, sizeHint int64, fn func(la ipld.ListAssembler)) {
	if *e != nil {
		return
	}
	la, err := na.BeginList(sizeHint)
	if err != nil {
		*e = err
		return
	}
	fn(la)
	if *e != nil {
		return
	}
	*e = la.Finish()
}

func AssembleListEntry(e *error, la ipld.ListAssembler, fn func(va ipld.NodeAssembler)) {
	if *e != nil {
		return
	}
	fn(la.AssembleValue())
}

func AssignNull(e *error, na ipld.NodeAssembler) {
	if *e != nil {
		return
	}
	*e = na.AssignNull()
}
func AssignBool(e *error, na ipld.NodeAssembler, x bool) {
	if *e != nil {
		return
	}
	*e = na.AssignBool(x)
}
func AssignInt(e *error, na ipld.NodeAssembler, x int64) {
	if *e != nil {
		return
	}
	*e = na.AssignInt(x)
}
func AssignFloat(e *error, na ipld.NodeAssembler, x float64) {
	if *e != nil {
		return
	}
	*e = na.AssignFloat(x)
}
func AssignString(e *error, na ipld.NodeAssembler, x string) {
	if *e != nil {
		return
	}
	*e = na.AssignString(x)
}
func AssignBytes(e *error, na ipld.NodeAssembler, x []byte) {
	if *e != nil {
		return
	}
	*e = na.AssignBytes(x)
}
func AssignLink(e *error, na ipld.NodeAssembler, x ipld.Link) {
	if *e != nil {
		return
	}
	*e = na.AssignLink(x)
}
func AssignNode(e *error, na ipld.NodeAssembler, x ipld.Node) {
	if *e != nil {
		return
	}
	*e = na.AssignNode(x)
}

// Assign takes any value and attempts to turn it into something we can reparse as Node-like,
// using the same logic as fluent.Reflect.
// It's not particularly performant, so use it only when convenience matters more than performance.
func Assign(e *error, na ipld.NodeAssembler, x interface{}) {
	if *e != nil {
		return
	}
	*e = fluent.ReflectIntoAssembler(na, x)
}

func AssignMapEntryNull(e *error, ma ipld.MapAssembler, k string) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = va.AssignNull()
}
func AssignMapEntryBool(e *error, ma ipld.MapAssembler, k string, v bool) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = va.AssignBool(v)
}
func AssignMapEntryInt(e *error, ma ipld.MapAssembler, k string, v int64) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = va.AssignInt(v)
}
func AssignMapEntryFloat(e *error, ma ipld.MapAssembler, k string, v float64) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = va.AssignFloat(v)
}
func AssignMapEntryString(e *error, ma ipld.MapAssembler, k string, v string) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = va.AssignString(v)
}
func AssignMapEntryBytes(e *error, ma ipld.MapAssembler, k string, v []byte) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = va.AssignBytes(v)
}
func AssignMapEntryLink(e *error, ma ipld.MapAssembler, k string, v ipld.Link) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = va.AssignLink(v)
}
func AssignMapEntryNode(e *error, ma ipld.MapAssembler, k string, v ipld.Node) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = va.AssignNode(v)
}

// AssignMapEntry takes any value and attempts to turn it into something we can reparse as Node-like,
// using the same logic as fluent.Reflect.
// It's not particularly performant, so use it only when convenience matters more than performance.
func AssignMapEntry(e *error, ma ipld.MapAssembler, k string, x interface{}) {
	if *e != nil {
		return
	}
	va, err := ma.AssembleEntry(k)
	if err != nil {
		*e = err
		return
	}
	*e = fluent.ReflectIntoAssembler(va, x)
}

func AssignListEntryNull(e *error, la ipld.ListAssembler) {
	if *e != nil {
		return
	}
	*e = la.AssembleValue().AssignNull()
}
func AssignListEntryBool(e *error, la ipld.ListAssembler, v bool) {
	if *e != nil {
		return
	}
	*e = la.AssembleValue().AssignBool(v)
}
func AssignListEntryInt(e *error, la ipld.ListAssembler, v int64) {
	if *e != nil {
		return
	}
	*e = la.AssembleValue().AssignInt(v)
}
func AssignListEntryFloat(e *error, la ipld.ListAssembler, v float64) {
	if *e != nil {
		return
	}
	*e = la.AssembleValue().AssignFloat(v)
}
func AssignListEntryString(e *error, la ipld.ListAssembler, v string) {
	if *e != nil {
		return
	}
	*e = la.AssembleValue().AssignString(v)
}
func AssignListEntryBytes(e *error, la ipld.ListAssembler, v []byte) {
	if *e != nil {
		return
	}
	*e = la.AssembleValue().AssignBytes(v)
}
func AssignListEntryLink(e *error, la ipld.ListAssembler, v ipld.Link) {
	if *e != nil {
		return
	}
	*e = la.AssembleValue().AssignLink(v)
}
func AssignListEntryNode(e *error, la ipld.ListAssembler, v ipld.Node) {
	if *e != nil {
		return
	}
	*e = la.AssembleValue().AssignNode(v)
}

// AssignListEntry takes any value and attempts to turn it into something we can reparse as Node-like,
// using the same logic as fluent.Reflect.
// It's not particularly performant, so use it only when convenience matters more than performance.
func AssignListEntry(e *error, la ipld.ListAssembler, x interface{}) {
	if *e != nil {
		return
	}
	*e = fluent.ReflectIntoAssembler(la.AssembleValue(), x)
}
