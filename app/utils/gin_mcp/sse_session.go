package gin_mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"sync"
	"sync/atomic"
)

// sseSession represents an active SSE connection.
type sseSession struct {
	done                chan struct{}
	eventQueue          chan string // Channel for queuing events
	sessionID           string
	requestID           atomic.Int64
	notificationChannel chan mcp.JSONRPCNotification
	initialized         atomic.Bool
	loggingLevel        atomic.Value
	tools               sync.Map     // stores session-specific tools
	clientInfo          atomic.Value // stores session-specific client info
}

func (s *sseSession) SessionID() string {
	return s.sessionID
}

func (s *sseSession) NotificationChannel() chan<- mcp.JSONRPCNotification {
	return s.notificationChannel
}

func (s *sseSession) Initialize() {
	// set default logging level
	s.loggingLevel.Store(mcp.LoggingLevelError)
	s.initialized.Store(true)
}

func (s *sseSession) Initialized() bool {
	return s.initialized.Load()
}

func (s *sseSession) SetLogLevel(level mcp.LoggingLevel) {
	s.loggingLevel.Store(level)
}

func (s *sseSession) GetLogLevel() mcp.LoggingLevel {
	level := s.loggingLevel.Load()
	if level == nil {
		return mcp.LoggingLevelError
	}
	return level.(mcp.LoggingLevel)
}

func (s *sseSession) GetSessionTools() map[string]server.ServerTool {
	tools := make(map[string]server.ServerTool)
	s.tools.Range(func(key, value any) bool {
		if tool, ok := value.(server.ServerTool); ok {
			tools[key.(string)] = tool
		}
		return true
	})
	return tools
}

func (s *sseSession) SetSessionTools(tools map[string]server.ServerTool) {
	// Clear existing tools
	s.tools.Clear()

	// Set new tools
	for name, tool := range tools {
		s.tools.Store(name, tool)
	}
}

func (s *sseSession) GetClientInfo() mcp.Implementation {
	if value := s.clientInfo.Load(); value != nil {
		if clientInfo, ok := value.(mcp.Implementation); ok {
			return clientInfo
		}
	}
	return mcp.Implementation{}
}

func (s *sseSession) SetClientInfo(clientInfo mcp.Implementation) {
	s.clientInfo.Store(clientInfo)
}
