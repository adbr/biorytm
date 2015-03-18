// 2015-02-25 Adam Bryt

package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	zeroMark  = "Z" // dzień zerowy biorytmu
	critMark  = "K" // dzień krytyczny
	maxMark   = "M" // dzień maksymalny
	minMark   = "m" // dzień minimalny
	dateMark  = "*" // aktualna data
	emptyMark = " " // zwykły dzień
)

// marks zwraca znaczniki jeśli dzień biorytmu jest wyróżniony, czyli
// jeśli jest dniem zerowym, lub dniem maksymum lub minimum biorytmu.
func marks(f, p, i int) (fmark, pmark, imark string) {
	switch f {
	case 0:
		fmark = zeroMark
	case 6:
		fmark = maxMark
	case 11, 12:
		fmark = critMark
	case 17:
		fmark = minMark
	default:
		fmark = emptyMark
	}

	switch p {
	case 0:
		pmark = zeroMark
	case 7:
		pmark = maxMark
	case 14:
		pmark = critMark
	case 21:
		pmark = minMark
	default:
		pmark = emptyMark
	}

	switch i {
	case 0:
		imark = zeroMark
	case 8:
		imark = maxMark
	case 16, 17:
		imark = critMark
	case 25:
		imark = minMark
	default:
		imark = emptyMark
	}

	return
}

// printBiorytm drukuje nd dni biorytmu na w.
func printBiorytm(w io.Writer, born, date time.Time, nd int) {
	fmt.Fprintf(w, "Data urodzenia: %s\n", born.Format(dateFmt))
	fmt.Fprintf(w, "Data biorytmu:  %s\n", date.Format(dateFmt))
	fmt.Fprintf(w, "Liczba dni:     %d\n", ndays(born, date))
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%-13s%-19s%-19s%-19s\n", "Data", "Fizyczny", "Psychiczny", "Intelektualny")

	const day = 24 * time.Hour

	d1 := date.Add(-time.Duration((nd-1)/2) * day) // data początku zakresu
	if d1.Before(born) {
		d1 = born
	}
	d2 := d1.Add(time.Duration(nd-1) * day) // data końca zakresu

	for {
		n := ndays(born, d1)
		f, p, i := bioDay(n)
		fv, pv, iv := bioVal(n)

		dmark := emptyMark
		if d1 == date {
			dmark = dateMark
		}
		fmark, pmark, imark := marks(f, p, i)

		fmt.Fprintf(w, "%s%s ", d1.Format(dateFmt), dmark)
		fmt.Fprintf(w, " F: %+5.2f (%2d/%d)%s ", fv, f, fPeriod, fmark)
		fmt.Fprintf(w, " P: %+5.2f (%2d/%d)%s ", pv, p, pPeriod, pmark)
		fmt.Fprintf(w, " I: %+5.2f (%2d/%d)%s \n", iv, i, iPeriod, imark)

		d1 = d1.Add(day)
		if d1.After(d2) {
			break
		}
	}
}

// biorytmCli drukuje biorytm na stdout.
func biorytmCli() {
	var (
		born time.Time // data urodzenia
		date time.Time // data docelowa biorytmu
		err  error
	)

	if *bornFlag == "" {
		fmt.Fprintln(os.Stderr, "brak opcji -born")
		usage()
	}

	born, err = time.Parse(dateFmt, *bornFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage()
	}

	if *dateFlag == "" {
		n := time.Now()
		date = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		date, err = time.Parse(dateFmt, *dateFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
		}
	}

	if date.Before(born) {
		fmt.Fprintln(os.Stderr, "data biorytmu wcześniejsza niż data urodzenia")
		usage()
	}

	if *daysFlag <= 0 {
		fmt.Fprintf(os.Stderr, "liczba dni musi być większa od 0, jest: %d\n", *daysFlag)
		usage()
	}

	printBiorytm(os.Stdout, born, date, *daysFlag)
}
