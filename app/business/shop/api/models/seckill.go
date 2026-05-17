package models

// Seckill 秒杀商品
type Seckill struct {
	ID           int64   `json:"id,string" gorm:"column:id"`
	ActivityID   int64   `json:"activityId,string" gorm:"column:activity_id"`
	ProductID    int64   `json:"productId,string" gorm:"column:product_id"`
	Image        string  `json:"image" gorm:"column:image"`
	Images       string  `json:"images" gorm:"column:images"`
	Title        string  `json:"title" gorm:"column:title"`
	Info         string  `json:"info" gorm:"column:info"`
	Price        float64 `json:"price" gorm:"column:price"`
	Cost         float64 `json:"cost" gorm:"column:cost"`
	OtPrice      float64 `json:"otPrice" gorm:"column:ot_price"`
	GiveIntegral float64 `json:"giveIntegral" gorm:"column:give_integral"`
	Sort         int64   `json:"sort" gorm:"column:sort"`
	Stock        int64   `json:"stock" gorm:"column:stock"`
	Sales        int64   `json:"sales" gorm:"column:sales"`
	UnitName     string  `json:"unitName" gorm:"column:unit_name"`
	Postage      float64 `json:"postage" gorm:"column:postage"`
	StartTime    string  `json:"startTime" gorm:"column:start_time"`
	StopTime     string  `json:"stopTime" gorm:"column:stop_time"`
	Status       int32   `json:"status" gorm:"column:status"`
	IsPostage    int32   `json:"isPostage" gorm:"column:is_postage"`
	IsHot        int32   `json:"isHot" gorm:"column:is_hot"`
	Num          int32   `json:"num" gorm:"column:num"`
	IsShow       int32   `json:"isShow" gorm:"column:is_show"`
	TimeID       string  `json:"timeId" gorm:"column:time_id"`
	Quota        int64   `json:"quota" gorm:"column:quota"`
	QuotaShow    int64   `json:"quotaShow" gorm:"column:quota_show"`
}

// SeckillQuery 秒杀商品查询参数
type SeckillQuery struct {
	Title      string `form:"title"`
	ActivityID int64  `form:"activityId"`
	ProductID  int64  `form:"productId"`
	Status     *int32 `form:"status"`
	IsShow     *int32 `form:"isShow"`
	IsHot      *int32 `form:"isHot"`
	TimeID     int64  `form:"timeId"`
	Page       int64  `form:"page"`
	Size       int64  `form:"size"`
}

// SeckillListData 秒杀商品列表数据
type SeckillListData struct {
	Rows  []*Seckill `json:"rows"`
	Total int64      `json:"total"`
}
