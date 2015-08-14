// 2015-02-25 Adam Bryt

package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/adbr/biorytm/cycle"
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

// fmarks zawiera znaczniki wyróżnionych dni cyklu fizycznego.
var fmarks = map[int]string{
	0:  zeroMark,
	6:  maxMark,
	11: critMark,
	12: critMark,
	17: minMark,
}

// fmark zwraca znacznik dla dnia day cyklu fizycznego.
func fmark(day int) string {
	s, ok := fmarks[day]
	if ok {
		return s
	}
	return emptyMark
}

// pmarks zawiera znaczniki wyróżnionych dni cyklu psychicznego.
var pmarks = map[int]string{
	0:  zeroMark,
	7:  maxMark,
	14: critMark,
	21: minMark,
}

// pmark zwraca znacznik dla dnia day cyklu psychicznego.
func pmark(day int) string {
	s, ok := pmarks[day]
	if ok {
		return s
	}
	return emptyMark
}

// imarks zawiera znaczniki wyróżnionych dni cyklu intelektualnego.
var imarks = map[int]string{
	0:  zeroMark,
	8:  maxMark,
	16: critMark,
	17: critMark,
	25: minMark,
}

// imark zwraca znacznik dla dnia day cyklu intelektualnego.
func imark(day int) string {
	s, ok := imarks[day]
	if ok {
		return s
	}
	return emptyMark
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

		dm := emptyMark
		if f.Date == date {
			dm = dateMark
		}
		fm := fmark(f.Day)
		pm := pmark(p.Day)
		im := imark(i.Day)

		fmt.Fprintf(w, "%s%s ", f.Date.Format(dateFmt), dm)
		fmt.Fprintf(w, " F: %+5.2f (%2d/%d)%s ", f.Val, f.Day, cycle.F, fm)
		fmt.Fprintf(w, " P: %+5.2f (%2d/%d)%s ", p.Val, p.Day, cycle.P, pm)
		fmt.Fprintf(w, " I: %+5.2f (%2d/%d)%s ", i.Val, i.Day, cycle.I, im)
		fmt.Fprintln(w)
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
