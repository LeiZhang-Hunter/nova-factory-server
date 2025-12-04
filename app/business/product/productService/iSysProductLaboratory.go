package productService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/product/productModels"
)

type ISysProductLaboratoryService interface {
	SelectLaboratoryList(c *gin.Context, dql *productModels.SysProductLaboratoryDQL) (list *productModels.SysProductLaboratoryList, err error)
	SelectLaboratoryById(c *gin.Context, id int64) (*productModels.SysProductLaboratoryVo, error)
	Set(c *gin.Context, data *productModels.SysProductLaboratoryVo) (*productModels.SysProductLaboratory, error)
	DeleteLaboratoryByIds(c *gin.Context, ids []int64) error
}
