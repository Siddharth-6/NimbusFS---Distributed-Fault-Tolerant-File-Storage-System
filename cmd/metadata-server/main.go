package main

import (
	"fmt"
	"net/http"

	"github.com/user/distributed-file-management-system/internal/api"
	"github.com/user/distributed-file-management-system/internal/client"
	"github.com/user/distributed-file-management-system/internal/metadata"
	"github.com/user/distributed-file-management-system/internal/node"
	"github.com/user/distributed-file-management-system/internal/repair"
	"github.com/user/distributed-file-management-system/internal/service"
)

func main() {

	metadataStore := metadata.NewJSONStore("data/metadata")

	nodes := []node.StorageNode{
		{
			ID:      "node1",
			Address: "localhost:9001",
		},
		{
			ID:      "node2",
			Address: "localhost:9002",
		},
		{
			ID:      "node3",
			Address: "localhost:9003",
		},
	}

	manager := node.NewManager(nodes)

	storageClient := client.NewStorageClient()

	fileService := service.NewFileService(
		metadataStore,
		manager,
		storageClient,
	)

	uploadHandler := api.NewUploadHandler(fileService)

	http.HandleFunc("/upload", uploadHandler.Upload)
	http.HandleFunc("/files", uploadHandler.List)
	http.HandleFunc("/download", uploadHandler.Download)

	heartbeat := node.NewHeartbeat(manager)
	go heartbeat.Start()

	repair := repair.NewService(
		metadataStore,
		manager,
		storageClient,
	)
	go repair.Start()

	fmt.Println("Metadata Server running on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
