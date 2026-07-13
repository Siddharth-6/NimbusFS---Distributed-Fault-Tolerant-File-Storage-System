package metadata

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/user/distributed-file-management-system/internal/models"
)

type JSONStore struct {
	BasePath string
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{
		BasePath: path,
	}
}

func (j *JSONStore) Save(meta models.FileMetadata) error {

	path := filepath.Join(j.BasePath, meta.ID+".json")

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(meta)
}

func (j *JSONStore) Get(id string) (models.FileMetadata, error) {

	path := filepath.Join(j.BasePath, id+".json")

	file, err := os.Open(path)
	if err != nil {
		return models.FileMetadata{}, err
	}
	defer file.Close()

	var meta models.FileMetadata

	err = json.NewDecoder(file).Decode(&meta)

	return meta, err
}

func (j *JSONStore) Delete(id string) error {

	path := filepath.Join(j.BasePath, id+".json")

	return os.Remove(path)
}

func (j *JSONStore) List() ([]models.FileMetadata, error) {

	files, err := os.ReadDir(j.BasePath)
	if err != nil {
		return nil, err
	}

	var result []models.FileMetadata

	for _, file := range files {

		f, err := os.Open(filepath.Join(j.BasePath, file.Name()))
		if err != nil {
			continue
		}

		var meta models.FileMetadata

		if err := json.NewDecoder(f).Decode(&meta); err == nil {
			result = append(result, meta)
		}

		f.Close()
	}

	return result, nil
}
