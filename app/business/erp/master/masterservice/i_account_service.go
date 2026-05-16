package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IAccountService ERP 结算账户服务接口
type IAccountService interface {
	Create(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error)
	Update(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Account, error)
	List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error)
}
