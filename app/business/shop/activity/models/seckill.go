package models

import "nova-factory-server/app/baize"

type SeckillMainInfo struct {
	ID           int64   `json:"id,string" gorm:"column:id"`                  // 秒杀商品ID
	ActivityID   int64   `json:"activityId,string" gorm:"column:activity_id"` // 活动ID
	ProductID    int64   `json:"productId,string" gorm:"column:product_id"`   // 商品ID
	Price        float64 `json:"price" gorm:"column:price"`                   // 秒杀价
	Cost         float64 `json:"cost" gorm:"column:cost"`                     // 成本价
	OtPrice      float64 `json:"otPrice" gorm:"column:ot_price"`              // 原价
	Stock        int64   `json:"stock" gorm:"column:stock"`                   // 库存
	StartTime    string  `json:"startTime" gorm:"column:start_time"`          // 开始时间
	StopTime     string  `json:"stopTime" gorm:"column:stop_time"`            // 结束时间
	AddTime      string  `json:"addTime" gorm:"column:add_time"`              // 添加时间
	Status       int32   `json:"status" gorm:"column:status"`                 // 商品状态
	IsPostage    int32   `json:"isPostage" gorm:"column:is_postage"`          // 是否包邮
	IsHot        int32   `json:"isHot" gorm:"column:is_hot"`                  // 是否热门推荐
	IsDel        int32   `json:"isDel" gorm:"column:is_del"`                  // 是否删除
	Num          int32   `json:"num" gorm:"column:num"`                       // 最多秒杀数量
	IsShow       int32   `json:"isShow" gorm:"column:is_show"`                // 是否显示
	TimeID       string  `json:"timeId,string" gorm:"column:time_id"`         // 时间段ID
	TempID       int32   `json:"tempId" gorm:"column:temp_id"`                // 运费模板ID
	Weight       float64 `json:"weight" gorm:"column:weight"`                 // 商品重量
	Volume       float64 `json:"volume" gorm:"column:volume"`                 // 商品体积
	Quota        int64   `json:"quota" gorm:"column:quota"`                   // 限购总数
	QuotaShow    int64   `json:"quotaShow" gorm:"column:quota_show"`          // 限购总数显示
	OnceNum      int32   `json:"onceNum" gorm:"column:once_num"`              // 单次购买数量
	Logistics    string  `json:"logistics" gorm:"column:logistics"`           // 物流类型
	Freight      int32   `json:"freight" gorm:"column:freight"`               // 运费设置
	CustomForm   string  `json:"customForm" gorm:"column:custom_form"`        // 自定义表单
	VirtualType  int32   `json:"virtualType" gorm:"column:virtual_type"`      // 商品类型
	IsCommission int32   `json:"isCommission" gorm:"column:is_commission"`    // 是否参与返佣
}

// FromatSeckillMainInfo 提取秒杀商品主信息。
func FromatSeckillMainInfo(sec *Seckill) *SeckillMainInfo {
	if sec == nil {
		return nil
	}
	return &SeckillMainInfo{
		ID:           sec.ID,
		ActivityID:   sec.ActivityID,
		ProductID:    sec.ProductID,
		Price:        sec.Price,
		Cost:         sec.Cost,
		OtPrice:      sec.OtPrice,
		Stock:        sec.Stock,
		StartTime:    sec.StartTime,
		StopTime:     sec.StopTime,
		AddTime:      sec.AddTime,
		Status:       sec.Status,
		IsPostage:    sec.IsPostage,
		IsHot:        sec.IsHot,
		IsDel:        sec.IsDel,
		Num:          sec.Num,
		IsShow:       sec.IsShow,
		TimeID:       sec.TimeID,
		TempID:       sec.TempID,
		Weight:       sec.Weight,
		Volume:       sec.Volume,
		Quota:        sec.Quota,
		QuotaShow:    sec.QuotaShow,
		OnceNum:      sec.OnceNum,
		Logistics:    sec.Logistics,
		Freight:      sec.Freight,
		CustomForm:   sec.CustomForm,
		VirtualType:  sec.VirtualType,
		IsCommission: sec.IsCommission,
	}
}

