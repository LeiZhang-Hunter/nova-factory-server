package dao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/config/models"
)

// ILogisticsCompanyDao ERP物流公司数据访问接口
type ILogisticsCompanyDao interface {
	Create(c *gin.Context, req *models.LogisticsCompanyUpsert) (*models.LogisticsCompany, error)
	Update(c *gin.Context, req *models.LogisticsCompanyUpsert) (*models.LogisticsCompany, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.LogisticsCompany, error)
	GetByCode(c *gin.Context, code string) (*models.LogisticsCompany, error)
	GetByName(c *gin.Context, name string) (*models.LogisticsCompany, error)
	List(c *gin.Context, req *models.LogisticsCompanyQuery) (*models.LogisticsCompanyListData, error)
	ListByCodes(c *gin.Context, codes []string) ([]*models.LogisticsCompany, error)
}
