package chunker

import "io"

type Chunk struct {
	ID    string
	Index int
	Data  []byte
}

type Chunker struct {
	reader io.Reader
	size   int
	index  int
}

func NewChunker(
	reader io.Reader,
	size int,
) *Chunker {

	return &Chunker{
		reader: reader,
		size:   size,
		index:  0,
	}
}

func (c *Chunker) Next() (*Chunk, error) {

	buffer := make([]byte, c.size)

	n, err := c.reader.Read(buffer)

	if err == io.EOF {
		return nil, io.EOF
	}

	if err != nil {
		return nil, err
	}

	data := make([]byte, n)

	copy(data, buffer[:n])

	chunk := &Chunk{
		Index: c.index,
		Data:  data,
	}

	c.index++

	return chunk, nil
}
