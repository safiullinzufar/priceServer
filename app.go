package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"priceServer/src"
	"time"
)

const (
	RETRIES_MAX = 5
)

func main() {
	db, err := src.SetDB()
	for cntRetries := 0; cntRetries != RETRIES_MAX; cntRetries++{
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}
	if err != nil {
		log.Fatal("Cannot set up db. Reason: ", err)
	}

	file, err := ioutil.ReadFile("config.txt")
	if err != nil {
		log.Fatal("Invalid config", err)
	}
	err = json.Unmarshal(file, db.Config)
	if err != nil {
		log.Fatal("Cannot read config.json")
	}

	//fmt.Println(db.Config)

	handler := http.NewServeMux()

	handler.HandleFunc("/subscribe", db.Subscribe)
	server := &http.Server{
		Addr:         ":9090",
		Handler:      handler,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	log.Println("Starting server on port 9090")
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("Error starting server: %s", err)
		os.Exit(1)
	}
}