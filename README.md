GoPanic
=======
Distributed panic.

Abstract
--------
GoPanic is daemon that runs and, upon an HTTP request, propagates a panic message across the network.

Usage
-----
First, check that daemon.go is setup correctly (perhaps you want to change the KEYWORD to something more secret?), then be sure to add to the function `do_panic()` Compile the daemon with `go build daemon.go`.  Run the executable on all required hosts.  To start the panic propagation, send an HTTP GET to any host on the network, on the chosen port.
