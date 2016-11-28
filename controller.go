package main

import (
	"fmt"
	"github.com/MooreGuy/waterapp/device"
	"github.com/MooreGuy/waterapp/network"
	"log"
	"net"
	"net/http"
)

func StartController() {
	website := Website{}
	fmt.Println("Starting controller web management console.")
	go http.ListenAndServe(":8081", website)

	outgoing, _ := ConnectToAggregator()
	go device.ReadSensors(outgoing)
	//GetDeviceInfo()
	//devices := device.FindDevices()
	//log.Println(len(devices))

}

func ConnectToAggregator() (outgoing chan network.Message, incoming chan network.Message) {
	log.Println("Connecting to aggregator.")
	aggregatorAddress := "138.68.46.103"
	aggregatorPort := "6666"
	conn, err := net.Dial("tcp", aggregatorAddress+":"+aggregatorPort)
	if err != nil {
		panic(err.Error())
	}

	outgoing = make(chan network.Message, 4)
	go network.Outgoing(conn, outgoing)
	incoming = make(chan network.Message, 4)
	go network.Reading(conn, incoming)
	log.Println("Connected to aggregator.")

	return
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
