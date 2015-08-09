package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"time"
)

const KEYWORD string = "Panic!"
const HTTP_PORT string = ":9999"
const UDP_PORT int = 9998
const BUFFER_SIZE int = 64000

// Generic error handle.
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
// By default, we'll try to halt the system.
func do_panic() {
    proc := exec.Command("shutdown", "-H", "now")
		proc.Start()
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
