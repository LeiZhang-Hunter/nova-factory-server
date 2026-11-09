// Package catalogwire defines the Nova MQTT catalog request/reply protocol.
package catalogwire

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	Prefix   = "nova/v1/gateways/"
	MaxBytes = 16 << 20
	MaxAge   = time.Minute
)

type Request struct {
	RequestID string `json:"requestId"`
	GatewayID string `json:"gatewayId"`
	Username  string `json:"username"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

type Response struct {
	RequestID string          `json:"requestId"`
	GatewayID string          `json:"gatewayId"`
	Error     string          `json:"error,omitempty"`
	Catalog   json.RawMessage `json:"catalog,omitempty"`
	Signature string          `json:"signature"`
}

func RequestTopic(id string) string {
	return Prefix + id + "/catalog/request"
}

func ResponseTopic(id, request string) string {
	return Prefix + id + "/catalog/response/" + request
}

func ValidID(id string) bool {
	n, e := strconv.ParseInt(id, 10, 64)
	return e == nil && n > 0 && strconv.FormatInt(n, 10) == id
}

func ValidRequestID(id string) bool {
	b, e := hex.DecodeString(id)
	return e == nil && len(b) == 24
}

func TopicGateway(topic string) string {
	if !strings.HasPrefix(topic, Prefix) || !strings.HasSuffix(topic, "/catalog/request") {
		return ""
	}
	id := strings.TrimSuffix(strings.TrimPrefix(topic, Prefix), "/catalog/request")
	if !ValidID(id) {
		return ""
	}
	return id
}

func mac(key string, body []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

func (r Request) Sign(key string) string {
	r.Signature = ""
	b, _ := json.Marshal(r)
	return mac(key, b)
}

func (r Request) Verify(key string, now time.Time) bool {
	age := now.Sub(time.Unix(r.Timestamp, 0))
	return key != "" && ValidID(r.GatewayID) && ValidRequestID(r.RequestID) && r.Username != "" && age >= -MaxAge && age <= MaxAge &&
		hmac.Equal([]byte(r.Sign(key)), []byte(r.Signature))
}

func (r Response) Sign(key string) string {
	r.Signature = ""
	b, _ := json.Marshal(r)
	return mac(key, b)
}

func (r Response) Verify(key string, req Request) error {
	if r.RequestID != req.RequestID || r.GatewayID != req.GatewayID || !hmac.Equal([]byte(r.Sign(key)), []byte(r.Signature)) {
		return errors.New("catalog response identity or signature invalid")
	}
	if r.Error != "" {
		return errors.New("cloud catalog request failed; check cloud logs")
	}
	return nil
}

// Authentic validates identity and signature, including signed error replies.
func (r Response) Authentic(key string, req Request) bool {
	return key != "" && r.RequestID == req.RequestID && r.GatewayID == req.GatewayID &&
		hmac.Equal([]byte(r.Sign(key)), []byte(r.Signature))
}
