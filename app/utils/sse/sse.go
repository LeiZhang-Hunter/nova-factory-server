package sse

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ApplyHeaders 应用并补齐 SSE 响应头。
func ApplyHeaders(h http.Header, upstream map[string]string) {
	for k, v := range upstream {
		if v != "" {
			h.Set(k, v)
		}
	}
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", "text/event-stream; charset=utf-8")
	}
	if h.Get("Cache-Control") == "" {
		h.Set("Cache-Control", "no-cache")
	}
	if h.Get("Connection") == "" {
		h.Set("Connection", "keep-alive")
	}
	if h.Get("X-Accel-Buffering") == "" {
		h.Set("X-Accel-Buffering", "no")
	}
}

// WriteErrorEvent 向下游写入 SSE 错误事件。
func WriteErrorEvent(w http.ResponseWriter, message string) error {
	if message == "" {
		message = "sse stream failed"
	}
	if _, err := fmt.Fprintf(w, "event: error\ndata: %s\n\n", message); err != nil {
		return err
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

// Transport SSE透传给下游并支持空闲超时控制。
func Transport(w http.ResponseWriter, body io.ReadCloser, idleTimeout time.Duration) error {
	if idleTimeout <= 0 {
		idleTimeout = 60 * time.Second
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("response writer not support flush")
	}
	type streamResult struct {
		data []byte
		err  error
	}
	resultCh := make(chan streamResult)
	go func() {
		defer close(resultCh)
		buf := make([]byte, 32*1024)
		for {
			n, readErr := body.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				resultCh <- streamResult{data: chunk}
			}
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					return
				}
				resultCh <- streamResult{err: readErr}
				return
			}
		}
	}()
	timer := time.NewTimer(idleTimeout)
	defer timer.Stop()
	resetTimer := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(idleTimeout)
	}
	for {
		select {
		case <-timer.C:
			_ = body.Close()
			return errors.New("sse idle timeout")
		case result, ok := <-resultCh:
			if !ok {
				return nil
			}
			if result.err != nil {
				return result.err
			}
			resetTimer()
			if len(result.data) == 0 {
				continue
			}
			if _, writeErr := w.Write(result.data); writeErr != nil {
				return writeErr
			}
			flusher.Flush()
		}
	}
}
