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
	// SelectUserLaboratoryList 读取用户化验单
	SelectUserLaboratoryList(ctx *gin.Context, dql *productModels.SysProductLaboratoryDQL) (list *productModels.SysProductLaboratoryList, err error)
	// FirstLaboratoryInfo 读取最新信息
	FirstLaboratoryInfo(ctx *gin.Context) (*productModels.SysProductLaboratory, error)
	// FirstLaboratoryList 最新化验单列表
	FirstLaboratoryList(ctx *gin.Context, dql *productModels.SysProductLaboratoryDQL) (*productModels.SysProductLaboratoryList, error)
}
