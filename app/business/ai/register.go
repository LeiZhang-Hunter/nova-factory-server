//go:build ai
// +build ai

package ai

import (
	"nova-factory-server/app/business/ai/agent/aidatasetcontroller"
	"nova-factory-server/app/business/ai/agent/aidatasetdao/aiDataSetDaoImpl"
	"nova-factory-server/app/business/ai/agent/aidatasetservice/aidatasetserviceimpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	aiDataSetDaoImpl.ProviderSet,
	aidatasetserviceimpl.ProviderSet,
	aidatasetcontroller.ProviderSet,
	GinProviderSet,
)
