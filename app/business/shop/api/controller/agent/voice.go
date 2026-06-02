//go:build ai

package agent

import (
	"encoding/json"
	"errors"
	"net"
	"strings"
	"time"

	apiModels "nova-factory-server/app/business/shop/api/models"
	shopService "nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Voice struct {
	service shopService.IApiShopVoiceService
}

func NewVoice(service shopService.IApiShopVoiceService) *Voice {
	return &Voice{service: service}
}

func (voice *Voice) WsRegister(router *gin.RouterGroup) {
	router.GET("/ws/voice", voice.WsVoice)
}

func (voice *Voice) PublicRoutes(router *gin.RouterGroup) {
}

type voiceEmitter struct {
	conn *websocket.Conn
}

func (e *voiceEmitter) SendEvent(event *apiModels.ShopVoiceServerEvent) error {
	return e.conn.WriteJSON(event)
}

func (e *voiceEmitter) SendAudioChunk(chunk []byte) error {
	return e.conn.WriteMessage(websocket.BinaryMessage, chunk)
}

func isVoiceWebSocketClosedError(err error) bool {
	if err == nil {
		return false
	}
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		return true
	}
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	errText := err.Error()
	return strings.Contains(errText, "websocket: close sent") ||
		strings.Contains(errText, "use of closed network connection") ||
		strings.Contains(errText, "broken pipe") ||
		strings.Contains(errText, "connection reset by peer")
}

func (voice *Voice) WsVoice(c *gin.Context) {
	wsStart := time.Now()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("voice websocket upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()
	defer zap.L().Info("[voice-debug] websocket closed",
		zap.String("client_ip", c.ClientIP()),
		zap.Duration("ws_elapsed", time.Since(wsStart)),
	)
	zap.L().Info("[voice-debug] websocket connected", zap.String("client_ip", c.ClientIP()))

	emitter := &voiceEmitter{conn: conn}
	zap.L().Info("[voice-debug] send session_ready", zap.Duration("ws_elapsed", time.Since(wsStart)))
	_ = emitter.SendEvent(&apiModels.ShopVoiceServerEvent{Type: "session_ready"})

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			zap.L().Warn("voice websocket read failed", zap.Duration("ws_elapsed", time.Since(wsStart)), zap.Error(err))
			return
		}

		var req apiModels.ShopVoiceSubmitReq
		if err := json.Unmarshal(payload, &req); err != nil {
			_ = emitter.SendEvent(&apiModels.ShopVoiceServerEvent{Type: "error", Message: "invalid voice payload"})
			continue
		}
		if strings.TrimSpace(req.Type) == "" {
			req.Type = "submit_text"
		}
		zap.L().Info("[voice-debug] websocket message received",
			zap.String("type", req.Type),
			zap.Int64("conversation_id", req.ConversationID),
			zap.Int("payload_bytes", len(payload)),
			zap.Duration("ws_elapsed", time.Since(wsStart)),
		)
		if req.Type != "submit_text" {
			continue
		}
		turnStart := time.Now()
		if err := voice.service.ProcessTurn(c, emitter, &req); err != nil {
			if isVoiceWebSocketClosedError(err) {
				zap.L().Warn("voice websocket closed during turn",
					zap.Int64("conversation_id", req.ConversationID),
					zap.Duration("turn_elapsed", time.Since(turnStart)),
					zap.Duration("ws_elapsed", time.Since(wsStart)),
					zap.Error(err),
				)
				return
			}
			zap.L().Error("process voice turn failed",
				zap.Int64("conversation_id", req.ConversationID),
				zap.Duration("turn_elapsed", time.Since(turnStart)),
				zap.Duration("ws_elapsed", time.Since(wsStart)),
				zap.Error(err),
			)
			_ = emitter.SendEvent(&apiModels.ShopVoiceServerEvent{Type: "error", ConversationID: req.ConversationID, Message: err.Error()})
		}
	}
}
