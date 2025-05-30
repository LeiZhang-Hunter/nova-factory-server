package toolController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/tool/toolModels"
	"nova-factory-server/app/business/tool/toolService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strings"
)

type GenTable struct {
	gt toolService.IGenTableService
}

func NewGenTable(gt toolService.IGenTableService) *GenTable {
	return &GenTable{gt: gt}
}

func (gc *GenTable) PrivateRoutes(router *gin.RouterGroup) {
	genTable := router.Group("/tool/gen")
	genTable.GET("/list", middlewares.HasPermission("tool:gen:list"), gc.GenTableList)
	genTable.GET(":tableId", middlewares.HasPermission("tool:gen:query"), gc.GenTableGetInfo)
	genTable.GET("/db/list", middlewares.HasPermission("tool:gen:list"), gc.DataList)
	genTable.GET("/column/:talbleId", middlewares.HasPermission("tool:gen:list"), gc.ColumnList)
	genTable.POST("/importTable", middlewares.HasPermission("tool:gen:list"), gc.ImportTable)
	genTable.PUT("", middlewares.HasPermission("tool:gen:edit"), gc.EditSave)
	genTable.DELETE("/:tableIds", middlewares.HasPermission("tool:gen:remove"), gc.GenTableRemove)
	genTable.GET("/preview/:tableId", middlewares.HasPermission("tool:gen:code"), gc.Preview)
	genTable.GET("/genCode/:tableId", middlewares.HasPermission("tool:gen:code"), gc.GenCode)
}

func (gc *GenTable) GenTableList(c *gin.Context) {
	getTable := new(toolModels.GenTableDQL)
	c.ShouldBind(getTable)
	list, count := gc.gt.SelectGenTableList(c, getTable)
	baizeContext.SuccessListData(c, list, count)
}

func (gc *GenTable) GenTableGetInfo(c *gin.Context) {
	tableId := baizeContext.ParamInt64(c, "tableId")
	genTable := gc.gt.SelectGenTableById(c, tableId)
	tables := gc.gt.SelectGenTableAll(c)
	list := gc.gt.SelectGenTableColumnListByTableId(c, tableId)
	data := make(map[string]interface{})
	data["info"] = genTable
	data["rows"] = list
	data["tables"] = tables
	baizeContext.SuccessData(c, data)
}

func (gc *GenTable) DataList(c *gin.Context) {
	getTable := new(toolModels.GenTableDQL)
	_ = c.ShouldBind(getTable)
	list, count := gc.gt.SelectDbTableList(c, getTable)
	baizeContext.SuccessListData(c, list, count)

}
func (gc *GenTable) ColumnList(c *gin.Context) {
	tableId := baizeContext.ParamInt64(c, "tableId")
	list := gc.gt.SelectGenTableColumnListByTableId(c, tableId)
	total := int64(len(list))
	baizeContext.SuccessListData(c, list, total)
}
func (gc *GenTable) ImportTable(c *gin.Context) {
	gc.gt.ImportTableSave(c, strings.Split(c.Query("tables"), ","), "")
	baizeContext.SuccessMsg(c, "导入成功")
}
func (gc *GenTable) EditSave(c *gin.Context) {
	genTable := new(toolModels.GenTableDML)
	_ = c.ShouldBindJSON(genTable)
	genTable.SetUpdateBy(baizeContext.GetUserId(c))
	gc.gt.UpdateGenTable(c, genTable)
	baizeContext.Success(c)
}
func (gc *GenTable) GenTableRemove(c *gin.Context) {
	gc.gt.DeleteGenTableByIds(c, baizeContext.ParamInt64Array(c, "tableIds"))
	baizeContext.Success(c)
}
func (gc *GenTable) Preview(c *gin.Context) {
	s := gc.gt.PreviewCode(c, baizeContext.ParamInt64(c, "tableId"))
	baizeContext.SuccessData(c, s)
}

func (gc *GenTable) GenCode(c *gin.Context) {
	s := gc.gt.GenCode(c, baizeContext.ParamInt64(c, "tableId"))
	baizeContext.DataPackageZip(c, s)

}
