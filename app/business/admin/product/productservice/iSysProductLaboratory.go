package productservice

import (
	"nova-factory-server/app/business/admin/product/productmodels"

	"github.com/gin-gonic/gin"
)

type ISysProductLaboratoryService interface {
	SelectLaboratoryList(c *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (list *productmodels.SysProductLaboratoryList, err error)
	SelectLaboratoryById(c *gin.Context, id int64) (*productmodels.SysProductLaboratoryVo, error)
	Set(c *gin.Context, data *productmodels.SysProductLaboratoryVo) (*productmodels.SysProductLaboratory, error)
	DeleteLaboratoryByIds(c *gin.Context, ids []int64) error
	// SelectUserLaboratoryList 读取用户化验单
	SelectUserLaboratoryList(ctx *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (list *productmodels.SysProductLaboratoryList, err error)
	// FirstLaboratoryInfo 读取化验单
	FirstLaboratoryInfo(ctx *gin.Context, req *productmodels.SysProductLaboratoryInfoDQL) (*productmodels.SysProductLaboratory, error)
	// FirstLaboratoryList 最新化验单列表
	FirstLaboratoryList(ctx *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (*productmodels.SysProductLaboratoryList, error)
}
