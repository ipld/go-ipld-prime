package datamodel

type Map interface {
	Put(key string, value interface{}) bool
	Get(key string) (value Node, found bool)
	Remove(key string) bool
	Keys() []string

	Container
}
