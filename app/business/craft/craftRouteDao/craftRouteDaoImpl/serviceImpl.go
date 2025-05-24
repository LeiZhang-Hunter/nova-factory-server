package craftRouteDaoImpl

import (
	"github.com/google/wire"
	"nova-factory-server/app/business/gateway/gatewayDao/gatewayDaoImpl"
)

var ProviderSet = wire.NewSet(gatewayDaoImpl.NewISysGatewayInboundConfigDaoImpl)
