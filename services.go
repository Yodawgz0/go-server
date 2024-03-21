package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
)

type Year struct {
	Year string `json:"year"`
}

type CensusData struct {
	Key                  string            `json:"key"`
	ArrayContains        map[string]string `json:"array_contains"`
	ArraySize            *int              `json:"array_size,omitempty"`
	DocJSON              map[string]string `json:"doc_json"`
	ExistKeys            []string          `json:"exist_keys"`
	QueryBoolValues      *bool             `json:"query_bool_values,omitempty"`
	QueryDblValues       *float64          `json:"query_dbl_values,omitempty"`
	QueryNullValues      *bool             `json:"query_null_values,omitempty"`
	QueryTextValues      map[string]string `json:"query_text_values"`
	QueryTimestampValues *string           `json:"query_timestamp_values,omitempty"`
	TxID                 string            `json:"tx_id"`
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved GET reuqest from %s for %s", r.RemoteAddr, r.URL.Path)
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	// cluster.Authenticator = gocql.PasswordAuthenticator{
	// 	Username: "ashley",
	// 	Password: "bazzi",
	// }
	cluster.Keyspace = "datapop"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	iter := session.Query("SELECT * FROM census WHERE time = 2100 ALLOW FILTERING").Iter()
	var censusData CensusData
	censusRecords := []CensusData{}
	for iter.Scan(&censusData.ID, &censusData.Geo, &censusData.Name, &censusData.Population, &censusData.Time) {
		censusRecords = append(censusRecords, censusData)
	}
	jsonData, err := json.Marshal(censusRecords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)

}
func handleYearFilter(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received POST request from %s for %s", r.RemoteAddr, r.URL.Path)
	var Year Year
	err := json.NewDecoder(r.Body).Decode(&Year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var year int
	year, err = strconv.Atoi(Year.Year)
	if err != nil {
		panic(err)
	}
	log.Printf("Data received is year: %d", year)
	responseMessage := fmt.Sprintf("Data received is year: %d", year)
	fmt.Print(w, responseMessage)

}
