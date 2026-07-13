package node

import (
	"fmt"
	"sync"
)

type Manager struct {
	nodes []StorageNode
	mu    sync.RWMutex
}

func NewManager(nodes []StorageNode) *Manager {

	return &Manager{
		nodes: nodes,
	}
}

func (m *Manager) ChooseNodes(start int, replicas int) []StorageNode {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []StorageNode

	if len(m.nodes) == 0 {
		return result
	}

	for i := 0; i < len(m.nodes) && len(result) < replicas; i++ {

		index := (start + i) % len(m.nodes)

		if m.nodes[index].Alive {
			result = append(result, m.nodes[index])
		}
	}

	return result
}

func (m *Manager) GetNodes() []StorageNode {

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]StorageNode, len(m.nodes))
	copy(result, m.nodes)

	return result
}

func (m *Manager) SetAlive(id string, alive bool) {

	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.nodes {

		if m.nodes[i].ID == id {

			m.nodes[i].Alive = alive
			return
		}
	}
}

func (m *Manager) IsAlive(id string) bool {

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, node := range m.nodes {

		if node.ID == id {
			return node.Alive
		}
	}

	return false
}

func (m *Manager) FindReplacementNode(
	excluded map[string]bool,
) (StorageNode, error) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, node := range m.nodes {

		if !node.Alive {
			continue
		}

		if excluded[node.ID] {
			continue
		}

		return node, nil
	}

	return StorageNode{}, fmt.Errorf("no healthy node available")
}
