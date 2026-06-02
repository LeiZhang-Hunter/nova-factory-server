//go:build ai

package impl

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func buildXFYunAuthURL(rawURL, apiKey, apiSecret string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse websocket url: %w", err)
	}
	requestURI := u.EscapedPath()
	if u.RawQuery != "" {
		requestURI += "?" + u.RawQuery
	}
	date := time.Now().UTC().Format(http.TimeFormat)
	signatureOrigin := fmt.Sprintf("host: %s\ndate: %s\nGET %s HTTP/1.1", u.Host, date, requestURI)
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(signatureOrigin))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorizationOrigin := fmt.Sprintf("api_key=\"%s\", algorithm=\"hmac-sha256\", headers=\"host date request-line\", signature=\"%s\"", apiKey, signature)
	query := u.Query()
	query.Set("host", u.Host)
	query.Set("date", date)
	query.Set("authorization", base64.StdEncoding.EncodeToString([]byte(authorizationOrigin)))
	u.RawQuery = query.Encode()
	return u.String(), nil
}

type xfyunASRClient struct {
	appID     string
	apiKey    string
	apiSecret string
	url       string
}

func newXFYunASRClient() *xfyunASRClient {
	return &xfyunASRClient{
		appID:     strings.TrimSpace(viper.GetString("voice.xfyun.app_id")),
		apiKey:    strings.TrimSpace(viper.GetString("voice.xfyun.api_key")),
		apiSecret: strings.TrimSpace(viper.GetString("voice.xfyun.api_secret")),
		url:       defaultString(strings.TrimSpace(viper.GetString("voice.xfyun.asr_url")), "wss://iat-api.xfyun.cn/v2/iat"),
	}
}

type asrRequest struct {
	Common struct {
		AppID string `json:"app_id"`
	} `json:"common,omitempty"`
	Business map[string]any `json:"business,omitempty"`
	Data     *asrData       `json:"data,omitempty"`
}

type asrData struct {
	Status   int    `json:"status"`
	Format   string `json:"format"`
	Encoding string `json:"encoding"`
	Audio    string `json:"audio"`
}

type asrResponse struct {
	Action  string          `json:"action"`
	Code    asrInt          `json:"code"`
	Message string          `json:"message"`
	Desc    string          `json:"desc"`
	SID     string          `json:"sid"`
	Data    json.RawMessage `json:"data"`
}

type asrInt int

func (v *asrInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*v = 0
		return nil
	}

	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*v = asrInt(intValue)
		return nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		stringValue = strings.TrimSpace(stringValue)
		if stringValue == "" {
			*v = 0
			return nil
		}
		parsed, err := strconv.Atoi(stringValue)
		if err != nil {
			return err
		}
		*v = asrInt(parsed)
		return nil
	}

	return fmt.Errorf("unsupported asr int payload: %s", string(data))
}

