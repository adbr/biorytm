// 2015-02-24 Adam Bryt

package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	textFormTmpl    = template.New("textForm")
	textDisplayTmpl = template.New("textDisplay")
)

// biorytm grupuje parametry biorytmu.
type biorytm struct {
	Born time.Time
	Date time.Time
	Days int
}

func (b biorytm) DateString() string {
	return b.Date.Format(dateFmt)
}

func biorytmWeb() {
	initTemplates()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/text/form/", textFormHandler)
	http.HandleFunc("/text/display/", textDisplayHandler)
	//http.HandleFunc("/graph/form/", graphFormHandler)
	//http.HandleFunc("/graph/display/", graphDisplayHandler)

	log.Printf("biorytm: adres usługi HTTP: %s", *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func initTemplates() {
	var err error

	textFormTmpl, err = textFormTmpl.Parse(textFormHTML)
	if err != nil {
		log.Fatalf("błąd parsowania template 'textFormHTML': %s", err)
	}

	textDisplayTmpl, err = textDisplayTmpl.Parse(textDisplayHTML)
	if err != nil {
		log.Fatalf("błąd parsowania template 'textDisplayHTML': %s", err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/text/form/", http.StatusFound)
}

// textFormHandler wyświetla formatkę dla wprowadzania danych biorytmu
// w postaci tekstowej.
func textFormHandler(w http.ResponseWriter, r *http.Request) {
	// ustaw domyślne dane dla formatki
	b := biorytm{
		Date: time.Now(),
		Days: 30,
	}

	err := textFormTmpl.Execute(w, b)
	if err != nil {
		log.Printf("błąd wykonania template 'textFormTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'textFormTmpl': %s\n", err)
	}
}

// textDisplayHandler wyświetla biorytm w postaci tekstowej.
func textDisplayHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getTextFormData(r)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	buf := new(bytes.Buffer)
	printBiorytm(buf, b.Born, b.Date, b.Days)

	err = textDisplayTmpl.Execute(w, buf.String())
	if err != nil {
		log.Printf("błąd wykonania template 'textDisplayTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'textDisplayTmpl': %s\n", err)
	}
}

// getTextFormData pobiera z requestu i parsuje parametry biorytmu.
func getTextFormData(r *http.Request) (biorytm, error) {
	b := biorytm{}

	err := r.ParseForm()
	if err != nil {
		return b, err
	}

	// pobierz datę urodzenia

	bornPar, ok := r.Form["born"]
	if !ok {
		return b, errors.New("brak parametru 'born' (data urodzenia)")
	}

	b.Born, err = time.Parse(dateFmt, bornPar[0])
	if err != nil {
		return b, fmt.Errorf("błędna data urodzenia: %s", err)
	}

	// pobierz datę aktualną

	datePar, ok := r.Form["date"]
	if !ok {
		return b, errors.New("brak parametru 'date' (data aktualna)")
	}

	b.Date, err = time.Parse(dateFmt, datePar[0])
	if err != nil {
		return b, fmt.Errorf("błędna data aktualna: %s", err)
	}

	// pobierz parametr range

	rangePar, ok := r.Form["range"]
	if !ok {
		return b, errors.New("brak parametru 'range' (liczba dni biorytmu)")
	}

	n, err := strconv.Atoi(rangePar[0])
	if err != nil {
		return b, fmt.Errorf("błędny zakres dni: %s", err)
	}

	const maxRange = 1000
	if n > maxRange {
		return b, fmt.Errorf("za duża wartość zakresu (max: %d): %d", maxRange, n)
	}
	b.Days = n

	return b, nil
}
