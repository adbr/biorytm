// 2015-02-24 Adam Bryt

package main

import (
	"bytes"
	"fmt"
	"html/template"
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

// template dla formatki
const formTmplStr = `
<html>
<head>
	<title>Biorytm</title>
</head>
<body>
	<p>
	Podaj datę urodzenia i datę biorytmu w formacie yyyy-mm-dd
	</p>

	<form action="/">
		<table>
			<tr>
				<td style="text-align:right">data urodzenia:</td>
				<td><input type="text" name="born"></td>
			</tr>
			<tr>
				<td style="text-align:right">data aktualna:</td>
				<td><input type="text" name="date" value="{{ . }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">ilość dni:</td>
				<td><input type="text" name="range" value="30"></td>
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

// template dla wyniku
const outputTmplStr = `
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
	formTmpl   *template.Template
	outputTmpl *template.Template
)

func biorytmHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm(): %s\n", err)
		return
	}

	// sprawdź czy podano parametr 'born'
	_, ok := r.Form["born"]
	if !ok {
		// nie ma parametru 'born': wyświetl formatkę
		now := time.Now().Format(dateFmt)
		err := formTmpl.Execute(w, now)
		if err != nil {
			log.Printf("execute form template: %s", err)
		}
		return
	}

	// jest parametr 'born': parsuj argumenty i wyświetl wynik

	var (
		born time.Time // data urodzenia
		date time.Time // data aktualna
	)

	// pobierz datę urodzenia
	bornPar, ok := r.Form["born"]
	born, err = time.Parse(dateFmt, bornPar[0])
	if err != nil {
		fmt.Fprintf(w, "błędna data urodzenia: %s\n", err)
		return
	}

	// pobierz datę aktualną
	datePar, ok := r.Form["date"]
	if ok {
		date, err = time.Parse(dateFmt, datePar[0])
		if err != nil {
			fmt.Fprintf(w, "błędna data aktualna: %s\n", err)
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
			fmt.Fprintf(w, "błędny zakres: %s\n", err)
			return
		}

		const maxRange = 1000
		if n > maxRange {
			fmt.Fprintf(w, "za duża wartość zakresu (max: %d): %d\n", maxRange, n)
			return
		}
		*rangeFlag = n
	}

	// wyświetl biorytm
	buf := new(bytes.Buffer)
	printBiorytm(buf, born, date, *rangeFlag)

	err = outputTmpl.Execute(w, buf.String())
	if err != nil {
		log.Printf("execute output template: %s", err)
	}
}

func biorytmWeb() {
	var err error

	formTmpl = template.New("form")
	formTmpl, err = formTmpl.Parse(formTmplStr)
	if err != nil {
		log.Fatalf("template 'form' parse: %s", err)
	}

	outputTmpl = template.New("output")
	outputTmpl, err = outputTmpl.Parse(outputTmplStr)
	if err != nil {
		log.Fatalf("template 'output' parse: %s", err)
	}

	http.HandleFunc("/", biorytmHandler)

	log.Printf("adres usługi: %s", *httpAddr)
	err = http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
