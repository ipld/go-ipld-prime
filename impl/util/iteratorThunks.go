package nodeutil

import (
	ipld "github.com/ipld/go-ipld-prime"
)

func MapIteratorErrorThunk(err error) ipld.MapIterator {
	return mapIteratorErrorThunk{err}
}
func ListIteratorErrorThunk(err error) ipld.ListIterator {
	return listIteratorErrorThunk{err}
}

type mapIteratorErrorThunk struct{ err error }
type listIteratorErrorThunk struct{ err error }

func (itr mapIteratorErrorThunk) Next() (ipld.Node, ipld.Node, error) { return nil, nil, itr.err }
func (itr mapIteratorErrorThunk) Done() bool                          { return false }

func (itr listIteratorErrorThunk) Next() (int, ipld.Node, error) { return -1, nil, itr.err }
func (itr listIteratorErrorThunk) Done() bool                    { return false }
