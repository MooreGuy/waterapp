package network

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"
)

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
