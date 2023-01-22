package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"github.com/olegelantsev/linkshortener-go/shortener"
	"github.com/olegelantsev/linkshortener-go/store"
)

var urlStore store.UrlStore
var baseUrl string

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL shortener website\n")
}

var harvester *telemetry.Harvester

func GetShortLink(w http.ResponseWriter, r *http.Request) {
	const shortPath = "/x/"
	if strings.HasPrefix(r.URL.Path, shortPath) {
		slug := r.URL.Path[len(shortPath):]
		fullUrl, err := urlStore.GetUrl(slug)
		if err == nil && fullUrl != "" {
			http.Redirect(w, r, fullUrl, 301)
		} else {
			io.WriteString(w, "Not found!")
			w.WriteHeader(404)
		}
	}
}

func AddShortLink(w http.ResponseWriter, r *http.Request) {
	var createUrlRequest shortener.CreateUrlRequest
	err := json.NewDecoder(r.Body).Decode(&createUrlRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortener := shortener.NewLinkShortener(baseUrl, func(s string) (bool, error) {
		fullUrl, err := urlStore.GetUrl(s)
		if err != nil {
			return true, err
		}
		return len(fullUrl) > 0, nil
	})

	shortUrl, err := shortener.Shorten(createUrlRequest.Url)

	w.Write([]byte(fmt.Sprintf("{\"URL\": \"%s\"}", shortUrl)))
}

func Init(baseUrlParam string) {
	baseUrl = baseUrlParam
	urlStore = store.NewUrlStore()
}

func main() {
	log.Print("Starting")
	var err error
	harvester, err = telemetry.NewHarvester(telemetry.ConfigAPIKey(os.Getenv("NEW_RELIC_LICENSE_KEY")))
	if err != nil {
		fmt.Println(err)
	}
	Init(os.Getenv("BASE_URL_PATH"))
	mux := http.NewServeMux()
	mux.HandleFunc("/x/", GetShortLink)
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/links", AddShortLink)

	err = http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
