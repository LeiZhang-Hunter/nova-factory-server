package purchasemodels

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

func PurchaseOrderUpsertToEntity(upsert *PurchaseOrderUpsert) *PurchaseOrder {
	if upsert == nil {
		return nil
	}
	entity := &PurchaseOrder{
		ID:              upsert.ID,
		No:              upsert.No,
		SupplierID:      upsert.SupplierID,
		AccountID:       upsert.AccountID,
		DiscountPercent: upsert.DiscountPercent,
		DepositPrice:    upsert.DepositPrice,
		FileURL:         upsert.FileURL,
		Remark:          upsert.Remark,
		Items:           make([]*PurchaseOrderItem, 0),
	}
	if upsert.OrderTime != "" {
		if t, err := parseTime(upsert.OrderTime); err == nil {
			entity.OrderTime = t
		}
	}
	for _, item := range upsert.Items {
		if item == nil {
			continue
		}
		entity.Items = append(entity.Items, PurchaseOrderItemUpsertToEntity(item))
	}
	return entity
}

func PurchaseOrderToUpsert(entity *PurchaseOrder) *PurchaseOrderUpsert {
	if entity == nil {
		return nil
	}
	upsert := &PurchaseOrderUpsert{
		ID:              entity.ID,
		No:              entity.No,
		SupplierID:      entity.SupplierID,
		AccountID:       entity.AccountID,
		DiscountPercent: entity.DiscountPercent,
		DepositPrice:    entity.DepositPrice,
		FileURL:         entity.FileURL,
		Remark:          entity.Remark,
		Items:           make([]*PurchaseOrderItemUpsert, 0),
	}
	if entity.OrderTime != nil {
		upsert.OrderTime = formatTime(entity.OrderTime)
	}
	for _, item := range entity.Items {
		if item == nil {
			continue
		}
		upsert.Items = append(upsert.Items, PurchaseOrderItemToUpsert(item))
	}
	return upsert
}

func PurchaseOrderClone(entity *PurchaseOrder) *PurchaseOrder {
	if entity == nil {
		return nil
	}
	clone := &PurchaseOrder{
		ID:                entity.ID,
		No:                entity.No,
		Status:            entity.Status,
		SupplierID:        entity.SupplierID,
		AccountID:         entity.AccountID,
		OrderTime:         entity.OrderTime,
		TotalCount:        entity.TotalCount,
		TotalPrice:        entity.TotalPrice,
		TotalProductPrice: entity.TotalProductPrice,
		TotalTaxPrice:     entity.TotalTaxPrice,
		DiscountPercent:   entity.DiscountPercent,
		DiscountPrice:     entity.DiscountPrice,
		DepositPrice:      entity.DepositPrice,
		FileURL:           entity.FileURL,
		Remark:            entity.Remark,
		InCount:           entity.InCount,
		ReturnCount:       entity.ReturnCount,
		DeptID:            entity.DeptID,
		State:             entity.State,
		Items:             make([]*PurchaseOrderItem, len(entity.Items)),
	}
	for i, item := range entity.Items {
		if item == nil {
			continue
		}
		cloneItem := *item
		clone.Items[i] = &cloneItem
	}
	return clone
}

func PurchaseOrderItemUpsertToEntity(upsert *PurchaseOrderItemUpsert) *PurchaseOrderItem {
	if upsert == nil {
		return nil
	}
	return &PurchaseOrderItem{
		ID:            upsert.ID,
		OrderID:       upsert.OrderID,
		ProductID:     upsert.ProductID,
		ProductUnitID: upsert.ProductUnitID,
		ProductPrice:  upsert.ProductPrice,
		Count:         upsert.Count,
		TotalPrice:    upsert.TotalPrice,
		TaxPercent:    upsert.TaxPercent,
		TaxPrice:      upsert.TaxPrice,
		Remark:        upsert.Remark,
		InCount:       upsert.InCount,
		ReturnCount:   upsert.ReturnCount,
	}
}

func PurchaseOrderItemToUpsert(entity *PurchaseOrderItem) *PurchaseOrderItemUpsert {
	if entity == nil {
		return nil
	}
	return &PurchaseOrderItemUpsert{
		ID:            entity.ID,
		OrderID:       entity.OrderID,
		ProductID:     entity.ProductID,
		ProductUnitID: entity.ProductUnitID,
		ProductPrice:  entity.ProductPrice,
		Count:         entity.Count,
		TotalPrice:    entity.TotalPrice,
		TaxPercent:    entity.TaxPercent,
		TaxPrice:      entity.TaxPrice,
		Remark:        entity.Remark,
		InCount:       entity.InCount,
		ReturnCount:   entity.ReturnCount,
	}
}

