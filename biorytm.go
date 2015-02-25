// 2015-02-16 Adam Bryt

package main

import (
	"math"
	"time"
)

const (
	fPeriod = 23
	pPeriod = 28
	iPeriod = 33
)

// julianDayNumber returns the time's Julian Day Number
// relative to the epoch 12:00 January 1, 4713 BC, Monday.
// Algorithm: http://en.wikipedia.org/wiki/Julian_day
func julianDayNumber(year, month, day int) int64 {
	a := int64(14-month) / 12
	y := int64(year) + 4800 - a
	m := int64(month) + 12*a - 3
	return int64(day) + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

// ndays zwraca liczbę dni między datami t1 i t2.
func ndays(t1, t2 time.Time) int64 {
	a := julianDayNumber(t1.Year(), int(t1.Month()), t1.Day())
	b := julianDayNumber(t2.Year(), int(t2.Month()), t2.Day())
	n := b - a
	if n < 0 {
		n = -n
	}
	return n
}

// bioDay zwraca numer cyklu fizycznego, psychicznego i intelektualnego
// dla dnia n od urodzenia.
func bioDay(n int64) (f, p, i int) {
	f = int(n%fPeriod) + 1
	p = int(n%pPeriod) + 1
	i = int(n%iPeriod) + 1
	return
}

// bioVal zwraca wartość sinusoidy dla cyklu fizycznego, psychicznego
// i intelektualnego dla dnia n od urodzenia.
func bioVal(n int64) (fv, pv, iv float64) {
	f, p, i := bioDay(n)
	fv = math.Sin(float64(f) * 2 * math.Pi / fPeriod)
	pv = math.Sin(float64(p) * 2 * math.Pi / pPeriod)
	iv = math.Sin(float64(i) * 2 * math.Pi / iPeriod)
	return
}
