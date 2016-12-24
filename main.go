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
	"fmt"
	"log"
)

const (
	cliMode        = "cli"
	controllerMode = "controller"
	aggregatorMode = "aggregator"
	relayMode      = "relayMode"

	signalField = "signal"
	dataField   = "data"

	signalHeartbeat         = "heartbeat"
	signalHeartbeatResponse = "heartbeat_response"
)

func main() {
	var mode = flag.String("mode", cliMode,
		"Run the program in either shell mode or server mode")
	var test = flag.Bool("test", false, "Test a feature")
	var username = flag.String("u", "", "Username")
	var password = flag.String("p", "", "Password")

	flag.Parse()

	if *test {
		Testdb()
		return
	}

	if *mode == cliMode {
		StartCLIShell(*username, *password)
	} else if *mode == controllerMode {
		StartController()
	} else if *mode == aggregatorMode {
		StartAggregator()
	} else if *mode == relayMode {
		StartRelay()
	} else {
		log.Fatal("Unknown mode: " + *mode)
	}

	log.Println("Startup process complete.")

	var something string
	fmt.Scanln(&something)
	if something == "\n" {
	}
}
