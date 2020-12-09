package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// register a path
	// instead of constructing and returning a response object, we write directly
	// to the response object (w)
	// because of this, Golang http is http2 and websockets compatible
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v", r.URL)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
