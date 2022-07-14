// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/go-redis/redis"
// )

// func main() {
// 	fmt.Println("Hello everyone!")

// 	// Setup the connection to redis. "redis" is the name of the container which lets
// 	// docker handle the networking. "mypassword" is the password used in docker-compose.yml.
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     "redis:6379",
// 		Password: "mypassword",
// 		DB:       0,
// 	})

// 	// Set a key and value for testing.
// 	time.Sleep(5 * time.Second)
// 	err := client.Set("key", "value", 0).Err()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("key is set")
// 	fmt.Println(client.Get("key"))
// 	// For testing purposes sleep for 10 mins to keep container alive. Should serve as a web app.
// 	time.Sleep(10 * time.Minute)

// }
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Main function creates a SQL handler and listens to routes
// route output will provide json data
// route coins will take the array of coins and transform them to proper objects
// and store them sqldb

func main() {
	_db, err := gorm.Open(sqlite.Open("./output.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	db = _db

	if err := db.AutoMigrate(&CoinOutput{}); err != nil {
		panic(err)
	}
	setInitTaskId()
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/output", GetAllCoins)
	r.HandleFunc("/coins", TransformCoins).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", r))
}
