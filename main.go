/*
	TODO:
		http server
		Authenticate connections
		Use TLS
		Add fake device
*/
package main

import (
	"flag"
	"log"
)

const (
	clientMode     = "client"
	controllerMode = "controller"
	aggregatorMode = "aggregator"

	signalField = "signal"
	dataField   = "data"

	signalHeartbeat         = "heartbeat"
	signalHeartbeatResponse = "heartbeat_response"
)

func main() {
	var mode = flag.String("mode", "client",
		"Run the program in either shell mode or server mode")
	var test = flag.Bool("test", false, "Test a feature")
	var username = flag.String("u", "", "Username")
	var password = flag.String("p", "", "Password")

	flag.Parse()

	if *test {
		Testdb()
		return
	}

	if *mode == clientMode {
		StartCLIShell(*username, *password)
	} else if *mode == controllerMode {
		StartController()
	} else if *mode == aggregatorMode {
		StartAggregator()
	} else {
		log.Fatal("Unknown mode: " + *mode)
	}

	for {
	}
}
