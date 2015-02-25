// 2015-02-25 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"os"
)

const usageStr = `usage: biorytm [flagi] data_urodzenia [data_biorytmu]
	data w formacie 'yyyy-mm-dd'
	flagi:
		-range=15: zakres dni biorytmu
		-http:"": adres usługi HTTP (np. ':5050')`

const dateFmt = "2006-01-02"

var (
	rangeFlag = flag.Int("range", 15, "zakres dni biorytmu")
	httpAddr  = flag.String("http", "", "adres usługi (np. ':5050')")
)

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *httpAddr != "" {
		biorytmWeb()
		return
	}

	biorytmCli()
}
