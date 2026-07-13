package models

import "github.com/user/distributed-file-management-system/internal/node"

type ChunkLocation struct {
	NodeID  string `json:"node_id"`
	Address string `json:"address"`
}

type ChunkMetadata struct {
	ID        string          `json:"id"`
	Index     int             `json:"index"`
	Size      int64           `json:"size"`
	Locations []ChunkLocation `json:"locations"`
}

func (c *ChunkMetadata) AliveLocations(
	manager *node.Manager,
) []ChunkLocation {

	var alive []ChunkLocation

	for _, location := range c.Locations {

		if manager.IsAlive(location.NodeID) {
			alive = append(alive, location)
		}
	}

	return alive
}
