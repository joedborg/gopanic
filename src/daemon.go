package main

import (
	"fmt"
	"log"
	"net"
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
	go http.ListenAndServe(":9999", nil)

	ln, err := net.ListenPacket("udp", ":9998")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
}
