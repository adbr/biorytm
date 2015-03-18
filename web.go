// 2015-02-24 Adam Bryt

package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// opts zawiera wartości flag programu.
var opts struct {
	born time.Time // data urodzenia
	date time.Time // data biorytmu
	days int       // liczba dni biorytmu
}

// params zawiera parametry biorytmu przekazane w request.
type params struct {
	born time.Time // data urodzenia
	date time.Time // data biorytmu
	days int       // liczba dni biorytmu
	look string    // jak prezentować biorytm [text|graph]
}

func biorytmWeb() {
	initTemplates()
	parseOpts()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/form/", formHandler)
	http.HandleFunc("/display/", displayHandler)

	log.Printf("biorytm: adres usługi HTTP: %s", *httpFlag)
	err := http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func initTemplates() {
	var err error

	formTmpl, err = formTmpl.Parse(formHTML)
	if err != nil {
		log.Fatalf("błąd parsowania template 'formHTML': %s", err)
	}

	textTmpl, err = textTmpl.Parse(textHTML)
	if err != nil {
		log.Fatalf("błąd parsowania template 'textHTML': %s", err)
	}

	graphTmpl, err = graphTmpl.Parse(graphHTML)
	if err != nil {
		log.Fatalf("błąd parsowania template 'graphHTML': %s", err)
	}
}

// parseOpts parsuje wartości flag programu i ustawia zmienną opts.
func parseOpts() {
	// born
	if *bornFlag != "" {
		b, err := time.Parse(dateFmt, *bornFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
		}
		opts.born = b
	}

	// date
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

	// days
	opts.days = *daysFlag
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/form/", http.StatusFound)
}

// formHandler wyświetla formatkę dla wprowadzania danych biorytmu.
func formHandler(w http.ResponseWriter, r *http.Request) {
	// ustaw początkowe dane dla formatki
	d := formData{}
	if !opts.born.IsZero() {
		d.Born = opts.born.Format(dateFmt)
	}
	if !opts.date.IsZero() {
		d.Date = opts.date.Format(dateFmt)
	}
	d.Days = opts.days

	err := formTmpl.Execute(w, d)
	if err != nil {
		log.Printf("błąd wykonania template 'formTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'formTmpl': %s\n", err)
	}
}

// displayHandler wyświetla biorytm w wybranej postaci.
func displayHandler(w http.ResponseWriter, r *http.Request) {
	p, err := getParams(r)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	if p.date.Before(p.born) {
		born := p.born.Format(dateFmt)
		date := p.date.Format(dateFmt)
		log.Printf("data (%s) wcześniejsza niż data urodzenia (%s)", date, born)
		fmt.Fprintf(w, "data (%s) wcześniejsza niż data urodzenia (%s)\n", date, born)
		return
	}

	switch p.look {
	case "text":
		displayText(w, p)
	case "graph":
		displayGraph(w, p)
	default:
		log.Printf("nie poprawna wartość parametru 'look': %s", p.look)
		fmt.Fprintf(w, "nie poprawna wartość parametru 'look': %s\n", p.look)
	}
}

// displayText wyświetla biorytm w postaci tekstowej.
func displayText(w http.ResponseWriter, p params) {
	var buf bytes.Buffer
	printBiorytm(&buf, p.born, p.date, p.days)

	d := textData{}
	d.Text = buf.String()

	err := textTmpl.Execute(w, d)
	if err != nil {
		log.Printf("błąd wykonania template 'textTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'textTmpl': %s\n", err)
	}
}

// displayGraph wyświetla biorytm w postaci graficznej.
func displayGraph(w http.ResponseWriter, p params) {
	img := biorytmImage()

	var err error
	d := graphData{}
	d.Born = p.born.Format(dateFmt)
	d.Date = p.date.Format(dateFmt)
	d.Image, err = encodeImage(img)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	err = graphTmpl.Execute(w, d)
	if err != nil {
		log.Printf("błąd wykonania template 'graphTmpl': %s", err)
		fmt.Fprintf(w, "błąd wykonania template 'graphTmpl': %s\n", err)
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

// encodeImage koduje img w base64.
func encodeImage(img image.Image) (string, error) {
	var buf bytes.Buffer
	b64 := base64.NewEncoder(base64.StdEncoding, &buf)

	err := png.Encode(b64, img)
	if err != nil {
		return "", err
	}
	b64.Close()

	s := buf.String()
	return s, nil
}

// getParams pobiera i parsuje parametry biorytmu zawarte w r.
func getParams(r *http.Request) (params, error) {
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

	p.born, err = time.Parse(dateFmt, bornPar[0])
	if err != nil {
		return p, fmt.Errorf("błędna data urodzenia: %s", err)
	}

	// pobierz datę aktualną

	datePar, ok := r.Form["date"]
	if !ok {
		return p, errors.New("brak parametru 'date' (data aktualna)")
	}

	p.date, err = time.Parse(dateFmt, datePar[0])
	if err != nil {
		return p, fmt.Errorf("błędna data aktualna: %s", err)
	}

	// pobierz parametr days

	daysPar, ok := r.Form["days"]
	if !ok {
		return p, errors.New("brak parametru 'days' (liczba dni biorytmu)")
	}

	n, err := strconv.Atoi(daysPar[0])
	if err != nil {
		return p, fmt.Errorf("błędna liczba dni: %s", err)
	}

	const maxDays = 1000
	if n > maxDays {
		return p, fmt.Errorf("za duża wartość liczby dni (max: %d): %d", maxDays, n)
	}

	if n <= 0 {
		return p, fmt.Errorf("liczba dni musi być większa od 0, jest: %d", n)
	}

	p.days = n

	// pobierz parametr look

	lookPar, ok := r.Form["look"]
	if !ok {
		return p, errors.New("brak parametru 'look' (sposób prezentacji)")
	}

	p.look = lookPar[0]
	if !(p.look == "text" || p.look == "graph") {
		return p, fmt.Errorf("błędna wartość parametru 'look': %s", p.look)
	}

	return p, nil
}
