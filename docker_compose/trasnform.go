package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	redis "github.com/go-redis/redis/v8"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

const APIKEY = "COINGECKO_APIKEY"

// Check if coin Id exists then ignore else
// Get data from coingecko and update the sqlDb
// Ignore the error in case we cannot resolve the coin
func updateCoins(coins []string) error {
	for _, id := range coins {
		var coinTable CoinOutput
		result := db.First(&coinTable, "id = ?", id)
		// Insert coin
		if result.Error != nil {
			resp, err := getDataFromCG(id)
			fmt.Println("Response ", resp, err)
			if err != nil {
				fmt.Printf("error reading info from CG %+v", err)
			} else {
				if len(resp.Exchanges) > 0 {
					db.Create(&resp)
				}
			}
		} else { // already present so ignore
			fmt.Println("Coin already present ", id)
		}
	}
	return nil
}

//TODO move the api key to secrets, URI to env variable.
// Data is paginaged to 100 records , need to get more data to represent it appropriately
// GraphQl would be appropriate as json payload is too much for the request, unfortunately CoinGecko doent provide it.
// TODO lot of map key checks need to be done.
func getDataFromCG(id string) (CoinOutput, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	client := retryClient.StandardClient() // *http.Client

	uri := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/tickers", id)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return CoinOutput{}, err
	}
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", APIKEY)
	resp, err := client.Do(req)
	if err != nil {
		return CoinOutput{}, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	var jsonRes map[string]interface{}
	json.Unmarshal(respBody, &jsonRes)
	if _, ok := jsonRes["tickers"]; !ok {
		return CoinOutput{}, err
	}
	vals := jsonRes["tickers"].([]interface{})
	var platforms []string
	for _, ids := range vals {
		resp := ids.(map[string]interface{})
		name := resp["market"].(map[string]interface{})
		exchange := fmt.Sprintf("%v", name["identifier"])
		if exchange != "" {
			platforms = append(platforms, exchange)
		}
	}
	//Generate taskId from Redis
	taskId := 0
	if len(platforms) > 0 {
		taskId, err = generateTaskId()
		if err != nil {
			fmt.Println("error getting taskid")
			return CoinOutput{}, err
		}
	}
	return CoinOutput{Id: id, Exchanges: platforms, TaskRun: taskId}, err
}

// TODO make this method coroutine safe
// move all the connections sessions to init ex: redis,sql etc.
// read passwords from env.
func generateTaskId() (int, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "mypassword", // no password set
		DB:       0,            // use default DB
	})
	val, err := rdb.Get(ctx, "index").Result()
	if err != nil {
		return 0, err
	}
	intVar, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	incVar := strconv.Itoa(intVar + 1)
	err = rdb.Set(ctx, "index", incVar, 0).Err()
	if err != nil {
		return 0, err
	}
	return intVar, nil
}

func setInitTaskId() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "mypassword",
		DB:       0, // use default DB
	})
	incVar := strconv.Itoa(0)
	err := rdb.Set(ctx, "index", incVar, 0).Err()
	fmt.Println("error setting default val", err)
}
