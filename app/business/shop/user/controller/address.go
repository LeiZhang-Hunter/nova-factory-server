package shopcontroller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/region"
	"strconv"
	"strings"
	"sync"
)

// Address 商城用户地址控制器
type Address struct {
	service service.IShopAddressService
	mtx     sync.Mutex
}

// NewAddress 创建商城用户地址控制器。
func NewAddress(service service.IShopAddressService) *Address {
	return &Address{service: service}
}

// PrivateRoutes 注册商城用户地址路由
func (s *Address) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/user/address")
	group.GET("/list", middlewares.HasPermission("shop:user:address:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:user:address:info"), s.GetByID)
	group.GET("/region", middlewares.HasPermission("shop:user:address:region"), s.Region)
	group.POST("/set", middlewares.HasPermission("shop:user:address:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:user:address:remove"), s.Remove)
}

// List 获取商城用户地址列表
// @Summary 获取商城用户地址列表
// @Description 获取商城用户地址列表
// @Tags 商城/用户地址
// @Param object query models.AddressQuery true "商城用户地址查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/address/list [get]
func (s *Address) List(c *gin.Context) {
	req := new(models.AddressQuery)
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

// GetByID 获取商城用户地址详情
// @Summary 获取商城用户地址详情
// @Description 根据ID获取商城用户地址详情
// @Tags 商城/用户地址
// @Param id path int true "商城用户地址ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/address/info/{id} [get]
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

// Set 新增或修改商城用户地址
// @Summary 新增或修改商城用户地址
// @Description 新增或修改商城用户地址
// @Tags 商城/用户地址
// @Param object body models.AddressSetReq true "商城用户地址设置参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /shop/user/address/set [post]
func (s *Address) Set(c *gin.Context) {
	req := new(models.AddressSetReq)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("set address error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除商城用户地址
// @Summary 删除商城用户地址
// @Description 根据ID删除商城用户地址
// @Tags 商城/用户地址
// @Param ids path string true "商城用户地址ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/user/address/remove/{ids} [delete]
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

// Region 获取省市区行政区列表
// @Summary 获取省市区行政区列表
// @Description 不传 parentCode 返回省级列表，传省编码返回市级，传市编码返回区级
// @Tags 商城/用户地址
// @Param parentCode query string false "父级行政区编码"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/address/region [get]
func (s *Address) Region(c *gin.Context) {
	req := new(models.AddressRegionQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		zap.L().Error("get address region error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}

	req.ParentCode = strings.TrimSpace(req.ParentCode)
	req.Type = strings.TrimSpace(strings.ToLower(req.Type))

	pid := 0
	currentLevel := -1
	if req.ParentCode != "" {
		regionID, err := strconv.Atoi(req.ParentCode)
		if err != nil {
			baizeContext.ParameterError(c)
			return
		}
		s.mtx.Lock()
		info := region.GetRegionInfo(regionID)
		defer s.mtx.Unlock()
		if info == nil {
			baizeContext.ParameterError(c)
			return
		}
		pid = regionID
		currentLevel = info.Level
	}

	rows := region.GetRegionList(pid)
	items := make([]*models.AddressRegionItem, 0, len(rows))
	for _, item := range rows {
		level := regionLevelName(item.Level)
		if req.Type != "" && req.Type != level {
			continue
		}
		items = append(items, &models.AddressRegionItem{
			Code:       strconv.FormatInt(item.ID, 10),
			Name:       item.Name,
			Level:      level,
			ParentCode: req.ParentCode,
		})
	}

	if req.Type != "" && req.ParentCode != "" {
		expectedLevel := regionLevelName(currentLevel + 1)
		if expectedLevel != "" && req.Type != expectedLevel {
			items = make([]*models.AddressRegionItem, 0)
		}
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
