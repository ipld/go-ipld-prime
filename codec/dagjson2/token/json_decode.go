package jsontoken

import (
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime/codec/codectools"
	"github.com/ipld/go-ipld-prime/codec/codectools/scratch"
)

type Decoder struct {
	r scratch.Reader

	phase decoderPhase   // current phase.
	stack []decoderPhase // stack of any phases that need to be popped back up to before we're done with a complete tree.
	some  bool           // true after first value in any context; use to decide if a comma must precede the next value.  (doesn't need a stack, because if you're popping, it's true again.)

	tok codectools.Token // we'll be yielding this repeatedly.

	DecoderConfig
}

type DecoderConfig struct {
	AllowDanglingComma  bool // normal json: false; strict: false.
	AllowWhitespace     bool // normal json: true;  strict: false.
	AllowEscapedUnicode bool // normal json: true;  strict: false.
	ParseUtf8C8         bool // normal json: false; dag-json: true.
}

func (d *Decoder) Init(r io.Reader) {
	d.r.Init(r)
	d.phase = decoderPhase_acceptValue
	d.stack = d.stack[0:0]
	d.some = false
}

func (d *Decoder) Step(budget *int) (next *codectools.Token, err error) {
	switch d.phase {
	case decoderPhase_acceptValue:
		err = d.step_acceptValue()
	case decoderPhase_acceptMapKeyOrEnd:
		err = d.step_acceptMapKeyOrEnd()
	case decoderPhase_acceptMapValue:
		err = d.step_acceptMapValue()
	case decoderPhase_acceptListValueOrEnd:
		err = d.step_acceptListValueOrEnd()
	}
	return &d.tok, err
}

func (d *Decoder) pushPhase(newPhase decoderPhase) {
	d.stack = append(d.stack, d.phase)
	d.phase = newPhase
	d.some = false
}

func (d *Decoder) popPhase() {
	d.phase = d.stack[len(d.stack)-1]
	d.stack = d.stack[:len(d.stack)-1]
	d.some = true
}

type decoderPhase uint8

const (
	decoderPhase_acceptValue decoderPhase = iota
	decoderPhase_acceptMapKeyOrEnd
	decoderPhase_acceptMapValue
	decoderPhase_acceptListValueOrEnd
)

func (d *Decoder) readn1skippingWhitespace() (majorByte byte, err error) {
	if d.DecoderConfig.AllowWhitespace {
		for {
			majorByte, err = d.r.Readn1()
			switch majorByte {
			case ' ', '\t', '\r', '\n': // continue
			default:
				return
			}
		}
	} else {
		for {
			majorByte, err = d.r.Readn1()
			switch majorByte {
			case ' ', '\t', '\r', '\n':
				return 0, fmt.Errorf("whitespace not allowed by decoder configured for strictness")
			default:
				return
			}
		}
	}
}

// The original step, where any value is accepted, and no terminators for recursives are valid.
// ONLY used in the original step; all other steps handle leaf nodes internally.
func (d *Decoder) step_acceptValue() error {
	majorByte, err := d.r.Readn1()
	if err != nil {
		return err
	}
	return d.stepHelper_acceptValue(majorByte)
}

// Step in midst of decoding a map, key expected up next, or end.
func (d *Decoder) step_acceptMapKeyOrEnd() error {
	majorByte, err := d.readn1skippingWhitespace()
	if err != nil {
		return err
	}
	if d.some {
		switch majorByte {
		case '}':
			d.tok.Kind = codectools.TokenKind_MapClose
			d.popPhase()
			return nil
		case ',':
			majorByte, err = d.readn1skippingWhitespace()
			if err != nil {
				return err
			}
			// and now fall through to the next switch
			// FIXME: AllowDanglingComma needs a check hereabouts
		}
	}
	switch majorByte {
	case '}':
		d.tok.Kind = codectools.TokenKind_MapClose
		d.popPhase()
		return nil
	default:
		// Consume a value for key.
		//  Given that this is JSON, this has to be a string.
		err := d.stepHelper_acceptValue(majorByte)
		if err != nil {
			return err
		}
		if d.tok.Kind != codectools.TokenKind_String {
			return fmt.Errorf("unexpected non-string token where expecting a map key")
		}
		// Now scan up to consume the colon as well, which is required next.
		majorByte, err = d.readn1skippingWhitespace()
		if err != nil {
			return err
		}
		if majorByte != ':' {
			return fmt.Errorf("expected colon after map key; got 0x%x", majorByte)
		}
		// Next up: expect a value.
		d.phase = decoderPhase_acceptMapValue
		d.some = true
		return nil
	}
}

// Step in midst of decoding a map, value expected up next.
func (d *Decoder) step_acceptMapValue() error {
	majorByte, err := d.readn1skippingWhitespace()
	if err != nil {
		return err
	}
	d.phase = decoderPhase_acceptMapKeyOrEnd
	return d.stepHelper_acceptValue(majorByte)
}

// Step in midst of decoding an array.
func (d *Decoder) step_acceptListValueOrEnd() error {
	majorByte, err := d.readn1skippingWhitespace()
	if err != nil {
		return err
	}
	if d.some {
		switch majorByte {
		case ']':
			d.tok.Kind = codectools.TokenKind_ListClose
			d.popPhase()
			return nil
		case ',':
			majorByte, err = d.readn1skippingWhitespace()
			if err != nil {
				return err
			}
			// and now fall through to the next switch
			// FIXME: AllowDanglingComma needs a check hereabouts
		}
	}
	switch majorByte {
	case ']':
		d.tok.Kind = codectools.TokenKind_ListClose
		d.popPhase()
		return nil
	default:
		d.some = true
		return d.stepHelper_acceptValue(majorByte)
	}
}

func (d *Decoder) stepHelper_acceptValue(majorByte byte) (err error) {
	switch majorByte {
	case '{':
		d.tok.Kind = codectools.TokenKind_MapOpen
		d.tok.Length = -1
		d.pushPhase(decoderPhase_acceptMapKeyOrEnd)
		return nil
	case '[':
		d.tok.Kind = codectools.TokenKind_ListOpen
		d.tok.Length = -1
		d.pushPhase(decoderPhase_acceptListValueOrEnd)
		return nil
	case 'n':
		d.r.Readnzc(3) // FIXME must check these equal "ull"!
		d.tok.Kind = codectools.TokenKind_Null
		return nil
	case '"':
		d.tok.Kind = codectools.TokenKind_String
		d.tok.Str, err = DecodeStringBody(&d.r)
		if err == nil {
			d.r.Readn1() // Swallow the trailing `"` (which DecodeStringBody has insured we have).
		}
		return err
	case 'f':
		d.r.Readnzc(4) // FIXME must check these equal "alse"!
		d.tok.Kind = codectools.TokenKind_Bool
		d.tok.Bool = false
		return nil
	case 't':
		d.r.Readnzc(3) // FIXME must check these equal "rue"!
		d.tok.Kind = codectools.TokenKind_Bool
		d.tok.Bool = true
		return nil
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Some kind of numeric... but in json, we can't tell if it's float or int.  At least, certainly not yet.
		// We'll have to look ahead quite a bit more to try to differentiate.  The decodeNumber function does this for us.
		d.r.Unreadn1()
		d.tok.Kind, d.tok.Int, d.tok.Float, err = DecodeNumber(&d.r)
		return err
	default:
		return fmt.Errorf("Invalid byte while expecting start of value: 0x%x", majorByte)
	}
}
