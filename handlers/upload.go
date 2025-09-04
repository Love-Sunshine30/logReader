package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	//"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/love-sunshine30/logReader/models"
)

// holds the response from the upload handler
type UploadResponse struct {
	Status   string `json:"status"`
	FileName string `json:"filename"`
	FileSize int64  `json:"filesize"`
	UploadId string `json:"uploadId"`
}

// holds each logline's attributes
type LogLine struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	Message   string `json:"message"`
}

// count for how many line was valid and invalid
var Valid = 0
var Failed = 0

// helpers
// generates id for log lifes
func generateId() string {
	rand.NewSource(time.Now().UnixNano())
	return fmt.Sprintf("log_%d_%04d", time.Now().UnixNano(), rand.Intn(1000))
}

// validate log line
func (ll *LogLine) valid() bool {
	lvl := strings.ToUpper(strings.TrimSpace(ll.Level))
	if ll.Timestamp == "" || ll.Message == "" || ll.Service == "" || ll.Level == "" || (lvl != "DEBUG" && lvl != "INFO" && lvl != "ERROR" && lvl != "WARN") {
		return false
	}
	return true
}

func Upload(w http.ResponseWriter, r *http.Request) {
	// parse the multipart form
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// get the file
	file, metadata, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get the file", http.StatusBadRequest)
		return
	}
	//close the file
	defer file.Close()

	// generate upload id, get filename and filesize from metadata
	uploadid := generateId()
	filename := metadata.Filename
	filesize := metadata.Size

	err = models.InsertUpload(uploadid, filename, filesize)
	if err != nil {
		http.Error(w, "insert database error", 500)
	}

	// Scan the file
	// scanner
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	// scan each line and insert into log_entries table
	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		var ll LogLine
		if err := json.Unmarshal(line, &ll); err != nil {
			Failed++
			continue
		}

		if ll.valid() {
			err = models.InsertLogEntry(uploadid, ll.Timestamp, ll.Level, ll.Service, ll.Message)
			if err != nil {
				http.Error(w, "Error inserting log line", http.StatusInternalServerError)
			}
			Valid++
		} else {
			Failed++
		}
	}

	// if any error happens reading the file, update status to "failed"
	if err := scanner.Err(); err != nil {
		_ = models.UpdateUploadStatus(uploadid, "failed")
		http.Error(w, "file read error", http.StatusInternalServerError)
		return
	}

	if Failed > 10 && Failed > Valid/4 {
		_ = models.UpdateUploadStatus(uploadid, "failed")
	} else {
		_ = models.UpdateUploadStatus(uploadid, "completed")
	}
	response := UploadResponse{
		Status:   "successful",
		FileName: filename,
		FileSize: filesize,
		UploadId: uploadid,
	}

	jsondata, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Json marshaling error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsondata)
}
