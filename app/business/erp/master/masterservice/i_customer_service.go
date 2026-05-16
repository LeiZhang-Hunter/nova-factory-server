package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// ICustomerService ERP 客户服务接口
type ICustomerService interface {
	Create(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error)
	Update(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Customer, error)
	List(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error)
}
