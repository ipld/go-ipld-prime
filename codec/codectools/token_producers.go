package codectools

import (
	"errors"
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
)

// TokenWalk walks an ipld Node and repeatedly calls the visitFn,
// calling it once for every "token" yielded by the walk.
// Every map and list is yielded as a token at their beginning,
// and another token when they're finished;
// every scalar value (strings, bools, bytes, ints, etc) is yielded as a single token.
//
// The token pointer given to the visitFn will be identical on every call,
// but the data it contains will vary.
// The token may contain invalid data that is leftover from previous calls
// in some of its union fields; correct behavior requires looking at the
// token's Kind field before handling any of its other fields.
//
// If any error is returned by the visitFn, it will cause the walk to halt,
// and TokenWalk will return that error.
// However, if the error is the value TokenWalkSkip, and it's been returned
// when visitFn was called with a MapOpen or ListOpen token, the walk will
// skip forward over that entire map or list, and continue (with the
// next token being the close token that complements the open token).
// Returning a TokenWalkSkip when the token was any of the scalar kinds
// (e.g. anything other than a MapOpen or a ListOpen) has no effect.
//
// TokenAssembler is the rough dual of TokenWalk.
func TokenWalk(n ipld.Node, visitFn func(tk *Token) error) error {
	// TokenWalk would be trivial to implement over NodeTokenizer,
	//  but we do a distinct implementation here because NodeTokenizer's resumable implementation means it needs a user-space stack,
	//  and to reuse that would require allocations which this method (since it's not resumable in the same way) can easily avoid (or at least, keep on the stack).

	var tk Token // For capture, once.
	return tokenWalk(&tk, n, visitFn)
}

func tokenWalk(tk *Token, n ipld.Node, visitFn func(*Token) error) error {
	switch n.ReprKind() {
	case ipld.ReprKind_Map:
		tk.Kind = TokenKind_MapOpen
		tk.Length = n.Length()
		tk.Node = n
		if err := visitFn(tk); err != nil {
			return err
		}
		mitr := n.MapIterator()
		for !mitr.Done() {
			k, v, err := mitr.Next()
			if err != nil {
				return err
			}
			if err := tokenWalk(tk, k, visitFn); err != nil {
				return err
			}
			if err := tokenWalk(tk, v, visitFn); err != nil {
				return err
			}
		}
		tk.Kind = TokenKind_MapClose
		tk.Node = n
		return visitFn(tk)
	case ipld.ReprKind_List:
		tk.Kind = TokenKind_ListOpen
		tk.Length = n.Length()
		tk.Node = n
		if err := visitFn(tk); err != nil {
			return err
		}
		litr := n.ListIterator()
		for !litr.Done() {
			_, v, err := litr.Next()
			if err != nil {
				return err
			}
			if err := tokenWalk(tk, v, visitFn); err != nil {
				return err
			}
		}
		tk.Kind = TokenKind_ListClose
		tk.Node = n
		return visitFn(tk)
	case ipld.ReprKind_Null:
		tk.Kind = TokenKind_Null
		return visitFn(tk)
	case ipld.ReprKind_Bool:
		tk.Kind = TokenKind_Bool
		tk.Bool, _ = n.AsBool()
		return visitFn(tk)
	case ipld.ReprKind_Int:
		tk.Kind = TokenKind_Int
		i, _ := n.AsInt()
		tk.Int = int64(i) // TODO: upgrade all of ipld to use high precision int consistently
		return visitFn(tk)
	case ipld.ReprKind_Float:
		tk.Kind = TokenKind_Float
		tk.Float, _ = n.AsFloat()
		return visitFn(tk)
	case ipld.ReprKind_String:
		tk.Kind = TokenKind_String
		tk.Str, _ = n.AsString()
		return visitFn(tk)
	case ipld.ReprKind_Bytes:
		tk.Kind = TokenKind_Bytes
		tk.Bytes, _ = n.AsBytes()
		return visitFn(tk)
	case ipld.ReprKind_Link:
		tk.Kind = TokenKind_Link
		tk.Link, _ = n.AsLink()
		return visitFn(tk)
	default:
		panic(fmt.Errorf("unrecognized node kind (%q?)", n.ReprKind()))
	}
	return nil
}

var TokenWalkSkip = errors.New("token walk: skip")

// --- the stepwise token producer system (more complicated; has a userland stack) is below -->

//
// A TokenReader can be produced from any ipld.Node using NodeTokenizer.
// TokenReader are also commonly implemented by codec packages,
// wherein they're created over a serial data stream and tokenize that stream when pumped.
//
// TokenReader implementations are encouraged to yield the same token pointer repeatedly,
// just varying the contents of the value, in order to avoid unnecessary allocations.
//
// TODO: as elegant as this is, it's not able to provide much help for error reporting:
//  if I hit an error while handling the token, I'd like to be able to ask this thing where it thinks it is,
//   and include that position info in my error report.
//  Maybe putting position info directly into the Token struct would solve this satisfactorily?
//   More comments can be found in the Token definition.
//
// TODO: this probably ought to take a budget parameter.
//  It doesn't make much sense if you're walking in-memory data,
//  but it's sure relevant if you're parsing serial data and want to pass down info for limiting how big of a string is allowed to be allocated.
type TokenReader func() (next *Token, err error)

