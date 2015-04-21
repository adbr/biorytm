// 2015-02-25 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"os"
)

const dateFmt = "2006-01-02"

var (
	bornFlag  = flag.String("born", "", "data urodzenia (yyyy-mm-dd)")
	dateFlag  = flag.String("date", "", "data biorytmu (yyyy-mm-dd)")
	httpFlag  = flag.String("http", "", "adres usługi HTTP (np. ':5050')")
	daysFlag  = flag.Int("days", 15, "liczba dni biorytmu")
	fontsFlag = flag.String("fonts", "", "katalog z fontami")
	helpFlag  = flag.Bool("help", false, "wyświetla help")
)

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func help() {
	fmt.Fprintln(os.Stdout, helpStr)
	os.Exit(0)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *helpFlag {
		help()
	}

	if *httpFlag != "" {
		webMain()
		return
	}

	cliMain()
}
