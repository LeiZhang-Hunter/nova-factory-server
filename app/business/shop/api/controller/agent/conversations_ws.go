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
// @Router /api/v1/app/shop/agent/conversations/ws/chat [ws]
func (conversations *Conversations) WsChat(c *gin.Context) {
	// 升级为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("websocket upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	// 读取客户端第一条消息（发起聊天请求）
	_, messageBytes, err := conn.ReadMessage()
	if err != nil {
		zap.L().Warn("ws read init message failed", zap.Error(err))
		return
	}

	var initReq aidatasetmodels.SendMessageInput
	if err := json.Unmarshal(messageBytes, &initReq); err != nil {
		zap.L().Error("ws init message unmarshal failed", zap.Error(err))
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"invalid JSON"}`))
		return
	}

	if initReq.ConversationID == 0 {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"conversation_id is required"}`))
		return
	}
	if strings.TrimSpace(initReq.Content) == "" {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"content is required"}`))
		return
	}
	if strings.TrimSpace(initReq.TabID) == "" {
		initReq.TabID = "team"
	}
	var ret bool
	ret = true
	initReq.EnableThinking = &ret

	// 调用网关获取 SSE 流
	sseResp, err := conversations.gatewayService.Chat(c, &initReq)
	if err != nil {
		zap.L().Error("ws gateway chat failed", zap.Error(err))
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"`+err.Error()+`"}`))
		return
	}
	if sseResp == nil {
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"chat service无响应"}`))
		return
	}
	if sseResp.StatusCode != 200 {
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
		wsWriterDone := conversations.forwardSSEToWS(sseResp.Body, conn)
		<-wsWriterDone
		close(writerDone)
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "complete"), time.Now().Add(wsWriteWait))
	}
}

// forwardSSEToWS 读取 SSE 流，将每个事件转发为 WebSocket JSON 消息
// 返回一个 chan，在流结束后 close
func (conversations *Conversations) forwardSSEToWS(body io.Reader, conn *websocket.Conn) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		br := bufio.NewReaderSize(body, 4096)

		var currentEvent string
		var currentData []string

		flush := func() {
			if len(currentData) == 0 {
				return
			}
			wsMsg := map[string]interface{}{
				"event": currentEvent,
				"data":  strings.Join(currentData, "\n"),
			}
			jsonBytes, _ := json.Marshal(wsMsg)
			if err := conn.WriteMessage(websocket.TextMessage, jsonBytes); err != nil {
				return
			}
			currentEvent = ""
			currentData = nil
		}

		for {
			line, err := br.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					flush()
				}
				break
			}

			line = strings.TrimSuffix(line, "\n")
			if strings.HasPrefix(line, "event:") {
				flush()
				currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			} else if strings.HasPrefix(line, "data:") {
				currentData = append(currentData, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
			} else if line == "" {
				flush()
			}
		}
	}()
	return done
}

// WsChatRegister 注册 WebSocket 路由（与小程序 HTTP 路由保持同一前缀）
func (conversations *Conversations) WsChatRegister(router *gin.RouterGroup) {
	router.GET("/ws/chat", conversations.WsChat)
}
