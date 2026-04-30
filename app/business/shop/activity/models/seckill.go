package models

import "nova-factory-server/app/baize"

// Seckill 秒杀商品
type Seckill struct {
	ID           int64   `json:"id,string" db:"id"`                  // 秒杀商品ID
	ActivityID   int64   `json:"activityId,string" db:"activity_id"` // 活动ID
	ProductID    int64   `json:"productId,string" db:"product_id"`   // 商品ID
	Image        string  `json:"image" db:"image"`                   // 推荐图
	Images       string  `json:"images" db:"images"`                 // 轮播图，逗号分隔
	Title        string  `json:"title" db:"title"`                   // 活动标题
	Info         string  `json:"info" db:"info"`                     // 简介
	Price        float64 `json:"price" db:"price"`                   // 秒杀价
	Cost         float64 `json:"cost" db:"cost"`                     // 成本价
	OtPrice      float64 `json:"otPrice" db:"ot_price"`              // 原价
	GiveIntegral float64 `json:"giveIntegral" db:"give_integral"`    // 赠送积分
	Sort         int64   `json:"sort" db:"sort"`                     // 排序值
	Stock        int64   `json:"stock" db:"stock"`                   // 库存
	Sales        int64   `json:"sales" db:"sales"`                   // 销量
	UnitName     string  `json:"unitName" db:"unit_name"`            // 单位名称
	Postage      float64 `json:"postage" db:"postage"`               // 邮费
	StartTime    string  `json:"startTime" db:"start_time"`          // 开始时间
	StopTime     string  `json:"stopTime" db:"stop_time"`            // 结束时间
	AddTime      string  `json:"addTime" db:"add_time"`              // 添加时间
	Status       int32   `json:"status" db:"status"`                 // 商品状态
	IsPostage    int32   `json:"isPostage" db:"is_postage"`          // 是否包邮
	IsHot        int32   `json:"isHot" db:"is_hot"`                  // 是否热门推荐
	IsDel        int32   `json:"isDel" db:"is_del"`                  // 是否删除
	Num          int32   `json:"num" db:"num"`                       // 最多秒杀数量
	IsShow       int32   `json:"isShow" db:"is_show"`                // 是否显示
	TimeID       string  `json:"timeId" db:"time_id"`                // 时间段ID
	TempID       int32   `json:"tempId" db:"temp_id"`                // 运费模板ID
	Weight       float64 `json:"weight" db:"weight"`                 // 商品重量
	Volume       float64 `json:"volume" db:"volume"`                 // 商品体积
	Quota        int64   `json:"quota" db:"quota"`                   // 限购总数
	QuotaShow    int64   `json:"quotaShow" db:"quota_show"`          // 限购总数显示
	OnceNum      int32   `json:"onceNum" db:"once_num"`              // 单次购买数量
	Logistics    string  `json:"logistics" db:"logistics"`           // 物流类型
	Freight      int32   `json:"freight" db:"freight"`               // 运费设置
	CustomForm   string  `json:"customForm" db:"custom_form"`        // 自定义表单
	VirtualType  int32   `json:"virtualType" db:"virtual_type"`      // 商品类型
	IsCommission int32   `json:"isCommission" db:"is_commission"`    // 是否参与返佣
	DeptID       int64   `json:"deptId" db:"dept_id"`                // 部门ID
	baize.BaseEntity
	State int32 `json:"state" db:"state"` // 操作状态
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
	Page       int64  `form:"page"`       // 页码
	Size       int64  `form:"size"`       // 每页数量
}

// SeckillListData 秒杀商品列表数据
type SeckillListData struct {
	Rows  []*Seckill `json:"rows"`  // 列表数据
	Total int64      `json:"total"` // 总数
}
