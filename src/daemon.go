package main

import (
    "fmt"
    "net/http"
)

func handle_all(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "To trigger panic, browse to /panic")
}

func handle_panic(w http.ResponseWriter, r *http.Request) {
	print("Panic!")
}

func main() {
    http.HandleFunc("/", handle_all)
    http.HandleFunc("/panic", handle_panic)
    http.ListenAndServe(":9999", nil)
}