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

//	cluster.Authenticator = gocql.PasswordAuthenticator{
//		Username: "YWTfzyivjkJufvUzSRNKiNvZ",
//		Password: "bMnajGxe6ZGj2nPycO6d4m8MS+4FMPtsQc831uJ02zJoHHSq0pGuMO_52kuZSA5rv8gY-.e8DxiSghOf60Zca2ME-JS.0--z_2imZL5tFrnGZ8LqZn+aO5ZCFjxJzhln",
//	}
func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved GET reuqest from %s for %s", r.RemoteAddr, r.URL.Path)
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Keyspace = "census_data"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Execute the query
	query := "SELECT * FROM census_data WHERE query_text_values['time'] = '1841' LIMIT 5 ALLOW FILTERING"
	iter := session.Query(query).Iter()

	// Slice to hold the results
	var results []CensusData

	// Map to hold the row data
	var row map[string]interface{}

	// Iterate through the results
	for iter.MapScan(row) {
		// Convert the row to CensusData and append to results

		if !iter.MapScan(row) {
			break
		}
		data := CensusData{
			Key:             row["key"].(string),
			ArrayContains:   row["array_contains"].(map[string]string),
			DocJSON:         row["doc_json"].(map[string]string),
			ExistKeys:       row["exist_keys"].([]string),
			QueryTextValues: row["query_text_values"].(map[string]string),
			TxID:            row["tx_id"].(string),
		}
		results = append(results, data)
	}

	// Check for errors in iteration
	if err := iter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the results to JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
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