func PurchaseOrderItemClone(entity *PurchaseOrderItem) *PurchaseOrderItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func PurchaseInUpsertToEntity(upsert *PurchaseInUpsert) *PurchaseIn {
	if upsert == nil {
		return nil
	}
	entity := &PurchaseIn{
		ID:                upsert.ID,
		No:                upsert.No,
		Status:            upsert.Status,
		SupplierID:        upsert.SupplierID,
		AccountID:         upsert.AccountID,
		OrderID:           upsert.OrderID,
		OrderNo:           upsert.OrderNo,
		TotalCount:        upsert.TotalCount,
		TotalPrice:        upsert.TotalPrice,
		PaymentPrice:      upsert.PaymentPrice,
		TotalProductPrice: upsert.TotalProductPrice,
		TotalTaxPrice:     upsert.TotalTaxPrice,
		DiscountPercent:   upsert.DiscountPercent,
		DiscountPrice:     upsert.DiscountPrice,
		OtherPrice:        upsert.OtherPrice,
		FileURL:           upsert.FileURL,
		Remark:            upsert.Remark,
	}
	if upsert.InTime != "" {
		if t, err := parseTime(upsert.InTime); err == nil {
			entity.InTime = t
		}
	}
	return entity
}

func PurchaseInToUpsert(entity *PurchaseIn) *PurchaseInUpsert {
	if entity == nil {
		return nil
	}
	upsert := &PurchaseInUpsert{
		ID:                entity.ID,
		No:                entity.No,
		Status:            entity.Status,
		SupplierID:        entity.SupplierID,
		AccountID:         entity.AccountID,
		OrderID:           entity.OrderID,
		OrderNo:           entity.OrderNo,
		TotalCount:        entity.TotalCount,
		TotalPrice:        entity.TotalPrice,
		PaymentPrice:      entity.PaymentPrice,
		TotalProductPrice: entity.TotalProductPrice,
		TotalTaxPrice:     entity.TotalTaxPrice,
		DiscountPercent:   entity.DiscountPercent,
		DiscountPrice:     entity.DiscountPrice,
		OtherPrice:        entity.OtherPrice,
		FileURL:           entity.FileURL,
		Remark:            entity.Remark,
	}
	if entity.InTime != nil {
		upsert.InTime = formatTime(entity.InTime)
	}
	return upsert
}

func PurchaseInClone(entity *PurchaseIn) *PurchaseIn {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func PurchaseInItemUpsertToEntity(upsert *PurchaseInItemUpsert) *PurchaseInItem {
	if upsert == nil {
		return nil
	}
	return &PurchaseInItem{
		ID:            upsert.ID,
		InID:          upsert.InID,
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

func PurchaseInItemToUpsert(entity *PurchaseInItem) *PurchaseInItemUpsert {
	if entity == nil {
		return nil
	}
	return &PurchaseInItemUpsert{
		ID:            entity.ID,
		InID:          entity.InID,
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

func PurchaseInItemClone(entity *PurchaseInItem) *PurchaseInItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func PurchaseReturnUpsertToEntity(upsert *PurchaseReturnUpsert) *PurchaseReturn {
	if upsert == nil {
		return nil
	}
	entity := &PurchaseReturn{
		ID:                upsert.ID,
		No:                upsert.No,
		Status:            upsert.Status,
		SupplierID:        upsert.SupplierID,
		AccountID:         upsert.AccountID,
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

func PurchaseReturnToUpsert(entity *PurchaseReturn) *PurchaseReturnUpsert {
	if entity == nil {
		return nil
	}
	upsert := &PurchaseReturnUpsert{
		ID:                entity.ID,
		No:                entity.No,
		Status:            entity.Status,
		SupplierID:        entity.SupplierID,
		AccountID:         entity.AccountID,
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

func PurchaseReturnClone(entity *PurchaseReturn) *PurchaseReturn {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func PurchaseReturnItemUpsertToEntity(upsert *PurchaseReturnItemUpsert) *PurchaseReturnItem {
	if upsert == nil {
		return nil
	}
	return &PurchaseReturnItem{
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

func PurchaseReturnItemToUpsert(entity *PurchaseReturnItem) *PurchaseReturnItemUpsert {
	if entity == nil {
		return nil
	}
	return &PurchaseReturnItemUpsert{
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

func PurchaseReturnItemClone(entity *PurchaseReturnItem) *PurchaseReturnItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}
