package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func StartCLIShell() {
	fmt.Println("Connecting to server daemon")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error opening connection.", err.Error())
	}
	defer conn.Close()

	fmt.Println("Connected.")

	messageChannel := make(chan Message, 10)
	go handleConnection(conn)
	go outgoingChannel(conn, messageChannel)

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

		data, err = strconv.Atoi(input)
		if err != nil {
			log.Println("Bad input, enter number.")
			data = -1
		}

		message := Message{"signal": signalHeartbeat, "data": data}
		messageChannel <- message
	}
}
