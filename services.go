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
type GDPData struct {
	// Key      string            `json:"key"`
	QueryTextValues map[string]string
	QueryDblValues  *inf.Dec
}
type YearBoundaryValues struct {
	MinPopulation *inf.Dec
	MaxPopulation *inf.Dec
}

type gdpFilterBoundaryValues struct {
	MinValue *inf.Dec
	MaxValue *inf.Dec
}

type filteredGDPData struct {
	MinValue *inf.Dec
	MaxValue *inf.Dec
	AllData  []GDPData
}

type filteredCensusData struct {
	MinValue *inf.Dec
	MaxValue *inf.Dec
	AllData  []CensusData
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
			MinPopulation: row["system.min(query_dbl_values['Population'])"].(*inf.Dec),
			MaxPopulation: row["system.max(query_dbl_values['Population'])"].(*inf.Dec),
		}
		minPopulation = data.MinPopulation
		maxPopulation = data.MaxPopulation
	}
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
		MinValue: minPopulation,
		MaxValue: maxPopulation,
		AllData:  results,
	}

	jsonData, err := json.Marshal(finalData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func handleGdpFilterRequest(w http.ResponseWriter, r *http.Request, year string, typeIncomefilter string) {
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
	queryMaxMin := fmt.Sprintf("SELECT MAX(query_dbl_values['%s']), MIN(query_dbl_values['%s']) from gdp_data WHERE query_dbl_values['time'] = %s", typeIncomefilter, typeIncomefilter, year)
	var minValue, maxValue *inf.Dec
	yearIter := session.Query(queryMaxMin).Iter()
	for {
		row := make(map[string]interface{})
		if !yearIter.MapScan(row) {
			break
		}
		data := gdpFilterBoundaryValues{
			MinValue: row[fmt.Sprintf("system.min(query_dbl_values['%s'])", typeIncomefilter)].(*inf.Dec),
			MaxValue: row[fmt.Sprintf("system.max(query_dbl_values['%s'])", typeIncomefilter)].(*inf.Dec),
		}
		minValue = data.MinValue
		maxValue = data.MaxValue
	}

	query := fmt.Sprintf("SELECT query_dbl_values['%s'], query_text_values  FROM gdp_data WHERE query_dbl_values['time'] = %s", typeIncomefilter, year)
	iter := session.Query(query).Iter()
	var results []GDPData
	for {
		row := make(map[string]interface{})
		if !iter.MapScan(row) {
			break
		}
		data := GDPData{
			// Key:      row["key"].(string),
			QueryTextValues: row["query_text_values"].(map[string]string),
			QueryDblValues:  row[fmt.Sprintf("query_dbl_values['%s']", typeIncomefilter)].(*inf.Dec),
		}
		results = append(results, data)
	}

	finalData := filteredGDPData{
		MinValue: minValue,
		MaxValue: maxValue,
		AllData:  results,
	}
	if err := iter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(finalData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
