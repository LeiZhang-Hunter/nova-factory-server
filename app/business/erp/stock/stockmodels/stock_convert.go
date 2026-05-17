package stockmodels

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

func StockCheckUpsertToEntity(upsert *StockCheckUpsert) *StockCheck {
	if upsert == nil {
		return nil
	}
	entity := &StockCheck{
		ID:         upsert.ID,
		No:         upsert.No,
		TotalCount: upsert.TotalCount,
		TotalPrice: upsert.TotalPrice,
		Status:     upsert.Status,
		Remark:     upsert.Remark,
		FileURL:    upsert.FileURL,
	}
	if upsert.CheckTime != "" {
		if t, err := parseTime(upsert.CheckTime); err == nil {
			entity.CheckTime = t
		}
	}
	return entity
}

func StockCheckToUpsert(entity *StockCheck) *StockCheckUpsert {
	if entity == nil {
		return nil
	}
	upsert := &StockCheckUpsert{
		ID:         entity.ID,
		No:         entity.No,
		TotalCount: entity.TotalCount,
		TotalPrice: entity.TotalPrice,
		Status:     entity.Status,
		Remark:     entity.Remark,
		FileURL:    entity.FileURL,
	}
	if entity.CheckTime != nil {
		upsert.CheckTime = formatTime(entity.CheckTime)
	}
	return upsert
}

