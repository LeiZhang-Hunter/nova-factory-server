//go:build erp
// +build erp

package erp

import (
	"nova-factory-server/app/business/erp/finance/financecontroller"
	"nova-factory-server/app/business/erp/finance/financedao/financedaoimpl"
	"nova-factory-server/app/business/erp/finance/financeservice/financeserviceimpl"
	"nova-factory-server/app/business/erp/master/mastercontroller"
	"nova-factory-server/app/business/erp/master/masterdao/masterdaoimpl"
	"nova-factory-server/app/business/erp/master/masterservice/masterserviceimpl"
	erpObserver "nova-factory-server/app/business/erp/observer"
	"nova-factory-server/app/business/erp/purchase/purchasecontroller"
	"nova-factory-server/app/business/erp/purchase/purchasedao/purchasedaoimpl"
	"nova-factory-server/app/business/erp/purchase/purchaseservice/purchaseserviceimpl"
	"nova-factory-server/app/business/erp/sale/salecontroller"
	"nova-factory-server/app/business/erp/sale/saledao/saledaoimpl"
	"nova-factory-server/app/business/erp/sale/saleservice/saleserviceimpl"
	"nova-factory-server/app/business/erp/setting/settingcontroller"
	"nova-factory-server/app/business/erp/setting/settingdao/settingdaoimpl"
	"nova-factory-server/app/business/erp/setting/settingservice/settingserviceimpl"
	"nova-factory-server/app/business/erp/stock/stockcontroller"
	"nova-factory-server/app/business/erp/stock/stockdao/stockdaoimpl"
	"nova-factory-server/app/business/erp/stock/stockservice/stockserviceimpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	masterdaoimpl.ProviderSet,
	masterserviceimpl.ProviderSet,
	mastercontroller.ProviderSet,
	financedaoimpl.ProviderSet,
	financeserviceimpl.ProviderSet,
	financecontroller.ProviderSet,
	purchasedaoimpl.ProviderSet,
	purchaseserviceimpl.ProviderSet,
	purchasecontroller.ProviderSet,
	saledaoimpl.ProviderSet,
	saleserviceimpl.ProviderSet,
	salecontroller.ProviderSet,
	stockdaoimpl.ProviderSet,
	stockserviceimpl.ProviderSet,
	stockcontroller.ProviderSet,
	settingdaoimpl.ProviderSet,
	settingserviceimpl.ProviderSet,
	settingcontroller.ProviderSet,
	erpObserver.ProviderSet,
	GinProviderSet,
)