// Seckill 秒杀商品
type Seckill struct {
	ID           int64   `json:"id,string" gorm:"column:id"`                  // 秒杀商品ID
	ActivityID   int64   `json:"activityId,string" gorm:"column:activity_id"` // 活动ID
	ProductID    int64   `json:"productId,string" gorm:"column:product_id"`   // 商品ID
	Image        string  `json:"image" gorm:"column:image"`                   // 推荐图
	Images       string  `json:"images" gorm:"column:images"`                 // 轮播图，逗号分隔
	Title        string  `json:"title" gorm:"column:title"`                   // 活动标题
	Info         string  `json:"info" gorm:"column:info"`                     // 简介
	Price        float64 `json:"price" gorm:"column:price"`                   // 秒杀价
	Cost         float64 `json:"cost" gorm:"column:cost"`                     // 成本价
	OtPrice      float64 `json:"otPrice" gorm:"column:ot_price"`              // 原价
	GiveIntegral float64 `json:"giveIntegral" gorm:"column:give_integral"`    // 赠送积分
	Sort         int64   `json:"sort" gorm:"column:sort"`                     // 排序值
	Stock        int64   `json:"stock" gorm:"column:stock"`                   // 库存
	Sales        int64   `json:"sales" gorm:"column:sales"`                   // 销量
	UnitName     string  `json:"unitName" gorm:"column:unit_name"`            // 单位名称
	Postage      float64 `json:"postage" gorm:"column:postage"`               // 邮费
	StartTime    string  `json:"startTime" gorm:"column:start_time"`          // 开始时间
	StopTime     string  `json:"stopTime" gorm:"column:stop_time"`            // 结束时间
	AddTime      string  `json:"addTime" gorm:"column:add_time"`              // 添加时间
	Status       int32   `json:"status" gorm:"column:status"`                 // 商品状态
	IsPostage    int32   `json:"isPostage" gorm:"column:is_postage"`          // 是否包邮
	IsHot        int32   `json:"isHot" gorm:"column:is_hot"`                  // 是否热门推荐
	IsDel        int32   `json:"isDel" gorm:"column:is_del"`                  // 是否删除
	Num          int32   `json:"num" gorm:"column:num"`                       // 最多秒杀数量
	IsShow       int32   `json:"isShow" gorm:"column:is_show"`                // 是否显示
	TimeID       string  `json:"timeId" gorm:"column:time_id"`                // 时间段ID
	TempID       int32   `json:"tempId" gorm:"column:temp_id"`                // 运费模板ID
	Weight       float64 `json:"weight" gorm:"column:weight"`                 // 商品重量
	Volume       float64 `json:"volume" gorm:"column:volume"`                 // 商品体积
	Quota        int64   `json:"quota" gorm:"column:quota"`                   // 限购总数
	QuotaShow    int64   `json:"quotaShow" gorm:"column:quota_show"`          // 限购总数显示
	OnceNum      int32   `json:"onceNum" gorm:"column:once_num"`              // 单次购买数量
	Logistics    string  `json:"logistics" gorm:"column:logistics"`           // 物流类型
	Freight      int32   `json:"freight" gorm:"column:freight"`               // 运费设置
	CustomForm   string  `json:"customForm" gorm:"column:custom_form"`        // 自定义表单
	VirtualType  int32   `json:"virtualType" gorm:"column:virtual_type"`      // 商品类型
	IsCommission int32   `json:"isCommission" gorm:"column:is_commission"`    // 是否参与返佣
	DeptID       int64   `json:"deptId" gorm:"column:dept_id"`                // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"` // 操作状态
}

// SeckillSet 秒杀商品保存参数
type SeckillSet struct {
	ID           int64   `json:"id,string"`                           // 秒杀商品ID，更新时传
	ActivityID   int64   `json:"activityId,string"`                   // 活动ID
	ProductID    int64   `json:"productId,string" binding:"required"` // 商品ID
	Image        string  `json:"image" binding:"required"`            // 推荐图
	Images       string  `json:"images"`                              // 轮播图，逗号分隔
	Title        string  `json:"title" binding:"required"`            // 活动标题
	Info         string  `json:"info"`                                // 简介
	Price        float64 `json:"price"`                               // 秒杀价
	Cost         float64 `json:"cost"`                                // 成本价
	OtPrice      float64 `json:"otPrice"`                             // 原价
	GiveIntegral float64 `json:"giveIntegral"`                        // 赠送积分
	Sort         int64   `json:"sort"`                                // 排序值
	Stock        int64   `json:"stock"`                               // 库存
	Sales        int64   `json:"sales"`                               // 销量
	UnitName     string  `json:"unitName"`                            // 单位名称
	Postage      float64 `json:"postage"`                             // 邮费
	StartTime    string  `json:"startTime"`                           // 开始时间
	StopTime     string  `json:"stopTime"`                            // 结束时间
	Status       int32   `json:"status"`                              // 商品状态
	IsPostage    int32   `json:"isPostage"`                           // 是否包邮
	IsHot        int32   `json:"isHot"`                               // 是否热门推荐
	Num          int32   `json:"num"`                                 // 最多秒杀数量
	IsShow       int32   `json:"isShow"`                              // 是否显示
	TimeID       string  `json:"timeId"`                              // 时间段ID
	TempID       int32   `json:"tempId"`                              // 运费模板ID
	Weight       float64 `json:"weight"`                              // 商品重量
	Volume       float64 `json:"volume"`                              // 商品体积
	Quota        int64   `json:"quota"`                               // 限购总数
	QuotaShow    int64   `json:"quotaShow"`                           // 限购总数显示
	OnceNum      int32   `json:"onceNum"`                             // 单次购买数量
	Logistics    string  `json:"logistics"`                           // 物流类型
	Freight      int32   `json:"freight"`                             // 运费设置
	CustomForm   string  `json:"customForm"`                          // 自定义表单
	VirtualType  int32   `json:"virtualType"`                         // 商品类型
	IsCommission int32   `json:"isCommission"`                        // 是否参与返佣
}

// SeckillQuery 秒杀商品查询参数
type SeckillQuery struct {
	Title      string `form:"title"`      // 活动标题
	ActivityID int64  `form:"activityId"` // 活动ID
	ProductID  int64  `form:"productId"`  // 商品ID
	Status     *int32 `form:"status"`     // 商品状态
	IsShow     *int32 `form:"isShow"`     // 是否显示
	IsHot      *int32 `form:"isHot"`      // 是否热门推荐
	TimeID     int64  `form:"timeId"`     // 时间段ID
	Page       int64  `form:"page"`       // 页码
	Size       int64  `form:"size"`       // 每页数量
}

// SeckillListData 秒杀商品列表数据
type SeckillListData struct {
	Rows  []*Seckill `json:"rows"`  // 列表数据
	Total int64      `json:"total"` // 总数
}