type NodeTokenizer struct {
	// This structure is designed to be embeddable.  Use Initialize when doing so.

	tk  Token // We embed this to avoid allocations; we'll be repeatedly yielding a pointer to this piece of memory.
	stk nodeTokenizerStack
}

func (nt *NodeTokenizer) Initialize(n ipld.Node) {
	if nt.stk == nil {
		nt.stk = make(nodeTokenizerStack, 0, 10)
	} else {
		nt.stk = nt.stk[0:0]
	}
	nt.stk.Push(n)
}

type nodeTokenizerStackRow struct {
	state uint8             // 0: start this node; 1: continue list; 2: continue map with key; 3: continue map with value.
	n     ipld.Node         // Always present.
	litr  ipld.ListIterator // At most one of these is present.
	mitr  ipld.MapIterator  // At most one of these is present.
	mval  ipld.Node         // The value to resume at when in state 3.

}
type nodeTokenizerStack []nodeTokenizerStackRow

func (stk nodeTokenizerStack) Tip() *nodeTokenizerStackRow {
	return &stk[len(stk)-1]
}
func (stk *nodeTokenizerStack) Push(n ipld.Node) {
	*stk = append(*stk, nodeTokenizerStackRow{n: n})
}
func (stk *nodeTokenizerStack) Pop() {
	if len(*stk) == 0 {
		return
	}
	*stk = (*stk)[0 : len(*stk)-1]
}

// ReadToken fits the TokenReader functional interface, and so may be used anywhere a TokenReader is required.
func (nt *NodeTokenizer) ReadToken() (next *Token, err error) {
	// How stack depth works:
	// - finding that you're starting to handle map or least leaves it the same;
	// - before recursing to handle a child key or value, push stack;
	// - any time you finish something, whether scalar or recursive, pop stack.
	// This could be written differently: in particular,
	//  scalar leaves could be handled without increasing stack depth by that last increment.
	//   However, doing so would make for more complicated code.
	//    Maybe worth it; PRs welcome; benchmarks first.
	if len(nt.stk) == 0 {
		return nil, io.EOF
	}
	tip := nt.stk.Tip()
	switch tip.state {
	case 0:
		switch tip.n.ReprKind() {
		case ipld.ReprKind_Map:
			nt.tk.Kind = TokenKind_MapOpen
			nt.tk.Length = tip.n.Length()
			nt.tk.Node = tip.n
			tip.state = 2
			tip.mitr = tip.n.MapIterator()
			return &nt.tk, nil
		case ipld.ReprKind_List:
			nt.tk.Kind = TokenKind_ListOpen
			nt.tk.Length = tip.n.Length()
			nt.tk.Node = tip.n
			tip.state = 1
			tip.litr = tip.n.ListIterator()
			return &nt.tk, nil
		case ipld.ReprKind_Null:
			nt.tk.Kind = TokenKind_Null
			nt.stk.Pop()
			return &nt.tk, nil
		case ipld.ReprKind_Bool:
			nt.tk.Kind = TokenKind_Bool
			nt.tk.Bool, _ = tip.n.AsBool()
			nt.stk.Pop()
			return &nt.tk, nil
		case ipld.ReprKind_Int:
			nt.tk.Kind = TokenKind_Int
			i, _ := tip.n.AsInt()
			nt.tk.Int = int64(i) // TODO: upgrade all of ipld to use high precision int consistently
			nt.stk.Pop()
			return &nt.tk, nil
		case ipld.ReprKind_Float:
			nt.tk.Kind = TokenKind_Float
			nt.tk.Float, _ = tip.n.AsFloat()
			nt.stk.Pop()
			return &nt.tk, nil
		case ipld.ReprKind_String:
			nt.tk.Kind = TokenKind_String
			nt.tk.Str, _ = tip.n.AsString()
			nt.stk.Pop()
			return &nt.tk, nil
		case ipld.ReprKind_Bytes:
			nt.tk.Kind = TokenKind_Bytes
			nt.tk.Bytes, _ = tip.n.AsBytes()
			nt.stk.Pop()
			return &nt.tk, nil
		case ipld.ReprKind_Link:
			nt.tk.Kind = TokenKind_Link
			nt.tk.Link, _ = tip.n.AsLink()
			nt.stk.Pop()
			return &nt.tk, nil
		default:
			panic(fmt.Errorf("unrecognized node kind (%q?)", tip.n.ReprKind()))
		}
	case 1:
		if tip.litr.Done() {
			nt.tk.Kind = TokenKind_ListClose
			nt.tk.Node = tip.n
			nt.stk.Pop()
			return &nt.tk, nil
		}
		_, v, err := tip.litr.Next()
		if err != nil {
			return nil, err
		}
		nt.stk.Push(v)
		return nt.ReadToken()
	case 2:
		if tip.mitr.Done() {
			nt.tk.Kind = TokenKind_MapClose
			nt.tk.Node = tip.n
			nt.stk.Pop()
			return &nt.tk, nil
		}
		k, v, err := tip.mitr.Next()
		if err != nil {
			return nil, err
		}
		tip.mval = v
		tip.state = 3
		nt.stk.Push(k)
		return nt.ReadToken()
	case 3:
		tip.state = 2
		nt.stk.Push(tip.mval)
		return nt.ReadToken()
	default:
		panic("unreachable")
	}
}
