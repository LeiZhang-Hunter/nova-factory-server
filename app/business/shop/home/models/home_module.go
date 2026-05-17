package models

import "nova-factory-server/app/baize"

// HomeModule 首页模块。
type HomeModule struct {
	ID         int64  `json:"id,string" gorm:"id"`           // 主键ID
	PageKey    string `json:"pageKey" gorm:"page_key"`       // 页面标识
	ModuleType string `json:"moduleType" gorm:"module_type"` // 模块类型
	ModuleName string `json:"moduleName" gorm:"module_name"` // 模块名称
	Title      string `json:"title" gorm:"title"`            // 展示标题
	SubTitle   string `json:"subTitle" gorm:"sub_title"`     // 展示副标题
	SourceType int8   `json:"sourceType" gorm:"source_type"` // 数据来源
	LimitNum   int64  `json:"limitNum" gorm:"limit_num"`     // 展示数量
	Sort       int64  `json:"sort" gorm:"sort"`              // 排序值
	StartTime  int64  `json:"startTime" gorm:"start_time"`   // 生效开始时间
	EndTime    int64  `json:"endTime" gorm:"end_time"`       // 生效结束时间
	ShowMore   int8   `json:"showMore" gorm:"show_more"`     // 是否显示更多入口
	MoreLink   string `json:"moreLink" gorm:"more_link"`     // 更多跳转链接
	StyleJSON  string `json:"styleJson" gorm:"style_json"`   // 样式配置JSON
	RuleJSON   string `json:"ruleJson" gorm:"rule_json"`     // 规则配置JSON
	ExtJSON    string `json:"extJson" gorm:"ext_json"`       // 扩展配置JSON
	Status     int8   `json:"status" gorm:"status"`          // 状态
	DeptID     int64  `json:"deptId" gorm:"dept_id"`         // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"` // 操作状态
}

// HomeModuleItem 首页模块明细。
type HomeModuleItem struct {
	ID           int64  `json:"id,string" gorm:"id"`                // 主键ID
	ModuleID     int64  `json:"moduleId,string" gorm:"module_id"`   // 模块ID
	BusinessType string `json:"businessType" gorm:"business_type"`  // 业务类型
	LinkID       int64  `json:"linkId,string" gorm:"link_id"`       // 关联业务ID
	ItemName     string `json:"itemName" gorm:"item_name"`          // 内容项名称
	ItemSubTitle string `json:"itemSubTitle" gorm:"item_sub_title"` // 内容项副标题
	ItemImage    string `json:"itemImage" gorm:"item_image"`        // 内容项图片
	Sort         int64  `json:"sort" gorm:"sort"`                   // 排序值
	Status       int8   `json:"status" gorm:"status"`               // 状态
	ExtJSON      string `json:"extJson" gorm:"ext_json"`            // 扩展配置JSON
	DeptID       int64  `json:"deptId" gorm:"dept_id"`              // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"` // 操作状态
}

// HomeModuleSet 首页模块保存参数。
type HomeModuleSet struct {
	ID         int64  `json:"id,string"` // 主键ID，更新时传
	PageKey    string `json:"pageKey"`   // 页面标识
	ModuleType string `json:"moduleType"`
	ModuleName string `json:"moduleName"`
	Title      string `json:"title"`
	SubTitle   string `json:"subTitle"`
	SourceType int8   `json:"sourceType"`
	LimitNum   int64  `json:"limitNum"`
	Sort       int64  `json:"sort"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
	ShowMore   int8   `json:"showMore"`
	MoreLink   string `json:"moreLink"`
	StyleJSON  string `json:"styleJson"`
	RuleJSON   string `json:"ruleJson"`
	ExtJSON    string `json:"extJson"`
	Status     int8   `json:"status"`
}

// HomeModuleItemSet 首页模块明细保存参数。
type HomeModuleItemSet struct {
	ID           int64  `json:"id,string"`       // 主键ID，更新时传
	ModuleID     int64  `json:"moduleId,string"` // 模块ID
	BusinessType string `json:"businessType"`
	LinkID       int64  `json:"linkId,string"`
	ItemName     string `json:"itemName"`
	ItemSubTitle string `json:"itemSubTitle"`
	ItemImage    string `json:"itemImage"`
	Sort         int64  `json:"sort"`
	Status       int8   `json:"status"`
	ExtJSON      string `json:"extJson"`
}

// HomeModuleQuery 首页模块查询参数。
type HomeModuleQuery struct {
	PageKey    string `form:"pageKey"`    // 页面标识
	ModuleType string `form:"moduleType"` // 模块类型
	Title      string `form:"title"`      // 标题
	Status     *int8  `form:"status"`     // 状态
	Page       int64  `form:"page"`       // 页码
	Size       int64  `form:"size"`       // 每页数量
}

// HomeModuleListData 首页模块列表结果。
type HomeModuleListData struct {
	Rows  []*HomeModule `json:"rows"`  // 列表数据
	Total int64         `json:"total"` // 总数
}

// HomeModuleItemQuery 首页模块明细查询参数。
type HomeModuleItemQuery struct {
	ModuleID     int64  `form:"moduleId,string"` // 模块ID
	BusinessType string `form:"businessType"`    // 业务类型
	ItemName     string `form:"itemName"`        // 内容项名称
	Status       *int8  `form:"status"`          // 状态
	Page         int64  `form:"page"`            // 页码
	Size         int64  `form:"size"`            // 每页数量
}

// HomeModuleItemListData 首页模块明细列表结果。
type HomeModuleItemListData struct {
	Rows  []*HomeModuleItem `json:"rows"`  // 列表数据
	Total int64             `json:"total"` // 总数
}

// HomeModuleItemBusinessSync 首页模块明细业务同步参数。
type HomeModuleItemBusinessSync struct {
	BusinessType string   `json:"businessType"`  // 业务类型
	LinkID       int64    `json:"linkId,string"` // 业务主键ID
	ModuleIDs    []string `json:"moduleIds"`     // 首页模块ID集合
	ItemName     string   `json:"itemName"`      // 内容项名称
	ItemSubTitle string   `json:"itemSubTitle"`  // 内容项副标题
	ItemImage    string   `json:"itemImage"`     // 内容项图片
	Sort         int64    `json:"sort"`          // 排序值
	Status       int8     `json:"status"`        // 状态
	ExtJSON      string   `json:"extJson"`       // 扩展配置JSON
}
