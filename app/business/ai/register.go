//go:build ai
// +build ai

package ai

import (
	"nova-factory-server/app/business/ai/agent/aidatasetcontroller"
	"nova-factory-server/app/business/ai/agent/aidatasetdao/aiDataSetDaoImpl"
	"nova-factory-server/app/business/ai/agent/aidatasetservice/aidatasetserviceimpl"
	"nova-factory-server/app/business/ai/gateway/gatewaycontroller"
	"nova-factory-server/app/business/ai/gateway/gatewaydao/gatewaydaoimpl"
	"nova-factory-server/app/business/ai/gateway/gatewayservice/gatewayserviceimpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	aiDataSetDaoImpl.ProviderSet,
	aidatasetserviceimpl.ProviderSet,
	aidatasetcontroller.ProviderSet,
	gatewaydaoimpl.ProviderSet,
	gatewayserviceimpl.ProviderSet,
	gatewaycontroller.ProviderSet,
	GinProviderSet,
)
