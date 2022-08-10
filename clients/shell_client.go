package clients

import (
	"fmt"
	"io"
	"os"
)

type shellClient struct {
	w     io.Writer
	trace bool
}

// Handles responsibility of writing to **stdout**
func NewShellWriter(w io.Writer, trace bool) BaseWriter {
	return &shellClient{
		w:     w,
		trace: trace,
	}
}

func (s *shellClient) Write(content interface{}) {
	fmt.Fprintln(s.w, content)
}

func (s *shellClient) WriteError(content interface{}) {
	fmt.Fprint(s.w, content)
	if s.trace {
		panic(content)
	}
	os.Exit(1)
}
