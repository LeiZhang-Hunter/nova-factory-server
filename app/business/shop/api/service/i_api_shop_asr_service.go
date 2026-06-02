package service

import (
	"context"
)

type IApiShopASRService interface {
	Recognize(ctx context.Context, audioBase64 string) (string, error)
}
