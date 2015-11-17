package main

import (
	"fmt"
	"github.com/op/go-logging"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"time"
)

// Key word sent to trigger panic.
const Keyword string = "Panic!"

// HTTP port to listen on.
const HTTPPort int = 9999

// UDP port to listen on.
const UDPPort int = 9998

// Buffer size.
const BufferSize int = 64000

// Instantiate the logger.
var log = logging.MustGetLogger("gopanic")

// Generic error handle.
func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// Called when there is an HTTP request to the server, all
// we need to do is propagate the panic and then run the
// panic.
func handleHTTP(w http.ResponseWriter, r *http.Request) {
	propagatePanic()
	doPanic()
}

// Listen for UDP packets on the broadcast address, if
// this matches Keyword, then run the panic.
func handleUDP() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: UDPPort,
	})
	errorCheck(err)
	defer conn.Close()

	var buff = make([]byte, BufferSize)
	n, addr, err := conn.ReadFromUDP(buff)
	errorCheck(err)

	if string(buff[0:n]) == Keyword {
		log.Info("Accepting panic from %s", addr)
		doPanic()
	}
}

// Propagate the panic to the broadcast address, so that
// all hosts on this network can see.
func propagatePanic() {
	log.Info("Propagating panic to broadcast")
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: UDPPort,
	})
	errorCheck(err)
	defer conn.Close()
	_, err = conn.Write([]byte(Keyword))
	errorCheck(err)
}

// Template for what to do in the case of a panic.
// By default, we'll try to halt the system.
func doPanic() {
	proc := exec.Command("shutdown", "-H", "now")
	proc.Start()
}

// Setup the formatter and level of the logger.
func setupLogger() {
	format := logging.MustStringFormatter(
		"(%{time:2006/01/02 15:04:05.999 -07:00}) %{color}[%{level}]%{color:reset}: %{message}",
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backendLeveled, backendFormatter)
}

// Entry point.  Start the HTTP server and start
// listening for UDP.
func main() {
	setupLogger()
	log.Info("Starting gopanic daemon")
	currentUser, error := user.Current()
	errorCheck(error)
	if currentUser.Uid != "0" {
		log.Warning("Not running as root user (%s), this means no halt on panic.", currentUser.Username)
	}
	http.HandleFunc("/", handleHTTP)
	go http.ListenAndServe(fmt.Sprintf(":%d", HTTPPort), nil)
	for {
		time.Sleep(100 * time.Millisecond)
		handleUDP()
	}
}
