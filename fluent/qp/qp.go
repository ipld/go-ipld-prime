// qp is similar to fluent/quip, but with a bit more magic.
package qp

import (
	"github.com/ipld/go-ipld-prime"
)

type Assemble = func(ipld.NodeAssembler)

func BuildMap(np ipld.NodePrototype, sizeHint int64, fn func(ipld.MapAssembler)) (_ ipld.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	nb := np.NewBuilder()
	Map(sizeHint, fn)(nb)
	return nb.Build(), nil
}

type mapParams struct {
	sizeHint int64
	fn       func(ipld.MapAssembler)
}

func (mp mapParams) Assemble(na ipld.NodeAssembler) {
	ma, err := na.BeginMap(mp.sizeHint)
	if err != nil {
		panic(err)
	}
	mp.fn(ma)
	if err := ma.Finish(); err != nil {
		panic(err)
	}
}

func Map(sizeHint int64, fn func(ipld.MapAssembler)) Assemble {
	return mapParams{sizeHint, fn}.Assemble
}

func MapEntry(ma ipld.MapAssembler, k string, fn Assemble) {
	na, err := ma.AssembleEntry(k)
	if err != nil {
		panic(err)
	}
	fn(na)
}

func BuildList(np ipld.NodePrototype, sizeHint int64, fn func(ipld.ListAssembler)) (_ ipld.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	nb := np.NewBuilder()
	List(sizeHint, fn)(nb)
	return nb.Build(), nil
}

type listParams struct {
	sizeHint int64
	fn       func(ipld.ListAssembler)
}

func (lp listParams) Assemble(na ipld.NodeAssembler) {
	la, err := na.BeginList(lp.sizeHint)
	if err != nil {
		panic(err)
	}
	lp.fn(la)
	if err := la.Finish(); err != nil {
		panic(err)
	}
}

func List(sizeHint int64, fn func(ipld.ListAssembler)) Assemble {
	return listParams{sizeHint, fn}.Assemble
}

func ListEntry(la ipld.ListAssembler, fn Assemble) {
	fn(la.AssembleValue())
}

type stringParam string

func (s stringParam) Assemble(na ipld.NodeAssembler) {
	if err := na.AssignString(string(s)); err != nil {
		panic(err)
	}
}

func String(s string) Assemble {
	return stringParam(s).Assemble
}

type intParam int64

func (i intParam) Assemble(na ipld.NodeAssembler) {
	if err := na.AssignInt(int64(i)); err != nil {
		panic(err)
	}
}

func Int(i int64) Assemble {
	return intParam(i).Assemble
}
