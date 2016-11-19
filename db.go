package main

import (
	"fmt"
	"github.com/gocql/gocql"
)

func Testdb() (*gocql.Session, error) {
	// connect to the cluster
	cluster := gocql.NewCluster("138.197.192.41")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "gmoore", Password: "meatandpotatoes"}
	session, err := cluster.CreateSession()
	return session, err
}

func RecordHeartbeat(data float64) error {
	session, err := Testdb()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Query(`INSERT INTO heartbeat (id, data) values(?, ?)`,
		gocql.TimeUUID(), int(data)).Exec()

}

func GetAllHeartbeats() (heartbeats map[gocql.UUID]int, err error) {
	heartbeats = map[gocql.UUID]int{}
	session, err := Testdb()
	if err != nil {
		return
	}
	defer session.Close()

	var id gocql.UUID
	var data int

	iter := session.Query(`SELECT id, data FROM heartbeat`).Iter()
	for iter.Scan(&id, &data) {
		heartbeats[id] = data
		fmt.Println("Heartbeat: ", id, data)
	}
	err = iter.Close()
	return heartbeats, err
}
