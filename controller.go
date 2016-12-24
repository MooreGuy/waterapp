package main

import (
	"encoding/json"
	"github.com/MooreGuy/waterapp/device"
	"github.com/MooreGuy/waterapp/network"
	"github.com/gocql/gocql"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

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

func StartController() {
	website := Website{}
	log.Println("Starting controller web management console.")
	go http.ListenAndServe(":8081", website)

	log.Println("Connecting to aggregator.")
	outgoingAggregations := make(chan network.Message, 100)
	incomingAggregations := make(chan network.Message, 100)
	go ConnectToAggregator(outgoingAggregations, incomingAggregations)

	go device.ReadSensors(outgoingAggregations)

	outgoingControl := make(chan network.Message, 100)
	incomingControl := make(chan network.Message, 100)
	go ConnectToExternalController(outgoingControl, incomingControl)
	go network.ListenForConnections(":6667", outgoingControl, incomingControl)

	commandQueue := device.HandleDeviceSignal()
	go handleControlMessage(incomingControl, outgoingControl, commandQueue)
	//GetDeviceInfo()
	//devices := device.FindDevices()
	//log.Println(len(devices))
}

func ConnectToAggregator(outgoing chan network.Message, incoming chan network.Message) {
	aggregatorAddress := "138.68.46.103"
	aggregatorPort := "6666"
	conn, err := net.Dial("tcp", aggregatorAddress+":"+aggregatorPort)
	for err != nil {
		time.Sleep(2 * time.Second)
		log.Println("Retrying connection to aggregator.")
		conn, err = net.Dial("tcp", aggregatorAddress+":"+aggregatorPort)
	}
	log.Println("Connected to aggregator.")

	go network.Outgoing(conn, outgoing)
	go network.Reading(conn, incoming)

	return
}

// Starts a listening routine that listens for control requests.
func ConnectToExternalController(outgoing chan network.Message, incoming chan network.Message) {
	extConAddress := "138.68.46.103"
	extConPort := "6667"
	conn, err := net.Dial("tcp", extConAddress+":"+extConPort)
	for err != nil {
		time.Sleep(2 * time.Second)
		log.Println("Retrying connection to external controller.")
		conn, err = net.Dial("tcp", extConAddress+":"+extConAddress)
	}

	go network.Outgoing(conn, outgoing)
	go network.Reading(conn, incoming)
	log.Println("Connected to external controller.")

	commandQueue := device.HandleDeviceSignal()
	go handleControlMessage(outgoing, incoming, commandQueue)
}

// TODO: Use custom handlers to handle control, and aggregation.
func handleControlMessage(incomingControl chan network.Message, outgoingControl chan network.Message, commandQueue chan device.Command) {
	for {
		mes := <-incomingControl
		signal, ok := mes["signal"].(string)
		if !ok {
			log.Println("Missing signal")
			continue
		}

		if signal == "valve-turn" {
			deviceid, ok := mes["deviceid"].(string)
			if !ok {
				log.Println("Missing deviceid")
				continue
			}

			data, ok := mes["data"].(float64)
			if !ok {
				log.Println("Missing data")
				continue
			}
			log.Println("Control signal type: ", signal)
			log.Println("deviceid: ", deviceid)
			log.Println("Data payload: ", data)

			uuid, err := gocql.ParseUUID(deviceid)
			if err != nil {
				log.Println("Invalid id")
				continue
			}

			command := device.Command{Target: uuid, Data: int(data), Name: signal}
			commandQueue <- command
		} else if signal == "list-devices" {
			devices := device.GetDevices()
			outgoingControl <- network.Message{"devices": devices}
		}

		//_, deviceOnThisController := devices[deviceid]
		/*
			if signal == "valve-turn" && deviceOnThisController {
				log.Println("Sending valve turn to valve controller.")
			}
		*/
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
