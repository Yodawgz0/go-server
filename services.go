package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
)

type UserData struct {
	Username string `json:"username"`
	Age      int    `json:"age"`
}

type CensusData struct {
	ID         int
	Geo        string
	Name       string
	Time       int
	Population int
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved GET reuqest from %s for %s", r.RemoteAddr, r.URL.Path)

	responseMessage := fmt.Sprintf("Hello, this is your Go server!")
	fmt.Fprint(w, responseMessage)
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "ashley",
		Password: "bazzi",
	}
	cluster.Keyspace = "world_census" // Use system_schema keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Query table metadata from system_schema.tables
	iter := session.Query("SELECT * FROM census WHERE id < 20 ALLOW FILTERING").Iter()
	var censusData CensusData
	for iter.Scan(&censusData.ID, &censusData.Geo, &censusData.Name, &censusData.Time, &censusData.Population) {
		fmt.Println(censusData)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

}
func handleReadGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved GET reuqest from %s for %s", r.RemoteAddr, r.URL.Path)

	// Create a cluster configuration and a session using the keyspace
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Keyspace = "world_census"
	// other options
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	err = session.Query("describe census").Exec()

	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	responseMessage := fmt.Sprintf("Reading of the table name is done!")
	fmt.Fprint(w, responseMessage)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received POST request from %s for %s", r.RemoteAddr, r.URL.Path)
	var userData UserData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Data received: Username - %s, Age - %d", userData.Username, userData.Age)

	// Send a response to the client
	responseMessage := fmt.Sprintf("Data received is Username: %s and Age: %d",
		userData.Username, userData.Age)
	fmt.Fprintf(w, responseMessage)
}
