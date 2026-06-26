package impl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	logisticsCompanyDao "nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/business/shop/order/service"
)

type IShopOrderShipmentServiceImpl struct {
	orderShipmentDao    dao.IOrderShipmentDao
	logisticsCompanyDao logisticsCompanyDao.ILogisticsCompanyDao
}

// NewIShopOrderShipmentServiceImpl 初始化快递模块
func NewIShopOrderShipmentServiceImpl(orderShipmentDao dao.IOrderShipmentDao,
	logisticsCompanyDao logisticsCompanyDao.ILogisticsCompanyDao) service.IShopOrderShipmentService {
	return &IShopOrderShipmentServiceImpl{
		orderShipmentDao:    orderShipmentDao,
		logisticsCompanyDao: logisticsCompanyDao,
	}
}

func (i *IShopOrderShipmentServiceImpl) ListByOrderID(ctx *gin.Context, orderID uint64) ([]*models.OrderShipment, error) {
	list, err := i.orderShipmentDao.ListByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return list, nil
	}

	codes := make([]string, 0, len(list))
	seen := make(map[string]struct{}, len(list))
	for _, s := range list {
		code := s.Companycode
		if code == "" {
			continue
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, code)
	}
	if len(codes) > 0 {
		companies, err := i.logisticsCompanyDao.ListByCodes(ctx, codes)
		if err != nil {
			return nil, err
		}
		codeToName := make(map[string]string, len(companies))
		for _, c := range companies {
			codeToName[c.Code] = c.Name
		}
		for _, s := range list {
			if name, ok := codeToName[s.Companycode]; ok {
				s.CompanyName = name
			}
		}
	}

	return list, nil
}

func (i *IShopOrderShipmentServiceImpl) BatchInsert(tx *gorm.DB, shipments []*models.OrderShipmentSet) error {
	return i.orderShipmentDao.BatchInsert(tx, shipments)
}
