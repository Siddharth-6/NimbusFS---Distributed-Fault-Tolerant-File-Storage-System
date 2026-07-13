package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/user/distributed-file-management-system/internal/service"
)

type UploadHandler struct {
	FileService *service.FileService
}

func NewUploadHandler(fs *service.FileService) *UploadHandler {
	return &UploadHandler{
		FileService: fs,
	}
}

// POST /upload
func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	meta, err := h.FileService.Upload(file, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(meta)
}

// GET /files
func (h *UploadHandler) List(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := h.FileService.ListFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (h *UploadHandler) Download(w http.ResponseWriter, r *http.Request) {

	// Only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL format:
	// /download?id=<uuid>
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing file id", http.StatusBadRequest)
		return
	}

	reader, meta, err := h.FileService.Download(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer reader.Close()

	// Tell browser what we're sending
	w.Header().Set("Content-Type", meta.ContentType)

	// Force browser to download using original filename
	w.Header().Set(
		"Content-Disposition",
		`attachment; filename="`+meta.OriginalName+`"`,
	)

	// Stream file directly to client
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
