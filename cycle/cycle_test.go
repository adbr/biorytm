// 2015-02-17 Adam Bryt

package cycle

import (
	"testing"
	"time"
)

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

func TestNumDays(t *testing.T) {
	//layout := "2006-01-02"

	tests := []struct {
		t1 time.Time
		t2 time.Time
		n  int64 // liczba dni między t1 i t2
	}{
		{
			time.Now(),
			time.Now(),
			0,
		},
		{
			t1: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			t2: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			n:  0,
		},
		{
			t1: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			t2: time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC),
			n:  1,
		},
		{
			t1: time.Date(2015, 1, 5, 0, 0, 0, 0, time.UTC),
			t2: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			n:  4,
		},
	}

	for i, test := range tests {
		n := NumDays(test.t1, test.t2)
		if n != test.n {
			t.Errorf("#%d: NumDays(): oczekiwano: %d, jest: %d", i, test.n, n)
		}
	}
}
