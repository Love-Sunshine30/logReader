package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthStatus struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	TimeStamp string `json:"timestamp"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	response := HealthStatus{
		Status:    "Ok",
		Service:   "sucessful",
		TimeStamp: time.Now().Format(time.RFC1123),
	}

	jsodata, err := json.Marshal(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsodata)
}
