//go:build ai

package agent

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// WebSocket 超时配置
	wsWriteWait      = 10 * time.Second
	wsPongWait       = 60 * time.Second
	wsPingPeriod     = (wsPongWait * 9) / 10
	wsMaxMessageSize = 512
	wsReadBufSize    = 1024
	wsWriteBufSize   = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  wsReadBufSize,
	WriteBufferSize: wsWriteBufSize,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域，生产环境应做限制
	},
}

// WsChat WebSocket 流式聊天
// @Summary WebSocket 流式聊天
// @Description 建立 WebSocket 连接后，发送 JSON 消息体 { conversation_id, tab_id, content }，服务端将 SSE 事件实时推送回客户端
// @Tags app接口/商城/App智能体WebSocket
// @Param ws path string true "WebSocket 端点 /api/v1/app/shop/agent/conversations/ws/chat"
// @Router /api/v1/app/shop/agent/conversations/ws/chat [get]
func (conversations *Conversations) WsChat(c *gin.Context) {
	wsStart := time.Now()
	// 升级为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("websocket upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()
	defer func() {
		//zap.L().Info("[ws-chat-debug] websocket closed",
		//	zap.String("client_ip", c.ClientIP()),
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//)
	}()
	//zap.L().Info("[ws-chat-debug] websocket connected",
	//	zap.String("client_ip", c.ClientIP()),
	//	zap.Duration("ws_elapsed", time.Since(wsStart)),
	//)

	// 读取客户端第一条消息（发起聊天请求）
	_, messageBytes, err := conn.ReadMessage()
	if err != nil {
		//zap.L().Warn("[ws-chat-debug] read init message failed",
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//	zap.Error(err),
		//)
		return
	}
	//zap.L().Info("[ws-chat-debug] init message received",
	//	zap.Int("payload_bytes", len(messageBytes)),
	//	zap.Duration("ws_elapsed", time.Since(wsStart)),
	//)

	var initReq aidatasetmodels.SendMessageInput
	if err := json.Unmarshal(messageBytes, &initReq); err != nil {
		//zap.L().Error("[ws-chat-debug] init message unmarshal failed",
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//	zap.Error(err),
		//)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"invalid JSON"}`))
		return
	}
	//zap.L().Info("[ws-chat-debug] init message parsed",
	//	zap.Int64("conversation_id", initReq.ConversationID),
	//	zap.String("tab_id", initReq.TabID),
	//	zap.Int("content_runes", len([]rune(initReq.Content))),
	//	zap.Duration("ws_elapsed", time.Since(wsStart)),
	//)

	if initReq.ConversationID == 0 {
		//zap.L().Warn("[ws-chat-debug] missing conversation_id",
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"conversation_id is required"}`))
		return
	}
	if strings.TrimSpace(initReq.Content) == "" {
		//zap.L().Warn("[ws-chat-debug] missing content",
		//	zap.Int64("conversation_id", initReq.ConversationID),
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"content is required"}`))
		return
	}
	if strings.TrimSpace(initReq.TabID) == "" {
		initReq.TabID = "team"
	}
	var ret bool
	ret = false
	initReq.EnableThinking = &ret

	// 调用网关获取 SSE 流
	//gatewayStart := time.Now()
	//zap.L().Info("[ws-chat-debug] gateway chat start",
	//	zap.Int64("conversation_id", initReq.ConversationID),
	//	zap.String("tab_id", initReq.TabID),
	//	zap.Duration("ws_elapsed", time.Since(wsStart)),
	//)
	sseResp, err := conversations.gatewayService.Chat(c, &initReq)
	if err != nil {
		//zap.L().Error("[ws-chat-debug] gateway chat failed",
		//	zap.Int64("conversation_id", initReq.ConversationID),
		//	zap.Duration("gateway_elapsed", time.Since(gatewayStart)),
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//	zap.Error(err),
		//)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"`+err.Error()+`"}`))
		return
	}
	//zap.L().Info("[ws-chat-debug] gateway chat returned",
	//	zap.Int64("conversation_id", initReq.ConversationID),
	//	zap.Duration("gateway_elapsed", time.Since(gatewayStart)),
	//	zap.Duration("ws_elapsed", time.Since(wsStart)),
	//)
	if sseResp == nil {
		//zap.L().Warn("[ws-chat-debug] gateway chat nil response",
		//	zap.Int64("conversation_id", initReq.ConversationID),
		//	zap.Duration("gateway_elapsed", time.Since(gatewayStart)),
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"chat service无响应"}`))
		return
	}
	if sseResp.StatusCode != 200 {
		//zap.L().Warn("[ws-chat-debug] gateway chat non-200",
		//	zap.Int64("conversation_id", initReq.ConversationID),
		//	zap.Int("status_code", sseResp.StatusCode),
		//	zap.String("message", sseResp.Message),
		//	zap.Duration("gateway_elapsed", time.Since(gatewayStart)),
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"`+sseResp.Message+`"}`))
		return
	}

	// 启动 ping/pong 保活
	conn.SetReadLimit(wsMaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(wsPongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})

	// 启动 writer goroutine 处理 ping
	writerDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(wsPingPeriod)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			case <-writerDone:
				return
			}
		}
	}()

	// 将 SSE 流转发到 WebSocket 客户端
	if sseResp.Body != nil {
		defer sseResp.Body.Close()
		//zap.L().Info("[ws-chat-debug] forward stream start",
		//	zap.Int64("conversation_id", initReq.ConversationID),
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//)
		wsWriterDone := conversations.forwardSSEToWS(sseResp.Body, conn, initReq.ConversationID, wsStart)
		<-wsWriterDone
		//zap.L().Info("[ws-chat-debug] forward stream done",
		//	zap.Int64("conversation_id", initReq.ConversationID),
		//	zap.Duration("ws_elapsed", time.Since(wsStart)),
		//)
		close(writerDone)
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "complete"), time.Now().Add(wsWriteWait))
	}
}

