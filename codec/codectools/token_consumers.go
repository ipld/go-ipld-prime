package codectools

import (
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec"
)

// TokenAssemble takes an ipld.NodeAssembler and a TokenReader,
// and repeatedly pumps the TokenReader for tokens and feeds their data into the ipld.NodeAssembler
// until it finishes a complete value.
//
// To compare and contrast to other token oriented tools:
// TokenAssemble does the same direction of information transfer as the TokenAssembler gadget does,
// but TokenAssemble moves completely through a value in one step,
// whereas the TokenAssembler accepts tokens pumped into it one step at a time.
//
// TokenAssemble does not enforce the "map keys must be strings" rule which is present in the Data Model;
// it will also happily do even recursive structures in map keys,
// meaning it can be used when handling schema values like maps with complex keys.
func TokenAssemble(na ipld.NodeAssembler, tr TokenReader, budget int) error {
	tk, err := tr()
	if err != nil {
		return err
	}
	return tokenAssemble(na, tk, tr, &budget)
}

func tokenAssemble(na ipld.NodeAssembler, tk *Token, tr TokenReader, budget *int) error {
	if *budget < 0 {
		return codec.ErrBudgetExhausted{}
	}
	switch tk.Kind {
	case TokenKind_MapOpen:
		if tk.Length > 0 && *budget < tk.Length*2 { // Pre-check budget: at least two decrements estimated for each entry.
			return codec.ErrBudgetExhausted{}
		}
		ma, err := na.BeginMap(tk.Length)
		if err != nil {
			return err
		}
		for {
			// Peek one token.  We need to see if the map is about to end or not.
			tk, err = tr()
			if err != nil {
				return err
			}
			// If the map has ended, invoke the finish operation and check for any errors.
			if tk.Kind == TokenKind_MapClose {
				return ma.Finish()
			}
			// Recurse to assemble the key.
			*budget-- // Decrement budget by at least one for each key.  The key content may also cause further decrements.
			if err = tokenAssemble(ma.AssembleKey(), tk, tr, budget); err != nil {
				return err
			}
			// Recurse to assemble the value.
			//  (We don't really care to peek this token, but do so anyway to keep the calling convention regular.)
			tk, err = tr()
			if err != nil {
				return err
			}
			*budget-- // Decrement budget by at least one for each value.  The value content may also cause further decrements.
			if err = tokenAssemble(ma.AssembleValue(), tk, tr, budget); err != nil {
				return err
			}
			// Continue around the loop, to encounter either the next entry or the end of the map.
		}
	case TokenKind_MapClose:
		return ErrMalformedTokenSequence{"map close token encountered while not in the middle of a map"}
	case TokenKind_ListOpen:
		if tk.Length > 0 && *budget < tk.Length { // Pre-check budget: at least one decrement estimated for each entry.
			return codec.ErrBudgetExhausted{}
		}
		la, err := na.BeginList(tk.Length)
		if err != nil {
			return err
		}
		for {
			// Peek one token.  We need to see if the list is about to end or not.
			tk, err = tr()
			if err != nil {
				return err
			}
			// If the list has ended, invoke the finish operation and check for any errors.
			if tk.Kind == TokenKind_ListClose {
				return la.Finish()
			}
			// Recurse to assemble the value.
			*budget-- // Decrement budget by at least one for each value.  The value content may also cause further decrements.
			if err = tokenAssemble(la.AssembleValue(), tk, tr, budget); err != nil {
				return err
			}
			// Continue around the loop, to encounter either the next value or the end of the list.
		}
	case TokenKind_ListClose:
		return ErrMalformedTokenSequence{"list close token encountered while not in the middle of a list"}
	case TokenKind_Null:
		return na.AssignNull()
	case TokenKind_Bool:
		*budget--
		return na.AssignBool(tk.Bool)
	case TokenKind_Int:
		*budget--
		return na.AssignInt(int(tk.Int))
	case TokenKind_Float:
		*budget--
		return na.AssignFloat(tk.Float)
	case TokenKind_String:
		*budget -= len(tk.Str)
		return na.AssignString(tk.Str)
	case TokenKind_Bytes:
		*budget -= len(tk.Bytes)
		return na.AssignBytes(tk.Bytes)
	case TokenKind_Link:
		*budget--
		return na.AssignLink(tk.Link)
	default:
		panic(fmt.Errorf("unrecognized token kind (%q?)", tk.Kind))
	}
}

// --- the stepwise assembler system (more complicated; has a userland stack) is below -->

type TokenAssembler struct {
	// This structure is designed to be embeddable.  Use Initialize when doing so.

	stk    assemblerStack // this is going to end up being a stack you know
	budget int64
}

