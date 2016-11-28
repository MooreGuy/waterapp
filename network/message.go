package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	signalField = "signal"
	dataField   = "data"

	signalHeartbeat         = "heartbeat"
	signalHeartbeatResponse = "heartbeat_response"
)

// Probably want to implement this as a struct, but you know... woops
type Message map[string]interface{}

// Fixing my mistakes
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

func JSONMessage(rawMessage string) (mes Message, err error) {
	var jsonMessage interface{}
	decoder := json.NewDecoder(strings.NewReader(rawMessage))
	decoder.Decode(&jsonMessage)
	log.Println(jsonMessage)
	mes = map[string]interface{}{}

	mes, ok := jsonMessage.(map[string]interface{})
	if !ok {
		return mes, errors.New("Could not assert type of message.")
	}

	signal, ok := mes["signal"].(string)
	if !ok {
		return mes, errors.New("Could assert signal type to string.")
	}
	log.Println("Signal type: ", signal)

	sensorid, ok := mes["sensorid"].(string)
	if !ok {
		return mes, errors.New("Could assert signal type to string.")
	}
	log.Println("Sensorid: ", sensorid)

	fmt.Printf("%T\n", mes["data"])
	data, ok := mes["data"].(float64)
	if !ok {
		return mes, errors.New("Data field wasn't a number.")
	}
	log.Println("Data payload: ", data)

	/**
	if signal == signalHeartbeat {
		log.Println("Recording heartbeat")
		err := RecordHeartbeat(data)
		if err != nil {
			panic(err.Error())
		}
	}
	*/

	return mes, err
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
