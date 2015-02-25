// 2015-02-24 Adam Bryt

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func biorytmHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintln(w, err)
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
		return
	}
	born, err = time.Parse(dateFmt, bornPar[0])
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// pobierz datę aktualną
	datePar, ok := r.Form["date"]
	if ok {
		date, err = time.Parse(dateFmt, datePar[0])
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	} else {
		n := time.Now()
		date = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
	}

	printBiorytm(w, born, date)
}

func biorytmWeb() {
	http.HandleFunc("/", biorytmHandler)

	log.Printf("adres usługi: %s", *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
