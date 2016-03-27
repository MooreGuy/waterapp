package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	client = "client"
	server = "server"

	signalResponse       = "response"
	signalConnect        = "connect"
	signalConnectConfirm = "connect_confirm"

	signalFlow          = "flow"
	signalFlowResponse  = "flow_response"
	signalLevel         = "level"
	signalLevelResponse = "level_response"

	signalHeartbeat         = "heartbeat"
	signalHeartbeatResponse = "heartbeat_response"
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

		serverReadMessage(conn)
		sendHeartbeat(conn)
	}
}

func serverReadMessage(conn net.Conn) {
	message, ok := readJSON(conn).(map[string]interface{})
	if !ok {
		log.Println("Skipping")
		return
		log.Fatal("Could not assert message to map[string]interface{}")
	}

	signal, ok := message["signal"].(string)
	if !ok {
		log.Fatal("Could assert signal type to string.")
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
}

func sendHeartbeat(conn net.Conn) {
	data := 123
	message := Message{"signal": signalHeartbeat, "data": data}
	if err := writeJSON(conn, message); err != nil {
		log.Fatal(err.Error())
	}
}

func heartbeatResponse(conn net.Conn) {
	data := 123
	message := Message{"signal": signalHeartbeatResponse, "data": data}
	if err := writeJSON(conn, message); err != nil {
		log.Fatal(err.Error())
	}
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
	conn, err := net.Dial("tcp", "home.guywmoore.com:8080")
	if err != nil {
		log.Fatal("Error opening connection.", err.Error())
	}
	defer conn.Close()

	connectRequest(conn)

	for {
		clientReadMessage(conn)

		time.Sleep(time.Second)
	}
}

func clientReadMessage(conn net.Conn) {
	message, ok := readJSON(conn).(map[string]interface{})
	if !ok {
		log.Println("Skipping")
		return
		log.Fatal("Could not assert message to map[string]interface{}")
	}

	signal, ok := message["signal"].(string)
	if !ok {
		log.Fatal("Could assert signal type to string.")
	}

	if signal == signalFlow {
		log.Print("Flow request received")
	} else if signal == signalLevel {
		log.Print("Level request received")
	} else if signal == signalHeartbeat {
		log.Print("Heartbeat received")
		heartbeatResponse(conn)
	} else if signal == signalConnectConfirm {
		log.Println("Successfully connected to master server.")
	} else {
		log.Print("unknown signal")
	}
}

func readJSON(conn net.Conn) interface{} {
	// Set deadline for 3 seconds from now.
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	rawMessage, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("receiving", rawMessage)

	// Remove newline delimeter.
	jsonMessage := rawMessage[0 : len(rawMessage)-1]

	var decodedMessage interface{}
	decoder := json.NewDecoder(strings.NewReader(jsonMessage))
	decoder.Decode(&decodedMessage)

	return decodedMessage
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
