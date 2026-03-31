//go:build ai
// +build ai

package ai

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetController"
	"nova-factory-server/app/business/ai/agent/aiDataSetDao/aiDataSetDaoImpl"
	"nova-factory-server/app/business/ai/agent/aiDataSetService/aiDataSetServiceImpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	aiDataSetDaoImpl.ProviderSet,
	aiDataSetServiceImpl.ProviderSet,
	aiDataSetController.ProviderSet,
	GinProviderSet,
)
