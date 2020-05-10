/*
	"jst" -- JSON Table -- is a format that's parsable as JSON, while sprucing up the display to humans using the non-significant whitespace cleverly.
	Regular data can be piped into a JSON Table formatter, and some simple heuristics will attempt to detect table structure.

	Lists can be turned into column-aligned tables:

	  - Every time there's a list, and the first entry is a map, we'll try to treat it as a table.
	  - Every time a map in that list starts with the same first key as the first map did, it's a table row.
	  - Every thing that's a table row will be buffered, and we try to fit each key into a table column.
	    - (FUTURE) You can manually specify keys that should be excluded from columns; these will be shifted to the end and packed tightly.
	  - We'll store the size of the widest value for each column.  We'll need to do this over every row so we can align output spacing.
	  - If there's a map in the list that doesn't start with the same first key, it's a table exclusion, and gets formatted densely.
	  - If a map has a value that's a list, we attempt to apply this whole ruleset over again from the top.
	    - If a table is detected, we'll print the key on its own new line, slightly indented, together with the list open.  Then, emit the table onsubsequent further indented lines.

	Maps can also be turned into column-aligned tables:

	  - You have to use additional configuration to engage this: by default, only lists trigger table mode.
	  - The search for table rows begins anew with each map value.  The map keys form a defacto first column.
	  - Thereafter, all the rules for handling each table row is the same the rules described above for lists.

	Decorations can be applied:

	  - Defaults include colorizing map keys vs values, and optionally colorizing column names distinctly from other keys.
	  - These colorations operate by ANSI codes (e.g., they work in terminals).  The palette is accordingly limited.
	  - (FUTURE) You can override colors of specific keys and values.

	The overall nature of detecting traits of the data (particularly, size) means JSON Tables cannot be created streamingly;
	we have to process the entire structure first, and only then can we begin to output correctly aligned data.
	(It's for this reason that this is implemented over IPLD Nodes, and would be somewhat painful to implement using just refmt Tokens -- we'd end up buffering *all* the tokens anyway, and wanting to build an index, etc.)

	There's no unmarshal functions because unmarshal is just... regular JSON unmarshal.
*/
package jst

