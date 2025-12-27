package gin_mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
	"net/http"
	"sync"
)

type wrappedHTTP struct {
	http.Handler
	sessions    sync.Map
	server      *server.MCPServer
	contextFunc server.SSEContextFunc
}

func newWrappedHTTP(handler http.Handler) *wrappedHTTP {
	return &wrappedHTTP{
		Handler: handler,
	}
}

// writeJSONRPCError writes a JSON-RPC error response with the given error details.
func (wh *wrappedHTTP) writeJSONRPCError(
	w http.ResponseWriter,
	id any,
	code int,
	message string,
) {
	response := createErrorResponse(id, code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to encode response: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

func (wh *wrappedHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		wh.writeJSONRPCError(w, nil, mcp.INVALID_REQUEST, "Method not allowed")
		return
	}
	//wh.server.AddSessionTools()
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		wh.writeJSONRPCError(w, nil, mcp.INVALID_PARAMS, "Missing sessionId")
		return
	}
	sessionI, ok := wh.sessions.Load(sessionID)
	if !ok {
		wh.writeJSONRPCError(w, nil, mcp.INVALID_PARAMS, "Invalid session ID")
		return
	}
	session := sessionI.(*sseSession)

	// Set the client context before handling the message
	ctx := wh.server.WithContext(r.Context(), session)
	if wh.contextFunc != nil {
		ctx = wh.contextFunc(ctx, r)
	}

	// Parse message as raw JSON
	var rawMessage json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&rawMessage); err != nil {
		wh.writeJSONRPCError(w, nil, mcp.PARSE_ERROR, "Parse error")
		return
	}

	// Create a context that preserves all values from parent ctx but won't be canceled when the parent is canceled.
	// this is required because the http ctx will be canceled when the client disconnects
	detachedCtx := context.WithoutCancel(ctx)

	// quick return request, send 202 Accepted with no body, then deal the message and sent response via SSE
	w.WriteHeader(http.StatusAccepted)

	// Create a new context for handling the message that will be canceled when the message handling is done
	messageCtx, cancel := context.WithCancel(detachedCtx)

	go func(ctx context.Context) {
		defer cancel()
		// Use the context that will be canceled when session is done
		// Process message through MCPServer
		response := wh.server.HandleMessage(ctx, rawMessage)
		// Only send response if there is one (not for notifications)
		if response != nil {
			var message string
			if eventData, err := json.Marshal(response); err != nil {
				// If there is an error marshalling the response, send a generic error response
				log.Printf("failed to marshal response: %v", err)
				message = "event: message\ndata: {\"error\": \"internal error\",\"jsonrpc\": \"2.0\", \"id\": null}\n\n"
			} else {
				message = fmt.Sprintf("event: message\ndata: %s\n\n", eventData)
			}

			// Queue the event for sending via SSE
			select {
			case session.eventQueue <- message:
				// Event queued successfully
			case <-session.done:
				// Session is closed, don't try to queue
			default:
				// Queue is full, log this situation
				log.Printf("Event queue full for session %s", sessionID)
			}
		}
	}(messageCtx)
}
