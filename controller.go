package main

import (
	"fmt"
	"github.com/MooreGuy/waterapp/device"
	"github.com/MooreGuy/waterapp/network"
	"github.com/gocql/gocql"
	"log"
	"net"
	"net/http"
	"time"
)

func StartController() {
	website := Website{}
	fmt.Println("Starting controller web management console.")
	go http.ListenAndServe(":8081", website)

	log.Println("Connecting to aggregator.")
	outgoingAggregations := make(chan network.Message, 100)
	incomingAggregations := make(chan network.Message, 100)
	go ConnectToAggregator(outgoingAggregations, incomingAggregations)

	go device.ReadSensors(outgoingAggregations)

	outgoingControl := make(chan network.Message, 100)
	incomingControl := make(chan network.Message, 100)
	go ConnectToExtenrnalController(outgoingControl, incomingControl)
	go ListenForInternalControl(outgoingControl, incomingControl)

	commandQueue := device.HandleDeviceSignal()
	go handleControlMessage(incomingControl, commandQueue)
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
func ConnectToExtenrnalController(outgoing chan network.Message, incoming chan network.Message) {
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
}

func ListenForInternalControl(outgoing chan network.Message, incoming chan network.Message) {
	listener, err := net.Listen("tcp", ":6667")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go network.Outgoing(conn, outgoing)
		go network.Reading(conn, incoming)
		log.Println("Internal controller connected.")
	}
}

// TODO: Use custom handlers to handle control, and aggregation.
func handleControlMessage(incomingControl chan network.Message, commandQueue chan device.Command) {
	for {
		mes := <-incomingControl
		signal, ok := mes["signal"].(string)
		if !ok {
			log.Println("Missing signal")
			continue
		}

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
