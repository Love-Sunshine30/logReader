package models

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Setup test database connection
	connStr := "postgres://postgres:password123@localhost:5432/log_reader?sslmode=disable"
	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Wire the package-global DB used by model functions
	DB = testDB

	// Test the connection
	err = testDB.Ping()
	if err != nil {
		panic("Failed to ping test database: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	_ = testDB.Close()
	os.Exit(code)
}

func TestInsertUpload(t *testing.T) {
	// Test data
	uploadID := "test_upload_123"
	filename := "test.log"
	fileSize := int64(1024)

	// Test insert
	err := InsertUpload(uploadID, filename, fileSize)
	if err != nil {
		t.Errorf("InsertUpload failed: %v", err)
	}

	// Verify the insert by querying the database
	var id int
	var status string
	err = testDB.QueryRow("SELECT id, status FROM uploads WHERE upload_id = $1", uploadID).Scan(&id, &status)
	if err != nil {
		t.Errorf("Failed to query inserted upload: %v", err)
	}

	if status != "processing" {
		t.Errorf("Expected status 'processing', got '%s'", status)
	}

	// Cleanup
	testDB.Exec("DELETE FROM uploads WHERE upload_id = $1", uploadID)
}

func TestInsertLogEntry(t *testing.T) {
	// Test data
	uploadID := "test_log_upload_456"
	timestamp := "2025-01-01T12:00:00Z"
	level := "ERROR"
	service := "test-service"
	message := "Test error message"

	// First insert an upload record
	err := InsertUpload(uploadID, "test.log", 1024)
	if err != nil {
		t.Fatalf("Failed to insert upload for log entry test: %v", err)
	}

	// Test insert log entry
	err = InsertLogEntry(uploadID, timestamp, level, service, message)
	if err != nil {
		t.Errorf("InsertLogEntry failed: %v", err)
	}

	// Verify the insert
	var insertedMessage string
	err = testDB.QueryRow("SELECT message FROM log_entries WHERE upload_id = $1", uploadID).Scan(&insertedMessage)
	if err != nil {
		t.Errorf("Failed to query inserted log entry: %v", err)
	}

	if insertedMessage != message {
		t.Errorf("Expected message '%s', got '%s'", message, insertedMessage)
	}

	// Cleanup
	testDB.Exec("DELETE FROM log_entries WHERE upload_id = $1", uploadID)
	testDB.Exec("DELETE FROM uploads WHERE upload_id = $1", uploadID)
}

func TestUpdateUploadStatus(t *testing.T) {
	// Test data
	uploadID := "test_status_update_789"
	filename := "test.log"
	fileSize := int64(1024)

	// Insert an upload record
	err := InsertUpload(uploadID, filename, fileSize)
	if err != nil {
		t.Fatalf("Failed to insert upload for status update test: %v", err)
	}

	// Test status update
	newStatus := "completed"
	err = UpdateUploadStatus(uploadID, newStatus)
	if err != nil {
		t.Errorf("UpdateUploadStatus failed: %v", err)
	}

	// Verify the update
	var status string
	err = testDB.QueryRow("SELECT status FROM uploads WHERE upload_id = $1", uploadID).Scan(&status)
	if err != nil {
		t.Errorf("Failed to query updated status: %v", err)
	}

	if status != newStatus {
		t.Errorf("Expected status '%s', got '%s'", newStatus, status)
	}

	// Cleanup
	testDB.Exec("DELETE FROM uploads WHERE upload_id = $1", uploadID)
}

func TestUploadStruct(t *testing.T) {
	upload := Upload{
		ID:       1,
		UploadId: "test_123",
		FileName: "test.log",
		FileSize: 1024,
		UploadAt: time.Now(),
		Status:   "processing",
	}

	if upload.ID != 1 {
		t.Errorf("Expected ID 1, got %d", upload.ID)
	}
	if upload.UploadId != "test_123" {
		t.Errorf("Expected UploadId 'test_123', got '%s'", upload.UploadId)
	}
	if upload.FileName != "test.log" {
		t.Errorf("Expected FileName 'test.log', got '%s'", upload.FileName)
	}
}

func TestLogEntryStruct(t *testing.T) {
	logEntry := LogEntry{
		ID:        1,
		UploadID:  "test_123",
		TimeStamp: "2025-01-01T12:00:00Z",
		Level:     "ERROR",
		Service:   "test-service",
		Message:   "Test message",
		CreatedAt: time.Now(),
	}

	if logEntry.ID != 1 {
		t.Errorf("Expected ID 1, got %d", logEntry.ID)
	}
	if logEntry.UploadID != "test_123" {
		t.Errorf("Expected UploadID 'test_123', got '%s'", logEntry.UploadID)
	}
	if logEntry.Level != "ERROR" {
		t.Errorf("Expected Level 'ERROR', got '%s'", logEntry.Level)
	}
}
