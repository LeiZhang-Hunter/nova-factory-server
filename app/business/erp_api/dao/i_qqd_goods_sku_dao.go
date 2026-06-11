package dao

import (
	"context"

	"nova-factory-server/app/business/erp_api/models"
)

type IQQDGoodsSkuDao interface {
	ListByGoodsIDs(ctx context.Context, goodsIDs []string) ([]models.QQDProductSkuTable, error)
	UpdateQuantity(ctx context.Context, skuID string, quantity int64, lock bool) error
}
