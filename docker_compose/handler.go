package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get method returns all the coins
func GetAllCoins(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: returnAllCoins")
	output, err := findAllCoins()
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	response := []Response{}
	for _, obj := range output {
		response = append(response,
			Response{Id: obj.Id, Exchanges: obj.Exchanges, TaskRun: obj.TaskRun})
	}
	json.NewEncoder(w).Encode(response)
}

// Transform coins ETL logic
// Incase of errors in input handler should respond with 4XX
func TransformCoins(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var input []string
	json.Unmarshal(reqBody, &input)
	fmt.Println("Input coins are ", input)
	go updateCoins(input)
	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode("OK")
}

//Getall objects
func findAllCoins() ([]CoinOutput, error) {
	var coins []CoinOutput
	if result := db.Find(&coins); result.Error != nil {
		return nil, result.Error
	}
	return coins, nil
}
