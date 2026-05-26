package aidatasetservice

import (
	"context"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
)

type IConfigLoaderService interface {
	All(ctx context.Context) ([]*gatewaymodels.AIAgentConfigPublishHistory, error)
}
