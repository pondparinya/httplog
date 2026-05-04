package httplog

import "io"

const NoOpWriterKey = "noop"

type NoOpWriter struct {
	io.Writer
}

func NewNoOpWriter() NoOpWriter {
	return NoOpWriter{io.Discard}
}

func (w NoOpWriter) Close() error {
	return nil
}

func (w NoOpWriter) Sync() error {
	return nil
}
