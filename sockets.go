package main

import (
	"io"
	"log"
	"net"
)

func outgoing(conn net.Conn, outgoing chan Message) {
	for {
		currentMessage := <-outgoing

		err := writeJSON(conn, currentMessage)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func reading(conn net.Conn, outgoing chan Message) {
	for {
		err := readJSONMessage(conn, outgoing)
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

			log.Println(err.Error())
		}
	}
}
