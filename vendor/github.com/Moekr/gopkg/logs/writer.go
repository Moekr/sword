package logs

import (
	"io"
	"os"
)

type writer struct {
	w io.WriteCloser
}

func newWriter(w io.WriteCloser) *writer {
	return &writer{
		w: w,
	}
}

func (w *writer) Write(p []byte) (int, error) {
	if w.w != nil {
		_, _ = os.Stdout.Write(p)
		return w.w.Write(p)
	}
	return os.Stdout.Write(p)
}

func (w *writer) Close() error {
	if w.w != nil {
		return w.w.Close()
	}
	return nil
}
