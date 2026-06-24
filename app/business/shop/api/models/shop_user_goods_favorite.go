package models

import (
	"gorm.io/gorm"
)

// ShopUserGoodsFavorite 用户商品收藏
type ShopUserGoodsFavorite struct {
	gorm.Model
	ID      int64 `json:"id,string" gorm:"id"`            // 主键ID
	UserID  int64 `json:"userId,string" gorm:"user_id"`   // 用户ID
	GoodsID int64 `json:"goodsId,string" gorm:"goods_id"` // 商品ID (VARCHAR)
}

// FavoriteStatusResp 收藏状态查询响应
type FavoriteStatusResp struct {
	IsFavorite bool `json:"isFavorite"` // 是否已收藏
}

// FavoriteAddReq 添加收藏请求
type FavoriteAddReq struct {
	GoodsId int64 `json:"goodsId,string" binding:"required"` // 商品ID (VARCHAR)
}
