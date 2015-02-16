// 2015-02-16 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	dateFmt = "2006-01-02"
)

const usageStr = "usage: biorytm data_urodzenia [data_biorytmu]"

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	btime, err := time.Parse(dateFmt, flag.Arg(0)) // born time
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage()
	}

	dtime := time.Now() // destination time
	if flag.NArg() > 1 {
		dtime, err = time.Parse(dateFmt, flag.Arg(1))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			usage()
		}
	}

	fmt.Printf("data urodzenia: %s\n", btime.String())
	fmt.Printf("data docelowa:  %s\n", dtime.String())
}
