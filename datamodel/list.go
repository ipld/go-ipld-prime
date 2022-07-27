package datamodel

type List interface {
	Get(idx int64) (Node, bool)
	Remove(idx int64)
	Append(values ...interface{})
	Insert(idx int64, values ...interface{})
	Set(idx int64, value interface{})

	Container
}