// forwardSSEToWS 读取 SSE 流，将每个事件转发为 WebSocket JSON 消息
// 返回一个 chan，在流结束后 close
func (conversations *Conversations) forwardSSEToWS(body io.Reader, conn *websocket.Conn, conversationID int64, wsStart time.Time) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		forwardStart := time.Now()
		br := bufio.NewReaderSize(body, 4096)

		var currentEvent string
		var currentData []string
		var eventCount int
		var firstEventLogged bool
		var firstChunkLogged bool

		flush := func() bool {
			if len(currentData) == 0 {
				return true
			}
			wsMsg := map[string]interface{}{
				"event": currentEvent,
				"data":  strings.Join(currentData, "\n"),
			}
			jsonBytes, _ := json.Marshal(wsMsg)
			if !firstEventLogged {
				firstEventLogged = true
				//zap.L().Info("[ws-chat-debug] first upstream event",
				//	zap.Int64("conversation_id", conversationID),
				//	zap.String("event", currentEvent),
				//	zap.Int("payload_bytes", len(jsonBytes)),
				//	zap.Duration("first_event_elapsed", time.Since(forwardStart)),
				//	zap.Duration("ws_elapsed", time.Since(wsStart)),
				//)
			}
			eventCount++
			if currentEvent == "chat:chunk" && !firstChunkLogged {
				firstChunkLogged = true
				//zap.L().Info("[ws-chat-debug] first chunk event",
				//	zap.Int64("conversation_id", conversationID),
				//	zap.Int("event_count", eventCount),
				//	zap.Int("payload_bytes", len(jsonBytes)),
				//	zap.Duration("first_chunk_elapsed", time.Since(forwardStart)),
				//	zap.Duration("ws_elapsed", time.Since(wsStart)),
				//)
			}
			//if currentEvent == "chat:start" || currentEvent == "chat:complete" || currentEvent == "chat:stopped" || currentEvent == "chat:error" {
			//	zap.L().Info("[ws-chat-debug] upstream event",
			//		zap.Int64("conversation_id", conversationID),
			//		zap.String("event", currentEvent),
			//		zap.Int("event_count", eventCount),
			//		zap.Int("payload_bytes", len(jsonBytes)),
			//		zap.Duration("forward_elapsed", time.Since(forwardStart)),
			//		zap.Duration("ws_elapsed", time.Since(wsStart)),
			//	)
			//}
			//writeStart := time.Now()
			if err := conn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
				//zap.L().Warn("[ws-chat-debug] websocket write failed",
				//	zap.Int64("conversation_id", conversationID),
				//	zap.String("event", currentEvent),
				//	zap.Int("event_count", eventCount),
				//	zap.Duration("write_elapsed", time.Since(writeStart)),
				//	zap.Duration("forward_elapsed", time.Since(forwardStart)),
				//	zap.Duration("ws_elapsed", time.Since(wsStart)),
				//	zap.Error(err),
				//)
				return false
			}
			currentEvent = ""
			currentData = nil
			return true
		}

		for {
			line, err := br.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					flush()
					zap.L().Info("[ws-chat-debug] upstream stream eof",
						zap.Int64("conversation_id", conversationID),
						zap.Int("event_count", eventCount),
						zap.Duration("forward_elapsed", time.Since(forwardStart)),
						zap.Duration("ws_elapsed", time.Since(wsStart)),
					)
				} else {
					zap.L().Warn("[ws-chat-debug] upstream stream read failed",
						zap.Int64("conversation_id", conversationID),
						zap.Int("event_count", eventCount),
						zap.Duration("forward_elapsed", time.Since(forwardStart)),
						zap.Duration("ws_elapsed", time.Since(wsStart)),
						zap.Error(err),
					)
				}
				break
			}

			line = strings.TrimSuffix(line, "\n")
			if strings.HasPrefix(line, "event:") {
				if !flush() {
					break
				}
				currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			} else if strings.HasPrefix(line, "data:") {
				currentData = append(currentData, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
			} else if line == "" {
				if !flush() {
					break
				}
			}
		}
	}()
	return done
}

// WsChatRegister 注册 WebSocket 路由（与小程序 HTTP 路由保持同一前缀）
func (conversations *Conversations) WsChatRegister(router *gin.RouterGroup) {
	router.GET("/ws/chat", conversations.WsChat)
}
