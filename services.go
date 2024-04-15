package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
	"gopkg.in/inf.v0"
)

type Year struct {
	Year string `json:"year"`
}

type CensusData struct {
	// Key      string            `json:"key"`
	QueryTextValues map[string]string   `json:"query_text_values"`
	QueryDblValues  map[string]*inf.Dec `json:"query_dbl_values"`
}
type YearBoundaryValues struct {
	MinPopulation *inf.Dec
	MaxPopulation *inf.Dec
}

type filteredCensusData struct {
	MinPopulation *inf.Dec
	MaxPopulation *inf.Dec
	AllData       []CensusData
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
	queryMaxMin := fmt.Sprintf("SELECT MAX(query_dbl_values['Population']), MIN(query_dbl_values['Population']) from census_data WHERE query_dbl_values['time'] = %s", year)
	var minPopulation, maxPopulation *inf.Dec
	yearIter := session.Query(queryMaxMin).Iter()
	for {
		row := make(map[string]interface{})
		if !yearIter.MapScan(row) {
			break
		}
		data := YearBoundaryValues{
			MinPopulation: row["system.max(query_dbl_values['Population'])"].(*inf.Dec),
			MaxPopulation: row["system.min(query_dbl_values['Population'])"].(*inf.Dec),
		}
		minPopulation = data.MinPopulation
		maxPopulation = data.MaxPopulation
	}
	fmt.Printf("Min Year: %d, Max Year: %d\n", minPopulation, maxPopulation)
	query := fmt.Sprintf("SELECT * FROM census_data WHERE query_dbl_values['time'] = %s ALLOW FILTERING", year)
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
			QueryDblValues:  row["query_dbl_values"].(map[string]*inf.Dec),
		}
		results = append(results, data)
	}
	if err := iter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	finalData := filteredCensusData{
		MinPopulation: minPopulation,
		MaxPopulation: maxPopulation,
		AllData:       results,
	}

	fmt.Println(finalData)

	jsonData, err := json.Marshal(finalData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func handleGdpFilterRequest(w http.ResponseWriter, r *http.Request, year string) {
	log.Printf("Received GET request from %s for %s", r.RemoteAddr, r.URL.Path)
	cluster := gocql.NewCluster("localhost")
	cluster.Port = 9042
	cluster.Keyspace = "gdp_data"
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
	query := fmt.Sprintf("SELECT * FROM gdp_data WHERE query_dbl_values['time'] = %s ALLOW FILTERING", year)
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
			QueryDblValues:  row["query_dbl_values"].(map[string]*inf.Dec),
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
