package metadata

import "github.com/user/distributed-file-management-system/internal/models"

type Store interface {
	Save(meta models.FileMetadata) error

	Get(id string) (models.FileMetadata, error)

	List() ([]models.FileMetadata, error)

	Delete(id string) error
}
