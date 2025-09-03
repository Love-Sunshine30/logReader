package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealth(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Health)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect
	var response HealthStatus
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	// Check required fields
	if response.Status == "" {
		t.Error("handler returned empty status")
	}
	if response.Service == "" {
		t.Error("handler returned empty service")
	}
	if response.TimeStamp == "" {
		t.Error("handler returned empty timestamp")
	}

	// Validate timestamp format
	_, err = time.Parse(time.RFC1123, response.TimeStamp)
	if err != nil {
		t.Errorf("handler returned invalid timestamp format: %v", response.TimeStamp)
	}
}

func TestHealthResponseStructure(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Health)
	handler.ServeHTTP(rr, req)

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "application/json")
	}

	// Check response is valid JSON
	var response HealthStatus
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}
}
