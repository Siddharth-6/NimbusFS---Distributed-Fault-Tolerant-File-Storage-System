package service

import "io"

type MultiReadCloser struct {
	io.Reader
	closers []io.Closer
}

func NewMultiReadCloser(
	reader io.Reader,
	closers []io.Closer,
) *MultiReadCloser {

	return &MultiReadCloser{
		Reader:  reader,
		closers: closers,
	}
}

func (m *MultiReadCloser) Close() error {

	for _, c := range m.closers {
		c.Close()
	}

	return nil
}
