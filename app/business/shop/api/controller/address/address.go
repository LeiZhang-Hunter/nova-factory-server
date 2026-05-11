package address

import (
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/region"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Address 地址控制器
type Address struct {
	service service.IApiShopAddressService
}

// NewAddress 创建地址控制器
func NewAddress(service service.IApiShopAddressService) *Address {
	return &Address{service: service}
}

// PrivateRoutes 注册商城地址路由（商城模块只检查登录，不检查权限）
func (s *Address) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/address")
	group.GET("/list", s.List)
	group.GET("/info/:id", s.GetByID)
	group.POST("/set", s.Set)
	group.DELETE("/remove/:ids", s.Remove)
	group.GET("/region", s.Region)
	group.GET("/default", s.Default)
}

// List 获取地址列表
// @Summary 获取地址列表
// @Description 获取当前用户的地址列表
// @Tags 商城/用户地址
// @Param userId query int true "用户ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/address/list [get]
func (s *Address) List(c *gin.Context) {
	userId := baizeContext.GetUserId(c)
	data, err := s.service.List(c, userId)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取地址详情
// @Summary 获取地址详情
// @Description 根据ID获取地址详情
// @Tags 商城/用户地址
// @Param id path int true "地址ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/address/info/{id} [get]
func (s *Address) GetByID(c *gin.Context) {
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

// Set 新增或修改地址
// @Summary 新增或修改地址
// @Description 新增或修改地址
// @Tags 商城/用户地址
// @Param object body models.AddressSetReq true "地址设置参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /api/v1/app/shop/address/set [post]
func (s *Address) Set(c *gin.Context) {
	userId := baizeContext.GetUserId(c)
	req := new(models.AddressSetReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("set address error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	req.UserID = userId
	data, err := s.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除地址
// @Summary 删除地址
// @Description 根据ID删除地址
// @Tags 商城/用户地址
// @Param ids path string true "地址ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /api/v1/app/shop/address/remove/{ids} [delete]
func (s *Address) Remove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.Remove(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Region 获取省市区列表
// @Summary 获取省市区列表
// @Description 不传 parentCode 返回省级列表，传省编码返回市级，传市编码返回区级
// @Tags 商城/用户地址
// @Param parentCode query string false "父级行政区编码"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/address/region [get]
func (s *Address) Region(c *gin.Context) {
	parentCode := strings.TrimSpace(c.Query("parentCode"))

	pid := 0
	if parentCode != "" {
		id, err := strconv.Atoi(parentCode)
		if err != nil {
			baizeContext.ParameterError(c)
			return
		}
		info := region.GetRegionInfo(id)
		if info == nil {
			baizeContext.ParameterError(c)
			return
		}
		pid = id
	}

	rows := region.GetRegionList(pid)
	items := make([]*models.AddressRegionItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, &models.AddressRegionItem{
			Code:       strconv.FormatInt(item.ID, 10),
			Name:       item.Name,
			Level:      regionLevelName(item.Level),
			ParentCode: parentCode,
		})
	}
	baizeContext.SuccessData(c, items)
}

func regionLevelName(level int) string {
	switch level {
	case 0:
		return "province"
	case 1:
		return "city"
	case 2:
		return "district"
	case 3:
		return "street"
	default:
		return ""
	}
}

// Default 获取默认收货地址
// @Summary 获取默认收货地址
// @Description 获取当前用户的默认收货地址
// @Tags 商城/用户地址
// @Param userId query int true "用户ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/address/default [get]
func (s *Address) Default(c *gin.Context) {
	userId := baizeContext.GetUserId(c)
	data, err := s.service.Default(c, userId)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
