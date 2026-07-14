package requestlog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	StatusPending int32 = 0
	StatusSuccess int32 = 1
	StatusFailed  int32 = 2
)

const (
	TypeGjpQqdAPI          = "gjpqqd_api"
	TypeWechatPayNotify    = "shop_wechat_pay_notify"
	TypeWechatRefundNotify = "shop_wechat_refund_notify"
)

const (
	SourceGjpQqdAPI            = "datasyncapi/gjpqqd/api"
	SourceWechatPayNotify      = "shop/api/order/notify"
	SourceWechatRefundNotify   = "shop/api/order/refund_notify"
	DefaultRequestLogTableName = "sys_request_log"
)

type Entry struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id,string"`
	LogType       string `gorm:"column:log_type" json:"logType"`
	SourceModule  string `gorm:"column:source_module" json:"sourceModule"`
	Status        int32  `gorm:"column:status" json:"status"`
	RequestMethod string `gorm:"column:request_method" json:"requestMethod"`
	RequestPath   string `gorm:"column:request_path" json:"requestPath"`
	QueryString   string `gorm:"column:query_string" json:"queryString"`
	HeadersJSON   string `gorm:"column:headers_json" json:"headersJson"`
	BodyText      string `gorm:"column:body_text" json:"bodyText"`
	ClientIP      string `gorm:"column:client_ip" json:"clientIp"`
	UserAgent     string `gorm:"column:user_agent" json:"userAgent"`
	ErrorMessage  string `gorm:"column:error_message" json:"errorMessage"`
}

type Store interface {
	Create(ctx context.Context, entry *Entry) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status int32, errorMessage string) error
}

type Service interface {
	Create(ctx context.Context, entry *Entry) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status int32, errorMessage string) error
}

type emptyStore struct{}

func NewEmptyStore() Store {
	return emptyStore{}
}

func (emptyStore) Create(context.Context, *Entry) (int64, error) {
	return 0, nil
}

func (emptyStore) UpdateStatus(context.Context, int64, int32, string) error {
	return nil
}

const (
	CtxKeyRequestLogID    = "requestlog_id"
	CtxKeyRequestLogError = "requestlog_error"
)

// Middleware 创建请求日志记录中间件
// 自动在请求前创建日志条目，请求后根据 context 标记更新状态。
// handler 中通过 MarkError(c, msg) 标记失败；否则默认标记成功。
func Middleware(logType, sourceModule string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := RestoreBody(c)
		logID := CaptureGin(c, logType, sourceModule, body)
		c.Set(CtxKeyRequestLogID, logID)

		c.Next()

		if errMsg, exists := c.Get(CtxKeyRequestLogError); exists && errMsg.(string) != "" {
			MarkFailed(c, logID, errMsg.(string))
		} else {
			MarkSuccess(c, logID)
		}
	}
}

// MarkError 在 handler 中标记当前请求失败（配合 Middleware 使用）
func MarkError(c *gin.Context, msg string) {
	c.Set(CtxKeyRequestLogError, msg)
}

func CaptureGin(c *gin.Context, logType, sourceModule, body string) int64 {
	headers, _ := json.Marshal(headerMap(c.Request.Header))
	id, err := GetStore().Create(c, &Entry{
		LogType:       logType,
		SourceModule:  sourceModule,
		Status:        StatusPending,
		RequestMethod: c.Request.Method,
		RequestPath:   c.Request.URL.Path,
		QueryString:   c.Request.URL.RawQuery,
		HeadersJSON:   string(headers),
		BodyText:      body,
		ClientIP:      c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
	})
	if err != nil {
		zap.L().Error("create request log failed", zap.Error(err))
		return 0
	}
	return id
}

func RegisterService(s Service) {
	if s == nil {
		return
	}
	RegisterStore(serviceStore{s})
}

type serviceStore struct {
	Service
}

func (s serviceStore) Create(ctx context.Context, entry *Entry) (int64, error) {
	return s.Service.Create(ctx, entry)
}

func (s serviceStore) UpdateStatus(ctx context.Context, id int64, status int32, errorMessage string) error {
	return s.Service.UpdateStatus(ctx, id, status, errorMessage)
}

func MarkSuccess(ctx context.Context, id int64) {
	if id <= 0 {
		return
	}
	_ = GetStore().UpdateStatus(ctx, id, StatusSuccess, "")
}

func MarkFailed(ctx context.Context, id int64, message string) {
	if id <= 0 {
		return
	}
	_ = GetStore().UpdateStatus(ctx, id, StatusFailed, message)
}

func RestoreBody(c *gin.Context) (string, error) {
	if c == nil || c.Request == nil || c.Request.Body == nil {
		return "", nil
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return string(bodyBytes), nil
}

func headerMap(header http.Header) map[string][]string {
	m := make(map[string][]string, len(header))
	for k, v := range header {
		m[k] = v
	}
	return m
}
