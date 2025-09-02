package handlers

import (
	"encoding/json"
	"fmt"

	//"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/love-sunshine30/logReader/models"
)

type UploadResponse struct {
	Status   string `json:"status"`
	FileName string `json:"filename"`
	FileSize int64  `json:"filesize"`
	UploadId string `json:"uploadId"`
}

// generates id for log lifes
func generateId() string {
	rand.NewSource(time.Now().UnixNano())
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

	// generate upload id, get filename and filesize from metadata
	uploadid := generateId()
	filename := metadata.Filename
	filesize := metadata.Size

	err = models.InsertUpload(uploadid, filename, filesize)
	if err != nil {
		http.Error(w, "insert database error", 500)
	}
	response := UploadResponse{
		Status:   "successful",
		FileName: filename,
		FileSize: filesize,
		UploadId: uploadid,
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
