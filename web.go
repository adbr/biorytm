// 2015-02-24 Adam Bryt

package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	textFormTmpl     = template.New("textForm")
	textDisplayTmpl  = template.New("textDisplay")
	graphFormTmpl    = template.New("graphForm")
	graphDisplayTmpl = template.New("graphDisplay")
)

// opts zawiera wartości flag programu.
var opts struct {
	born   time.Time // data urodzenia
	date   time.Time // data biorytmu
	drange int       // liczba dni biorytmu (days range)
}

// params grupuje parametry biorytmu - dla prezentacji na stronach.
type params struct {
	Born        time.Time // data urodzenia
	Date        time.Time // data biorytmu
	Drange      int       // liczba dni biorytmu (days range)
	ImageString string    // obrazek zakodowany w base64
}

func (p params) DateString() string {
	if p.Date.IsZero() {
		return ""
	}
	return p.Date.Format(dateFmt)
}

func (p params) BornString() string {
	if p.Born.IsZero() {
		return ""
	}
	return p.Born.Format(dateFmt)
}

func biorytmWeb() {
	initTemplates()
	parseOpts()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/text/form/", textFormHandler)
	http.HandleFunc("/text/display/", textDisplayHandler)
	http.HandleFunc("/graph/form/", graphFormHandler)
	http.HandleFunc("/graph/display/", graphDisplayHandler)

	log.Printf("biorytm: adres usługi HTTP: %s", *httpFlag)
	err := http.ListenAndServe(*httpFlag, nil)
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

	graphFormTmpl, err = graphFormTmpl.Parse(graphFormHTML)
	if err != nil {
		log.Fatalf("błąd parsowania template 'graphFormHTML': %s", err)
	}

	graphDisplayTmpl, err = graphDisplayTmpl.Parse(graphDisplayHTML)
	if err != nil {
		log.Fatalf("błąd parsowania template 'graphDisplayHTML': %s", err)
	}
}

// parseOpts parsuje wartości flag programu i ustawia zmienną opts.
func parseOpts() {
	if *bornFlag != "" {
		b, err := time.Parse(dateFmt, *bornFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
		}
		opts.born = b
	}

	if *dateFlag != "" {
		d, err := time.Parse(dateFmt, *dateFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
		}
		opts.date = d
	} else {
		n := time.Now()
		d := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
		opts.date = d
	}

	opts.drange = *rangeFlag
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/graph/form/", http.StatusFound)
}

// textFormHandler wyświetla formatkę dla wprowadzania danych biorytmu
// w postaci tekstowej.
func textFormHandler(w http.ResponseWriter, r *http.Request) {
	// ustaw początkowe dane dla formatki
	p := params{
		Born:   opts.born,
		Date:   opts.date,
		Drange: opts.drange,
	}

	err := textFormTmpl.Execute(w, p)
	if err != nil {
		log.Printf("błąd wykonania template 'textFormTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'textFormTmpl': %s\n", err)
	}
}

// textDisplayHandler wyświetla biorytm w postaci tekstowej.
func textDisplayHandler(w http.ResponseWriter, r *http.Request) {
	p, err := getTextFormData(r)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	if p.Date.Before(p.Born) {
		log.Printf("data biorytmu (%s) wcześniejsza niż data urodzenia (%s)",
			p.DateString(), p.BornString())
		fmt.Fprintf(w, "data biorytmu (%s) wcześniejsza niż data urodzenia (%s)\n",
			p.DateString(), p.BornString())
		return
	}

	buf := new(bytes.Buffer)
	printBiorytm(buf, p.Born, p.Date, p.Drange)

	err = textDisplayTmpl.Execute(w, buf.String())
	if err != nil {
		log.Printf("błąd wykonania template 'textDisplayTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'textDisplayTmpl': %s\n", err)
	}
}

// graphFormHandler wyświetla formatkę dla wprowadzania danych biorytmu
// w postaci graficznej.
func graphFormHandler(w http.ResponseWriter, r *http.Request) {
	// ustaw początkowe dane dla formatki
	p := params{
		Born:   opts.born,
		Date:   opts.date,
		Drange: opts.drange,
	}

	err := graphFormTmpl.Execute(w, p)
	if err != nil {
		log.Printf("błąd wykonania template 'graphFormTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'graphFormTmpl': %s\n", err)
	}
}

// graphDisplayHandler wyświetla biorytm w postaci graficznej.
func graphDisplayHandler(w http.ResponseWriter, r *http.Request) {
	p, err := getTextFormData(r)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	img := biorytmImage()
	p.ImageString, err = encodeImage(img)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	err = graphDisplayTmpl.Execute(w, p)
	if err != nil {
		log.Printf("błąd wykonania template 'graphDisplayTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'graphDisplayTmpl': %s\n", err)
	}
}

// biorytmImage zwraca obrazek z wykresem biorytmu.
func biorytmImage() image.Image {
	r := image.Rect(0, 0, 800, 600)
	img := image.NewRGBA(r)

	back := color.RGBA{0, 50, 128, 128} // background color

	// ustaw kolor tła
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, back)
		}
	}

	return img
}

func encodeImage(img image.Image) (string, error) {
	buf := new(bytes.Buffer)
	b64 := base64.NewEncoder(base64.StdEncoding, buf)

	err := png.Encode(b64, img)
	if err != nil {
		return "", err
	}
	b64.Close()

	s := buf.String()
	return s, nil
}

// getTextFormData pobiera z requestu i parsuje parametry biorytmu.
func getTextFormData(r *http.Request) (params, error) {
	p := params{}

	err := r.ParseForm()
	if err != nil {
		return p, err
	}

	// pobierz datę urodzenia

	bornPar, ok := r.Form["born"]
	if !ok {
		return p, errors.New("brak parametru 'born' (data urodzenia)")
	}

	p.Born, err = time.Parse(dateFmt, bornPar[0])
	if err != nil {
		return p, fmt.Errorf("błędna data urodzenia: %s", err)
	}

	// pobierz datę aktualną

	datePar, ok := r.Form["date"]
	if !ok {
		return p, errors.New("brak parametru 'date' (data aktualna)")
	}

	p.Date, err = time.Parse(dateFmt, datePar[0])
	if err != nil {
		return p, fmt.Errorf("błędna data aktualna: %s", err)
	}

	// pobierz parametr range

	rangePar, ok := r.Form["range"]
	if !ok {
		return p, errors.New("brak parametru 'range' (liczba dni biorytmu)")
	}

	n, err := strconv.Atoi(rangePar[0])
	if err != nil {
		return p, fmt.Errorf("błędny zakres dni: %s", err)
	}

	const maxRange = 1000
	if n > maxRange {
		return p, fmt.Errorf("za duża wartość zakresu (max: %d): %d", maxRange, n)
	}
	p.Drange = n

	return p, nil
}
