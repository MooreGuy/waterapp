package main

import (
	"fmt"
	"log"
	"net"
)

func StartAggregator() {
	fmt.Println("Starting aggregator socket server.")
	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err.Error())
	}

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

		outgoingChan := make(chan Message, 100)
		go outgoing(conn, outgoingChan)
		incomingChan := make(chan Message, 100)
		go reading(conn, incomingChan)
		go handleIncoming(incomingChan)
	}
}

func aggregatorHandleIncoming(incoming chan Message) {
	for {
		currentMessage := <-incoming
		log.Println(currentMessage)
	}
}
