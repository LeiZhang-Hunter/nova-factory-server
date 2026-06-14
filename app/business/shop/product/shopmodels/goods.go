package shopmodels

import (
	"nova-factory-server/app/baize"
	goodsstore "nova-factory-server/app/utils/store/goods"
	"time"
)

// Goods 商品信息
type Goods struct {
	ID                 int64       `json:"id,string" gorm:"id"`                           // 主键ID
	GoodsID            string      `json:"goodsId" gorm:"goods_id"`                       // 商品业务ID
	ShopCategoryId     int64       `json:"shopCategoryId,string" gorm:"shop_category_id"` // 商品分类id
	ShopCategoryName   string      `json:"shopCategoryName" gorm:"-" gorm:"-"`            // 商品分类名称
	GoodsName          string      `json:"goodsName" gorm:"goods_name"`                   // 商品名称
	GoodsCode          string      `json:"goodsCode" gorm:"goods_code"`                   // 商品编码
	OuterID            string      `json:"outerId" gorm:"outer_id"`                       // 外部系统ID
	ImageURL           string      `json:"imageUrl" gorm:"image_url"`                     // 主图地址
	RetailPrice        float64     `json:"retailPrice" gorm:"retail_price"`               // 零售价
	GalleryImages      string      `json:"-" gorm:"gallery_images"`                       // 图集
	GalleryImagesArray []string    `json:"galleryImages" gorm:"-" gorm:"-"`               // 图集
	VideoURL           string      `json:"videoUrl" gorm:"video_url"`                     // 视频地址
	Description        string      `json:"description" gorm:"description"`                // 商品描述
	Weight             float64     `json:"weight" gorm:"weight"`                          // 重量
	WeightUnit         string      `json:"weightUnit" gorm:"weight_unit"`                 // 重量单位
	Unit               string      `json:"unit" gorm:"unit"`                              // 销售单位
	IsOnSale           int32       `json:"isOnSale" gorm:"is_on_sale"`                    // 是否上架
	Quantity           int64       `json:"quantity" gorm:"quantity"`                      // 库存数量
	HomeModuleIDs      string      `json:"homeModuleIds" gorm:"home_module_ids"`          // 推荐首页模块ID集合
	Skus               []*GoodsSku `json:"skus" gorm:"-" gorm:"-"`                        // 商品规格列表
	DeptID             int64       `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// GoodsUpsert 商品新增修改参数
type GoodsUpsert struct {
	ID             int64    `json:"id,string"`                                     // 主键ID
	GoodsID        string   `json:"goodsId" binding:"required"`                    // 商品业务ID
	ShopCategoryId int64    `json:"shopCategoryId,string" gorm:"shop_category_id"` // 商品分类id
	GoodsName      string   `json:"goodsName" binding:"required"`                  // 商品名称
	GoodsCode      string   `json:"goodsCode"`                                     // 商品编码
	OuterID        string   `json:"outerId"`                                       // 外部系统ID
	ImageURL       string   `json:"imageUrl"`                                      // 主图地址
	RetailPrice    float64  `json:"retailPrice"`                                   // 零售价
	GalleryImages  []string `json:"galleryImages"`                                 // 图集
	VideoURL       string   `json:"videoUrl"`                                      // 视频地址
	Description    string   `json:"description"`                                   // 商品描述
	Weight         float64  `json:"weight"`                                        // 重量
	WeightUnit     string   `json:"weightUnit"`                                    // 重量单位
	Unit           string   `json:"unit"`                                          // 销售单位
	IsOnSale       int32    `json:"isOnSale"`                                      // 是否上架
	Quantity       int64    `json:"quantity"`                                      // 库存数量
	HomeModuleIDs  []string `json:"homeModuleIds"`                                 // 推荐首页模块ID集合
	baize.BaseEntity
}

// GoodsQuery 商品查询参数
type GoodsQuery struct {
	ID            int64  `form:"id"`            // 主键ID
	GoodsName     string `form:"goodsName"`     // 商品名称
	GoodsCode     string `form:"goodsCode"`     // 商品编码
	OuterID       string `form:"outerId"`       // 外部系统ID
	IsOnSale      *bool  `form:"isOnSale"`      // 是否上架
	CategoryId    int64  `form:"categoryId"`    // 商品分类ID
	StartModified string `form:"startModified"` // 起始修改时间
	EndModified   string `form:"endModified"`   // 结束修改时间
	SortBy        string `form:"sortBy"`        // 排序字段
	SortOrder     string `form:"sortOrder"`     // 排序方向 asc/desc
	Page          int64  `form:"page"`          // 页码
	Size          int64  `form:"size"`          // 每页数量
}

// GoodsListData 商品列表结果
type GoodsListData struct {
	Rows  []*Goods `json:"rows"`  // 数据列表
	Total int64    `json:"total"` // 总数
}

type GoodsVectorResult struct {
	GoodsDBID  int64  `json:"goodsDbId"`
	GoodsID    string `json:"goodsId"`
	Collection string `json:"collection"`
	Dimension  int    `json:"dimension"`
	SkuCount   int    `json:"skuCount"`
}

// ImportGoodsList 商品导入列表
type ImportGoodsList struct {
	Count   int                 `json:"count"`
	Records []ImportGoodsRecord `json:"records"`
}

// ImportGoodsRecord 导入商品结果
type ImportGoodsRecord struct {
	ExternalID string             `json:"external_id"`
	Source     string             `json:"source"`
	Entity     string             `json:"entity"`
	Data       ImportGoodsRawData `json:"data"`
	SyncedAt   time.Time          `json:"synced_at"`
}

// ImportGoodsSkuRawData 导入的sku原始数据
type ImportGoodsSkuRawData struct {
	Barcode  string  `json:"barcode"`
	Lcmccode string  `json:"lcmccode"`
	Price    float64 `json:"price"`
	Price2   float64 `json:"price2"`
	Price3   float64 `json:"price3"`
	Price4   float64 `json:"price4"`
	Price5   float64 `json:"price5"`
	Size     float64 `json:"size"`
	Skucode  string  `json:"skucode"`
	Skuid    string  `json:"skuid"`
	Skuname  string  `json:"skuname"`
	Weight   float64 `json:"weight"`
}

// ImportGoodsRawData 导入的商品原始数据
type ImportGoodsRawData struct {
	ProductCode string                  `json:"product_code"`
	ProductName string                  `json:"product_name"`
	Remark      string                  `json:"remark"`
	Skus        []ImportGoodsSkuRawData `json:"skus"`
	UnitName    string                  `json:"unit_name"`
	Units       []Unit                  `json:"units"`
}

type Unit struct {
	Unitname string  `json:"unitname"`
	Barcode  string  `json:"barcode"`
	Rate     float64 `json:"rate"`
	Price    float64 `json:"price"`
	Price2   float64 `json:"price2"`
	Price3   float64 `json:"price3"`
	Price4   float64 `json:"price4"`
	Price5   float64 `json:"price5"`
}

// GoodsProductData 适配 shopmodels.Goods 为 goods.ProductData 接口，json tag 对齐管家婆 API 标准。
type GoodsProductData struct {
	Cid        int                     `json:"cid"`
	CatName    string                  `json:"catname"`
	ProductId  int64                   `json:"productid"`
	Name       string                  `json:"name"`
	OuterId    string                  `json:"outerid"`
	PicPath    string                  `json:"picpath"`
	Price      int                     `json:"price"`
	BarcodeStr string                  `json:"barcodestr"`
	Created    string                  `json:"created"`
	Desc       string                  `json:"desc"`
	Modified   string                  `json:"modified"`
	Status     string                  `json:"status"`
	Quantity   int                     `json:"quantity"`
	Skus       []goodsstore.ProductSku `json:"skus"`
}

func (g *GoodsProductData) GetCid() int                      { return g.Cid }
func (g *GoodsProductData) GetCatName() string               { return g.CatName }
func (g *GoodsProductData) GetProductId() int64              { return g.ProductId }
func (g *GoodsProductData) GetName() string                  { return g.Name }
func (g *GoodsProductData) GetOuterId() string               { return g.OuterId }
func (g *GoodsProductData) GetPicPath() string               { return g.PicPath }
func (g *GoodsProductData) GetPrice() int                    { return g.Price }
func (g *GoodsProductData) GetBarcodeStr() string            { return g.BarcodeStr }
func (g *GoodsProductData) GetCreated() string               { return g.Created }
func (g *GoodsProductData) GetDesc() string                  { return g.Desc }
func (g *GoodsProductData) GetModified() string              { return g.Modified }
func (g *GoodsProductData) GetStatus() string                { return g.Status }
func (g *GoodsProductData) GetQuantity() int                 { return g.Quantity }
func (g *GoodsProductData) GetSkus() []goodsstore.ProductSku { return g.Skus }

// GoodsProductSku 适配 shopmodels.GoodsSku 为 goods.ProductSku 接口，json tag 对齐管家婆 API 标准。
type GoodsProductSku struct {
	SkuId      int64  `json:"skuid"`
	SkuName    string `json:"skuname"`
	ProductId  int    `json:"productid"`
	OuterId    string `json:"outerid"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	Created    string `json:"created"`
	Modified   string `json:"modified"`
	Properties string `json:"properties"`
}

