package main

import (
	"fmt"
	"net/http"
)

type Website struct {
}

func (site Website) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	heartbeats, err := GetAllHeartbeats()
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Welcome to waterapp. Listing all heartbeats:\n")
	for id, data := range heartbeats {
		fmt.Fprintf(w, "ID: %s  Data: %d\n", id, data)
	}
}