func (c *xfyunASRClient) Transcribe(ctx context.Context, pcm []byte) (string, error) {
	if len(pcm) == 0 {
		return "", fmt.Errorf("pcm数据为空")
	}
	if c.appID == "" || c.apiKey == "" || c.apiSecret == "" {
		return "", fmt.Errorf("讯飞ASR配置缺失")
	}
	authURL, err := buildXFYunAuthURL(c.url, c.apiKey, c.apiSecret)
	if err != nil {
		return "", err
	}
	dialer := websocket.Dialer{HandshakeTimeout: 15 * time.Second}
	conn, _, err := dialer.DialContext(ctx, authURL, nil)
	if err != nil {
		return "", fmt.Errorf("dial asr websocket: %w", err)
	}
	defer conn.Close()
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.UnderlyingConn().SetDeadline(deadline)
	}
	frames := splitAudioFrames(pcm, 1280)
	zap.L().Info("[DEBUG-ASR-UPLOAD] xfyun asr connect",
		zap.String("url", c.url),
		zap.Int("pcm_bytes", len(pcm)),
		zap.Int("frame_count", len(frames)),
	)
	for i, frame := range frames {
		status := 1
		if i == 0 {
			status = 0
		}
		if i == len(frames)-1 {
			status = 1
		}
		req := asrRequest{}
		if status == 0 {
			req.Common.AppID = c.appID
			req.Business = map[string]any{"language": "zh_cn", "domain": "iat", "accent": "mandarin", "eos": 10000}
		}
		req.Data = &asrData{Status: status, Format: "audio/L16;rate=16000", Encoding: "raw", Audio: base64.StdEncoding.EncodeToString(frame)}
		zap.L().Info("[DEBUG-ASR-UPLOAD] xfyun asr write frame",
			zap.Int("frame_index", i),
			zap.Int("status", status),
			zap.Int("frame_bytes", len(frame)),
			zap.Int("audio_base64_len", len(req.Data.Audio)),
			zap.String("audio_base64_preview", previewDebugText(req.Data.Audio, 96)),
		)
		if err := conn.WriteJSON(req); err != nil {
			return "", fmt.Errorf("write asr frame: %w", err)
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(40 * time.Millisecond):
		}
	}
	if err := conn.WriteJSON(map[string]any{
		"data": map[string]any{
			"status": 2,
		},
	}); err != nil {
		return "", fmt.Errorf("write asr finish frame: %w", err)
	}
	var transcript strings.Builder
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				break
			}
			return "", fmt.Errorf("read asr response: %w", err)
		}
		zap.L().Info("[DEBUG-ASR-UPLOAD] xfyun asr raw response",
			zap.Int("bytes", len(message)),
			zap.String("payload", previewDebugText(string(message), 512)),
		)
		var resp asrResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			return "", fmt.Errorf("decode asr response: %w", err)
		}
		if int(resp.Code) != 0 {
			errorMessage := strings.TrimSpace(resp.Desc)
			if errorMessage == "" {
				errorMessage = strings.TrimSpace(resp.Message)
			}
			if errorMessage == "" {
				errorMessage = "unknown asr error"
			}
			if strings.TrimSpace(resp.Action) != "" || strings.TrimSpace(resp.SID) != "" {
				return "", fmt.Errorf("asr failed: code=%d action=%s sid=%s message=%s", int(resp.Code), resp.Action, resp.SID, errorMessage)
			}
			return "", fmt.Errorf("asr failed: code=%d message=%s", int(resp.Code), errorMessage)
		}

		var data struct {
			Status int `json:"status"`
			Result struct {
				Ws []struct {
					Cw []struct {
						W string `json:"w"`
					} `json:"cw"`
				} `json:"ws"`
				PGS string `json:"pgs"`
				RG  []int  `json:"rg"`
				SN  int    `json:"sn"`
			} `json:"result"`
		}
		if len(resp.Data) > 0 && string(resp.Data) != `""` && string(resp.Data) != "null" {
			if err := json.Unmarshal(resp.Data, &data); err != nil {
				return "", fmt.Errorf("decode asr data: %w", err)
			}
		}

		if data.Result.PGS == "rpl" {
			transcript.Reset()
		}
		for _, ws := range data.Result.Ws {
			for _, cw := range ws.Cw {
				transcript.WriteString(cw.W)
			}
		}
		if data.Status == 2 {
			break
		}
	}
	text := strings.TrimSpace(cleanASRText(transcript.String()))
	if text == "" {
		return "", fmt.Errorf("asr returned empty transcript")
	}
	return text, nil
}

func previewDebugText(text string, limit int) string {
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return text[:limit] + "...(truncated)"
}

func splitAudioFrames(data []byte, size int) [][]byte {
	if len(data) <= size {
		return [][]byte{data}
	}
	frames := make([][]byte, 0, (len(data)+size-1)/size)
	for start := 0; start < len(data); start += size {
		end := start + size
		if end > len(data) {
			end = len(data)
		}
		frames = append(frames, data[start:end])
	}
	return frames
}

func cleanASRText(text string) string {
	replacer := strings.NewReplacer("<s>", "", "</s>", "", "[noise]", "")
	return replacer.Replace(text)
}

type xfyunTTSClient struct {
	appID      string
	apiKey     string
	apiSecret  string
	url        string
	voice      string
	oralLevel  string
	speed      int
	volume     int
	pitch      int
	encoding   string
	sampleRate int
}

