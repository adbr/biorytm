// 2015-03-27 Adam Bryt

// Pakiet cycle grupuje funkcje związane z obliczaniem wartości biorytmów.
package cycle

import (
	"math"
	"time"
)

// Period jest liczbą dni cyklu biorytmu.
type Period int

// Długości cykli biorytmów w dniach.
const (
	F Period = 23 // biorytm fizyczny
	P Period = 28 // biorytm psychiczny
	I Period = 33 // biorytm intelektualny
)

const day = 24 * time.Hour

// Day zwraca wartość cyklu p [-1..1] w dniu date od daty urodzenia born.
func Val(p Period, born, date time.Time) float64 {
	d := Day(p, born, date)
	return math.Sin(float64(d) * 2 * math.Pi / float64(p))
}

// Day zwraca numer dnia cyklu p w dniu date od daty urodzenia born.
func Day(p Period, born, date time.Time) int {
	n := NumDays(born, date)
	return int(n % int64(p))
}

// NumDays zwraca liczbę dni między datami d1 i d2.
func NumDays(d1, d2 time.Time) int64 {
	a := julianDayNumber(d1.Year(), int(d1.Month()), d1.Day())
	b := julianDayNumber(d2.Year(), int(d2.Month()), d2.Day())
	n := b - a
	if n < 0 {
		n = -n
	}
	return n
}

// julianDayNumber returns the time's Julian Day Number
// relative to the epoch 12:00 January 1, 4713 BC, Monday.
// Algorithm: http://en.wikipedia.org/wiki/Julian_day
func julianDayNumber(year, month, day int) int64 {
	a := int64(14-month) / 12
	y := int64(year) + 4800 - a
	m := int64(month) + 12*a - 3
	return int64(day) + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

// Point reprezentuje wartość biorytmu [-1..1] dla danego dnia.
type Point struct {
	Date time.Time // data
	Val  float64   // wartość biorytmu w zakresie [-1..1]
}

// ValuesCenterDate zwraca slice wartości biorytmu dla zakresu days dni, gdzie
// date jest w środku zakresu dni.  Argument p jest długością cyklu biorytmu, a
// born to data urodzenia.
func ValuesCenterDate(p Period, born, date time.Time, days int) []Point {
	d1 := date.Add(-time.Duration(days/2) * day)
	if d1.Before(born) {
		d1 = born
	}
	return Values(p, born, d1, days)
}

// Values zwraca slice wartości biorytmu dla zakresu days dni od date.
// Argument p jest długością cyklu biorytmu, a born to data urodzenia.
func Values(p Period, born, date time.Time, days int) []Point {
	d1 := date
	d2 := d1.Add(time.Duration(days) * day)
	a := []Point{}
	for {
		v := Val(p, born, d1)
		pv := Point{
			Date: d1,
			Val:  v,
		}
		a = append(a, pv)

		d1 = d1.Add(day)
		if d1.After(d2) {
			break
		}
	}
	return a
}
