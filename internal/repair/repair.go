package repair

import (
	"io"
	"time"

	"github.com/user/distributed-file-management-system/internal/client"
	"github.com/user/distributed-file-management-system/internal/metadata"
	"github.com/user/distributed-file-management-system/internal/models"
	"github.com/user/distributed-file-management-system/internal/node"
)

const ReplicationFactor = 2

type Service struct {
	Metadata metadata.Store
	Nodes    *node.Manager
	Client   *client.StorageClient
}

func NewService(
	meta metadata.Store,
	nodes *node.Manager,
	client *client.StorageClient,
) *Service {

	return &Service{
		Metadata: meta,
		Nodes:    nodes,
		Client:   client,
	}
}

func (s *Service) Start() {

	for {

		s.Check()

		time.Sleep(10 * time.Second)
	}
}

func (s *Service) Check() {

	files, err := s.Metadata.List()
	if err != nil {
		return
	}

	for _, file := range files {
		s.checkFile(file)
	}
}

func (s *Service) checkFile(
	file models.FileMetadata,
) {

	changed := false

	for i := range file.Chunks {

		if s.repairChunk(
			&file.Chunks[i],
		) {

			changed = true
		}
	}

	if changed {

		s.Metadata.Save(file)
	}
}

func (s *Service) repairChunk(
	chunk *models.ChunkMetadata,
) bool {

	alive := chunk.AliveLocations(s.Nodes)

	if len(alive) >= ReplicationFactor {
		return false
	}

	excluded := make(map[string]bool)

	for _, loc := range chunk.Locations {
		excluded[loc.NodeID] = true
	}

	newNode, err := s.Nodes.FindReplacementNode(excluded)
	if err != nil {
		return false
	}

	reader, err := s.Client.DownloadChunk(
		alive[0].Address,
		chunk.ID,
	)

	if err != nil {
		return false
	}

	data, err := io.ReadAll(reader)
	reader.Close()

	if err != nil {
		return false
	}

	err = s.Client.UploadChunk(
		newNode.Address,
		chunk.ID,
		data,
	)

	if err != nil {
		return false
	}

	chunk.Locations = append(
		chunk.Locations,
		models.ChunkLocation{
			NodeID:  newNode.ID,
			Address: newNode.Address,
		},
	)

	var updated []models.ChunkLocation

	for _, loc := range chunk.Locations {

		if s.Nodes.IsAlive(loc.NodeID) {
			updated = append(updated, loc)
		}
	}

	chunk.Locations = updated

	return true
}
