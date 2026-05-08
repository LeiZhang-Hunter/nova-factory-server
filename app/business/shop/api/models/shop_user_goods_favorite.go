package models

import (
	"gorm.io/gorm"
)

// ShopUserGoodsFavorite 用户商品收藏
type ShopUserGoodsFavorite struct {
	gorm.Model
	ID      int64  `json:"id,string" db:"id"`     // 主键ID
	UserID  int64  `json:"userId" db:"user_id"`   // 用户ID
	GoodsID string `json:"goodsId" db:"goods_id"` // 商品ID (VARCHAR)
}

// FavoriteStatusResp 收藏状态查询响应
type FavoriteStatusResp struct {
	IsFavorite bool `json:"isFavorite"` // 是否已收藏
}

// FavoriteAddReq 添加收藏请求
type FavoriteAddReq struct {
	GoodsId string `json:"goodsId" binding:"required"` // 商品ID (VARCHAR)
}
