package logging

import (
	"bytes"
	"github.com/yangchenxing/go-string-mapformatter"
	"io"
)

// Writer wrap the io.Writer interface. It is useful for WriterFactory.
type Writer io.Writer

// Handler handles format the log message and dispatch it to writers
type Handler struct {
	// responsible levels
	Levels []string

	// format of log message
	Format string

	// writers
	Writers []Writer
}

func (handler *Handler) write(context ...map[string]interface{}) error {
	var buf bytes.Buffer
	buf.WriteString(mapformatter.Format(handler.Format, context...))
	buf.WriteRune('\n')
	text := buf.Bytes()
	for _, writer := range handler.Writers {
		if _, err := writer.Write(text); err != nil {
			return err
		}
	}
	return nil
}
