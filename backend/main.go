package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"github.com/olegelantsev/linkshortener-go/api"
)

var harvester *telemetry.Harvester

func main() {
	log.Print("Starting")
	var err error
	harvester, err = telemetry.NewHarvester(telemetry.ConfigAPIKey(os.Getenv("NEW_RELIC_LICENSE_KEY")))
	if err != nil {
		fmt.Println(err)
	}
	api.Init(os.Getenv("BASE_URL_PATH"))
	mux := http.NewServeMux()
	mux.HandleFunc("/x/", api.GetShortLink)
	mux.HandleFunc("/", api.GetRoot)
	mux.HandleFunc("/links", api.AddShortLink)

	err = http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
