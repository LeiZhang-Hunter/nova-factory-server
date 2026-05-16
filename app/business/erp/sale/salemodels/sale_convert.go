package salemodels

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

func SaleOutUpsertToEntity(upsert *SaleOutUpsert) *SaleOut {
	if upsert == nil {
		return nil
	}
	entity := &SaleOut{
		ID:                upsert.ID,
		No:                upsert.No,
		Status:            upsert.Status,
		CustomerID:        upsert.CustomerID,
		AccountID:         upsert.AccountID,
		SaleUserID:        upsert.SaleUserID,
		OrderID:           upsert.OrderID,
		OrderNo:           upsert.OrderNo,
		TotalCount:        upsert.TotalCount,
		TotalPrice:        upsert.TotalPrice,
		ReceiptPrice:      upsert.ReceiptPrice,
		TotalProductPrice: upsert.TotalProductPrice,
		TotalTaxPrice:     upsert.TotalTaxPrice,
		DiscountPercent:   upsert.DiscountPercent,
		DiscountPrice:     upsert.DiscountPrice,
		OtherPrice:        upsert.OtherPrice,
		FileURL:           upsert.FileURL,
		Remark:            upsert.Remark,
	}
	if upsert.OutTime != "" {
		if t, err := parseTime(upsert.OutTime); err == nil {
			entity.OutTime = t
		}
	}
	return entity
}

func SaleOutToUpsert(entity *SaleOut) *SaleOutUpsert {
	if entity == nil {
		return nil
	}
	upsert := &SaleOutUpsert{
		ID:                entity.ID,
		No:                entity.No,
		Status:            entity.Status,
		CustomerID:        entity.CustomerID,
		AccountID:         entity.AccountID,
		SaleUserID:        entity.SaleUserID,
		OrderID:           entity.OrderID,
		OrderNo:           entity.OrderNo,
		TotalCount:        entity.TotalCount,
		TotalPrice:        entity.TotalPrice,
		ReceiptPrice:      entity.ReceiptPrice,
		TotalProductPrice: entity.TotalProductPrice,
		TotalTaxPrice:     entity.TotalTaxPrice,
		DiscountPercent:   entity.DiscountPercent,
		DiscountPrice:     entity.DiscountPrice,
		OtherPrice:        entity.OtherPrice,
		FileURL:           entity.FileURL,
		Remark:            entity.Remark,
	}
	if entity.OutTime != nil {
		upsert.OutTime = formatTime(entity.OutTime)
	}
	return upsert
}

