// 2015-02-16 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	dateFmt = "2006-01-02"
)

const usageStr = `usage: biorytm data_urodzenia [data_biorytmu]
	argumenty data_urodzenia i data_biorytmu mają format 'yyyy-mm-dd'`

var (
	rFlag = flag.Int("r", 14, "zakres dni")
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

// ndays zwraca liczbę dni między t1 i t2.
func ndays(t1, t2 time.Time) int {
	return 1
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

	fmt.Printf("data urodzenia: %s\n", btime)
	jdn := julianDayNumber(btime.Year(), int(btime.Month()), btime.Day())
	fmt.Printf("JDN: %v\n", jdn)

	fmt.Printf("data docelowa:  %s\n", dtime)
	jdn = julianDayNumber(dtime.Year(), int(dtime.Month()), dtime.Day())
	fmt.Printf("JDN: %v\n", jdn)
}
