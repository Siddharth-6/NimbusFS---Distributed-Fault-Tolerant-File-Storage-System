package node

import (
	"io"
	"net/http"

	"github.com/user/distributed-file-management-system/internal/storage"
)

type Handler struct {
	store storage.Storage
}

func NewHandler(store storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) Chunk(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodPost:
		h.UploadChunk(w, r)

	case http.MethodGet:
		h.DownloadChunk(w, r)

	default:
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *Handler) UploadChunk(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing chunk id", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.store.SaveChunk(id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DownloadChunk(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.URL.Query().Get("id")

	reader, err := h.store.OpenChunk(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer reader.Close()

	io.Copy(w, reader)
}

func (h *Handler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}
