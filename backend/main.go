package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"github.com/olegelantsev/linkshortener-go/store"
)

var urlStore store.UrlStore

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL shortener website\n")
}

var harvester *telemetry.Harvester

func getShortLink(w http.ResponseWriter, r *http.Request) {
	const shortPath = "/x/"
	if strings.HasPrefix(r.URL.Path, shortPath) {
		slug := r.URL.Path[len(shortPath):]
		fullUrl, err := urlStore.GetUrl(slug)
		if err == nil && fullUrl != "" {
			http.Redirect(w, r, fullUrl, 301)
		} else {
			io.WriteString(w, "Not found!\n")
			w.WriteHeader(404)
		}
	}
}

func main() {
	log.Print("Starting")
	var err error
	harvester, err = telemetry.NewHarvester(telemetry.ConfigAPIKey(os.Getenv("NEW_RELIC_LICENSE_KEY")))
	if err != nil {
		fmt.Println(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/x/", getShortLink)
	mux.HandleFunc("/", getRoot)
	urlStore = store.NewUrlStore()

	err = http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
