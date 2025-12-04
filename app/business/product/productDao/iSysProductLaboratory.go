package productDao

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/product/productModels"
)

type ISysProductLaboratoryDao interface {
	SelectLaboratoryList(c *gin.Context, dql *productModels.SysProductLaboratoryDQL) (list *productModels.SysProductLaboratoryList, err error)
	SelectLaboratoryById(ctx context.Context, id int64) (*productModels.SysProductLaboratoryVo, error)
	Set(c *gin.Context, data *productModels.SysProductLaboratoryVo) (*productModels.SysProductLaboratory, error)
	DeleteLaboratoryByIds(ctx context.Context, ids []int64) error
}
