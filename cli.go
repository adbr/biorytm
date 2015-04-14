// 2015-02-25 Adam Bryt

package main

import (
	"biorytm/cycle"
	"fmt"
	"io"
	"os"
	"time"
)

// Oznaczenia wyróżnionych dni na wydruku tekstowym.
const (
	zeroMark  = "Z" // dzień zerowy biorytmu
	critMark  = "K" // dzień krytyczny
	maxMark   = "M" // dzień maksymalny
	minMark   = "m" // dzień minimalny
	dateMark  = "*" // aktualna data
	emptyMark = " " // zwykły dzień
)

// marks zwraca znaczniki jeśli dzień biorytmu jest wyróżniony.
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

// printBiorytm drukuje na w biorytm dla days dni.
func printBiorytm(w io.Writer, born, date time.Time, days int) {
	fmt.Fprintf(w, "Data urodzenia: %s\n", born.Format(dateFmt))
	fmt.Fprintf(w, "Data biorytmu:  %s\n", date.Format(dateFmt))
	fmt.Fprintf(w, "Liczba dni:     %d\n", cycle.NumDays(born, date))
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%-13s%-19s%-19s%-19s\n", "Data", "Fizyczny", "Psychiczny", "Intelektualny")

	fvs := cycle.ValuesCenter(cycle.F, born, date, days)
	pvs := cycle.ValuesCenter(cycle.P, born, date, days)
	ivs := cycle.ValuesCenter(cycle.I, born, date, days)

	for i := 0; i < len(fvs); i++ {
		f := fvs[i]
		p := pvs[i]
		i := ivs[i]

		dmark := emptyMark
		if f.Date == date {
			dmark = dateMark
		}
		fmark, pmark, imark := marks(f.Day, p.Day, i.Day)

		fmt.Fprintf(w, "%s%s ", f.Date.Format(dateFmt), dmark)
		fmt.Fprintf(w, " F: %+5.2f (%2d/%d)%s ", f.Val, f.Day, cycle.F, fmark)
		fmt.Fprintf(w, " P: %+5.2f (%2d/%d)%s ", p.Val, p.Day, cycle.P, pmark)
		fmt.Fprintf(w, " I: %+5.2f (%2d/%d)%s \n", i.Val, i.Day, cycle.I, imark)
	}
}

func cliMain() {
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
