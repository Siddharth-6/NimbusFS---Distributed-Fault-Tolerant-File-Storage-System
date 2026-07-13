package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/user/distributed-file-management-system/internal/chunker"
	"github.com/user/distributed-file-management-system/internal/client"
	"github.com/user/distributed-file-management-system/internal/metadata"
	"github.com/user/distributed-file-management-system/internal/models"
	"github.com/user/distributed-file-management-system/internal/node"
)

const (
	ChunkSize         = 1024 * 1024
	ReplicationFactor = 2
)

type FileService struct {
	Metadata metadata.Store
	Nodes    *node.Manager
	Client   *client.StorageClient
}

func NewFileService(
	metadataStore metadata.Store,
	manager *node.Manager,
	client *client.StorageClient,
) *FileService {

	return &FileService{
		Metadata: metadataStore,
		Nodes:    manager,
		Client:   client,
	}
}

func (f *FileService) Upload(
	file multipart.File,
	header *multipart.FileHeader,
) (models.FileMetadata, error) {

	fileID := uuid.New().String()

	meta := models.FileMetadata{
		ID:           fileID,
		OriginalName: header.Filename,
		ContentType:  header.Header.Get("Content-Type"),
		UploadedAt:   time.Now(),
	}

	c := chunker.NewChunker(file, ChunkSize)

	for {

		chunk, err := c.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return models.FileMetadata{}, err
		}

		chunkID := fmt.Sprintf("%s_%d", fileID, chunk.Index)

		nodes := f.Nodes.ChooseNodes(
			chunk.Index,
			ReplicationFactor,
		)

		chunkMeta := models.ChunkMetadata{
			ID:    chunkID,
			Index: chunk.Index,
			Size:  int64(len(chunk.Data)),
		}

		for _, nd := range nodes {

			err := f.Client.UploadChunk(
				nd.Address,
				chunkID,
				chunk.Data,
			)

			if err != nil {
				return models.FileMetadata{}, err
			}

			chunkMeta.Locations = append(
				chunkMeta.Locations,
				models.ChunkLocation{
					NodeID:  nd.ID,
					Address: nd.Address,
				},
			)
		}

		meta.Size += int64(len(chunk.Data))

		meta.Chunks = append(
			meta.Chunks,
			chunkMeta,
		)
	}

	if err := f.Metadata.Save(meta); err != nil {
		return models.FileMetadata{}, err
	}

	return meta, nil
}

func (f *FileService) ListFiles() ([]models.FileMetadata, error) {
	return f.Metadata.List()
}

func (f *FileService) Download(
	id string,
) (io.ReadCloser, models.FileMetadata, error) {

	meta, err := f.Metadata.Get(id)
	if err != nil {
		return nil, models.FileMetadata{}, err
	}

	var readers []io.Reader
	var closers []io.Closer

	for _, chunk := range meta.Chunks {

		var (
			reader io.ReadCloser
			ok     bool
		)

		for _, location := range chunk.Locations {

			reader, err = f.Client.DownloadChunk(
				location.Address,
				chunk.ID,
			)

			if err == nil {
				readers = append(readers, reader)
				closers = append(closers, reader)

				ok = true
				break
			}
		}

		if !ok {

			// Close all previously opened readers
			for _, c := range closers {
				c.Close()
			}

			return nil,
				models.FileMetadata{},
				fmt.Errorf(
					"unable to download chunk %s from any replica",
					chunk.ID,
				)
		}
	}

	reader := io.MultiReader(readers...)

	return NewMultiReadCloser(
		reader,
		closers,
	), meta, nil
}
