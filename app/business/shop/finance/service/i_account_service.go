package service

import (
	"nova-factory-server/app/business/shop/finance/models"

	"github.com/gin-gonic/gin"
)

// IAccountService ERP 结算账户服务接口
type IAccountService interface {
	Create(c *gin.Context, req *models.AccountUpsert) (*models.Account, error)
	Update(c *gin.Context, req *models.AccountUpsert) (*models.Account, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.Account, error)
	List(c *gin.Context, req *models.AccountQuery) (*models.AccountListData, error)
}
