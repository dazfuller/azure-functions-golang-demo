package main

import (
	"github.com/dazfuller/azure-functions-golang-demo/givenergy"
	"log"
	"net/http"
	"os"
	"strings"
)

func caseInsensitiveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiKey := os.Getenv("GivEnergyApiKey")
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	givEnergyManager := givenergy.GivEnergyManager{
		ApiKey: apiKey,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/meterdata", givEnergyManager.GivEnergyHandler)
	log.Printf("Internal go server about to list on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, caseInsensitiveMiddleware(mux)))
}
