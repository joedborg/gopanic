package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const KEYWORD string = "Panic!"
const UDP_PORT int = 9998

func error_check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handle_http(w http.ResponseWriter, r *http.Request) {
	propagate_panic()
	do_panic()
}

func handle_udp() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: UDP_PORT,
	})
	error_check(err)
	defer conn.Close()

	buff := []byte{}
	n, addr, err := conn.ReadFromUDP(buff)
	error_check(err)

	fmt.Println("From address:", addr, " Got message:", string(buff[0:n]), n)
}

func propagate_panic() {
	fmt.Println("Sending:", KEYWORD)
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: UDP_PORT,
	})
	error_check(err)
	defer conn.Close()
	_, err = conn.Write([]byte(KEYWORD))
	error_check(err)
}

func do_panic() {

}

func main() {
	http.HandleFunc("/", handle_http)
	go http.ListenAndServe(":9999", nil)
	for {
		time.Sleep(100 * time.Millisecond)
		handle_udp()
	}
}
