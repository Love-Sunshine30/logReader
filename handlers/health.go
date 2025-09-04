package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/love-sunshine30/logReader/models"
)

type HealthStatus struct {
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
	ErrorCount int       `json:"error_count"`
	ErrorRate  string    `json:"error_rate"`
}

// response variable for /health handler

func Health(w http.ResponseWriter, r *http.Request) {
	var response HealthStatus
	response.getDataFromDB()
	jsdt, err := json.Marshal(response)
	if err != nil {
		fmt.Print("Json marshaling error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsdt)

}
func (h *HealthStatus) getDataFromDB() {
	query1 := `SELECT uploaded_at, status FROM uploads ORDER BY uploaded_at DESC LIMIT 1;`

	var up time.Time
	var status string
	err := models.DB.QueryRow(query1).Scan(&up, &status)
	if err != nil {
		fmt.Println("DB query issue")
	}

	query2 := `SELECT COUNT(*) FILTER(WHERE level = 'ERROR'),
				COUNT(*)  FROM log_entries;`
	var failed, total int
	err = models.DB.QueryRow(query2).Scan(&failed, &total)
	if err != nil {
		fmt.Println("Query issue")
	}

	h.ErrorCount = failed
	h.ErrorRate = fmt.Sprintf("%.2f%%", float32(failed)/float32(total)*100)
	h.UploadedAt = up
	h.Status = status
}
