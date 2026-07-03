package impl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nova-factory-server/app/business/shop/order/models"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/utils/baizeContext"
)

// ReviewRefund 审核售后单。审核通过后由后台退款按钮手动触发支付退款。
func (o *OrderServiceImpl) ReviewRefund(c *gin.Context, req *models.RefundReviewReq) error {
	if req == nil || req.ID == 0 {
		return errors.New("售后单ID不能为空")
	}
	aftersale, err := o.orderRefundDao.GetByID(c, req.ID)
	if err != nil {
		return fmt.Errorf("查询售后单失败: %v", err)
	}
	if aftersale == nil {
		return errors.New("售后单不存在")
	}
	if aftersale.Status != orderConstant.AftersaleStatusPendingReview {
		return errors.New("当前售后状态不允许审核")
	}

	now := time.Now()
	userID := baizeContext.GetUserId(c)
	updates := map[string]any{
		"audit_by":     userID,
		"audit_time":   &now,
		"audit_remark": strings.TrimSpace(req.Remark),
		"update_by":    userID,
		"update_time":  &now,
	}
	targetStatus := orderConstant.AftersaleStatusRejected
	if req.Approved {
		targetStatus = orderConstant.AftersaleStatusApproved
		if !isReturnRefundAftersale(aftersale.PreviousStatus) {
			targetStatus = orderConstant.AftersaleStatusPendingRefund
			updates["erp_sync_status"] = orderConstant.AftersaleSyncSuccess
			updates["erp_sync_message"] = "仅退款无需同步ERP"
			updates["erp_sync_time"] = &now
		}
	}

	if err := o.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := o.orderRefundDao.UpdateStatusWithTx(tx, aftersale.ID, targetStatus, updates); err != nil {
			return err
		}
		if !req.Approved {
			return o.orderDao.UpdateByID(tx, uint64(aftersale.OrderID), map[string]any{
				"status":      aftersale.PreviousStatus,
				"update_by":   userID,
				"update_time": &now,
			})
		}
		return nil
	}); err != nil {
		return fmt.Errorf("审核售后单失败: %v", err)
	}

	aftersale.Status = targetStatus
	aftersale.AuditBy = userID
	aftersale.AuditTime = &now
	aftersale.AuditRemark = strings.TrimSpace(req.Remark)
	if !req.Approved {
		return nil
	}
	return nil
}

func isReturnRefundAftersale(previousStatus string) bool {
	switch strings.TrimSpace(previousStatus) {
	case orderConstant.ERPStatusSended, orderConstant.ERPStatusPartSend:
		return true
	default:
		return false
	}
}

// Refund 后台手动触发支付退款。
func (o *OrderServiceImpl) Refund(c *gin.Context, req *models.RefundManualReq) error {
	if req == nil || req.ID == 0 {
		return errors.New("售后单ID不能为空")
	}
	aftersale, err := o.orderRefundDao.GetByID(c, req.ID)
	if err != nil {
		return fmt.Errorf("查询售后单失败: %v", err)
	}
	if aftersale == nil {
		return errors.New("售后单不存在")
	}
	if aftersale.Status != orderConstant.AftersaleStatusPendingRefund &&
		aftersale.Status != orderConstant.AftersaleStatusRefundFailed {
		return errors.New("当前售后状态不允许退款")
	}
	return o.refundViaPaymentChannel(c, aftersale)
}
