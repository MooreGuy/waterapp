package main

import (
	"fmt"
	"github.com/MooreGuy/waterapp/network"
	"log"
	"net"
)

func StartAggregator() {
	fmt.Println("Starting aggregator socket server.")
	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Handle incoming connections.")
	go ListenAggregatorServer(listener)
}

func ListenAggregatorServer(listener net.Listener) {
	fmt.Println("Listening.")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		defer conn.Close()

		log.Println("Got connection.")

		// Sets up network communication channels.
		outgoingChan := make(chan network.Message, 100)
		go network.Outgoing(conn, outgoingChan)
		incomingChan := make(chan network.Message, 100)
		go network.Reading(conn, incomingChan)

		// Handles incoming messages.
		go aggregatorHandleIncoming(incomingChan)
	}
}

func aggregatorHandleIncoming(incoming chan network.Message) {
	for {
		currentMessage := <-incoming
		log.Println("Aggregator handling message", currentMessage)
	}
}
