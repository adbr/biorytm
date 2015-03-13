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
	"strconv"
	"time"
)

var (
	textFormTmpl     = template.New("textForm")
	textDisplayTmpl  = template.New("textDisplay")
	graphFormTmpl    = template.New("graphForm")
	graphDisplayTmpl = template.New("graphDisplay")
)

// biorytm grupuje parametry biorytmu.
type biorytm struct {
	Born        time.Time
	Date        time.Time
	Days        int
	ImageString string // obrazek zakodowany w base64
}

func (b biorytm) DateString() string {
	return b.Date.Format(dateFmt)
}

func (b biorytm) BornString() string {
	return b.Born.Format(dateFmt)
}

func biorytmWeb() {
	initTemplates()

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

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/graph/form/", http.StatusFound)
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

// graphFormHandler wyświetla formatkę dla wprowadzania danych biorytmu
// w postaci graficznej.
func graphFormHandler(w http.ResponseWriter, r *http.Request) {
	// ustaw domyślne dane dla formatki
	b := biorytm{
		Date: time.Now(),
		Days: 30,
	}

	err := graphFormTmpl.Execute(w, b)
	if err != nil {
		log.Printf("błąd wykonania template 'graphFormTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'graphFormTmpl': %s\n", err)
	}
}

// graphDisplayHandler wyświetla biorytm w postaci graficznej.
func graphDisplayHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getTextFormData(r)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	img := biorytmImage()
	b.ImageString, err = encodeImage(img)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	err = graphDisplayTmpl.Execute(w, b)
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
