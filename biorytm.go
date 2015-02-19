// 2015-02-16 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"
)

const (
	fPeriod = 23
	pPeriod = 28
	iPeriod = 33
)

const (
	dateFmt = "2006-01-02"
)

const usageStr = `usage: biorytm data_urodzenia [data_biorytmu]
	argumenty data_urodzenia i data_biorytmu mają format 'yyyy-mm-dd'`

var (
	rangeFlag = flag.Int("range", 15, "zakres dni")
)

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
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
	f = int(n % fPeriod)
	p = int(n % pPeriod)
	i = int(n % iPeriod)
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

func biorytm(btime, dtime time.Time) {
	fmt.Printf("data urodzenia: %s\n", btime.Format(dateFmt))
	fmt.Printf("data docelowa:  %s\n", dtime.Format(dateFmt))
	fmt.Printf("liczba dni:     %d\n", ndays(btime, dtime))

	day := 24 * time.Hour

	r := *rangeFlag / 2
	d := dtime.Add(-time.Duration(r) * day) // data początku zakresu
	n := ndays(btime, d)     // liczba dni od urodzenia do początku zakresu

	for i := 0; i < *rangeFlag; i++ {
		f, p, i := bioDay(n)
		fv, pv, iv := bioVal(n)

		fmt.Printf("%s", d.Format(dateFmt))
		fmt.Printf(" F: %+5.2f (%2d/%d)", fv, f, fPeriod)
		fmt.Printf(" P: %+5.2f (%2d/%d)", pv, p, pPeriod)
		fmt.Printf(" I: %+5.2f (%2d/%d)\n", iv, i, iPeriod)

		d = d.Add(day)
		n++
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var btime time.Time // data urodzenia
	var dtime time.Time // data docelowa biorytmu
	var err error

	switch flag.NArg() {
	case 1:
		btime, err = time.Parse(dateFmt, flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
		}
		n := time.Now()
		dtime = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
	case 2:
		btime, err = time.Parse(dateFmt, flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
		}
		dtime, err = time.Parse(dateFmt, flag.Arg(1))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
		}
	default:
		usage()
	}

	biorytm(btime, dtime)
}
