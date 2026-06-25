package service

import (
	"nova-factory-server/app/business/shop/config/models"

	"github.com/gin-gonic/gin"
)

// ILogisticsCompanyService ERP物流公司服务接口
type ILogisticsCompanyService interface {
	Create(c *gin.Context, req *models.LogisticsCompanyUpsert) (*models.LogisticsCompany, error)
	Update(c *gin.Context, req *models.LogisticsCompanyUpsert) (*models.LogisticsCompany, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.LogisticsCompany, error)
	GetByCode(c *gin.Context, code string) (*models.LogisticsCompany, error)
	GetByName(c *gin.Context, name string) (*models.LogisticsCompany, error)
	List(c *gin.Context, req *models.LogisticsCompanyQuery) (*models.LogisticsCompanyListData, error)
}
