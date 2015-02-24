// 2015-02-24 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":5050", "adres usługi")
)

func biorytmHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("form: %v\n", r.Form)
	fmt.Printf("request: %s\n", r)
	fmt.Fprintln(w, "biorytmHandler: test")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", biorytmHandler)

	log.Printf("adres usługi: %s", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
