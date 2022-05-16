package terminal

import (
	"bytes"
	"io"
	"os"
)

type Terminal struct {
	outBuff bytes.Buffer
	multi   io.Writer
}

// New creates a new Terminal instance
func New() Terminal {
	var t Terminal
	t.multi = io.MultiWriter(os.Stdout, &t.outBuff)
	return t
}

// Write writes to the Terminal internal buffer and to os.Stdout
func (t Terminal) Write(p []byte) (n int, err error) {
	return t.multi.Write(p)
}

// Read reads from the internal buffer of the Terminal
func (t Terminal) Read(p []byte) (n int, err error) {
	return t.outBuff.Read(p)
}

func (t Terminal) Len() int {
	return t.outBuff.Len()
}