import (
	"bytes"
	"io"

	"github.com/polydawn/refmt/json"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func Marshal(n ipld.Node, w io.Writer) error {
	ctx := state{
		cfg: Config{
			Indent: []byte{' ', ' '},
		},
	}
	// Stride first -- see how much spacing we need.
	err := stride(&ctx, n)
	if err != nil {
		return err
	}
	// Marshal -- using the spacing nodes from our stride.
	return marshal(&ctx, n, w)
}

type state struct {
	cfg    Config
	path   []ipld.PathSegment // TODO replace with PathBuffer... once you, you know, write it.
	stack  []step
	tables map[tableGroupID]*table
	spare  bytes.Buffer
	indent int
}

type step struct {
}

type tableGroupID string
type columnName string

type table struct {
	entryStyles map[columnName]entryStyle
	keySize     map[columnName]int // size of columnName rendered
	cols        []columnName
	colSize     map[columnName]int // max rendered value width
	ownLine     []columnName
}

type entryStyle uint8

const (
	entryStyle_column   entryStyle = 1
	entryStyle_trailing entryStyle = 2
	entryStyle_ownLine  entryStyle = 3
)

type Config struct {
	Indent []byte

	// FUTURE: selectors and other forms of specification can override where tables appear, what their tableGroupID is, and so on.
	// FUTURE: whether to emit trailing commas unconditionally, even on the last elements of maps and lists.
	// FUTURE: list of keys to exclude from column alignment efforts.
	// FUTURE: fixed column widths (would even potentially enable streaming operation!  (probably won't on first draft though; makes many codepaths diverge)).
	// FUTURE: additional coloration cues (could be from selectors, or take cues from schema types).
	// ..... etc ......
}

func (ctx *state) Table(tgid tableGroupID) *table {
	if tab, exists := ctx.tables[tgid]; exists {
		return tab
	}
	tab := &table{
		entryStyles: make(map[columnName]entryStyle),
		keySize:     make(map[columnName]int),
		colSize:     make(map[columnName]int),
	}
	if ctx.tables == nil {
		ctx.tables = make(map[tableGroupID]*table)
	}
	ctx.tables[tgid] = tab
	return tab
}

func (tab *table) ColumnObserved(cn columnName, size int) {
	switch tab.entryStyles[cn] {
	case entryStyle_trailing, entryStyle_ownLine:
		return
	}
	prev, exists := tab.colSize[cn]
	if !exists {
		tab.cols = append(tab.cols, cn)
		tab.entryStyles[cn] = entryStyle_column
		prev = -1
	}
	tab.colSize[cn] = max(size, prev)
}
func (tab *table) GetsOwnLine(cn columnName) {
	switch tab.entryStyles[cn] {
	case entryStyle_ownLine:
		return
	}
	tab.entryStyles[cn] = entryStyle_ownLine
	tab.ownLine = append(tab.ownLine, cn)
}
func (tab *table) Finalize(cn columnName) {
	// TODO drop all entries in tab.cols that ended up with a different entrystyle.
	//  (This happens when something gets observed as a column first, but forced into ownLine mode by a subtable in a subsequent row.)
}

func (tab *table) IsRow(n ipld.Node) bool {
	// FUTURE: this entire function's behavior might be *heavily* redirected by config.
	//  Having it attached to the table struct might not even be sensible by the end of the day.
	//  Alternately: unclear if we need this function at all, if the "trailing" and "ownLine" entryStyle can simply carry the day for all userstories like comments and etc.
	switch n.ReprKind() {
	case ipld.ReprKind_Map:
		if n.Length() < 1 {
			return false
		}
		ks := mustFirstKeyAsString(n)
		if len(tab.cols) < 1 {
			// FUTURE: may want to check for `ks == "comment"` or other configured values, and then say "no".
			return true
		}
		return columnName(ks) == tab.cols[0]
	case ipld.ReprKind_List:
		// FUTURE: maybe this could be 'true', but it requires very different logic.  And it's not in my first-draft pareto priority choices.
		return false
	default:
		return false
	}
}

// Discerns if the given node starts a table, and if so, what string to use as tableGroupID for cross-table alignment.
// By default, the tableGroupID is just the first key in the first row.
// (We might extend the tableGroupID heuristic later, perhaps to also include the last key we saw on the path here, etc).
func peekMightBeTable(ctx *state, n ipld.Node) (bool, tableGroupID) {
	// FUTURE: might need to apply a selector or other rules from ctx.cfg to say additonal "no"s.
	// FUTURE: the ctx.cfg can also override what the tableGroupID is.
	switch n.ReprKind() {
	case ipld.ReprKind_Map:
		// TODO: maps can definitely be tables!  but gonna come back to this.  and by default, they're not.
		return false, ""
	case ipld.ReprKind_List:
		if n.Length() < 1 {
			return false, ""
		}
		n0, _ := n.LookupIndex(0)
		if n0.ReprKind() != ipld.ReprKind_Map {
			return false, ""
		}
		if n0.Length() < 1 {
			return false, ""
		}
		return true, tableGroupID(mustFirstKeyAsString(n0))
	default:
		return false, ""
	}
}

// Fills out state.tables.
// Applies itself and other stride functions recursively (have to, because:
//  some row values might themselves be tables,
//   which removes them from the column flow and changes our size counting).
func stride(ctx *state, n ipld.Node) error {
	switch n.ReprKind() {
	case ipld.ReprKind_Map:
		panic("todo")
	case ipld.ReprKind_List:
		return strideList(ctx, n)
	default:
		// There's never anything we need to record for scalars that their parents don't already note.
		return nil
	}
}

func strideList(ctx *state, listNode ipld.Node) error {
	isTable, tgid := peekMightBeTable(ctx, listNode)
	if !isTable {
		return nil
	}
	tab := ctx.Table(tgid)
	listItr := listNode.ListIterator()
	for !listItr.Done() {
		_, row, err := listItr.Next()
		// TODO grow ctx.path
		if err != nil {
			return recordErrorPosition(ctx, err)
		}
		if !tab.IsRow(row) {
			continue
		}
		rowItr := row.MapIterator()
		for !rowItr.Done() {
			k, v, err := rowItr.Next()
			// TODO grow ctx.path
			if err != nil {
				return recordErrorPosition(ctx, err)
			}
			ks, _ := k.AsString()
			if vIsTable, _ := peekMightBeTable(ctx, v); vIsTable {
				tab.GetsOwnLine(columnName(ks))
				stride(ctx, v) // i do believe this results in calling peekMightBeTable repeatedly; could refactor to improve; but doesn't affect correctness.
			} else {
				if tab.entryStyles[columnName(ks)] != entryStyle_trailing {
					ctx.spare.Reset()
					if err := marshalPlain(ctx, v, &ctx.spare); err != nil {
						return err
					}
					computedSize := ctx.spare.Len() // FIXME this is ignoring charsets, renderable glyphs, etc at present.
					tab.ColumnObserved(columnName(ks), computedSize)
				}
			}
		}
	}
	return nil
}

func marshal(ctx *state, n ipld.Node, w io.Writer) error {
	switch n.ReprKind() {
	case ipld.ReprKind_Map:
		panic("todo")
	case ipld.ReprKind_List:
		return marshalList(ctx, n, w)
	default:
		return marshalPlain(ctx, n, w)
	}
}

// this function is probably a placeholder.
// It doesn't colorize or anything else.  To replace it with something clever that does,
// we'll have to tear deeper into the plumbing level of json serializers; will, but later.
func marshalPlain(ctx *state, n ipld.Node, w io.Writer) error {
	err := dagjson.Marshal(n, json.NewEncoder(w, json.EncodeOptions{
		// never indent here: these values will always end up being emitted mid-line.
	}))
	if err != nil {
		return recordErrorPosition(ctx, err)
	}
	return nil
}

func marshalList(ctx *state, listNode ipld.Node, w io.Writer) error {
	isTab, tgid := peekMightBeTable(ctx, listNode)
	if !isTab {
		return marshalPlain(ctx, listNode, w)
	}
	tab := ctx.Table(tgid)
	ctx.indent++
	w.Write([]byte{'[', '\n'})
	listItr := listNode.ListIterator()
	for !listItr.Done() {
		_, row, err := listItr.Next()
		// TODO grow ctx.path
		if err != nil {
			return recordErrorPosition(ctx, err)
		}
		if err := marshalListValue(ctx, tab, row, w); err != nil {
			return err
		}
	}
	ctx.indent--
	w.Write(bytes.Repeat(ctx.cfg.Indent, ctx.indent))
	w.Write([]byte{']'}) // TODO: linebreak?  not sure.  probably not: if it's a sub-table, a close-brace might be coming.
	return nil
}
func marshalListValue(ctx *state, tab *table, row ipld.Node, w io.Writer) error {
	if row.ReprKind() != ipld.ReprKind_Map {
		// TODO make this a lot more open... scalars aren't exactly "rows" for example but we can surely print them just fine.
		panic("table rows can only be maps at present")
	}
	w.Write(bytes.Repeat(ctx.cfg.Indent, ctx.indent))
	w.Write([]byte{'{'})

	// Flow here goes by the table notes rather than the data!  Mostly.
	//  0. Figure out what the last column is that we have a value for.
	//      We end lines early rather than adding tons of extraneous spaces if we can.
	//      Skip this if the row has any elements that are in the trailing style.
	//       FIXME figure out if there's a cheaper way than iterating to sort that out; it's kinda gross, considering then we iterate again in stage 2.
	//  1. First all the columns are emitted, in order.
	//      If there is no entry for that column, and there's more stuff coming, we have to emit space for both the column key and the value.
	//  2. Next all the trailing entries are emitted.
	//  3. Finally we linebreak, indent further, and start dealing with ownLine stuff (which may include sub-tables).

	// Stage 0 -- looking ahead for where we can rest.
	for _, _ = range tab.cols {
		// TODO figure out who's the last one with presence.
	}

	// Stage 1 -- emitting regular columns.
	for _, col := range tab.cols {
		v, err := row.LookupString(string(col))
		if v == nil {
			w.Write(bytes.Repeat([]byte{' '}, tab.keySize[col]+1 /* plus one for the colon */))
			w.Write(bytes.Repeat([]byte{' '}, tab.colSize[col]))
			continue
		}
		if err != nil {
			return recordErrorPosition(ctx, err)
		}
		if err := emitKey(ctx, basicnode.NewString(string(col)), w); err != nil { // FIXME this would be a lot less irritating if we had more plumbing access to the json encoding
			return err
		}
		// Rather hackily, marshal to an intermediate buffer, so we can count print size.  Would rather stream, but needs more code to do so.
		ctx.spare.Reset()
		if err := marshalPlain(ctx, v, &ctx.spare); err != nil {
			return err
		}
		computedSize := ctx.spare.Len() // FIXME this is ignoring charsets, renderable glyphs, etc at present.
		w.Write(ctx.spare.Bytes())
		w.Write([]byte{','}) // FIXME: you know, the maybe-trailing shenanigans.
		w.Write(bytes.Repeat([]byte{' '}, tab.colSize[col]-computedSize+1))
	}

	// Stage 2 -- emitting trailing entries.
	rowItr := row.MapIterator()
	for !rowItr.Done() {
		k, v, err := rowItr.Next()
		// TODO grow ctx.path
		if err != nil {
			return recordErrorPosition(ctx, err)
		}
		ks, _ := k.AsString()
		switch tab.entryStyles[columnName(ks)] {
		case entryStyle_column, entryStyle_ownLine:
			continue
		}
		if err := emitKey(ctx, k, w); err != nil {
			return err
		}
		if err := marshal(ctx, v, w); err != nil {
			return err
		}
		w.Write([]byte{','}) // FIXME: you know, the maybe-trailing shenanigans.
	}

	// Stage 3 -- emitting ownLine entries.
	if len(tab.ownLine) > 0 { // FIXME: wrong condition; this should be based on if we *have* any (needs stage0 precount).
		w.Write([]byte{'\n'})
		ctx.indent++
		w.Write(bytes.Repeat(ctx.cfg.Indent, ctx.indent))
		defer func() { ctx.indent-- }()
	}
	for _, col := range tab.ownLine {
		v, err := row.LookupString(string(col))
		if v == nil {
			continue
		}
		if err := emitKey(ctx, basicnode.NewString(string(col)), w); err != nil { // FIXME this would be a lot less irritating if we had more plumbing access to the json encoding
			return err
		}
		if err != nil {
			return recordErrorPosition(ctx, err)
		}
		if err := marshal(ctx, v, w); err != nil { // whole recursion.  can even have sub-tables.
			return err
		}
		w.Write([]byte{','}) // FIXME: you know, the maybe-trailing shenanigans.
	}

	// End of entries.  Close the row.
	w.Write([]byte{'}'})
	w.Write([]byte{','}) // FIXME: you know, the maybe-trailing shenanigans.
	w.Write([]byte{'\n'})
	return nil
}

func emitKey(ctx *state, k ipld.Node, w io.Writer) error {
	// FUTURE: ansi color goes around here too (hence the need to have ctx in hand).
	if err := dagjson.Marshal(k, json.NewEncoder(w, json.EncodeOptions{})); err != nil {
		return recordErrorPosition(ctx, err)
	}
	w.Write([]byte{':'})
	w.Write([]byte{' '}) // FUTURE: this should be configurable
	return nil
}
