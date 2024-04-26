// seehuhn.de/go/float - decimal representation of floating point numbers
// Copyright (C) 2024  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package float provides functions for the decimal representation of floating
// point numbers.
//
// Numbers formatted by this package consist of the following components in
// order, all optional: a sign, an integer part consisting of a sequence of
// decimal digits, and a fractional part consisting of a decimal point followed
// by a sequence of decimal digits.  At least one of the integer part or the
// fractional part must be present.
//
// The package guarantees that
//
//	Round(x, k) == Parse(Format(x, k), k)
//
// for all non-NaN values of x and all k ∈ {0, 1, ..., 10}.
package float

import (
	"errors"
	"io"
	"math"
	"strconv"
)

// Parse parses the decimal representation of a floating point number.  The
// function also accepts numbers in scientific notation.  The result is rounded
// to the nearest number with at most k digits after the decimal point.
func Parse(s string, digits int) (float64, error) {
	x, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if x != x {
		return 0, errors.New("NaN")
	}
	return Round(x, digits), nil
}

// Round rounds x to the nearest number with at most k digits after the decimal
// point.  Very large (positive or negative) numbers are clamped to the range
// of representable numbers.
func Round(x float64, digits int) float64 {
	// TODO(voss): I have found the value 15.42 by trial and error.  Make sure
	// that this is correct.
	limit := math.Pow(10, 15.42-float64(digits))
	if x > limit {
		x = limit
	} else if x < -limit {
		x = -limit
	}

	base := math.Pow(10, float64(digits))
	return math.Round(x*base) / base
}

// Format formats x as a decimal number with at most k digits after the decimal
// point.  Trailing zeros are omitted.  The output never uses scientific
// notation.
func Format(x float64, digits int) string {
	return string(doFormat(x, digits))
}

var zero = []byte{'0'}

func doFormat(x float64, digits int) []byte {
	z := int64(math.Round(x * math.Pow(10, float64(digits))))

	if z == 0 {
		return zero
	}

	var neg bool
	if z < 0 {
		neg = true
		z = -z
	}

	var cc []byte
	i := 0
	for z != 0 || i <= digits {
		if i == digits {
			cc = append(cc, '.')
		}
		cc = append(cc, byte(z%10)+'0')
		z /= 10
		i++
	}
	if i == digits+1 && cc[digits+1] == '0' {
		cc = cc[:digits+1]
	}

	firstNonZero := 0
	for cc[firstNonZero] == '0' {
		firstNonZero++
	}
	if firstNonZero == digits {
		firstNonZero++
	}
	cc = cc[firstNonZero:]

	if neg {
		cc = append(cc, '-')
	}

	for i, j := 0, len(cc)-1; i < j; i, j = i+1, j-1 {
		cc[i], cc[j] = cc[j], cc[i]
	}

	return cc
}

// Write is like [Format], but writes the result to w.
func Write(w io.Writer, x float64, digits int) error {
	_, err := w.Write(doFormat(x, digits))
	return err
}
