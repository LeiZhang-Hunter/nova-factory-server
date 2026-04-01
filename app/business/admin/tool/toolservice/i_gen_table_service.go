package toolservice

import (
	toolModels2 "nova-factory-server/app/business/admin/tool/toolmodels"

	"github.com/gin-gonic/gin"
)

type IGenTableService interface {
	SelectGenTableList(c *gin.Context, getTable *toolModels2.GenTableDQL) (list []*toolModels2.GenTableVo, total int64)
	SelectDbTableList(c *gin.Context, getTable *toolModels2.GenTableDQL) (list []*toolModels2.DBTableVo, total int64)
	SelectGenTableAll(c *gin.Context) (list []*toolModels2.GenTableVo)
	SelectGenTableById(c *gin.Context, id int64) (genTable *toolModels2.GenTableVo)
	ImportTableSave(c *gin.Context, table []string, userName string)
	UpdateGenTable(c *gin.Context, genTable *toolModels2.GenTableDML)
	DeleteGenTableByIds(c *gin.Context, ids []int64)
	PreviewCode(c *gin.Context, tableId int64) (m map[string]string)
	GenCode(c *gin.Context, tableId int64) []byte
	SelectGenTableColumnListByTableId(c *gin.Context, tableId int64) (list []*toolModels2.GenTableColumnVo)
}
