package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

type Year struct {
	Year string `json:"year"`
}

type CensusData struct {
	// Key      string            `json:"key"`
	QueryTextValues map[string]string `json:"query_text_values"`
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func handleYearFilterRequest(w http.ResponseWriter, r *http.Request, year string) {
	log.Printf("Received GET request from %s for %s", r.RemoteAddr, r.URL.Path)
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Keyspace = "census_data"
	Username := goDotEnvVariable("CLUSTER_USERNAME")
	Password := goDotEnvVariable("CLUSTER_PASSWORD")
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: Username,
		Password: Password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	query := fmt.Sprintf("SELECT * FROM census_data WHERE query_text_values['time'] = '%s' LIMIT 10 ALLOW FILTERING", year)
	iter := session.Query(query).Iter()
	var results []CensusData
	for {
		row := make(map[string]interface{})
		if !iter.MapScan(row) {
			break
		}
		data := CensusData{
			// Key:      row["key"].(string),
			QueryTextValues: row["query_text_values"].(map[string]string),
		}
		results = append(results, data)
	}
	if err := iter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
