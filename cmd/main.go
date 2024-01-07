package main

import (
	"github.com/dazfuller/azure-functions-golang-demo/givenergy"
	"log"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("GivEnergyApiKey")
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	givEnergyManager := givenergy.GivEnergyManager{
		ApiKey: apiKey,
	}

	http.HandleFunc("/api/MeterData", givEnergyManager.GivEnergyHandler)
	log.Printf("Internal go server about to list on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
