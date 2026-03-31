package systemservice

import (
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type IDictDataService interface {
	SelectDictDataByType(c *gin.Context, dictType string) (data []byte)
	SelectDictDataList(c *gin.Context, dictData *systemmodels.SysDictDataDQL) (list []*systemmodels.SysDictDataVo, total int64)
	ExportDictData(c *gin.Context, dictData *systemmodels.SysDictDataDQL) (data []byte)
	SelectDictDataById(c *gin.Context, dictCode int64) (dictData *systemmodels.SysDictDataVo)
	InsertDictData(c *gin.Context, dictData *systemmodels.SysDictDataVo)
	UpdateDictData(c *gin.Context, dictData *systemmodels.SysDictDataVo)
	DeleteDictDataByIds(c *gin.Context, dictCodes []int64)
	CheckDictDataByTypes(c *gin.Context, dictType []string) bool
}
