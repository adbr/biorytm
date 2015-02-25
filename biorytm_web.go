// 2015-02-24 Adam Bryt

package main

import (
	"fmt"
	"log"
	"net/http"
)

func biorytmHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("form: %v\n", r.Form)
	fmt.Printf("request: %v\n", r)
	fmt.Fprintln(w, "biorytmHandler: test")
}

func biorytmWeb() {
	http.HandleFunc("/", biorytmHandler)

	log.Printf("adres usługi: %s", *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
