package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type StorageClient struct {
	httpClient *http.Client
}

func NewStorageClient() *StorageClient {

	return &StorageClient{
		httpClient: &http.Client{},
	}
}

func (c *StorageClient) UploadChunk(
	address string,
	id string,
	data []byte,
) error {

	url := fmt.Sprintf(
		"http://%s/chunk?id=%s",
		address,
		id,
	)

	resp, err := c.httpClient.Post(
		url,
		"application/octet-stream",
		bytes.NewReader(data),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed")
	}

	return nil
}

func (c *StorageClient) DownloadChunk(
	address string,
	id string,
) (io.ReadCloser, error) {

	url := fmt.Sprintf(
		"http://%s/chunk?id=%s",
		address,
		id,
	)

	resp, err := c.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("chunk not found")
	}

	return resp.Body, nil
}
