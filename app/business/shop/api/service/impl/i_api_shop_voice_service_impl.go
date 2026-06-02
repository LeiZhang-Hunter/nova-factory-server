//go:build ai

package impl

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	apiModels "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const maxVoiceSynthesisRunes = 600

type IApiShopVoiceServiceImpl struct {
	gatewayService aidatasetservice.IAIGatewayService
}

func NewIApiShopVoiceServiceImpl(gatewayService aidatasetservice.IAIGatewayService) service.IApiShopVoiceService {
	return &IApiShopVoiceServiceImpl{
		gatewayService: gatewayService,
	}
}

func (s *IApiShopVoiceServiceImpl) ProcessTurn(c *gin.Context, emitter service.IApiShopVoiceEmitter, req *apiModels.ShopVoiceSubmitReq) error {
	turnStart := time.Now()
	if emitter == nil {
		return errors.New("voice emitter is nil")
	}
	if req == nil {
		return errors.New("voice request is nil")
	}
	if req.ConversationID == 0 {
		return errors.New("conversation_id不能为空")
	}
	if strings.TrimSpace(req.Text) == "" {
		return errors.New("text不能为空")
	}
	zap.L().Info("[voice-debug] turn start",
		zap.Int64("conversation_id", req.ConversationID),
		zap.String("tab_id", req.TabID),
		zap.Int("text_runes", utf8.RuneCountInString(req.Text)),
	)
	transcript := strings.TrimSpace(req.Text)

	if strings.TrimSpace(req.TabID) == "" {
		req.TabID = "team"
	}
	if err := emitter.SendEvent(&apiModels.ShopVoiceServerEvent{Type: "chat_start", ConversationID: req.ConversationID}); err != nil {
		zap.L().Warn("[voice-debug] send chat_start failed", zap.Int64("conversation_id", req.ConversationID), zap.Error(err))
		return err
	}
	var ret bool
	chatStart := time.Now()
	chatResp, err := s.gatewayService.Chat(c, &aidatasetmodels.SendMessageInput{
		ConversationID: req.ConversationID,
		Content:        transcript,
		TabID:          req.TabID,
		EnableThinking: &ret,
	})
	if err != nil {
		zap.L().Warn("[voice-debug] chat request failed", zap.Int64("conversation_id", req.ConversationID), zap.Duration("elapsed", time.Since(chatStart)), zap.Error(err))
		return err
	}
	if chatResp == nil {
		return errors.New("聊天服务无响应")
	}
	if chatResp.StatusCode != http.StatusOK {
		return errors.New(chatResp.Message)
	}
	if !chatResp.IsStream || chatResp.Body == nil {
		return errors.New("聊天服务未返回流式响应")
	}
	defer chatResp.Body.Close()

	ttsClient := newXFYunTTSClient()
	if err := emitter.SendEvent(&apiModels.ShopVoiceServerEvent{
		Type:           "tts_start",
		ConversationID: req.ConversationID,
		Mime:           "audio/pcm",
		SampleRate:     ttsClient.sampleRate,
		Channels:       1,
		BitDepth:       16,
	}); err != nil {
		zap.L().Warn("[voice-debug] send tts_start failed", zap.Int64("conversation_id", req.ConversationID), zap.Error(err))
		return err
	}
	zap.L().Info("[voice-debug] tts start",
		zap.Int64("conversation_id", req.ConversationID),
		zap.Int("sample_rate", ttsClient.sampleRate),
		zap.Int("max_voice_synthesis_runes", maxVoiceSynthesisRunes),
		zap.Duration("elapsed", time.Since(turnStart)),
	)

	streamCtx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	audioSeq := 0
	var assistantText strings.Builder

	ttsSession, err := ttsClient.OpenStream(streamCtx, func(chunk []byte) error {
		if len(chunk) == 0 {
			return nil
		}
		audioSeq++
		zap.L().Info("[voice-debug] tts audio chunk",
			zap.Int64("conversation_id", req.ConversationID),
			zap.Int("audio_seq", audioSeq),
			zap.Int("bytes", len(chunk)),
			zap.Duration("elapsed", time.Since(turnStart)),
		)
		if err := emitter.SendEvent(&apiModels.ShopVoiceServerEvent{Type: "tts_part", ConversationID: req.ConversationID, Mime: "audio/pcm", SampleRate: ttsClient.sampleRate, Channels: 1, BitDepth: 16, Seq: audioSeq}); err != nil {
			zap.L().Warn("[voice-debug] send tts_part failed", zap.Int64("conversation_id", req.ConversationID), zap.Int("audio_seq", audioSeq), zap.Error(err))
			return err
		}
		if err := emitter.SendAudioChunk(chunk); err != nil {
			zap.L().Warn("[voice-debug] send audio chunk failed", zap.Int64("conversation_id", req.ConversationID), zap.Int("audio_seq", audioSeq), zap.Int("bytes", len(chunk)), zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		zap.L().Warn("[voice-debug] open tts stream failed", zap.Int64("conversation_id", req.ConversationID), zap.Duration("elapsed", time.Since(turnStart)), zap.Error(err))
		return err
	}
	defer ttsSession.Close()

	chunker := newVoiceTextChunker(16, 11)
	textSeq := 0
	textStarted := false
	synthesizedRunes := 0

	flushChunk := func(piece string) error {
		piece = strings.TrimSpace(piece)
		if piece == "" {
			return nil
		}
		remainingRunes := maxVoiceSynthesisRunes - synthesizedRunes
		if remainingRunes <= 0 {
			zap.L().Info("[voice-debug] tts text skipped after limit",
				zap.Int64("conversation_id", req.ConversationID),
				zap.Int("synthesized_runes", synthesizedRunes),
				zap.Int("piece_runes", utf8.RuneCountInString(piece)),
			)
			return nil
		}
		pieceRunes := utf8.RuneCountInString(piece)
		if pieceRunes > remainingRunes {
			piece = trimRunes(piece, remainingRunes)
			pieceRunes = remainingRunes
		}
		status := 1
		if !textStarted {
			status = 0
			textStarted = true
		}
		textSeq++
		zap.L().Info("[voice-debug] tts send text",
			zap.Int64("conversation_id", req.ConversationID),
			zap.Int("text_seq", textSeq),
			zap.Int("status", status),
			zap.Int("piece_runes", pieceRunes),
			zap.Int("synthesized_runes_before", synthesizedRunes),
			zap.Duration("elapsed", time.Since(turnStart)),
		)
		if err := ttsSession.SendText(piece, status, textSeq); err != nil {
			zap.L().Warn("[voice-debug] tts send text failed",
				zap.Int64("conversation_id", req.ConversationID),
				zap.Int("text_seq", textSeq),
				zap.Int("status", status),
				zap.Int("piece_runes", pieceRunes),
				zap.Int("synthesized_runes_before", synthesizedRunes),
				zap.Duration("elapsed", time.Since(turnStart)),
				zap.Error(err),
			)
			return err
		}
		synthesizedRunes += pieceRunes
		return nil
	}

	br := bufio.NewReaderSize(chatResp.Body, 4096)
	var currentEvent string
	var currentData []string

	flushSSE := func() error {
		if len(currentData) == 0 {
			return nil
		}
		event := currentEvent
		dataText := strings.Join(currentData, "\n")
		currentEvent = ""
		currentData = nil

		var payload map[string]any
		if err := json.Unmarshal([]byte(dataText), &payload); err != nil {
			return nil
		}

		switch event {
		case "chat:chunk":
			delta := stringValue(payload["delta"])
			if delta == "" {
				return nil
			}
			zap.L().Info("[voice-debug] chat chunk",
				zap.Int64("conversation_id", req.ConversationID),
				zap.Int("delta_runes", utf8.RuneCountInString(delta)),
				zap.Int("assistant_runes_before", utf8.RuneCountInString(assistantText.String())),
				zap.Duration("elapsed", time.Since(turnStart)),
			)
			if err := emitter.SendEvent(&apiModels.ShopVoiceServerEvent{Type: "chat_chunk", ConversationID: req.ConversationID, Text: delta}); err != nil {
				zap.L().Warn("[voice-debug] send chat_chunk failed", zap.Int64("conversation_id", req.ConversationID), zap.Error(err))
				return err
			}
			assistantText.WriteString(delta)
			chunker.Append(cleanTextForVoiceSynthesis(delta))
			ready := chunker.FlushReady(false)
			for _, piece := range ready {
				if err := flushChunk(piece); err != nil {
					return err
				}
			}
		case "chat:complete":
			zap.L().Info("[voice-debug] chat complete event",
				zap.Int64("conversation_id", req.ConversationID),
				zap.Int("assistant_runes", utf8.RuneCountInString(assistantText.String())),
				zap.Int("synthesized_runes", synthesizedRunes),
				zap.Duration("elapsed", time.Since(turnStart)),
			)
			ready := chunker.FlushReady(true)
			for _, piece := range ready {
				if err := flushChunk(piece); err != nil {
					return err
				}
			}
		}
		return nil
	}

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				if err := flushSSE(); err != nil {
					return err
				}
				break
			}
			return err
		}
		line = strings.TrimSuffix(line, "\n")
		if strings.HasPrefix(line, "event:") {
			if err := flushSSE(); err != nil {
				return err
			}
			currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}
		if strings.HasPrefix(line, "data:") {
			currentData = append(currentData, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
			continue
		}
		if line == "" {
			if err := flushSSE(); err != nil {
				return err
			}
		}
	}

	if err := emitter.SendEvent(&apiModels.ShopVoiceServerEvent{Type: "chat_complete", ConversationID: req.ConversationID, Text: assistantText.String()}); err != nil {
		zap.L().Warn("[voice-debug] send chat_complete failed", zap.Int64("conversation_id", req.ConversationID), zap.Error(err))
		return err
	}
	textSeq++
	zap.L().Info("[voice-debug] tts finish send",
		zap.Int64("conversation_id", req.ConversationID),
		zap.Int("text_seq", textSeq),
		zap.Int("assistant_runes", utf8.RuneCountInString(assistantText.String())),
		zap.Int("synthesized_runes", synthesizedRunes),
		zap.Duration("elapsed", time.Since(turnStart)),
	)
	if err := ttsSession.Finish(textSeq); err != nil {
		zap.L().Warn("[voice-debug] tts finish failed", zap.Int64("conversation_id", req.ConversationID), zap.Duration("elapsed", time.Since(turnStart)), zap.Error(err))
		return err
	}
	if err := emitter.SendEvent(&apiModels.ShopVoiceServerEvent{Type: "tts_end", ConversationID: req.ConversationID}); err != nil {
		zap.L().Warn("[voice-debug] send tts_end failed", zap.Int64("conversation_id", req.ConversationID), zap.Error(err))
		return err
	}
	zap.L().Info("[voice-debug] turn complete",
		zap.Int64("conversation_id", req.ConversationID),
		zap.Int("assistant_runes", utf8.RuneCountInString(assistantText.String())),
		zap.Int("synthesized_runes", synthesizedRunes),
		zap.Int("audio_chunks", audioSeq),
		zap.Duration("elapsed", time.Since(turnStart)),
	)
	return nil
}

func stringValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func trimRunes(text string, limit int) string {
	if limit <= 0 {
		return ""
	}
	if utf8.RuneCountInString(text) <= limit {
		return text
	}
	runes := []rune(text)
	return string(runes[:limit])
}

var (
	voiceHTMLTagPattern      = regexp.MustCompile(`<[^>]+>`)
	voiceMarkdownLinkPattern = regexp.MustCompile(`!?\[([^\]]*)\]\([^)]+\)`)
	voiceMarkdownCodeFence   = regexp.MustCompile(`(?m)^\s*` + "```" + `.*$`)
	voiceMarkdownHeading     = regexp.MustCompile(`(?m)^\s{0,3}#{1,6}\s*`)
	voiceMarkdownListMarker  = regexp.MustCompile(`(?m)^\s*(?:[-+*]|\d+[.)])\s+`)
	voiceMarkdownQuote       = regexp.MustCompile(`(?m)^\s{0,3}>\s?`)
	voiceMarkdownRule        = regexp.MustCompile(`(?m)^\s*[-*_]{3,}\s*$`)
	voiceMarkdownTableAlign  = regexp.MustCompile(`(?m)^\s*\|?\s*:?-{3,}:?\s*(?:\|\s*:?-{3,}:?\s*)+\|?\s*$`)
	voiceWhitespace          = regexp.MustCompile(`[ \t\r\n]+`)
)

func cleanTextForVoiceSynthesis(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = voiceHTMLTagPattern.ReplaceAllString(text, "")
	text = voiceMarkdownLinkPattern.ReplaceAllString(text, "$1")
	text = voiceMarkdownCodeFence.ReplaceAllString(text, "")
	text = voiceMarkdownHeading.ReplaceAllString(text, "")
	text = voiceMarkdownListMarker.ReplaceAllString(text, "")
	text = voiceMarkdownQuote.ReplaceAllString(text, "")
	text = voiceMarkdownRule.ReplaceAllString(text, "")
	text = voiceMarkdownTableAlign.ReplaceAllString(text, "")
	text = strings.NewReplacer(
		"```", " ",
		"`", "",
		"**", "",
		"__", "",
		"~~", "",
		"*", "",
		"_", "",
		"|", "，",
		"[", "",
		"]", "",
		"(", "",
		")", "",
	).Replace(text)
	text = voiceWhitespace.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

type voiceTextChunker struct {
	flushChars      int
	firstFlushChars int
	buffer          strings.Builder
	firstFlushed    bool
}

func newVoiceTextChunker(flushChars, firstFlushChars int) *voiceTextChunker {
	return &voiceTextChunker{flushChars: flushChars, firstFlushChars: firstFlushChars}
}

func (c *voiceTextChunker) Append(delta string) {
	c.buffer.WriteString(delta)
}

func (c *voiceTextChunker) FlushReady(force bool) []string {
	text := c.buffer.String()
	if strings.TrimSpace(text) == "" {
		return nil
	}
	if force {
		c.buffer.Reset()
		c.firstFlushed = true
		return []string{text}
	}
	threshold := c.flushChars
	if !c.firstFlushed && c.firstFlushChars > 0 {
		threshold = c.firstFlushChars
	}
	var ready []string
	for {
		cut := findVoiceChunkBoundary(text, threshold)
		if cut <= 0 {
			break
		}
		ready = append(ready, text[:cut])
		text = text[cut:]
		c.firstFlushed = true
		threshold = c.flushChars
	}
	if len(ready) > 0 {
		c.buffer.Reset()
		c.buffer.WriteString(text)
	}
	return ready
}

func findVoiceChunkBoundary(text string, flushChars int) int {
	runes := []rune(text)
	if len(runes) == 0 {
		return 0
	}
	for i, r := range runes {
		if strings.ContainsRune("。！？!?；;", r) {
			return len(string(runes[:i+1]))
		}
	}
	for i, r := range runes {
		if i+1 >= flushChars && strings.ContainsRune("，,：:", r) {
			return len(string(runes[:i+1]))
		}
	}
	if len(runes) >= flushChars {
		return len(string(runes[:flushChars]))
	}
	return 0
}

func normalizePCM(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("pcm音频数据为空")
	}
	return data, nil
}
