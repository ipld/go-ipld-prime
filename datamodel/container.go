package datamodel

type Container interface {
	Empty() bool
	Length() int64
	Clear()
	Values() []Node
}
