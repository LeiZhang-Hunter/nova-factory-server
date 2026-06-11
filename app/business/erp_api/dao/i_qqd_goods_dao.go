package dao

import (
	"context"

	"nova-factory-server/app/business/erp_api/models"
)

type IQQDGoodsDao interface {
	List(ctx context.Context, pageNo, pageSize int) ([]models.QQDProductTable, int64, error)
	UpdateQuantity(ctx context.Context, goodsID string, quantity int64) error
}
