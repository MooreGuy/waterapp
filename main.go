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
	"github.com/MooreGuy/waterapp/device"
	"log"
	"net"
	"net/http"
)

const (
	clientMode = "client"
	serverMode = "server"

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
	} else if *mode == serverMode {
		startDaemon()
	} else {
		log.Fatal("Unknown mode: " + *mode)
	}
}

func startDaemon() {
	fmt.Println("Starting socket server.")
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err.Error())
	}

	GetDeviceInfo()
	go ListenSocketServer(listener)

	devices := device.FindDevices()
	log.Println(len(devices))

	website := Website{}
	fmt.Println("Starting web server.")
	http.ListenAndServe(":8081", website)
}

func ListenSocketServer(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		defer conn.Close()

		outgoingChan := make(chan Message, 10)
		go outgoing(conn, outgoingChan)
		go reading(conn, outgoingChan)
	}
}

func GetDeviceInfo() {
	for {
		log.Println("Reading devices")
		devices := device.FindDevices()
		for _, device := range devices {
			readBuf := make([]byte, 2, 2)
			device.Read(readBuf)
			log.Println(readBuf)
		}
	}
}
