package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

/* Probably want to implement this as a struct, but you know... woops */
type Message map[string]interface{}

/* Fixing my mistakes */
// type MessageReal struct

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