func StockCheckClone(entity *StockCheck) *StockCheck {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockCheckItemUpsertToEntity(upsert *StockCheckItemUpsert) *StockCheckItem {
	if upsert == nil {
		return nil
	}
	return &StockCheckItem{
		ID:            upsert.ID,
		CheckID:       upsert.CheckID,
		WarehouseID:   upsert.WarehouseID,
		ProductID:     upsert.ProductID,
		ProductUnitID: upsert.ProductUnitID,
		ProductPrice:  upsert.ProductPrice,
		StockCount:    upsert.StockCount,
		ActualCount:   upsert.ActualCount,
		Count:         upsert.Count,
		TotalPrice:    upsert.TotalPrice,
		Remark:        upsert.Remark,
	}
}

func StockCheckItemToUpsert(entity *StockCheckItem) *StockCheckItemUpsert {
	if entity == nil {
		return nil
	}
	return &StockCheckItemUpsert{
		ID:            entity.ID,
		CheckID:       entity.CheckID,
		WarehouseID:   entity.WarehouseID,
		ProductID:     entity.ProductID,
		ProductUnitID: entity.ProductUnitID,
		ProductPrice:  entity.ProductPrice,
		StockCount:    entity.StockCount,
		ActualCount:   entity.ActualCount,
		Count:         entity.Count,
		TotalPrice:    entity.TotalPrice,
		Remark:        entity.Remark,
	}
}

func StockCheckItemClone(entity *StockCheckItem) *StockCheckItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockInUpsertToEntity(upsert *StockInUpsert) *StockIn {
	if upsert == nil {
		return nil
	}
	entity := &StockIn{
		ID:         upsert.ID,
		No:         upsert.No,
		SupplierID: upsert.SupplierID,
		TotalCount: upsert.TotalCount,
		TotalPrice: upsert.TotalPrice,
		Status:     upsert.Status,
		Remark:     upsert.Remark,
		FileURL:    upsert.FileURL,
	}
	if upsert.InTime != "" {
		if t, err := parseTime(upsert.InTime); err == nil {
			entity.InTime = t
		}
	}
	return entity
}

func StockInToUpsert(entity *StockIn) *StockInUpsert {
	if entity == nil {
		return nil
	}
	upsert := &StockInUpsert{
		ID:         entity.ID,
		No:         entity.No,
		SupplierID: entity.SupplierID,
		TotalCount: entity.TotalCount,
		TotalPrice: entity.TotalPrice,
		Status:     entity.Status,
		Remark:     entity.Remark,
		FileURL:    entity.FileURL,
	}
	if entity.InTime != nil {
		upsert.InTime = formatTime(entity.InTime)
	}
	return upsert
}

func StockInClone(entity *StockIn) *StockIn {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockInItemUpsertToEntity(upsert *StockInItemUpsert) *StockInItem {
	if upsert == nil {
		return nil
	}
	return &StockInItem{
		ID:            upsert.ID,
		InID:          upsert.InID,
		WarehouseID:   upsert.WarehouseID,
		ProductID:     upsert.ProductID,
		ProductUnitID: upsert.ProductUnitID,
		ProductPrice:  upsert.ProductPrice,
		Count:         upsert.Count,
		TotalPrice:    upsert.TotalPrice,
		Remark:        upsert.Remark,
	}
}

func StockInItemToUpsert(entity *StockInItem) *StockInItemUpsert {
	if entity == nil {
		return nil
	}
	return &StockInItemUpsert{
		ID:            entity.ID,
		InID:          entity.InID,
		WarehouseID:   entity.WarehouseID,
		ProductID:     entity.ProductID,
		ProductUnitID: entity.ProductUnitID,
		ProductPrice:  entity.ProductPrice,
		Count:         entity.Count,
		TotalPrice:    entity.TotalPrice,
		Remark:        entity.Remark,
	}
}

func StockInItemClone(entity *StockInItem) *StockInItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockMoveUpsertToEntity(upsert *StockMoveUpsert) *StockMove {
	if upsert == nil {
		return nil
	}
	entity := &StockMove{
		ID:         upsert.ID,
		No:         upsert.No,
		TotalCount: upsert.TotalCount,
		TotalPrice: upsert.TotalPrice,
		Status:     upsert.Status,
		Remark:     upsert.Remark,
		FileURL:    upsert.FileURL,
	}
	if upsert.MoveTime != "" {
		if t, err := parseTime(upsert.MoveTime); err == nil {
			entity.MoveTime = t
		}
	}
	return entity
}

func StockMoveToUpsert(entity *StockMove) *StockMoveUpsert {
	if entity == nil {
		return nil
	}
	upsert := &StockMoveUpsert{
		ID:         entity.ID,
		No:         entity.No,
		TotalCount: entity.TotalCount,
		TotalPrice: entity.TotalPrice,
		Status:     entity.Status,
		Remark:     entity.Remark,
		FileURL:    entity.FileURL,
	}
	if entity.MoveTime != nil {
		upsert.MoveTime = formatTime(entity.MoveTime)
	}
	return upsert
}

func StockMoveClone(entity *StockMove) *StockMove {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockMoveItemUpsertToEntity(upsert *StockMoveItemUpsert) *StockMoveItem {
	if upsert == nil {
		return nil
	}
	return &StockMoveItem{
		ID:              upsert.ID,
		MoveID:          upsert.MoveID,
		FromWarehouseID: upsert.FromWarehouseID,
		ToWarehouseID:   upsert.ToWarehouseID,
		ProductID:       upsert.ProductID,
		ProductUnitID:   upsert.ProductUnitID,
		ProductPrice:    upsert.ProductPrice,
		Count:           upsert.Count,
		TotalPrice:      upsert.TotalPrice,
		Remark:          upsert.Remark,
	}
}

func StockMoveItemToUpsert(entity *StockMoveItem) *StockMoveItemUpsert {
	if entity == nil {
		return nil
	}
	return &StockMoveItemUpsert{
		ID:              entity.ID,
		MoveID:          entity.MoveID,
		FromWarehouseID: entity.FromWarehouseID,
		ToWarehouseID:   entity.ToWarehouseID,
		ProductID:       entity.ProductID,
		ProductUnitID:   entity.ProductUnitID,
		ProductPrice:    entity.ProductPrice,
		Count:           entity.Count,
		TotalPrice:      entity.TotalPrice,
		Remark:          entity.Remark,
	}
}

func StockMoveItemClone(entity *StockMoveItem) *StockMoveItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockOutUpsertToEntity(upsert *StockOutUpsert) *StockOut {
	if upsert == nil {
		return nil
	}
	entity := &StockOut{
		ID:         upsert.ID,
		No:         upsert.No,
		CustomerID: upsert.CustomerID,
		TotalCount: upsert.TotalCount,
		TotalPrice: upsert.TotalPrice,
		Status:     upsert.Status,
		Remark:     upsert.Remark,
		FileURL:    upsert.FileURL,
	}
	if upsert.OutTime != "" {
		if t, err := parseTime(upsert.OutTime); err == nil {
			entity.OutTime = t
		}
	}
	return entity
}

func StockOutToUpsert(entity *StockOut) *StockOutUpsert {
	if entity == nil {
		return nil
	}
	upsert := &StockOutUpsert{
		ID:         entity.ID,
		No:         entity.No,
		CustomerID: entity.CustomerID,
		TotalCount: entity.TotalCount,
		TotalPrice: entity.TotalPrice,
		Status:     entity.Status,
		Remark:     entity.Remark,
		FileURL:    entity.FileURL,
	}
	if entity.OutTime != nil {
		upsert.OutTime = formatTime(entity.OutTime)
	}
	return upsert
}

func StockOutClone(entity *StockOut) *StockOut {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockOutItemUpsertToEntity(upsert *StockOutItemUpsert) *StockOutItem {
	if upsert == nil {
		return nil
	}
	return &StockOutItem{
		ID:            upsert.ID,
		OutID:         upsert.OutID,
		WarehouseID:   upsert.WarehouseID,
		ProductID:     upsert.ProductID,
		ProductUnitID: upsert.ProductUnitID,
		ProductPrice:  upsert.ProductPrice,
		Count:         upsert.Count,
		TotalPrice:    upsert.TotalPrice,
		Remark:        upsert.Remark,
	}
}

func StockOutItemToUpsert(entity *StockOutItem) *StockOutItemUpsert {
	if entity == nil {
		return nil
	}
	return &StockOutItemUpsert{
		ID:            entity.ID,
		OutID:         entity.OutID,
		WarehouseID:   entity.WarehouseID,
		ProductID:     entity.ProductID,
		ProductUnitID: entity.ProductUnitID,
		ProductPrice:  entity.ProductPrice,
		Count:         entity.Count,
		TotalPrice:    entity.TotalPrice,
		Remark:        entity.Remark,
	}
}

func StockOutItemClone(entity *StockOutItem) *StockOutItem {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockRecordUpsertToEntity(upsert *StockRecordUpsert) *StockRecord {
	if upsert == nil {
		return nil
	}
	return &StockRecord{
		ID:          upsert.ID,
		ProductID:   upsert.ProductID,
		WarehouseID: upsert.WarehouseID,
		Count:       upsert.Count,
		TotalCount:  upsert.TotalCount,
		BizType:     upsert.BizType,
		BizID:       upsert.BizID,
		BizItemId:   upsert.BizItemId,
		BizNo:       upsert.BizNo,
	}
}

func StockRecordToUpsert(entity *StockRecord) *StockRecordUpsert {
	if entity == nil {
		return nil
	}
	return &StockRecordUpsert{
		ID:          entity.ID,
		ProductID:   entity.ProductID,
		WarehouseID: entity.WarehouseID,
		Count:       entity.Count,
		TotalCount:  entity.TotalCount,
		BizType:     entity.BizType,
		BizID:       entity.BizID,
		BizItemId:   entity.BizItemId,
		BizNo:       entity.BizNo,
	}
}

func StockRecordClone(entity *StockRecord) *StockRecord {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func StockUpsertToEntity(upsert *StockUpsert) *Stock {
	if upsert == nil {
		return nil
	}
	return &Stock{
		ID:          upsert.ID,
		ProductID:   upsert.ProductID,
		WarehouseID: upsert.WarehouseID,
		Count:       upsert.Count,
	}
}

func StockToUpsert(entity *Stock) *StockUpsert {
	if entity == nil {
		return nil
	}
	return &StockUpsert{
		ID:          entity.ID,
		ProductID:   entity.ProductID,
		WarehouseID: entity.WarehouseID,
		Count:       entity.Count,
	}
}

func StockClone(entity *Stock) *Stock {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}
