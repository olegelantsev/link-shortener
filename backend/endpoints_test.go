package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetShortLink_NotFound(t *testing.T) {
	Init("https://example.com")
	req, err := http.NewRequest("GET", "/x/abcdef", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetShortLink)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `Not found!`
	body := rr.Body.String()
	if body != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetShortLink_CreateShortLink(t *testing.T) {
	Init("https://example.com")
	var jsonData = []byte(`{
		"URL": "https://example.com/full-url",
		"Title": "TestUrl"
	}`)

	req, err := http.NewRequest("POST", "/links", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddShortLink)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"URL": "https://example.com/bOvh7"}`
	body := rr.Body.String()
	if body != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
