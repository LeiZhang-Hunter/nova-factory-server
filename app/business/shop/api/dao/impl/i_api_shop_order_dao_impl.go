package impl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/shop/api/dao"
)

type IApiShopOrderDaoImpl struct{}

// NewIApiShopOrderDaoImpl 初始化空order，没有erp的时候会走这个
func NewIApiShopOrderDaoImpl() dao.IApiShopOrderDao {
	return &IApiShopOrderDaoImpl{}
}

// Set 新增或修改 ERP 订单及其子表。
func (i *IApiShopOrderDaoImpl) Set(c *gin.Context, req *salemodels.OrderSet) (*salemodels.Order, error) {
	return nil, errors.New("not implemented")
}

// SetWithTx 新增或修改 ERP 订单及其子表（带事务）。
func (i *IApiShopOrderDaoImpl) SetWithTx(c *gin.Context, tx *gorm.DB, req *salemodels.OrderSet) (*salemodels.Order, error) {
	return nil, errors.New("not implemented")
}

// GetByID 查询 ERP 订单详情。
func (i *IApiShopOrderDaoImpl) GetByID(c *gin.Context, id uint64) (*salemodels.Order, error) {
	return nil, errors.New("not implemented")
}

// List 分页查询 ERP 订单。
func (i *IApiShopOrderDaoImpl) List(c *gin.Context, req *salemodels.OrderQuery) (*salemodels.OrderListData, error) {
	return nil, errors.New("not implemented")
}

// DeleteByIDs 删除 ERP 订单。
func (i *IApiShopOrderDaoImpl) DeleteByIDs(c *gin.Context, ids []uint64) error {
	return nil
}
