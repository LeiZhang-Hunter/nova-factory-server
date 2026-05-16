package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IAccountDao ERP 结算账户数据访问接口
type IAccountDao interface {
	Create(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error)
	Update(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Account, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Account, error)
	ListPage(c *gin.Context, req *mastermodels.AccountQuery) (*erpbiz.PageResult[mastermodels.Account], error)
	List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error)
}
