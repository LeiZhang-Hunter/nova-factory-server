package gatewaycontroller

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Skills 已安装技能控制器。
type Skills struct {
	service gatewayservice.IInstalledSkillService
}

// NewSkills 创建已安装技能控制器。
func NewSkills(service gatewayservice.IInstalledSkillService) *Skills {
	return &Skills{service: service}
}

// PrivateRoutes 注册已安装技能路由。
func (s *Skills) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/skills")
	group.GET("/list", middlewares.HasPermission("ai:skills:list"), s.List)
	group.GET("/query/:id", middlewares.HasPermission("ai:skills:query"), s.GetByID)
	group.POST("/set", middlewares.HasPermission("ai:skills:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:skills:remove"), s.Delete)
}

// List 获取已安装技能列表
// @Summary 获取已安装技能列表
// @Description 获取已安装技能列表
// @Tags 工业智能体/已安装技能
// @Param object query gatewaymodels.InstalledSkillQuery true "已安装技能查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/skills/list [get]
func (s *Skills) List(c *gin.Context) {
	req := new(gatewaymodels.InstalledSkillQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取已安装技能详情
// @Summary 获取已安装技能详情
// @Description 根据ID获取已安装技能详情
// @Tags 工业智能体/已安装技能
// @Param id path int true "技能ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/skills/query/{id} [get]
func (s *Skills) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 保存已安装技能
// @Summary 保存已安装技能
// @Description 保存已安装技能，id为空时新增，不为空时修改
// @Tags 工业智能体/已安装技能
// @Param object body gatewaymodels.InstalledSkillUpsert true "已安装技能保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /ai/skills/set [post]
func (s *Skills) Set(c *gin.Context) {
	req := new(gatewaymodels.InstalledSkillUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *gatewaymodels.InstalledSkill
		err  error
	)
	if req.ID > 0 {
		data, err = s.service.Update(c, req)
	} else {
		data, err = s.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除已安装技能
// @Summary 删除已安装技能
// @Description 根据ID删除已安装技能
// @Tags 工业智能体/已安装技能
// @Param ids path string true "技能ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/skills/remove/{ids} [delete]
func (s *Skills) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
