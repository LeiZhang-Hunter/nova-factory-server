package systemController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

// Shift 班次配置
type Shift struct {
	service systemService.ISysShiftService
}

func NewShift(service systemService.ISysShiftService) *Shift {
	return &Shift{
		service: service,
	}
}

func (s *Shift) PrivateRoutes(router *gin.RouterGroup) {
	ele := router.Group("/system/shift")
	ele.GET("/list", middlewares.HasPermission("system:shift:list"), s.List)
	ele.POST("/set", middlewares.HasPermissions([]string{"system:shift:set"}), s.Set)
	ele.DELETE("/remove/:ids", middlewares.HasPermission("system:shift:remove"), s.Remove)
	return
}

// Set 设置班次配置
// @Summary 设置班次配置
// @Description 设置班次配置
// @Tags 系统管理/班次配置
// @Param  object body systemModels.SysWorkShiftSettingVO true "班次配置参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /system/shift/set [post]
func (s *Shift) Set(c *gin.Context) {
	info := new(systemModels.SysWorkShiftSettingVO)
	err := c.ShouldBindJSON(info)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	setting, err := systemModels.ToSysWorkShiftSetting(info)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	data := s.service.Check(c, setting.ID, setting.BeginTime, setting.EndTime)
	if data != nil {
		if info.ID != data.ID {
			baizeContext.Waring(c, fmt.Sprintf("存在重复班次:%s", data.Name))
			return
		}

	}

	value, err := s.service.Set(c, info)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, value)
}

// List 班次配置列表
// @Summary 班次配置列表
// @Description 班次配置列表
// @Tags 系统管理/班次配置
// @Param  object query systemModels.SysWorkShiftSettingReq true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /system/shift/list [get]
func (s *Shift) List(c *gin.Context) {
	req := new(systemModels.SysWorkShiftSettingReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := s.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Remove 删除班次配置
// @Summary 删除班次配置
// @Description 删除告警AI推理发送配置
// @Tags 告警管理/班次配置
// @Param  ids path string true "ids"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object}  response.ResponseData "成功"
// @Router /system/shift/remove/{ids}  [delete]
func (s *Shift) Remove(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.Waring(c, "请选择删除选项")
		return
	}
	err := s.service.Remove(c, ids)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}

	baizeContext.Success(c)
}
