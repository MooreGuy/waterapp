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
	"strconv"
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
	flag.Parse()

	if *mode == clientMode {
		cliShell()
	} else if *mode == serverMode {
		startDaemon()
	} else {
		log.Fatal("Unknown mode: " + *mode)
	}
}

func startDaemon() {
	log.Println("Starting server.")
	listener, err := net.Listen("tcp", ":8080")
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

func cliShell() {
	log.Println("Connecting to server daemon")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error opening connection.", err.Error())
	}
	defer conn.Close()

	messageChannel := make(chan Message, 10)
	go handleConnection(conn)
	go outgoingChannel(conn, messageChannel)

	var data int = -1
	for {
		var input string
		_, err = fmt.Scanln(&input)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		data, err = strconv.Atoi(input)
		if err != nil {
			log.Println("Bad input, enter number.")
			data = -1
		}

		log.Println("Sending: ", data)

		message := Message{"signal": signalHeartbeat, "data": data}
		messageChannel <- message
	}
}

func outgoingChannel(conn net.Conn, outgoing chan Message) {
	defer conn.Close()

	for {
		currentMessage := <-outgoing

		log.Println("sent")
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
	log.Debug("Signal type: ", signal)

	data, ok := typedMessage["data"].(float64)
	if !ok {
		return errors.New("Data field wasn't a number.")
	}
	log.Debug("Data payload: ", data)

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
