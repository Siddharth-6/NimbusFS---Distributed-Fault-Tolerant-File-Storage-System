package storage

import "io"

type Storage interface {
	SaveChunk(id string, data []byte) error

	OpenChunk(id string) (io.ReadCloser, error)

	DeleteChunk(id string) error
}
