package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello,"))
	fmt.Fprintf(w, r.RemoteAddr + "\n")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Printf("Starts server at ;8080")

	http.ListenAndServe(":8080", nil)
}
