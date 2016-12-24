package main

import (
	"log"
)

type relayAPI struct {
	outgoing chan network.Message
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

		this.handleValveSwitch(w, r)
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

	log.Println("Relaying message.", commandMessage)
	this.outgoing <- commandMessage
}

// TODO: Actually confirm a deviceid is valid to operate on for a given user
//       and deviceid.
// TODO: Use userid in this confirmation.
func validDeviceid(stringDeviceid string) bool {
	return true
}

// TODO: This should be renamed to "relay" since all it is doing is checking to
//       make sure that they are valid messages going to valid controllers from
//       valid clients from valid users.
func StartRelay() {
	outgoingControl := make(chan network.Message, 100)
	incomingControl := make(chan network.Message, 100)
	log.Println("Starting master http api server")
	go http.ListenAndServe(":8080", relayAPI{outgoingControl})

	log.Println("Starting master controller socket server")
	go ListenForConnections(outgoingControl, incomingControl)
}

func HandleMasterControllerMessages(incoming chan network.Message, outgoing chan network.Message) {
	for {
		mes := <-incoming
		signal, ok := mes["signal"].(string)
		if !ok {
			log.Println("Missing signal")
			continue
		}

		id, ok := mes["deviceid"].(string)
		if !ok {
			log.Println("Missing deviceid")
			continue
		}

		if signal == "valve-turn" {
			log.Println("Turing valve for " + id)
		}
	}
}
