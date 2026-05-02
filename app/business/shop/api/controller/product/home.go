package product

import (
	"time"

	homeModels "nova-factory-server/app/business/shop/home/models"
	homeService "nova-factory-server/app/business/shop/home/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Home 前台首页模块控制器。
type Home struct {
	moduleService homeService.IShopHomeModuleService
	itemService   homeService.IShopHomeModuleItemService
}

// HomeListQuery 前台首页模块查询参数。
type HomeListQuery struct {
	PageKey string `form:"pageKey"` // 页面标识
}

// HomeModuleData 前台首页模块数据。
type HomeModuleData struct {
	ID         int64                 `json:"id,string"`  // 主键ID
	PageKey    string                `json:"pageKey"`    // 页面标识
	ModuleType string                `json:"moduleType"` // 模块类型
	ModuleName string                `json:"moduleName"` // 模块名称
	Title      string                `json:"title"`      // 展示标题
	SubTitle   string                `json:"subTitle"`   // 展示副标题
	SourceType int8                  `json:"sourceType"` // 数据来源
	LimitNum   int64                 `json:"limitNum"`   // 展示数量
	Sort       int64                 `json:"sort"`       // 排序值
	StartTime  int64                 `json:"startTime"`  // 生效开始时间
	EndTime    int64                 `json:"endTime"`    // 生效结束时间
	ShowMore   int8                  `json:"showMore"`   // 是否显示更多入口
	MoreLink   string                `json:"moreLink"`   // 更多跳转链接
	StyleJSON  string                `json:"styleJson"`  // 样式配置JSON
	RuleJSON   string                `json:"ruleJson"`   // 规则配置JSON
	ExtJSON    string                `json:"extJson"`    // 扩展配置JSON
	Items      []*HomeModuleItemData `json:"items"`      // 模块明细
}

// HomeModuleItemData 前台首页模块明细数据。
type HomeModuleItemData struct {
	ID           int64  `json:"id,string"`       // 主键ID
	ModuleID     int64  `json:"moduleId,string"` // 模块ID
	BusinessType string `json:"businessType"`    // 业务类型
	LinkID       int64  `json:"linkId,string"`   // 关联业务ID
	ItemName     string `json:"itemName"`        // 内容项名称
	ItemSubTitle string `json:"itemSubTitle"`    // 内容项副标题
	ItemImage    string `json:"itemImage"`       // 内容项图片
	Sort         int64  `json:"sort"`            // 排序值
	ExtJSON      string `json:"extJson"`         // 扩展配置JSON
}

// HomeModuleListData 前台首页模块列表结果。
type HomeModuleListData struct {
	Rows  []*HomeModuleData `json:"rows"`  // 模块列表
	Total int64             `json:"total"` // 模块数量
}

// NewHome 创建前台首页模块控制器。
func NewHome(moduleService homeService.IShopHomeModuleService, itemService homeService.IShopHomeModuleItemService) *Home {
	return &Home{
		moduleService: moduleService,
		itemService:   itemService,
	}
}

// PublicRoutes 注册前台首页模块路由。
func (h *Home) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/home")
	group.GET("/list", h.List)
}

func (h *Home) PrivateRoutes(router *gin.RouterGroup) {

}

// List 获取前台首页模块列表。
// @Summary 获取前台首页模块列表
// @Description 读取前台首页可展示的模块及其明细项
// @Tags app接口/商城/App首页
// @Produce application/json
// @Param pageKey query string false "页面标识，默认 index"
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/home/list [get]
func (h *Home) List(ctx *gin.Context) {
	req := new(HomeListQuery)
	if err := ctx.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}

	status := int8(1)
	moduleData, err := h.moduleService.List(ctx, &homeModels.HomeModuleQuery{
		Status: &status,
		Page:   1,
		Size:   200,
	})
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}

	now := time.Now().Unix()
	modules := make([]*homeModels.HomeModule, 0, len(moduleData.Rows))
	moduleIDs := make([]int64, 0, len(moduleData.Rows))
	for _, row := range moduleData.Rows {
		if row == nil || !isHomeModuleAvailable(row, now) {
			continue
		}
		modules = append(modules, row)
		moduleIDs = append(moduleIDs, row.ID)
	}

	items, err := h.itemService.ListByModuleIDs(ctx, moduleIDs)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	itemMap := make(map[int64][]*HomeModuleItemData, len(moduleIDs))
	for _, item := range items {
		if item == nil {
			continue
		}
		itemMap[item.ModuleID] = append(itemMap[item.ModuleID], buildHomeModuleItemData(item))
	}

	rows := make([]*HomeModuleData, 0, len(modules))
	for _, module := range modules {
		rows = append(rows, &HomeModuleData{
			ID:         module.ID,
			PageKey:    module.PageKey,
			ModuleType: module.ModuleType,
			ModuleName: module.ModuleName,
			Title:      module.Title,
			SubTitle:   module.SubTitle,
			SourceType: module.SourceType,
			LimitNum:   module.LimitNum,
			Sort:       module.Sort,
			StartTime:  module.StartTime,
			EndTime:    module.EndTime,
			ShowMore:   module.ShowMore,
			MoreLink:   module.MoreLink,
			StyleJSON:  module.StyleJSON,
			RuleJSON:   module.RuleJSON,
			ExtJSON:    module.ExtJSON,
			Items:      itemMap[module.ID],
		})
	}
	baizeContext.SuccessData(ctx, &HomeModuleListData{
		Rows:  rows,
		Total: int64(len(rows)),
	})
}

func isHomeModuleAvailable(module *homeModels.HomeModule, now int64) bool {
	if module == nil || module.Status != 1 {
		return false
	}
	if module.StartTime > 0 && now < module.StartTime {
		return false
	}
	if module.EndTime > 0 && now > module.EndTime {
		return false
	}
	return true
}

func buildHomeModuleItemData(item *homeModels.HomeModuleItem) *HomeModuleItemData {
	if item == nil {
		return nil
	}
	return &HomeModuleItemData{
		ID:           item.ID,
		ModuleID:     item.ModuleID,
		BusinessType: item.BusinessType,
		LinkID:       item.LinkID,
		ItemName:     item.ItemName,
		ItemSubTitle: item.ItemSubTitle,
		ItemImage:    item.ItemImage,
		Sort:         item.Sort,
		ExtJSON:      item.ExtJSON,
	}
}
