package impl

import (
	"errors"
	"fmt"
	"nova-factory-server/app/business/shop/order/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func (o *OrderServiceImpl) fillOrderSyncRequestFromDB(c *gin.Context, req *models.OrderSyncRequest) error {
	if req == nil || len(req.Orders) == 0 {
		return errors.New("orders不能为空")
	}

	orders := make([]*models.OrderSyncOrder, 0, len(req.Orders))
	seen := make(map[string]struct{}, len(req.Orders))
	for _, info := range req.Orders {
		if info == nil {
			continue
		}
		tid := strings.TrimSpace(info.Tid)
		info.Tid = tid
		if tid == "" {
			return errors.New("tid不能为空")
		}
		if _, ok := seen[tid]; ok {
			continue
		}
		seen[tid] = struct{}{}

		orderInfo, err := o.orderDao.GetByTid(c, tid)
		if err != nil {
			return err
		}
		if orderInfo == nil {
			return fmt.Errorf("订单不存在: %s", tid)
		}
		orders = append(orders, models.ToOrderSyncOrder(orderInfo, nil))
	}
	if len(orders) == 0 {
		return errors.New("orders不能为空")
	}
	req.Orders = orders
	return nil
}