func newXFYunTTSClient() *xfyunTTSClient {
	return &xfyunTTSClient{
		appID:      strings.TrimSpace(viper.GetString("voice.xfyun.app_id")),
		apiKey:     strings.TrimSpace(viper.GetString("voice.xfyun.api_key")),
		apiSecret:  strings.TrimSpace(viper.GetString("voice.xfyun.api_secret")),
		url:        strings.TrimSpace(viper.GetString("voice.xfyun.tts_url")),
		voice:      defaultString(strings.TrimSpace(viper.GetString("voice.xfyun.tts_voice")), "x5_lingxiaoxuan_flow"),
		oralLevel:  strings.TrimSpace(viper.GetString("voice.xfyun.tts_oral_level")),
		speed:      defaultInt(viper.GetInt("voice.xfyun.tts_speed"), 50),
		volume:     defaultInt(viper.GetInt("voice.xfyun.tts_volume"), 50),
		pitch:      defaultInt(viper.GetInt("voice.xfyun.tts_pitch"), 50),
		encoding:   defaultString(strings.TrimSpace(viper.GetString("voice.xfyun.tts_encoding")), "raw"),
		sampleRate: defaultInt(viper.GetInt("voice.xfyun.tts_sample_rate"), 24000),
	}
}

type xfyunTTSStreamSession struct {
	client       *xfyunTTSClient
	conn         *websocket.Conn
	onAudioChunk func([]byte) error
	errCh        chan error
	closeOnce    sync.Once
	openedAt     time.Time
}

type ttsFrame struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  int    `json:"status"`
		SID     string `json:"sid"`
	} `json:"header"`
	Payload struct {
		Audio *struct {
			Status int    `json:"status"`
			Seq    int    `json:"seq"`
			Audio  string `json:"audio"`
		} `json:"audio"`
	} `json:"payload"`
}

func (c *xfyunTTSClient) buildStreamMessage(text string, status int, seq int) map[string]any {
	message := map[string]any{
		"header": map[string]any{"app_id": c.appID, "status": status},
		"parameter": map[string]any{"tts": map[string]any{
			"vcn": c.voice, "speed": c.speed, "volume": c.volume, "pitch": c.pitch, "bgs": 0, "reg": 0, "rdn": 0, "rhy": 0,
			"audio": map[string]any{"encoding": c.encoding, "sample_rate": c.sampleRate, "channels": 1, "bit_depth": 16, "frame_size": 0},
		}},
		"payload": map[string]any{"text": map[string]any{"encoding": "utf8", "compress": "raw", "format": "plain", "status": status, "seq": seq, "text": base64.StdEncoding.EncodeToString([]byte(text))}},
	}
	if c.oralLevel != "" {
		message["parameter"].(map[string]any)["oral"] = map[string]any{"oral_level": c.oralLevel, "spark_assist": 1, "stop_split": 0, "remain": 0}
	}
	return message
}

func (c *xfyunTTSClient) OpenStream(ctx context.Context, onAudioChunk func([]byte) error) (*xfyunTTSStreamSession, error) {
	if c.appID == "" || c.apiKey == "" || c.apiSecret == "" || c.url == "" {
		return nil, fmt.Errorf("讯飞TTS配置缺失")
	}
	authURL, err := buildXFYunAuthURL(c.url, c.apiKey, c.apiSecret)
	if err != nil {
		return nil, err
	}
	dialer := websocket.Dialer{HandshakeTimeout: 15 * time.Second}
	conn, _, err := dialer.DialContext(ctx, authURL, nil)
	if err != nil {
		return nil, fmt.Errorf("dial tts websocket: %w", err)
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.UnderlyingConn().SetDeadline(deadline)
	}
	session := &xfyunTTSStreamSession{client: c, conn: conn, onAudioChunk: onAudioChunk, errCh: make(chan error, 1), openedAt: time.Now()}
	zap.L().Info("[voice-debug] xfyun tts stream opened",
		zap.String("encoding", c.encoding),
		zap.Int("sample_rate", c.sampleRate),
		zap.String("voice", c.voice),
	)
	go session.readLoop()
	return session, nil
}

func (s *xfyunTTSStreamSession) SendText(text string, status int, seq int) error {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	zap.L().Info("[voice-debug] xfyun tts write text",
		zap.Int("seq", seq),
		zap.Int("status", status),
		zap.Int("text_runes", utf8.RuneCountInString(text)),
		zap.Duration("stream_elapsed", time.Since(s.openedAt)),
	)
	if err := s.conn.WriteJSON(s.client.buildStreamMessage(text, status, seq)); err != nil {
		zap.L().Warn("[voice-debug] xfyun tts write text failed",
			zap.Int("seq", seq),
			zap.Int("status", status),
			zap.Int("text_runes", utf8.RuneCountInString(text)),
			zap.Duration("stream_elapsed", time.Since(s.openedAt)),
			zap.Error(err),
		)
		return fmt.Errorf("write tts stream chunk: %w", err)
	}
	return nil
}

