package handlers

import (
	"encoding/json"
	"fmt"

	//"io"
	"math/rand"
	"net/http"
	"time"
)

type UploadResponse struct {
	Status   string `json:"status"`
	FileName string `json:"filename"`
	FileSize int    `json:"filesize"`
	UploadId string `json:"uploadId"`
}

// generates id for log lifes
func generateId() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("log_%d_%04d", time.Now().UnixNano(), rand.Intn(1000))
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

	response := UploadResponse{
		Status:   "successful",
		FileName: metadata.Filename,
		FileSize: int(metadata.Size),
		UploadId: generateId(),
	}

	// filebytes, err := io.ReadAll(file)
	// if err != nil {
	// 	http.Error(w, "Unable to read file", http.StatusInternalServerError)
	// }

	jsondata, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Json marshaling error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsondata)
}
