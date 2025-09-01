package models

import (
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
