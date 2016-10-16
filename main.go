/*
	TODO:
		Authenticate connections
		CLI shell
*/
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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

type Message map[string]interface{}

func main() {
	var mode = flag.String("mode", "client",
		"Run the program in either shell mode or server mode")
	var test = flag.Bool("test", false, "Test a feature")

	flag.Parse()

	if *test {
		Testdb()
		return
	}

	if *mode == clientMode {
		StartCLIShell()
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

		// Talking to clients happens her.
		go handleConnection(conn)
	}
}

func outgoingChannel(conn net.Conn, outgoing chan Message) {
	defer conn.Close()
	for {
		currentMessage := <-outgoing

		err := writeJSON(conn, currentMessage)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

// Add CLI here.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		err := readJSONMessage(conn)
		if err != nil {
			opErr, ok := err.(net.Error)
			if ok {
				if !opErr.Temporary() {
					return
				} else {
					continue
				}

			}

			if err == io.EOF {
				return
			}

			log.Println(err.Error())
		}

	}
}

func sendHeartbeat(conn net.Conn) {
	data := 0
	message := Message{"signal": signalHeartbeat, "data": data}
	if err := writeJSON(conn, message); err != nil {
		log.Print("Sending heartbeat failed. ")
		log.Println(err.Error())
	}
}

func heartbeatResponse(conn net.Conn) error {
	data := 0
	message := Message{"signal": signalHeartbeatResponse, "data": data}
	if err := writeJSON(conn, message); err != nil {
		return err
	}

	return nil
}

func readJSONMessage(conn net.Conn) (err error) {
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
