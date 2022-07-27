package traversal

import "github.com/ipld/go-ipld-prime/datamodel"

var (
	_ datamodel.Node = &anyAmender{}
	_ Amender        = &anyAmender{}
)

type anyAmender struct {
	amendCfg
}

func (opts AmendOptions) newAnyAmender(base datamodel.Node, parent Amender, create bool) Amender {
	// If the base node is already an any-amender, reuse it but reset `parent` and `created`.
	if amd, castOk := base.(*anyAmender); castOk {
		return &anyAmender{amendCfg{&opts, amd.base, parent, create}}
	} else {
		return &anyAmender{amendCfg{&opts, base, parent, create}}
	}
}

// -- Node -->

func (a *anyAmender) Kind() datamodel.Kind {
	return a.base.Kind()
}

func (a *anyAmender) LookupByString(key string) (datamodel.Node, error) {
	return a.base.LookupByString(key)
}

func (a *anyAmender) LookupByNode(key datamodel.Node) (datamodel.Node, error) {
	return a.base.LookupByNode(key)
}

func (a *anyAmender) LookupByIndex(idx int64) (datamodel.Node, error) {
	return a.base.LookupByIndex(idx)
}

func (a *anyAmender) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	return a.base.LookupBySegment(seg)
}

func (a *anyAmender) MapIterator() datamodel.MapIterator {
	return a.base.MapIterator()
}

func (a *anyAmender) ListIterator() datamodel.ListIterator {
	return a.base.ListIterator()
}

func (a *anyAmender) Length() int64 {
	return a.base.Length()
}

func (a *anyAmender) IsAbsent() bool {
	return a.base.IsAbsent()
}

func (a *anyAmender) IsNull() bool {
	return a.base.IsNull()
}

func (a *anyAmender) AsBool() (bool, error) {
	return a.base.AsBool()
}

func (a *anyAmender) AsInt() (int64, error) {
	return a.base.AsInt()
}

func (a *anyAmender) AsFloat() (float64, error) {
	return a.base.AsFloat()
}

func (a *anyAmender) AsString() (string, error) {
	return a.base.AsString()
}

func (a *anyAmender) AsBytes() ([]byte, error) {
	return a.base.AsBytes()
}

func (a *anyAmender) AsLink() (datamodel.Link, error) {
	return a.base.AsLink()
}

func (a *anyAmender) Prototype() datamodel.NodePrototype {
	return a.base.Prototype()
}

// -- Amender -->

func (a *anyAmender) Fetch(prog *Progress, path datamodel.Path, trackProgress bool) (datamodel.Node, error) {
	// If the base node is an amender, use it, otherwise return the base node.
	if amd, castOk := a.base.(Amender); castOk {
		return amd.Fetch(prog, path, trackProgress)
	}
	return a.base, nil
}

func (a *anyAmender) Transform(prog *Progress, path datamodel.Path, fn TransformFn, createParents bool) (datamodel.Node, error) {
	// Allow the base node to be replaced.
	if path.Len() == 0 {
		prevNode := a.Build()
		if newNode, err := fn(*prog, prevNode); err != nil {
			return nil, err
		} else {
			// Go through `newAnyAmender` in case `newNode` is already an any-amender.
			*a = *a.opts.newAnyAmender(newNode, a.parent, a.created).(*anyAmender)
			return prevNode, nil
		}
	}
	// If the base node is an amender, use it, otherwise panic.
	if amd, castOk := a.base.(Amender); castOk {
		return amd.Transform(prog, path, fn, createParents)
	}
	panic("misuse")
}

func (a *anyAmender) Build() datamodel.Node {
	// `anyAmender` is also a `Node`.
	return (datamodel.Node)(a)
}

func (a *anyAmender) isCreated() bool {
	return a.created
}
