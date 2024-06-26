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

package float

import (
	"math"
	"strconv"
	"testing"
)

func TestFormat(t *testing.T) {
	type testCase struct {
		in     float64
		digits int
		out    string
	}
	testCases := []testCase{
		{0, 0, "0"},
		{1, 0, "1"},
		{-1, 0, "-1"},
		{0, 1, "0"},
		{1, 1, "1"},
		{1, 5, "1"},
		{-1, 1, "-1"},
		{0.1, 0, "0"},
		{0.1, 1, ".1"},
		{0.1, 2, ".1"},
		{0.9, 0, "1"},
		{0.9, 1, ".9"},
		{0.9, 2, ".9"},
		{0.19, 1, ".2"},
		{-0.19, 1, "-.2"},
		{math.Pi, 0, "3"},
		{math.Pi, 1, "3.1"},
		{math.Pi, 2, "3.14"},
		{math.Pi, 4, "3.1416"},
		{math.Pi, 5, "3.14159"},
	}
	for _, tc := range testCases {
		got := Format(tc.in, tc.digits)
		if got != tc.out {
			t.Errorf("Format(%g, %d) = %q, want %q", tc.in, tc.digits, got, tc.out)
		}
	}
}

// FuzzRound verifies that Round(x) = Format(Parse(Round(x))).
func FuzzRound(f *testing.F) {
	f.Fuzz(func(t *testing.T, x float64, digits int) {
		if digits < 0 || digits > 10 {
			t.Skip()
		}

		x = Round(x, digits)
		s := Format(x, digits)
		y, err := Parse(s, digits)

		if err != nil {
			t.Fatalf("%s -> %q -> %v [%d]",
				strconv.FormatFloat(x, 'f', -1, 64), s, err, digits)
		}
		if x != y {
			t.Errorf("%s -> %q -> %s [%d]",
				strconv.FormatFloat(x, 'f', -1, 64), s, strconv.FormatFloat(y, 'f', -1, 64), digits)
		}
	})
}