type assemblerStackRow struct {
	state uint8              // 0: assign this node; 1: continue list; 2: continue map with key; 3: continue map with value.
	na    ipld.NodeAssembler // Always present.
	la    ipld.ListAssembler // At most one of these is present.
	ma    ipld.MapAssembler  // At most one of these is present.
}
type assemblerStack []assemblerStackRow

func (stk assemblerStack) Tip() *assemblerStackRow {
	return &stk[len(stk)-1]
}
func (stk *assemblerStack) Push(na ipld.NodeAssembler) {
	*stk = append(*stk, assemblerStackRow{na: na})
}
func (stk *assemblerStack) Pop() {
	if len(*stk) == 0 {
		return
	}
	*stk = (*stk)[0 : len(*stk)-1]
}

func (ta *TokenAssembler) Initialize(na ipld.NodeAssembler, budget int64) {
	if ta.stk == nil {
		ta.stk = make(assemblerStack, 0, 10)
	} else {
		ta.stk = ta.stk[0:0]
	}
	ta.stk.Push(na)
	ta.budget = budget
}

// Process takes a Token pointer as an argument.
// (Notice how this function happens to match the definition of the visitFn that's usable as an argument to TokenWalk.)
// The token argument can be understood to be "borrowed" for the duration of the Process call, but will not be mutated.
// The use of a pointer here is so that a single Token can be reused by multiple calls, avoiding unnecessary allocations.
//
// Note that Process does very little sanity checking of token sequences itself,
// mostly handing information to the NodeAssemblers directly,
// which presumably will reject the data if it is out of line.
// The NodeAssembler this TokenAssembler is wrapping should already be enforcing the relevant logical rules,
// so it is not useful for TokenAssembler.Process to attempt to duplicate those checks;
// TokenAssembler.Process will also return any errors from the NodeAssembler without attempting to enforce a pattern on those errors.
// In particular, TokenAssembler.Process does not check if every MapOpen is paired with a MapClose;
// it does not check if every ListOpen is paired with a ListClose;
// and it does not check if the token stream is continuing after all open recursives have been closed.
// TODO: review this documentation; more of these checks turn out necessary anyway than originally expected.
func (ta *TokenAssembler) Process(tk *Token) (err error) {
	if len(ta.stk) == 0 {
		return io.EOF
	}
	tip := ta.stk.Tip()
	switch tip.state {
	case 0:
		switch tk.Kind {
		case TokenKind_MapOpen:
			tip.ma, err = tip.na.BeginMap(tk.Length)
			tip.state = 2
			return err
		case TokenKind_MapClose:
			// Mostly we try to just forward things, but can't not check this one: tip.ma would be nil; there's reasonable target for forwarding.
			return ErrMalformedTokenSequence{"map close token encountered while not in the middle of a map"}
		case TokenKind_ListOpen:
			tip.la, err = tip.na.BeginList(tk.Length)
			tip.state = 1
			return err
		case TokenKind_ListClose:
			// Mostly we try to just forward things, but can't not check this one: tip.la would be nil; there's reasonable target for forwarding.
			return ErrMalformedTokenSequence{"list close token encountered while not in the middle of a list"}
		case TokenKind_Null:
			err = tip.na.AssignNull()
			ta.stk.Pop()
			return err
		case TokenKind_Bool:
			err = tip.na.AssignBool(tk.Bool)
			ta.stk.Pop()
			return err
		case TokenKind_Int:
			err = tip.na.AssignInt(int(tk.Int)) // TODO: upgrade all of ipld to use high precision int consistently
			ta.stk.Pop()
			return err
		case TokenKind_Float:
			err = tip.na.AssignFloat(tk.Float)
			ta.stk.Pop()
			return err
		case TokenKind_String:
			err = tip.na.AssignString(tk.Str)
			ta.stk.Pop()
			return err
		case TokenKind_Bytes:
			err = tip.na.AssignBytes(tk.Bytes)
			ta.stk.Pop()
			return err
		case TokenKind_Link:
			err = tip.na.AssignLink(tk.Link)
			ta.stk.Pop()
			return err
		default:
			panic(fmt.Errorf("unrecognized token kind (%q?)", tk.Kind))
		}
		return nil
	case 1:
		if tk.Kind == TokenKind_ListClose {
			err = tip.la.Finish()
			ta.stk.Pop()
			return err
		}
		ta.stk.Push(tip.la.AssembleValue())
		return ta.Process(tk)
	case 2:
		if tk.Kind == TokenKind_MapClose {
			err = tip.ma.Finish()
			ta.stk.Pop()
			return err
		}
		tip.state = 3
		ta.stk.Push(tip.ma.AssembleKey())
		return ta.Process(tk)
	case 3:
		tip.state = 2
		ta.stk.Push(tip.ma.AssembleValue())
		return ta.Process(tk)
	default:
		panic("unreachable")
	}
}
