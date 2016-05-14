package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	client = "client"
	server = "server"

	signalField = "signal"
	dataField   = "data"

	signalResponse       = "response"
	signalConnect        = "connect"
	signalConnectConfirm = "connect_confirm"

	signalFlow          = "flow"
	signalFlowResponse  = "flow_response"
	signalLevel         = "level"
	signalLevelResponse = "level_response"

	signalHeartbeat         = "heartbeat"
	signalHeartbeatResponse = "heartbeat_response"

	connectI2CDevice   = 0xFE
	connectI2CResponse = 0xFF
)

type Message map[string]interface{}

func main() {
	var mode = flag.String("mode", "client",
		"Run the program in either client mode")
	flag.Parse()

	if *mode == client {
		clientFunc()
	} else if *mode == server {
		serverFunc()
	} else {
		log.Fatal("Unknown mode: " + *mode)
	}
}

func serverFunc() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		time.Sleep(time.Second)

		err := serverReadMessage(conn)
		if err != nil {
			log.Println(err.Error())
		}

		sendHeartbeat(conn)
	}
}

func serverReadMessage(conn net.Conn) error {
	json, err := readJSON(conn)
	if err != nil {
		log.Print("Error reading json: ")
		return err
	}
	mapMessage, ok := json.(map[string]interface{})
	if !ok {
		return errors.New("Could not assert message type to map[string]interface{}")
	}

	signal, ok := mapMessage["signal"].(string)
	if !ok {
		return errors.New("Could assert signal type to string.")
	}

	if signal == signalConnect {
		log.Print("Request to connect received")
		connectResponse(conn)
	} else if signal == signalFlowResponse {
		log.Print("Flow response received")
	} else if signal == signalLevelResponse {
		log.Print("Level response received")
	} else if signal == signalHeartbeatResponse {
		log.Print("Heartbeat response received")
	} else {
		log.Print("unknown signal")
	}

	return nil
}

func sendDummyFlow(conn net.Conn) {
	data := 123
	message := Message{signalField: signalFlowResponse, dataField: data}
	if err := writeJSON(conn, message); err != nil {
		log.Print("Sending flow failed.")
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

func connectRequest(conn net.Conn) {
	data := 123
	message := Message{"signal": signalConnect, "data": data}
	if err := writeJSON(conn, message); err != nil {
		log.Fatal(err.Error())
	}
}

func connectResponse(conn net.Conn) {
	data := 123
	message := Message{"signal": signalConnectConfirm, "data": data}
	if err := writeJSON(conn, message); err != nil {
		log.Fatal(err.Error())
	}
}

func clientFunc() {
	//devices := findDevices()
	//sendTestMessage(devices)

	conn, err := net.Dial("tcp", "home.guywmoore.com:8080")
	if err != nil {
		log.Fatal("Error opening connection.", err.Error())
	}
	defer conn.Close()

	connectRequest(conn)

	for {
		if err := clientReadMessage(conn); err != nil {
			log.Println("Reading from sockets has its problems.")
		}

		sendDummyFlow(conn)

		time.Sleep(time.Second)
	}
}

func clientReadMessage(conn net.Conn) error {
	message, err := readJSON(conn)
	if err != nil {
		log.Println("Error reading json.")
		return err
	}

	typedMessage, ok := message.(map[string]interface{})
	if !ok {
		return errors.New("Could not assert type of message.")
	}

	signal, ok := typedMessage["signal"].(string)
	if !ok {
		return errors.New("Could assert signal type to string.")
	}

	data, ok := typedMessage["data"].(int)
	if !ok {
		return errors.New("Data field wasn't a number.")
	}

	log.Println("data: " + strconv.Itoa(data))
	if signal == signalFlow {
		log.Print("Flow request received")
	} else if signal == signalLevel {
		log.Print("Level request received")
	} else if signal == signalHeartbeat {
		log.Print("Heartbeat received")
		if err := heartbeatResponse(conn); err != nil {
			log.Print("Error getting heartbeat response.")
			log.Println(err.Error())
			return err
		}
	} else if signal == signalConnectConfirm {
		log.Println("Successfully connected to master server.")
	} else {
		log.Print("unknown signal")
	}

	return nil
}

func readJSON(conn net.Conn) (interface{}, error) {
	// Set deadline for 3 seconds from now.
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	rawMessage, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err
	}
	log.Println("receiving", rawMessage)

	// Remove newline delimeter.
	jsonMessage := rawMessage[0 : len(rawMessage)-1]

	var decodedMessage interface{}
	decoder := json.NewDecoder(strings.NewReader(jsonMessage))
	decoder.Decode(&decodedMessage)

	return decodedMessage, nil
}

func writeJSON(conn net.Conn, jsonObject interface{}) error {
	// Set deadline for 3 seconds from now.
	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	log.Println("sending", jsonObject)
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
