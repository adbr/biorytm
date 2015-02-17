// 2015-02-17 Adam Bryt

package main

import "testing"

func TestJulianDayNumber(t *testing.T) {
	tests := []struct {
		y, m, d int   // year, month, day
		jd      int64 // julian day number
	}{
		{
			0, 1, 1,
			1721060,
		},
		{
			500, 6, 5,
			1903837,
		},
		{
			2000, 5, 2,
			2451667,
		},
		{
			4321, 1, 2,
			3299274,
		},
	}

	for i, test := range tests {
		jd := julianDayNumber(test.y, test.m, test.d)
		if jd != test.jd {
			t.Errorf("#%d: julianDayNumber(): oczekiwano: %d, jest: %d", i, test.jd, jd)
		}
	}
}
