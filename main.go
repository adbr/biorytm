// 2015-02-25 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"os"
)

const usageStr = `Sposób użycia:
	biorytm [flagi] -born <data urodzenia>
	biorytm [flagi] -http <host:port>

Flagi:
	-born="": data urodzenia w formacie yyyy-mm-dd
	-date="": data biorytmu w formacie yyyy-mm-dd (domyślnie: dzisiaj)
	-http="": adres usługi HTTP (np. ':5050')
	-days=15: liczba dni biorytmu
	-fonts="": katalog z fontami (domyślnie ./fonts lub $VGFONTPATH)`

const dateFmt = "2006-01-02"

var (
	bornFlag  = flag.String("born", "", "data urodzenia (yyyy-mm-dd)")
	dateFlag  = flag.String("date", "", "data biorytmu (yyyy-mm-dd)")
	httpFlag  = flag.String("http", "", "adres usługi HTTP (np. ':5050')")
	daysFlag  = flag.Int("days", 15, "liczba dni biorytmu")
	fontsFlag = flag.String("fonts", "", "katalog z fontami")
)

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *httpFlag != "" {
		webMain()
		return
	}

	cliMain()
}
