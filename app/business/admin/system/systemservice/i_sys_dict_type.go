package systemservice

import (
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type IDictTypeService interface {
	SelectDictTypeList(c *gin.Context, dictType *modelquery.SysDictTypeDQL) (list []*modelresponse.SysDictTypeVo, total int64)
	ExportDictType(c *gin.Context, dictType *modelquery.SysDictTypeDQL) (data []byte)
	SelectDictTypeById(c *gin.Context, dictId int64) (dictType *modelresponse.SysDictTypeVo)
	SelectDictTypeByIds(c *gin.Context, dictId []int64) (dictTypes []string)
	InsertDictType(c *gin.Context, dictType *modelrequest.SysDictTypeDML)
	UpdateDictType(c *gin.Context, dictType *modelrequest.SysDictTypeDML)
	DeleteDictTypeByIds(c *gin.Context, dictIds []int64)
	CheckDictTypeUnique(c *gin.Context, id int64, dictType string) bool
	DictTypeClearCache(c *gin.Context)
	SelectDictTypeAll(c *gin.Context) (list []*modelresponse.SysDictTypeVo)
}
