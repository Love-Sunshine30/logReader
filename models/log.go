package models

import (
	"fmt"
	"time"
)

// Upload hold file uploads data
type Upload struct {
	ID       int       `json:"id"`
	UploadId string    `json:"upload_id"`
	FileName string    `json:"filename"`
	FileSize int64     `json:"file_size"`
	UploadAt time.Time `json:"uploaded_at"`
	Status   string    `json:"status"`
}

// LogEntry holds log entry data
type LogEntry struct {
	ID        int       `json:"id"`
	UploadID  string    `json:"upload_id"`
	TimeStamp string    `json:"timestamp"`
	Level     string    `json:"level"`
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// updates status for uploaded file after reading the file
func UpdateUploadStatus(uploadid, status string) error {
	_, err := DB.Exec(`UPDATE uploads SET status=$2 WHERE upload_id=$1`, uploadid, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %v", err)
	}
	return nil
}

func InsertUpload(upload_id, filename string, filesize int64) error {
	query := `
	INSERT INTO uploads(upload_id, filename, file_size, status)
	VALUES($1, $2, $3, 'processing')
	RETURNING id
	`
	var id int
	err := DB.QueryRow(query, upload_id, filename, filesize).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to insert upload: %v", err)
	}
	fmt.Println("upload inserted sucessfully!")
	return nil
}

func InsertLogEntry(upload_id, timestamp, level, service, message string) error {
	query := `
	INSERT INTO log_entries(upload_id, timestamp, level, service, message)
	VALUES($1, $2, $3, $4, $5) 
	`
	_, err := DB.Exec(query, upload_id, timestamp, level, service, message)
	if err != nil {
		return fmt.Errorf("failed to insert log entries")
	}
	return nil
}
