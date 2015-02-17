package main

import (
	"fmt"
	"html"
	"net/http"
)

func main() {
	http.HandleFunc("/bar", handleFunction)
	http.ListenAndServe(":8080", nil)
}

func handleFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
