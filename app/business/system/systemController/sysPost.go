package systemController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Post struct {
	ps systemService.IPostService
}

func NewPost(ps systemService.IPostService) *Post {
	return &Post{ps: ps}
}
func (pc *Post) PrivateRoutes(router *gin.RouterGroup) {
	systemPost := router.Group("/system/post")
	systemPost.GET("/list", middlewares.HasPermission("system:post"), pc.PostList)
	systemPost.POST("/export", middlewares.HasPermission("system:post:export"), pc.PostExport)
	systemPost.GET("/:postId", middlewares.HasPermission("system:post:query"), pc.PostGetInfo)
	systemPost.POST("", middlewares.SetLog("岗位管理", middlewares.Insert), middlewares.HasPermission("system:post:add"), pc.PostAdd)
	systemPost.PUT("", middlewares.SetLog("岗位管理", middlewares.Update), middlewares.HasPermission("system:post:edit"), pc.PostEdit)
	systemPost.DELETE("/:postIds", middlewares.SetLog("岗位管理", middlewares.Delete), middlewares.HasPermission("system:post:remove"), pc.PostRemove)

}

// PostList 查询岗位列表查询
// @Summary 查询岗位列表查询
// @Description 查询岗位列表查询
// @Tags 岗位相关
// @Param  object query systemModels.SysPostDQL true "查询信息"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData{data=response.ListData{Rows=[]systemModels.SysPostVo}}  "成功"
// @Router /system/post/list  [get]
func (pc *Post) PostList(c *gin.Context) {
	post := new(systemModels.SysPostDQL)
	_ = c.ShouldBind(post)
	list, count := pc.ps.SelectPostList(c, post)
	baizeContext.SuccessListData(c, list, count)

}

// PostExport 导出岗位
// @Summary 导出岗位
// @Description 导出岗位
// @Tags 岗位相关
// @Param  object query systemModels.SysPostDQL true "查询信息"
// @Security BearerAuth
// @Produce application/octet-stream
// @Success 200 {object} []byte
// @Router /system/post/export [post]
func (pc *Post) PostExport(c *gin.Context) {
	post := new(systemModels.SysPostDQL)
	_ = c.ShouldBind(post)
	data := pc.ps.PostExport(c, post)
	baizeContext.DataPackageExcel(c, data)
}

// PostGetInfo 根据岗位ID获取岗位信息
// @Summary 根据岗位ID获取岗位信息
// @Description 根据岗位ID获取岗位信息
// @Tags 岗位相关
// @Param id path string true "PostId"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData{data=systemModels.SysPostVo}  "成功"
// @Router /system/post/{postId}  [get]
func (pc *Post) PostGetInfo(c *gin.Context) {
	postId := baizeContext.ParamInt64(c, "postId")
	if postId == 0 {
		baizeContext.ParameterError(c)
		return
	}
	baizeContext.SuccessData(c, pc.ps.SelectPostById(c, postId))
}

// PostAdd 添加岗位
// @Summary 添加岗位
// @Description 添加岗位
// @Tags 岗位相关
// @Param  object body systemModels.SysPostVo true "公司信息"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /system/post  [post]
func (pc *Post) PostAdd(c *gin.Context) {
	sysPost := new(systemModels.SysPostVo)
	if err := c.ShouldBindJSON(sysPost); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	sysPost.SetCreateBy(baizeContext.GetUserId(c))
	pc.ps.InsertPost(c, sysPost)
	baizeContext.Success(c)
}

// PostEdit 修改岗位
// @Summary 修改岗位
// @Description 修改岗位
// @Tags 岗位相关
// @Param  object body systemModels.SysPostVo true "公司信息"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /system/post  [put]
func (pc *Post) PostEdit(c *gin.Context) {
	post := new(systemModels.SysPostVo)
	if err := c.ShouldBindJSON(post); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	post.SetUpdateBy(baizeContext.GetUserId(c))
	pc.ps.UpdatePost(c, post)
	baizeContext.Success(c)
}

// PostRemove 删除岗位
// @Summary 删除岗位
// @Description 删除岗位
// @Tags 岗位相关
// @Param ids path []string true "postId"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /system/post/{postIds}  [delete]
func (pc *Post) PostRemove(c *gin.Context) {
	pc.ps.DeletePostByIds(c, baizeContext.ParamInt64Array(c, "postIds"))
	baizeContext.Success(c)
}
