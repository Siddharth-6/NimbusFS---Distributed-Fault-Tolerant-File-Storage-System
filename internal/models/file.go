package models

import "time"

type FileMetadata struct {
	ID           string
	OriginalName string
	Size         int64
	ContentType  string
	UploadedAt   time.Time
	Chunks       []ChunkMetadata `json:"chunks"`
}
