// 2015-02-25 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	zeroMark  = "Z" // zaznacza dzień zerowy biorytmu
	maxMark   = "M" // zaznacza dzień maksymalny lub minimalny
	dateMark  = "*" // aktualna data
	emptyMark = " " // dla zwykłego dnia
)

// marks zwraca znaczniki jeśli dzień biorytmu jest wyróżniony, czyli
// jeśli jest dniem zerowym, lub dniem maksymum lub minimum biorytmu.
func marks(f, p, i int) (fmark, pmark, imark string) {
	switch f {
	case 23, 11, 12:
		fmark = zeroMark
	case 6, 17:
		fmark = maxMark
	default:
		fmark = emptyMark
	}

	switch p {
	case 28, 14:
		pmark = zeroMark
	case 7, 21:
		pmark = maxMark
	default:
		pmark = emptyMark
	}

	switch i {
	case 33, 16, 17:
		imark = zeroMark
	case 8, 25:
		imark = maxMark
	default:
		imark = emptyMark
	}

	return
}

// printBiorytm drukuje nd dni biorytmu na w.
func printBiorytm(w io.Writer, btime, dtime time.Time, nd int) {
	fmt.Fprintf(w, "data urodzenia: %s\n", btime.Format(dateFmt))
	fmt.Fprintf(w, "data docelowa:  %s\n", dtime.Format(dateFmt))
	fmt.Fprintf(w, "liczba dni:     %d\n", ndays(btime, dtime))

	day := 24 * time.Hour

	r := nd / 2
	d := dtime.Add(-time.Duration(r) * day) // data początku zakresu
	n := ndays(btime, d)                    // ilość dni do początku zakresu

	for i := 0; i < nd; i++ {
		f, p, i := bioDay(n)
		fv, pv, iv := bioVal(n)

		dmark := emptyMark
		if d == dtime {
			dmark = dateMark
		}
		fmark, pmark, imark := marks(f, p, i)

		fmt.Fprintf(w, "%s%s ", d.Format(dateFmt), dmark)
		fmt.Fprintf(w, " F: %+5.2f (%2d/%d)%s ", fv, f, fPeriod, fmark)
		fmt.Fprintf(w, " P: %+5.2f (%2d/%d)%s ", pv, p, pPeriod, pmark)
		fmt.Fprintf(w, " I: %+5.2f (%2d/%d)%s \n", iv, i, iPeriod, imark)

		d = d.Add(day)
		n++
	}
}

// biorytmCli drukuje biorytm na stdout.
func biorytmCli() {
	var (
		btime time.Time // data urodzenia
		dtime time.Time // data docelowa biorytmu
		err   error
	)

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

	printBiorytm(os.Stdout, btime, dtime, *rangeFlag)
}