func (s *xfyunTTSStreamSession) Finish(seq int) error {
	zap.L().Info("[voice-debug] xfyun tts write finish",
		zap.Int("seq", seq),
		zap.Duration("stream_elapsed", time.Since(s.openedAt)),
	)
	if err := s.conn.WriteJSON(s.client.buildStreamMessage("", 2, seq)); err != nil {
		zap.L().Warn("[voice-debug] xfyun tts write finish failed", zap.Int("seq", seq), zap.Duration("stream_elapsed", time.Since(s.openedAt)), zap.Error(err))
		return fmt.Errorf("write tts finish chunk: %w", err)
	}
	return <-s.errCh
}

func (s *xfyunTTSStreamSession) readLoop() {
	defer close(s.errCh)
	for {
		_, raw, err := s.conn.ReadMessage()
		if err != nil {
			zap.L().Warn("[voice-debug] xfyun tts read failed", zap.Duration("stream_elapsed", time.Since(s.openedAt)), zap.Error(err))
			s.errCh <- fmt.Errorf("read tts response: %w", err)
			return
		}
		var frame ttsFrame
		if err := json.Unmarshal(raw, &frame); err != nil {
			s.errCh <- fmt.Errorf("decode tts response: %w", err)
			return
		}
		if frame.Header.Code != 0 {
			zap.L().Warn("[voice-debug] xfyun tts frame failed",
				zap.Int("code", frame.Header.Code),
				zap.String("message", frame.Header.Message),
				zap.String("sid", frame.Header.SID),
				zap.Int("status", frame.Header.Status),
				zap.Duration("stream_elapsed", time.Since(s.openedAt)),
			)
			s.errCh <- fmt.Errorf("tts failed: code=%d message=%s sid=%s", frame.Header.Code, frame.Header.Message, frame.Header.SID)
			return
		}
		if frame.Payload.Audio != nil && frame.Payload.Audio.Audio != "" {
			decoded, err := base64.StdEncoding.DecodeString(frame.Payload.Audio.Audio)
			if err != nil {
				s.errCh <- fmt.Errorf("decode tts audio: %w", err)
				return
			}
			zap.L().Info("[voice-debug] xfyun tts read audio",
				zap.Int("seq", frame.Payload.Audio.Seq),
				zap.Int("status", frame.Payload.Audio.Status),
				zap.Int("bytes", len(decoded)),
				zap.String("sid", frame.Header.SID),
				zap.Duration("stream_elapsed", time.Since(s.openedAt)),
			)
			if s.onAudioChunk != nil {
				if err := s.onAudioChunk(decoded); err != nil {
					zap.L().Warn("[voice-debug] xfyun tts audio callback failed", zap.Int("seq", frame.Payload.Audio.Seq), zap.Error(err))
					s.errCh <- err
					return
				}
			}
		}
		if frame.Header.Status == 2 || (frame.Payload.Audio != nil && frame.Payload.Audio.Status == 2) {
			zap.L().Info("[voice-debug] xfyun tts stream complete",
				zap.String("sid", frame.Header.SID),
				zap.Int("header_status", frame.Header.Status),
				zap.Duration("stream_elapsed", time.Since(s.openedAt)),
			)
			s.errCh <- nil
			return
		}
	}
}

func (s *xfyunTTSStreamSession) Close() error {
	var err error
	s.closeOnce.Do(func() {
		if s.conn != nil {
			zap.L().Info("[voice-debug] xfyun tts stream close", zap.Duration("stream_elapsed", time.Since(s.openedAt)))
			err = s.conn.Close()
		}
	})
	return err
}

func defaultInt(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func synthesizeOrdered(_ context.Context, _ *xfyunTTSClient, _ string) ([]byte, error) {
	return nil, nil
}

func mergeOrderedChunks(chunks map[int][]byte) []byte {
	order := make([]int, 0, len(chunks))
	for seq := range chunks {
		order = append(order, seq)
	}
	sort.Ints(order)
	var out []byte
	for _, seq := range order {
		out = append(out, chunks[seq]...)
	}
	return out
}
