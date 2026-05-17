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
