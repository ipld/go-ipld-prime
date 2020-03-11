package fluent

import (
	ipld "github.com/ipld/go-ipld-prime"
)

func Build(ns ipld.NodeStyle, fn func(NodeAssembler)) (ipld.Node, error) {
	nb := ns.NewBuilder()
	fna := WrapAssembler(nb)
	err := Recover(func() {
		fn(fna)
	})
	return nb.Build(), err
}

func MustBuild(ns ipld.NodeStyle, fn func(NodeAssembler)) ipld.Node {
	nb := ns.NewBuilder()
	fn(WrapAssembler(nb))
	return nb.Build()
}
func MustBuildMap(ns ipld.NodeStyle, sizeHint int, fn func(MapNodeAssembler)) ipld.Node {
	return MustBuild(ns, func(fna NodeAssembler) { fna.CreateMap(sizeHint, fn) })
}
func MustBuildList(ns ipld.NodeStyle, sizeHint int, fn func(ListNodeAssembler)) ipld.Node {
	return MustBuild(ns, func(fna NodeAssembler) { fna.CreateList(sizeHint, fn) })
}

func WrapAssembler(na ipld.NodeAssembler) NodeAssembler {
	return &nodeAssembler{na}
}

// NodeAssembler is the same as the interface in the core package, except:
// instead of returning errors, any error will cause panic
// (and you can collect these with `fluent.Recover`);
// and all recursive operations take a function as a parameter,
// within which you will receive another {Map,List,}NodeAssembler.
type NodeAssembler interface {
	CreateMap(sizeHint int, fn func(MapNodeAssembler))
	CreateList(sizeHint int, fn func(ListNodeAssembler))
	AssignNull()
	AssignBool(bool)
	AssignInt(int)
	AssignFloat(float64)
	AssignString(string)
	AssignBytes([]byte)
	AssignLink(ipld.Link)
	AssignNode(ipld.Node)

	Style() ipld.NodeStyle
}

// MapNodeAssembler is the same as the interface in the core package, except:
// instead of returning errors, any error will cause panic
// (and you can collect these with `fluent.Recover`);
// and all recursive operations take a function as a parameter,
// within which you will receive another {Map,List,}NodeAssembler.
type MapNodeAssembler interface {
	AssembleKey() NodeAssembler
	AssembleValue() NodeAssembler

	AssembleDirectly(k string) NodeAssembler

	KeyStyle() ipld.NodeStyle
	ValueStyle(k string) ipld.NodeStyle
}

// ListNodeAssembler is the same as the interface in the core package, except:
// instead of returning errors, any error will cause panic
// (and you can collect these with `fluent.Recover`);
// and all recursive operations take a function as a parameter,
// within which you will receive another {Map,List,}NodeAssembler.
type ListNodeAssembler interface {
	AssembleValue() NodeAssembler

	ValueStyle() ipld.NodeStyle
}

type nodeAssembler struct {
	na ipld.NodeAssembler
}

func (fna *nodeAssembler) CreateMap(sizeHint int, fn func(MapNodeAssembler)) {
	if ma, err := fna.na.BeginMap(sizeHint); err != nil {
		panic(Error{err})
	} else {
		fn(&mapNodeAssembler{ma})
		if err := ma.Finish(); err != nil {
			panic(Error{err})
		}
	}
}
func (fna *nodeAssembler) CreateList(sizeHint int, fn func(ListNodeAssembler)) {
	if la, err := fna.na.BeginList(sizeHint); err != nil {
		panic(Error{err})
	} else {
		fn(&listNodeAssembler{la})
		if err := la.Finish(); err != nil {
			panic(Error{err})
		}
	}
}
func (fna *nodeAssembler) AssignNull() {
	if err := fna.na.AssignNull(); err != nil {
		panic(Error{err})
	}
}
func (fna *nodeAssembler) AssignBool(v bool) {
	if err := fna.na.AssignBool(v); err != nil {
		panic(Error{err})
	}
}
func (fna *nodeAssembler) AssignInt(v int) {
	if err := fna.na.AssignInt(v); err != nil {
		panic(Error{err})
	}
}
func (fna *nodeAssembler) AssignFloat(v float64) {
	if err := fna.na.AssignFloat(v); err != nil {
		panic(Error{err})
	}
}
func (fna *nodeAssembler) AssignString(v string) {
	if err := fna.na.AssignString(v); err != nil {
		panic(Error{err})
	}
}
func (fna *nodeAssembler) AssignBytes(v []byte) {
	if err := fna.na.AssignBytes(v); err != nil {
		panic(Error{err})
	}
}
func (fna *nodeAssembler) AssignLink(v ipld.Link) {
	if err := fna.na.AssignLink(v); err != nil {
		panic(Error{err})
	}
}
func (fna *nodeAssembler) AssignNode(v ipld.Node) {
	if err := fna.na.AssignNode(v); err != nil {
		panic(Error{err})
	}
}
func (fna *nodeAssembler) Style() ipld.NodeStyle {
	return fna.na.Style()
}

type mapNodeAssembler struct {
	ma ipld.MapNodeAssembler
}

func (fma *mapNodeAssembler) AssembleKey() NodeAssembler {
	return &nodeAssembler{fma.ma.AssembleKey()}
}
func (fma *mapNodeAssembler) AssembleValue() NodeAssembler {
	return &nodeAssembler{fma.ma.AssembleValue()}
}
func (fma *mapNodeAssembler) AssembleDirectly(k string) NodeAssembler {
	va, err := fma.ma.AssembleDirectly(k)
	if err != nil {
		panic(Error{err})
	}
	return &nodeAssembler{va}
}
func (fma *mapNodeAssembler) KeyStyle() ipld.NodeStyle {
	return fma.ma.KeyStyle()
}
func (fma *mapNodeAssembler) ValueStyle(k string) ipld.NodeStyle {
	return fma.ma.ValueStyle(k)
}

type listNodeAssembler struct {
	la ipld.ListNodeAssembler
}

func (fla *listNodeAssembler) AssembleValue() NodeAssembler {
	return &nodeAssembler{fla.la.AssembleValue()}
}
func (fla *listNodeAssembler) ValueStyle() ipld.NodeStyle {
	return fla.la.ValueStyle()
}
