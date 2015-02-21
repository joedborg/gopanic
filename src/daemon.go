package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const KEYWORD string = "Panic!"
const HTTP_PORT string = ":9999"
const UDP_PORT int = 9998
const BUFFER_SIZE = 64000

func error_check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Called when there is an HTTP request to the server, all
// we need to do is propagate the panic and then run the
// panic.
func handle_http(w http.ResponseWriter, r *http.Request) {
	propagate_panic()
	do_panic()
}

// Listen for UDP packets on the broadcast address, if
// this matches KEYWORD, then run the panic.
func handle_udp() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: UDP_PORT,
	})
	error_check(err)
	defer conn.Close()

	var buff []byte = make([]byte, BUFFER_SIZE)
	n, addr, err := conn.ReadFromUDP(buff)
	error_check(err)

	if string(buff[0:n]) == KEYWORD {
		fmt.Println("Accepting panic from ", addr)
		do_panic()
	}
}

// Propagate the panic to the broadcast address, so that
// all hosts on this network can see.
func propagate_panic() {
	fmt.Println("Propagating panic to broadcast")
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: UDP_PORT,
	})
	error_check(err)
	defer conn.Close()
	_, err = conn.Write([]byte(KEYWORD))
	error_check(err)
}

// Template for what to do in the case of a panic.
func do_panic() {

}

// Entry point.  Start the HTTP server and start
// listening for UDP.
func main() {
	http.HandleFunc("/", handle_http)
	go http.ListenAndServe(HTTP_PORT, nil)
	for {
		time.Sleep(100 * time.Millisecond)
		handle_udp()
	}
}
