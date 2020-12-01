package jsontoken

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ipld/go-ipld-prime/codec/codectools"
	"github.com/ipld/go-ipld-prime/codec/codectools/scratch"
)

// License note: the string and numeric parsers here borrow
// heavily from the golang stdlib json parser scanner.
// That code is originally Copyright 2010 The Go Authors,
// and is governed by a BSD-style license.

// DecodeNumber will attempt to decode data in the format of a JSON numer from the reader.
// JSON is somewhat ambiguous about numbers: we'll return an int if we can, and a float if there's any decimal point involved.
// The boolean return indicates which kind of number we have:
// if true, we have an int (and the float return is invalid);
// if false, we have a float (and the int return is invalid).
func DecodeNumber(r *scratch.Reader) (codectools.TokenKind, int64, float64, error) {
	r.Track()
	// Scan until scanner tells us end of numeric.
	// Pick the first scanner stepfunc based on the leading byte.
	majorByte, err := r.Readn1()
	if err != nil {
		return codectools.TokenKind_Null, 0, 0, err
	}
	var step numscanStep
	switch majorByte {
	case '-':
		step = numscan_neg
	case '0':
		step = numscan_0
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		step = numscan_1
	default:
		panic("unreachable") // FIXME not anymore it ain't, this is exported
	}
	for {
		b, err := r.Readn1()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, 0, err
		}
		step, err = step(b)
		if step == nil {
			// Unread one.  The scan loop consumed one char beyond the end (this is unavoidable in json!),
			//  and that might be part of what whatever is going to be decoded from this stream next.
			r.Unreadn1()
			break
		}
		if err != nil {
			return 0, 0, 0, err
		}
	}
	// Parse!
	// *This is not a fast parse*.
	// Try int first; if it fails try float; if that fails return the float error.
	s := string(r.StopTrack())
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return codectools.TokenKind_Int, i, 0, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	return codectools.TokenKind_Float, 0, f, err
}

// Scan steps are looped over the stream to find how long the number is.
// A nil step func is returned to indicate the string is done.
// Actually parsing the string is done by 'parseString()'.
type numscanStep func(c byte) (numscanStep, error)

// numscan_neg is the state after reading `-` during a number.
func numscan_neg(c byte) (numscanStep, error) {
	if c == '0' {
		return numscan_0, nil
	}
	if '1' <= c && c <= '9' {
		return numscan_1, nil
	}
	return nil, fmt.Errorf("invalid byte in numeric literal: 0x%x", c)
}

// numscan_1 is the state after reading a non-zero integer during a number,
// such as after reading `1` or `100` but not `0`.
func numscan_1(c byte) (numscanStep, error) {
	if '0' <= c && c <= '9' {
		return numscan_1, nil
	}
	return numscan_0(c)
}

// numscan_0 is the state after reading `0` during a number.
func numscan_0(c byte) (numscanStep, error) {
	if c == '.' {
		return numscan_dot, nil
	}
	if c == 'e' || c == 'E' {
		return numscan_e, nil
	}
	return nil, nil
}

// numscan_dot is the state after reading the integer and decimal point in a number,
// such as after reading `1.`.
func numscan_dot(c byte) (numscanStep, error) {
	if '0' <= c && c <= '9' {
		return numscan_dot0, nil
	}
	return nil, fmt.Errorf("invalid byte after decimal in numeric literal: 0x%x", c)
}

// numscan_dot0 is the state after reading the integer, decimal point, and subsequent
// digits of a number, such as after reading `3.14`.
func numscan_dot0(c byte) (numscanStep, error) {
	if '0' <= c && c <= '9' {
		return numscan_dot0, nil
	}
	if c == 'e' || c == 'E' {
		return numscan_e, nil
	}
	return nil, nil
}

// numscan_e is the state after reading the mantissa and e in a number,
// such as after reading `314e` or `0.314e`.
func numscan_e(c byte) (numscanStep, error) {
	if c == '+' || c == '-' {
		return numscan_eSign, nil
	}
	return numscan_eSign(c)
}

// numscan_eSign is the state after reading the mantissa, e, and sign in a number,
// such as after reading `314e-` or `0.314e+`.
func numscan_eSign(c byte) (numscanStep, error) {
	if '0' <= c && c <= '9' {
		return numscan_e0, nil
	}
	return nil, fmt.Errorf("invalid byte in exponent of numeric literal: 0x%x", c)
}

// numscan_e0 is the state after reading the mantissa, e, optional sign,
// and at least one digit of the exponent in a number,
// such as after reading `314e-2` or `0.314e+1` or `3.14e0`.
func numscan_e0(c byte) (numscanStep, error) {
	if '0' <= c && c <= '9' {
		return numscan_e0, nil
	}
	return nil, nil
}
