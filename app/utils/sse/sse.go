package sse

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	gosse "github.com/tmaxmax/go-sse"
)

// Transport SSE透传给下游
func Transport(w http.ResponseWriter, body io.Reader) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("response writer not support flush")
	}
	for ev, readErr := range gosse.Read(body, nil) {
		if readErr != nil {
			fmt.Fprintf(w, "event: error\n")
			fmt.Fprintf(w, "data: %s\n\n", readErr.Error())
			flusher.Flush()
			return readErr
		}
		var builder strings.Builder
		if ev.Type != "" {
			builder.WriteString("event: ")
			builder.WriteString(ev.Type)
			builder.WriteString("\n")
		}
		if ev.LastEventID != "" {
			builder.WriteString("id: ")
			builder.WriteString(ev.LastEventID)
			builder.WriteString("\n")
		}
		if ev.Data == "" {
			builder.WriteString("data:\n")
		} else {
			for _, line := range strings.Split(ev.Data, "\n") {
				builder.WriteString("data: ")
				builder.WriteString(line)
				builder.WriteString("\n")
			}
		}
		builder.WriteString("\n")
		if _, writeErr := w.Write([]byte(builder.String())); writeErr != nil {
			fmt.Fprintf(w, "event: error\n")
			fmt.Fprintf(w, "data: %s\n\n", writeErr.Error())
			flusher.Flush()
			return writeErr
		}
		flusher.Flush()
	}
	return nil
}
