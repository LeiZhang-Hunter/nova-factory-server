package dao

import (
	"nova-factory-server/app/business/shop/finance/models"

	"github.com/gin-gonic/gin"
)

// IAccountDao ERP 结算账户数据访问接口
type IAccountDao interface {
	Create(c *gin.Context, req *models.AccountUpsert) (*models.Account, error)
	Update(c *gin.Context, req *models.AccountUpsert) (*models.Account, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.Account, error)
	GetByColumn(c *gin.Context, column string, value any) (*models.Account, error)
	ListPage(c *gin.Context, req *models.AccountQuery) (*models.AccountListData, error)
	List(c *gin.Context, req *models.AccountQuery) (*models.AccountListData, error)
	ResetDefaultStatus(c *gin.Context, excludeID int64) error
	GetDefaultFromCache(c *gin.Context) *models.Account
	// GetDefaultAccountNo 从缓存读取默认结算账户编码，缓存未命中返回空字符串。
	GetDefaultAccountNo(c *gin.Context) string
}
