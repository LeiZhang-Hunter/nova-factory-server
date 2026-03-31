package systemservice

import (
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type IDictTypeService interface {
	SelectDictTypeList(c *gin.Context, dictType *systemmodels.SysDictTypeDQL) (list []*systemmodels.SysDictTypeVo, total int64)
	ExportDictType(c *gin.Context, dictType *systemmodels.SysDictTypeDQL) (data []byte)
	SelectDictTypeById(c *gin.Context, dictId int64) (dictType *systemmodels.SysDictTypeVo)
	SelectDictTypeByIds(c *gin.Context, dictId []int64) (dictTypes []string)
	InsertDictType(c *gin.Context, dictType *systemmodels.SysDictTypeVo)
	UpdateDictType(c *gin.Context, dictType *systemmodels.SysDictTypeVo)
	DeleteDictTypeByIds(c *gin.Context, dictIds []int64)
	CheckDictTypeUnique(c *gin.Context, id int64, dictType string) bool
	DictTypeClearCache(c *gin.Context)
	SelectDictTypeAll(c *gin.Context) (list []*systemmodels.SysDictTypeVo)
}
