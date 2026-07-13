package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/user/distributed-file-management-system/internal/node"
	"github.com/user/distributed-file-management-system/internal/storage"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <port>")
		return
	}

	port := os.Args[1]

	store := storage.NewLocalStorage("data/node" + port)

	handler := node.NewHandler(store)

	http.HandleFunc("/chunk", handler.Chunk)
	http.HandleFunc("/health", handler.Health)

	fmt.Println("Storage Node running on :" + port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
