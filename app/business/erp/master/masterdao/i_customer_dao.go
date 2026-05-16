package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// ICustomerDao ERP 客户数据访问接口
type ICustomerDao interface {
	Create(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error)
	Update(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Customer, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Customer, error)
	ListPage(c *gin.Context, req *mastermodels.CustomerQuery) (*erpbiz.PageResult[mastermodels.Customer], error)
	List(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error)
}
