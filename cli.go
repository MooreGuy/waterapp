package main

import (
	"fmt"
	"github.com/MooreGuy/waterapp/device"
	"github.com/MooreGuy/waterapp/network"
	"log"
	"net"
	"strconv"
)

func StartCLIShell(username string, password string) {
	fmt.Println("Connecting to server daemon")
	conn, err := net.Dial("tcp", "localhost:6667")
	if err != nil {
		log.Fatal("Error opening connection.", err.Error())
	}
	defer conn.Close()

	fmt.Println("Connected.")

	messageChannel := make(chan network.Message, 10)
	go network.Reading(conn, messageChannel)
	go network.Outgoing(conn, messageChannel)

	var data int = -1
	for {
		fmt.Print("all-heartbeats or enter a number to send to the server: ")
		var input string
		_, err = fmt.Scanln(&input)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if input == "all-heartbeats" {
			GetAllHeartbeats()
			continue
		}

		if input == "list-devices" {
			fmt.Println(device.GetFakeDevices())
			continue
		}

		data, err = strconv.Atoi(input)
		for err != nil {
			log.Println("Bad input, enter number.")
			data, err = strconv.Atoi(input)
		}

		var dev device.Device
		for _, dev = range device.GetFakeDevices() {
			break
		}

		message := network.Message{"deviceid": dev.UUID(), "signal": "valve-turn", "data": data}
		messageChannel <- message
	}
}
