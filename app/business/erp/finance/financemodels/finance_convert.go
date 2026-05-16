package financemodels

import (
	"errors"
	"time"
)

func parseTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	layouts := []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &parsed, nil
		}
	}
	return nil, errors.New("时间格式错误")
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func FinanceReceiptUpsertToEntity(upsert *FinanceReceiptUpsert) *FinanceReceipt {
	if upsert == nil {
		return nil
	}
	entity := &FinanceReceipt{
		ID:            upsert.ID,
		No:            upsert.No,
		Status:        upsert.Status,
		FinanceUserID: upsert.FinanceUserID,
		CustomerID:    upsert.CustomerID,
		AccountID:     upsert.AccountID,
		TotalPrice:    upsert.TotalPrice,
		DiscountPrice: upsert.DiscountPrice,
		ReceiptPrice:  upsert.ReceiptPrice,
		Remark:        upsert.Remark,
	}
	if upsert.ReceiptTime != "" {
		if t, err := parseTime(upsert.ReceiptTime); err == nil {
			entity.ReceiptTime = t
		}
	}
	return entity
}

func FinanceReceiptToUpsert(entity *FinanceReceipt) *FinanceReceiptUpsert {
	if entity == nil {
		return nil
	}
	upsert := &FinanceReceiptUpsert{
		ID:            entity.ID,
		No:            entity.No,
		Status:        entity.Status,
		FinanceUserID: entity.FinanceUserID,
		CustomerID:    entity.CustomerID,
		AccountID:     entity.AccountID,
		TotalPrice:    entity.TotalPrice,
		DiscountPrice: entity.DiscountPrice,
		ReceiptPrice:  entity.ReceiptPrice,
		Remark:        entity.Remark,
	}
	if entity.ReceiptTime != nil {
		upsert.ReceiptTime = formatTime(entity.ReceiptTime)
	}
	return upsert
}

func FinanceReceiptClone(entity *FinanceReceipt) *FinanceReceipt {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func FinanceReceiptItemUpsertToEntity(upsert *FinanceReceiptItemUpsert) *FinanceReceiptItem {
	if upsert == nil {
		return nil
	}
	return &FinanceReceiptItem{
		ID:             upsert.ID,
		ReceiptID:      upsert.ReceiptID,
		BizType:        upsert.BizType,
		BizID:          upsert.BizID,
		BizNo:          upsert.BizNo,
		TotalPrice:     upsert.TotalPrice,
		ReceiptedPrice: upsert.ReceiptedPrice,
		ReceiptPrice:   upsert.ReceiptPrice,
		Remark:         upsert.Remark,
	}
}

func FinanceReceiptItemToUpsert(entity *FinanceReceiptItem) *FinanceReceiptItemUpsert {
	if entity == nil {
		return nil
	}
	return &FinanceReceiptItemUpsert{
		ID:             entity.ID,
		ReceiptID:      entity.ReceiptID,
		BizType:        entity.BizType,
		BizID:          entity.BizID,
		BizNo:          entity.BizNo,
		TotalPrice:     entity.TotalPrice,
		ReceiptedPrice: entity.ReceiptedPrice,
		ReceiptPrice:   entity.ReceiptPrice,
		Remark:         entity.Remark,
	}
}

func FinanceReceiptItemClone(entity *FinanceReceiptItem) *FinanceReceiptItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func FinancePaymentUpsertToEntity(upsert *FinancePaymentUpsert) *FinancePayment {
	if upsert == nil {
		return nil
	}
	entity := &FinancePayment{
		ID:            upsert.ID,
		No:            upsert.No,
		Status:        upsert.Status,
		FinanceUserID: upsert.FinanceUserID,
		SupplierID:    upsert.SupplierID,
		AccountID:     upsert.AccountID,
		TotalPrice:    upsert.TotalPrice,
		DiscountPrice: upsert.DiscountPrice,
		PaymentPrice:  upsert.PaymentPrice,
		Remark:        upsert.Remark,
	}
	if upsert.PaymentTime != "" {
		if t, err := parseTime(upsert.PaymentTime); err == nil {
			entity.PaymentTime = t
		}
	}
	return entity
}

func FinancePaymentToUpsert(entity *FinancePayment) *FinancePaymentUpsert {
	if entity == nil {
		return nil
	}
	upsert := &FinancePaymentUpsert{
		ID:            entity.ID,
		No:            entity.No,
		Status:        entity.Status,
		FinanceUserID: entity.FinanceUserID,
		SupplierID:    entity.SupplierID,
		AccountID:     entity.AccountID,
		TotalPrice:    entity.TotalPrice,
		DiscountPrice: entity.DiscountPrice,
		PaymentPrice:  entity.PaymentPrice,
		Remark:        entity.Remark,
	}
	if entity.PaymentTime != nil {
		upsert.PaymentTime = formatTime(entity.PaymentTime)
	}
	return upsert
}

func FinancePaymentClone(entity *FinancePayment) *FinancePayment {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func FinancePaymentItemUpsertToEntity(upsert *FinancePaymentItemUpsert) *FinancePaymentItem {
	if upsert == nil {
		return nil
	}
	return &FinancePaymentItem{
		ID:           upsert.ID,
		PaymentID:    upsert.PaymentID,
		BizType:      upsert.BizType,
		BizID:        upsert.BizID,
		BizNo:        upsert.BizNo,
		TotalPrice:   upsert.TotalPrice,
		PaidPrice:    upsert.PaidPrice,
		PaymentPrice: upsert.PaymentPrice,
		Remark:       upsert.Remark,
	}
}

func FinancePaymentItemToUpsert(entity *FinancePaymentItem) *FinancePaymentItemUpsert {
	if entity == nil {
		return nil
	}
	return &FinancePaymentItemUpsert{
		ID:           entity.ID,
		PaymentID:    entity.PaymentID,
		BizType:      entity.BizType,
		BizID:        entity.BizID,
		BizNo:        entity.BizNo,
		TotalPrice:   entity.TotalPrice,
		PaidPrice:    entity.PaidPrice,
		PaymentPrice: entity.PaymentPrice,
		Remark:       entity.Remark,
	}
}

func FinancePaymentItemClone(entity *FinancePaymentItem) *FinancePaymentItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}
