//go:build ai

package impl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"nova-factory-server/app/business/shop/api/service"
)

type IApiShopASRServiceImpl struct{}

func NewIApiShopASRServiceImpl() service.IApiShopASRService {
	return &IApiShopASRServiceImpl{}
}

type xfyunLFASRClient struct {
	appID     string
	secretKey string
	language  string
	client    *http.Client
}

type xfyunLFASRBaseResp struct {
	Code     string          `json:"code"`
	DescInfo string          `json:"descInfo"`
	Content  json.RawMessage `json:"content"`
}

type xfyunLFASRUploadContent struct {
	OrderID          string `json:"orderId"`
	TaskEstimateTime int64  `json:"taskEstimateTime"`
}

type xfyunLFASRGetResultContent struct {
	OrderInfo struct {
		OrderID          string `json:"orderId"`
		Status           int    `json:"status"`
		FailType         int    `json:"failType"`
		OriginalDuration int64  `json:"originalDuration"`
		RealDuration     int64  `json:"realDuration"`
	} `json:"orderInfo"`
	OrderResult      string `json:"orderResult"`
	TaskEstimateTime int64  `json:"taskEstimateTime"`
}

type xfyunLFASROrderResult struct {
	Lattice []struct {
		JSON1Best string `json:"json_1best"`
	} `json:"lattice"`
}

type xfyunLFASRBestResult struct {
	ST struct {
		RT []struct {
			WS []struct {
				CW []struct {
					W  string `json:"w"`
					WP string `json:"wp"`
				} `json:"cw"`
			} `json:"ws"`
		} `json:"rt"`
	} `json:"st"`
}

func (s *IApiShopASRServiceImpl) Recognize(ctx context.Context, audioBase64 string) (string, error) {
	audioBase64 = strings.TrimSpace(audioBase64)
	if audioBase64 == "" {
		return "", errors.New("audio base64 is empty")
	}
	pcm, err := base64.StdEncoding.DecodeString(audioBase64)
	if err != nil {
		return "", err
	}
	client := newXFYunASRClient()
	return client.Transcribe(ctx, pcm)
}
