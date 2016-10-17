/*
	TODO:
		http server
		Authenticate connections
		Use TLS
		Add fake device
*/
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
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
	fmt.Println("Starting server.")
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err.Error())
	}

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

func readJSONMessage(conn net.Conn, outgoing chan Message) (err error) {
	rawMessage, err := readMessage(conn)
	if err != nil {
		return err
	}

	var jsonMessage interface{}
	decoder := json.NewDecoder(strings.NewReader(rawMessage))
	decoder.Decode(&jsonMessage)
	log.Println(jsonMessage)

	typedMessage, ok := jsonMessage.(map[string]interface{})
	if !ok {
		return errors.New("Could not assert type of message.")
	}

	signal, ok := typedMessage["signal"].(string)
	if !ok {
		return errors.New("Could assert signal type to string.")
	}
	log.Println("Signal type: ", signal)

	fmt.Printf("%T\n", typedMessage["data"])
	data, ok := typedMessage["data"].(float64)
	if !ok {
		return errors.New("Data field wasn't a number.")
	}
	log.Println("Data payload: ", data)

	if signal == signalHeartbeat {
		log.Println("Recording heartbeat")
		err := RecordHeartbeat(data)
		if err != nil {
			panic(err.Error())
		}
	}

	return
}

func readMessage(conn net.Conn) (string, error) {
	// Set deadline for 3 seconds from now.
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	rawMessage, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	// Remove newline delimeter.
	return rawMessage[0 : len(rawMessage)-1], nil
}

func writeJSON(conn net.Conn, jsonObject interface{}) error {
	// Set deadline for 3 seconds from now.
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(jsonObject)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(conn, "\n")
	if err != nil {
		return err
	}

	return nil
}
