package models

// CombinationSku 拼团商品SKU
type CombinationSku struct {
	ID             uint64   `json:"id,string"`
	SkuID          int64    `json:"skuId,string"`
	SkuName        string   `json:"skuName"`
	ImageURL       string   `json:"imageUrl"`
	GalleryImages  []string `json:"galleryImages" gorm:"-"`
	VideoURL       string   `json:"videoUrl"`
	Price          float64  `json:"price"`          // 拼团价
	OriginalPrice  float64  `json:"originalPrice"`  // 单买价
	Quantity       int64    `json:"quantity"`       // SKU总库存
	AvailableStock int64    `json:"availableStock"` // 拼团可用库存
	Unit           string   `json:"unit"`
	Weight         float64  `json:"weight"`
	WeightUnit     string   `json:"weightUnit"`
}

// Combination 拼团商品
type Combination struct {
	ID            int64   `json:"id,string" gorm:"column:id"`
	ProductID     int64   `json:"productId,string" gorm:"column:product_id"`
	MerID         int64   `json:"merId,string" gorm:"column:mer_id"`
	Image         string  `json:"image" gorm:"column:image"`
	Images        string  `json:"images" gorm:"column:images"`
	Title         string  `json:"title" gorm:"column:title"`
	GoodsName     string  `json:"goodsName" gorm:"column:goods_name"`
	Attr          string  `json:"attr" gorm:"column:attr"`
	People        int32   `json:"people" gorm:"column:people"`
	Info          string  `json:"info" gorm:"column:info"`
	Price         float64 `json:"price" gorm:"column:price"`
	OtPrice       float64 `json:"otPrice" gorm:"column:ot_price"`
	Sort          int32   `json:"sort" gorm:"column:sort"`
	Sales         int64   `json:"sales" gorm:"column:sales"`
	Stock         int64   `json:"stock" gorm:"column:stock"`
	IsHost        int32   `json:"isHost" gorm:"column:is_host"`
	IsShow        int32   `json:"isShow" gorm:"column:is_show"`
	IsPostage     int32   `json:"isPostage" gorm:"column:is_postage"`
	Postage       float64 `json:"postage" gorm:"column:postage"`
	StartTime     int64   `json:"startTime" gorm:"column:start_time"`
	StopTime      int64   `json:"stopTime" gorm:"column:stop_time"`
	EffectiveTime int64   `json:"effectiveTime" gorm:"column:effective_time"`
	Browse        int64   `json:"browse" gorm:"column:browse"`
	UnitName      string  `json:"unitName" gorm:"column:unit_name"`
	Weight        float64 `json:"weight" gorm:"column:weight"`
	Volume        float64 `json:"volume" gorm:"column:volume"`
	Num           int64   `json:"num" gorm:"column:num"`
	OnceNum       int64   `json:"onceNum" gorm:"column:once_num"`
	Quota         int64   `json:"quota" gorm:"column:quota"`
	QuotaShow     int64   `json:"quotaShow" gorm:"column:quota_show"`
	Virtual       int64   `json:"virtual" gorm:"column:virtual"`
	HomeModuleIDs string  `json:"homeModuleIds" gorm:"column:home_module_ids"`
	PinkCount     int64   `json:"pinkCount" gorm:"column:pink_count"`
	// 扩展字段（非数据库映射）
	GoodsID  int64             `json:"goodsId" gorm:"-"`
	Gallery  []string          `json:"gallery" gorm:"-"`
	VideoURL string            `json:"videoUrl" gorm:"-"`
	Skus     []*CombinationSku `json:"skus" gorm:"-"`
}

// CombinationQuery 拼团商品查询参数
type CombinationQuery struct {
	Title     string `form:"title"`
	ProductID int64  `form:"productId"`
	IsShow    *int32 `form:"isShow"`
	IsHost    *int32 `form:"isHost"`
	Page      int64  `form:"page"`
	Size      int64  `form:"size"`
}

// CombinationListData 拼团商品列表数据
type CombinationListData struct {
	Rows  []*Combination `json:"rows"`
	Total int64          `json:"total"`
}
