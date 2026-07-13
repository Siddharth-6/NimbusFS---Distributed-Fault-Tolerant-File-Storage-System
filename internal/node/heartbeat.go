package node

import (
	"net/http"
	"time"
)

type Heartbeat struct {
	manager *Manager
	client  *http.Client
}

func NewHeartbeat(
	manager *Manager,
) *Heartbeat {

	return &Heartbeat{
		manager: manager,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (h *Heartbeat) Start() {

	for {

		h.check()

		time.Sleep(5 * time.Second)
	}
}

func (h *Heartbeat) check() {

	nodes := h.manager.GetNodes()

	for _, n := range nodes {

		resp, err := h.client.Get(
			"http://" + n.Address + "/health",
		)

		if err != nil {

			h.manager.SetAlive(n.ID, false)

			continue
		}

		resp.Body.Close()

		h.manager.SetAlive(n.ID, true)
	}
}
