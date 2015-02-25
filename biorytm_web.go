// 2015-02-24 Adam Bryt

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const webUsageStr = `usage: URL/?born=<date>&date=<date>&range=<ndays>
	parametry:
		born: data urodzenia
		date: data aktualna (jeśli nie występuje, to time.Now())
		range: liczba dni (domyślnie 15)
	example:
		localhost:5050/?born=1970-01-02&date=2015-01-01&range=20`

func biorytmHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln(w, err)
		fmt.Fprintln(w, webUsageStr)
		return
	}

	var (
		born time.Time // data urodzenia
		date time.Time // data aktualna
	)

	// pobierz datę urodzenia
	bornPar, ok := r.Form["born"]
	if !ok {
		fmt.Fprintln(w, "brak parametru 'born' (data urodzenia)")
		fmt.Fprintln(w, webUsageStr)
		return
	}
	born, err = time.Parse(dateFmt, bornPar[0])
	if err != nil {
		fmt.Fprintln(w, err)
		fmt.Fprintln(w, webUsageStr)
		return
	}

	// pobierz datę aktualną
	datePar, ok := r.Form["date"]
	if ok {
		date, err = time.Parse(dateFmt, datePar[0])
		if err != nil {
			fmt.Fprintln(w, err)
			fmt.Fprintln(w, webUsageStr)
			return
		}
	} else {
		n := time.Now()
		date = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
	}

	// pobierz parametr range
	rangePar, ok := r.Form["range"]
	if ok {
		n, err := strconv.Atoi(rangePar[0])
		if err != nil {
			fmt.Fprintln(w, err)
			fmt.Fprintln(w, webUsageStr)
			return
		}
		*rangeFlag = n
	}

	printBiorytm(w, born, date, *rangeFlag)
}

func biorytmWeb() {
	http.HandleFunc("/", biorytmHandler)

	log.Printf("adres usługi: %s", *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
