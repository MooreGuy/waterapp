package main

import (
	"encoding/json"
	"github.com/MooreGuy/waterapp/network"
	"io/ioutil"
	"log"
	"net/http"
)

type relayAPI struct {
	commands chan RelayCommand
}

func (this relayAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.EscapedPath()[len("/api/"):]
	method := r.Method
	log.Println("Got api request.", endpoint)

	if endpoint == "device/command" {
		if method != "POST" {
			http.NotFound(w, r)
			return
		}

		this.relayRequest(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (this relayAPI) relayRequest(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read response of device/command")
		log.Fatal(err)
	}

	commandMap := map[string]string{}
	err = json.Unmarshal(bodyBytes, &commandMap)
	if err != nil {
		log.Println("Failed to parse command.")
		return
	}

	data, ok := commandMap["data"]
	if !ok {
		log.Println("No data in command Message")
		return
	}

	command, ok := commandMap["command"]
	if !ok {
		log.Println("No command in command Message")
		return
	}
	if command != "valve-switch" {
		log.Println("Invalid command, only valve-turn is supported.")
		return
	}

	deviceidString, ok := commandMap["deviceid"]
	if !ok {
		log.Println("No deviceid in command Message")
		return
	}
	if !validDeviceid(deviceidString) {
		log.Println("Invalid deviceid")
		return
	}

	commandMessage := network.Message{
		"signal":   command,
		"deviceid": deviceidString,
		"data":     data,
	}
	relayCommand := RelayCommand{
		connectionid: 1,
		conn:         network.NetConn{},
		commandName:  "something",
		mes:          commandMessage,
	}

	log.Println("Relaying message.", relayCommand)
	this.commands <- relayCommand
}

// TODO: Actually confirm a deviceid is valid to operate on for a given user
//       and deviceid.
// TODO: Use userid in this confirmation.
func validDeviceid(stringDeviceid string) bool {
	return true
}

// TODO: Relay requests should be authenticated.
func StartRelay() {
	//	in := make(chan network.Message, 100)
	//	out := make(chan network.Message, 100)
	//	log.Println("Starting relay socket server")
	//	go network.RelayRouter(":*8081", in, out)
	relayCommands := make(chan RelayCommand, 10)
	log.Println("Starting relay router.")
	go RelayRouter(relayCommands)
	log.Println("Starting relay http api server")
	go http.ListenAndServe(":8080", relayAPI{relayCommands})
}

// Command to update a connection.
// connectionid uniquely identifies a connection.
// conn are channels for the connection
// command is what we should do to the identified connection
// Message is a message that should be sent out over the connection
type RelayCommand struct {
	connectionid int
	conn         network.NetConn
	commandName  string
	mes          network.Message
}

// Takes in relay requests then routes them to the
func RelayRouter(relays chan RelayCommand) {
	connPool := map[int]network.NetConn{}
	for {
		relay := <-relays
		switch relay.commandName {
		case "add":
			connPool[relay.connectionid] = relay.conn
			break
		case "remove":
			// TODO: Implement connection removal.
			break
		case "route":
			connPool[relay.connectionid].Outgoing <- relay.mes
		}
	}
}