func SaleOutClone(entity *SaleOut) *SaleOut {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func SaleOutItemUpsertToEntity(upsert *SaleOutItemUpsert) *SaleOutItem {
	if upsert == nil {
		return nil
	}
	return &SaleOutItem{
		ID:            upsert.ID,
		OutID:         upsert.OutID,
		OrderItemID:   upsert.OrderItemID,
		WarehouseID:   upsert.WarehouseID,
		ProductID:     upsert.ProductID,
		ProductUnitID: upsert.ProductUnitID,
		ProductPrice:  upsert.ProductPrice,
		Count:         upsert.Count,
		TotalPrice:    upsert.TotalPrice,
		TaxPercent:    upsert.TaxPercent,
		TaxPrice:      upsert.TaxPrice,
		Remark:        upsert.Remark,
	}
}

func SaleOutItemToUpsert(entity *SaleOutItem) *SaleOutItemUpsert {
	if entity == nil {
		return nil
	}
	return &SaleOutItemUpsert{
		ID:            entity.ID,
		OutID:         entity.OutID,
		OrderItemID:   entity.OrderItemID,
		WarehouseID:   entity.WarehouseID,
		ProductID:     entity.ProductID,
		ProductUnitID: entity.ProductUnitID,
		ProductPrice:  entity.ProductPrice,
		Count:         entity.Count,
		TotalPrice:    entity.TotalPrice,
		TaxPercent:    entity.TaxPercent,
		TaxPrice:      entity.TaxPrice,
		Remark:        entity.Remark,
	}
}

func SaleOutItemClone(entity *SaleOutItem) *SaleOutItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func SaleReturnUpsertToEntity(upsert *SaleReturnUpsert) *SaleReturn {
	if upsert == nil {
		return nil
	}
	entity := &SaleReturn{
		ID:                upsert.ID,
		No:                upsert.No,
		Status:            upsert.Status,
		CustomerID:        upsert.CustomerID,
		AccountID:         upsert.AccountID,
		SaleUserID:        upsert.SaleUserID,
		OrderID:           upsert.OrderID,
		OrderNo:           upsert.OrderNo,
		TotalCount:        upsert.TotalCount,
		TotalPrice:        upsert.TotalPrice,
		RefundPrice:       upsert.RefundPrice,
		TotalProductPrice: upsert.TotalProductPrice,
		TotalTaxPrice:     upsert.TotalTaxPrice,
		DiscountPercent:   upsert.DiscountPercent,
		DiscountPrice:     upsert.DiscountPrice,
		OtherPrice:        upsert.OtherPrice,
		FileURL:           upsert.FileURL,
		Remark:            upsert.Remark,
	}
	if upsert.ReturnTime != "" {
		if t, err := parseTime(upsert.ReturnTime); err == nil {
			entity.ReturnTime = t
		}
	}
	return entity
}

func SaleReturnToUpsert(entity *SaleReturn) *SaleReturnUpsert {
	if entity == nil {
		return nil
	}
	upsert := &SaleReturnUpsert{
		ID:                entity.ID,
		No:                entity.No,
		Status:            entity.Status,
		CustomerID:        entity.CustomerID,
		AccountID:         entity.AccountID,
		SaleUserID:        entity.SaleUserID,
		OrderID:           entity.OrderID,
		OrderNo:           entity.OrderNo,
		TotalCount:        entity.TotalCount,
		TotalPrice:        entity.TotalPrice,
		RefundPrice:       entity.RefundPrice,
		TotalProductPrice: entity.TotalProductPrice,
		TotalTaxPrice:     entity.TotalTaxPrice,
		DiscountPercent:   entity.DiscountPercent,
		DiscountPrice:     entity.DiscountPrice,
		OtherPrice:        entity.OtherPrice,
		FileURL:           entity.FileURL,
		Remark:            entity.Remark,
	}
	if entity.ReturnTime != nil {
		upsert.ReturnTime = formatTime(entity.ReturnTime)
	}
	return upsert
}

func SaleReturnClone(entity *SaleReturn) *SaleReturn {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func SaleReturnItemUpsertToEntity(upsert *SaleReturnItemUpsert) *SaleReturnItem {
	if upsert == nil {
		return nil
	}
	return &SaleReturnItem{
		ID:            upsert.ID,
		ReturnID:      upsert.ReturnID,
		OrderItemID:   upsert.OrderItemID,
		WarehouseID:   upsert.WarehouseID,
		ProductID:     upsert.ProductID,
		ProductUnitID: upsert.ProductUnitID,
		ProductPrice:  upsert.ProductPrice,
		Count:         upsert.Count,
		TotalPrice:    upsert.TotalPrice,
		TaxPercent:    upsert.TaxPercent,
		TaxPrice:      upsert.TaxPrice,
		Remark:        upsert.Remark,
	}
}

func SaleReturnItemToUpsert(entity *SaleReturnItem) *SaleReturnItemUpsert {
	if entity == nil {
		return nil
	}
	return &SaleReturnItemUpsert{
		ID:            entity.ID,
		ReturnID:      entity.ReturnID,
		OrderItemID:   entity.OrderItemID,
		WarehouseID:   entity.WarehouseID,
		ProductID:     entity.ProductID,
		ProductUnitID: entity.ProductUnitID,
		ProductPrice:  entity.ProductPrice,
		Count:         entity.Count,
		TotalPrice:    entity.TotalPrice,
		TaxPercent:    entity.TaxPercent,
		TaxPrice:      entity.TaxPrice,
		Remark:        entity.Remark,
	}
}

func SaleReturnItemClone(entity *SaleReturnItem) *SaleReturnItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}
