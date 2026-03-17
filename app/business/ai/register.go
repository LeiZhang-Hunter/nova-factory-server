//go:build shop
// +build shop

package ai

import (
	"nova-factory-server/app/business/ai/aiDataSetController"
	"nova-factory-server/app/business/ai/aiDataSetDao/aiDataSetDaoImpl"
	"nova-factory-server/app/business/ai/aiDataSetService/aiDataSetServiceImpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	aiDataSetDaoImpl.ProviderSet,
	aiDataSetServiceImpl.ProviderSet,
	aiDataSetController.ProviderSet,
)