func (g *GoodsProductSku) GetSkuId() int64       { return g.SkuId }
func (g *GoodsProductSku) GetSkuName() string    { return g.SkuName }
func (g *GoodsProductSku) GetProductId() int     { return g.ProductId }
func (g *GoodsProductSku) GetOuterId() string    { return g.OuterId }
func (g *GoodsProductSku) GetPrice() int         { return g.Price }
func (g *GoodsProductSku) GetQuantity() int      { return g.Quantity }
func (g *GoodsProductSku) GetCreated() string    { return g.Created }
func (g *GoodsProductSku) GetModified() string   { return g.Modified }
func (g *GoodsProductSku) GetProperties() string { return g.Properties }

// GoodsDataResult 实现 goods.DataResult 接口，json tag 对齐管家婆 API 标准。
type GoodsDataResult struct {
	IsError      bool                     `json:"iserror"`
	ErrorMsg     string                   `json:"errormsg"`
	TotalResults int                      `json:"totalresults"`
	ProductInfo  []goodsstore.ProductData `json:"productinfo"`
}

func (r *GoodsDataResult) GetIsError() bool                         { return r.IsError }
func (r *GoodsDataResult) GetErrorMsg() string                      { return r.ErrorMsg }
func (r *GoodsDataResult) GetTotalResults() int                     { return r.TotalResults }
func (r *GoodsDataResult) GetProductInfo() []goodsstore.ProductData { return r.ProductInfo }
