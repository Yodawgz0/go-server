package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func main() {
	// Connect to the Cassandra cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "System"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Create a new keyspace
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS mykeyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}").Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Switch to the new keyspace
	session.Close()
	cluster.Keyspace = "subKeySpace"
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Create a table with user and gender columns
	err = session.Query("CREATE TABLE IF NOT EXISTS user_info (user_id uuid PRIMARY KEY, user_name text, gender text)").Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Insert a row
	id := gocql.TimeUUID()
	if err := session.Query("INSERT INTO user_info (user_id, user_name, gender) VALUES (?, ?, ?)", id, "John", "Wick").Exec(); err != nil {
		log.Fatal(err)
	}

	// Select a row
	iter := session.Query("SELECT user_id, user_name, gender FROM user_info").Iter()
	var userID gocql.UUID
	var userName, gender string
	fmt.Println("User_info table contents:")
	for iter.Scan(&userID, &userName, &gender) {
		fmt.Printf("UserID: %s, UserName: %s, Gender: %s\n", userID, userName, gender)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	// Delete the table
	err = session.Query("DROP TABLE IF EXISTS user_info").Exec()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table 'user_info' deleted.")
}
