package network

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"
)

type NetConn struct {
	incoming chan Message
	outgoing chan Message
}

func ListenForConnections(address string, incoming chan Message, outgoing chan Message) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go Outgoing(conn, outgoing)
		go Reading(conn, incoming)
		log.Println("new connection")
	}
}

func Outgoing(conn net.Conn, outgoing chan Message) {
	for {
		currentMessage := <-outgoing

		err := writeJSON(conn, currentMessage)
		if err != nil {
			log.Println("ERROR WRITING JSON", err.Error())
		}
	}
}

// Reads from the socket connection and puts read messages through json decoding
// into messages, and then puts them in the incoming message channel.
func Reading(conn net.Conn, incoming chan Message) {
	for {
		// Set deadline for 3 seconds from now.
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		rawMessage, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			opErr, ok := err.(net.Error)
			if ok {
				if !opErr.Temporary() {
					return
				} else {
					continue
				}

			}

			if err == io.EOF {
				return
			}
		}

		// Remove newline delimeter.
		mes, err := JSONMessage(rawMessage[0 : len(rawMessage)-1])
		if err != nil {
			log.Println(err)
		}

		incoming <- mes
	}
}

// Starts a socket server that listens for incoming connections and manages
// them based on commands.
// Receives a channel for all messages that should be routed, and then another
// channel to
func SocketServer(address string, in chan Message, out chan Message) {
	connManager := make(chan connCommand)
	go ManageConns(connManager)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// Start in/out connection goroutines
		net := NetConn{make(chan Message, 100), make(chan Message, 100)}
		go Outgoing(conn, net.outgoing)
		go Reading(conn, in)

		addConn := connCommand{i, net, "add", nil}
		i = i + 1
		connManager <- addConn
		log.Println("sent add command")
	}
}


// Stores and manages the state of the relay connections
func PoolConns(commands chan ConnCommand) {
	for {
		command := <-commands
		switch command.commandName {
		case "add":
			log.Println("Added new connection")
			conns[command.connectionid] = conn
			break
		case "delete":
			log.Println("TODO, actually remove the connection.")
			break
		case "get":
			log.Println("TODO , return the pointer to this connection.")
			break
		}
	}
}
