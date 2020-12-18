package main

import (
	"fmt"
	"log"
	"net/http"
	"ledger/pkg/mytemplate"
)

func main() {
	// register a path
	// instead of constructing and returning a response object, we write directly
	// to the response object (w)
	// because of this, Golang http is http2 and websockets compatible
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", r.URL)
	})
	// http.HandleFunc("/insert", mytemplate.insertHandler)
	http.HandleFunc("/ledger", mytemplate.LedgerHandler)
	http.HandleFunc("/dailyledger", mytemplate.DailyLedgerHandler)
	http.HandleFunc("/insert", mytemplate.InsertHandler)
	http.HandleFunc("/save", mytemplate.SaveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
