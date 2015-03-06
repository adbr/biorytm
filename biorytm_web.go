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

// template dla formatki biorytmu tekstowego
const formTextTmplStr = `
<html>
<head>
	<title>Biorytm</title>
</head>
<body>
	<p>
	Podaj datę urodzenia i datę biorytmu w formacie yyyy-mm-dd
	</p>

	<form action="/text/biorytm/">
		<table>
			<tr>
				<td style="text-align:right">data urodzenia:</td>
				<td><input type="text" name="born"></td>
			</tr>
			<tr>
				<td style="text-align:right">data aktualna:</td>
				<td><input type="text" name="date" value="{{ .DateString }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">ilość dni:</td>
				<td><input type="text" name="range" value="{{ .Days }}"></td>
			</tr>
			<tr>
				<td></td>
				<td style="text-align:right"><input type="submit"></td>
			</tr>
		</table>
	</form>
</body>
</html>
`

// template dla wyniku biorytmu tekstowego
const outputTextTmplStr = `
<html>
<head>
	<title>Biorytm</title>
	<!--
	<style>
		body {background-color:black; color:#FFC200}
	</style>
	-->
</head>
<body>
	<p>
	<pre>{{ . }}</pre>
	</p>
</body>
</html>
`

var (
	formTextTmpl   = template.New("formText")
	outputTextTmpl = template.New("outputText")
)

// parametry biorytmu
type biorytm struct {
	Born time.Time
	Date time.Time
	Days int
}

func (b biorytm) DateString() string {
	return b.Date.Format(dateFmt)
}

func biorytmWeb() {
	var err error

	formTextTmpl, err = formTextTmpl.Parse(formTextTmplStr)
	if err != nil {
		log.Fatalf("template 'formText' parse: %s", err)
	}

	outputTextTmpl, err = outputTextTmpl.Parse(outputTextTmplStr)
	if err != nil {
		log.Fatalf("template 'outputText' parse: %s", err)
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/text/", textHandler)
	http.HandleFunc("/text/biorytm/", textBiorytmHandler)
	//http.HandleFunc("/graph/", graphHandler)
	//http.HandleFunc("/graph/biorytm/", graphBiorytmHandler)

	log.Printf("biorytm: adres usługi HTTP: %s", *httpAddr)
	err = http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/text/", http.StatusFound)
}

// textHandler wyświetla formatkę dla danych wejściowych.
func textHandler(w http.ResponseWriter, r *http.Request) {
	// ustaw domyślne dane dla formatki
	b := biorytm{
		Date: time.Now(),
		Days: 30,
	}

	err := formTextTmpl.Execute(w, b)
	if err != nil {
		log.Printf("execute formText template: %s", err)
		fmt.Fprintf(w, "execute formText template: %s\n", err)
	}
	return
}

// textBiorytmHandler wyświetla biorytm w postaci tekstowej.
func textBiorytmHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getFormTextData(r)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	buf := new(bytes.Buffer)
	printBiorytm(buf, b.Born, b.Date, b.Days)

	err = outputTextTmpl.Execute(w, buf.String())
	if err != nil {
		log.Printf("execute outputText template: %s", err)
		fmt.Fprintf(w, "execute outputText template: %s\n", err)
	}
}

// getFormTextData pobiera z requestu parsuje parametry biorytmu.
func getFormTextData(r *http.Request) (biorytm, error) {
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
