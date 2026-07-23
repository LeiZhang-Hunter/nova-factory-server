package systemservice

import (
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type IDictDataService interface {
	SelectDictDataByType(c *gin.Context, dictType string) (data []byte)
	SelectDictDataList(c *gin.Context, dictData *modelquery.SysDictDataDQL) (list []*modelresponse.SysDictDataVo, total int64)
	ExportDictData(c *gin.Context, dictData *modelquery.SysDictDataDQL) (data []byte)
	SelectDictDataById(c *gin.Context, dictCode int64) (dictData *modelresponse.SysDictDataVo)
	InsertDictData(c *gin.Context, dictData *modelrequest.SysDictDataDML)
	UpdateDictData(c *gin.Context, dictData *modelrequest.SysDictDataDML)
	DeleteDictDataByIds(c *gin.Context, dictCodes []int64)
	CheckDictDataByTypes(c *gin.Context, dictType []string) bool
}
