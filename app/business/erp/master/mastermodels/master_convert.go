package mastermodels

import "nova-factory-server/app/utils/code"

func ProductUpsertToEntity(upsert *ProductUpsert) *Product {
	if upsert == nil {
		return nil
	}
	return &Product{
		ID:            upsert.ID,
		Name:          upsert.Name,
		BarCode:       upsert.BarCode,
		CategoryId:    upsert.CategoryId,
		UnitId:        upsert.UnitId,
		Status:        upsert.Status,
		Standard:      upsert.Standard,
		Remark:        upsert.Remark,
		ExpiryDay:     upsert.ExpiryDay,
		Weight:        upsert.Weight,
		PurchasePrice: upsert.PurchasePrice,
		SalePrice:     upsert.SalePrice,
		MinPrice:      upsert.MinPrice,
	}
}

func ProductToUpsert(entity *Product) *ProductUpsert {
	if entity == nil {
		return nil
	}
	return &ProductUpsert{
		ID:            entity.ID,
		Name:          entity.Name,
		BarCode:       entity.BarCode,
		CategoryId:    entity.CategoryId,
		UnitId:        entity.UnitId,
		Status:        entity.Status,
		Standard:      entity.Standard,
		Remark:        entity.Remark,
		ExpiryDay:     entity.ExpiryDay,
		Weight:        entity.Weight,
		PurchasePrice: entity.PurchasePrice,
		SalePrice:     entity.SalePrice,
		MinPrice:      entity.MinPrice,
	}
}

func ProductClone(entity *Product) *Product {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func ProductCategoryUpsertToEntity(upsert *ProductCategoryUpsert) *ProductCategory {
	if upsert == nil {
		return nil
	}
	return &ProductCategory{
		ID:       upsert.ID,
		ParentID: upsert.ParentID,
		Name:     upsert.Name,
		Code:     upsert.Code,
		Sort:     upsert.Sort,
		Status:   upsert.Status,
	}
}

func ProductCategoryToUpsert(entity *ProductCategory) *ProductCategoryUpsert {
	if entity == nil {
		return nil
	}
	return &ProductCategoryUpsert{
		ID:       entity.ID,
		ParentID: entity.ParentID,
		Name:     entity.Name,
		Code:     entity.Code,
		Sort:     entity.Sort,
		Status:   entity.Status,
	}
}

func ProductCategoryClone(entity *ProductCategory) *ProductCategory {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func ProductUnitUpsertToEntity(upsert *ProductUnitUpsert) *ProductUnit {
	if upsert == nil {
		return nil
	}
	return &ProductUnit{
		ID:     upsert.ID,
		Name:   upsert.Name,
		Status: upsert.Status,
	}
}

func ProductUnitToUpsert(entity *ProductUnit) *ProductUnitUpsert {
	if entity == nil {
		return nil
	}
	return &ProductUnitUpsert{
		ID:     entity.ID,
		Name:   entity.Name,
		Status: entity.Status,
	}
}

func ProductUnitClone(entity *ProductUnit) *ProductUnit {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func CustomerUpsertToEntity(upsert *CustomerUpsert) *Customer {
	if upsert == nil {
		return nil
	}
	return &Customer{
		ID:          upsert.ID,
		Name:        upsert.Name,
		Code:        upsert.Code,
		Contact:     upsert.Contact,
		Mobile:      upsert.Mobile,
		Telephone:   upsert.Telephone,
		Email:       upsert.Email,
		Fax:         upsert.Fax,
		Remark:      upsert.Remark,
		Status:      upsert.Status,
		Sort:        upsert.Sort,
		TaxNo:       upsert.TaxNo,
		TaxPercent:  upsert.TaxPercent,
		BankName:    upsert.BankName,
		BankAccount: upsert.BankAccount,
		BankAddress: upsert.BankAddress,
	}
}

func CustomerToUpsert(entity *Customer) *CustomerUpsert {
	if entity == nil {
		return nil
	}
	return &CustomerUpsert{
		ID:          entity.ID,
		Name:        entity.Name,
		Code:        entity.Code,
		Contact:     entity.Contact,
		Mobile:      entity.Mobile,
		Telephone:   entity.Telephone,
		Email:       entity.Email,
		Fax:         entity.Fax,
		Remark:      entity.Remark,
		Status:      entity.Status,
		Sort:        entity.Sort,
		TaxNo:       entity.TaxNo,
		TaxPercent:  entity.TaxPercent,
		BankName:    entity.BankName,
		BankAccount: entity.BankAccount,
		BankAddress: entity.BankAddress,
	}
}

func CustomerClone(entity *Customer) *Customer {
	if entity == nil {
		return nil
	}
	clone := *entity
	return &clone
}

func SupplierUpsertToEntity(upsert *SupplierUpsert) *Supplier {
	if upsert == nil {
		return nil
	}
	return &Supplier{
		ID:          upsert.ID,
		Name:        upsert.Name,
		Code:        upsert.Code,
		Contact:     upsert.Contact,
		Mobile:      upsert.Mobile,
		Telephone:   upsert.Telephone,
		Email:       upsert.Email,
		Fax:         upsert.Fax,
		Remark:      upsert.Remark,
		Status:      upsert.Status,
		Sort:        upsert.Sort,
		TaxNo:       upsert.TaxNo,
		TaxPercent:  upsert.TaxPercent,
		BankName:    upsert.BankName,
		BankAccount: upsert.BankAccount,
		BankAddress: upsert.BankAddress,
	}
}

func SupplierToUpsert(entity *Supplier) *SupplierUpsert {
	if entity == nil {
		return nil
	}
	return &SupplierUpsert{
		ID:          entity.ID,
		Name:        entity.Name,
		Code:        entity.Code,
		Contact:     entity.Contact,
		Mobile:      entity.Mobile,
		Telephone:   entity.Telephone,
		Email:       entity.Email,
		Fax:         entity.Fax,
		Remark:      entity.Remark,
		Status:      entity.Status,
		Sort:        entity.Sort,
		TaxNo:       entity.TaxNo,
		TaxPercent:  entity.TaxPercent,
		BankName:    entity.BankName,
		BankAccount: entity.BankAccount,
		BankAddress: entity.BankAddress,
	}
}

func WarehouseUpsertToEntity(upsert *WarehouseUpsert) *Warehouse {
	if upsert == nil {
		return nil
	}
	w := &Warehouse{
		ID:             upsert.ID,
		Name:           upsert.Name,
		Code:           upsert.Code,
		Address:        upsert.Address,
		Sort:           upsert.Sort,
		Remark:         upsert.Remark,
		Principal:      upsert.Principal,
		WarehousePrice: upsert.WarehousePrice,
		TruckagePrice:  upsert.TruckagePrice,
		Status:         upsert.Status,
		DefaultStatus:  upsert.DefaultStatus,
	}
	if w.Code == "" {
		w.Code = code.GenerateWarehouseCode()
	}
	if w.DefaultStatus == nil {
		var defaultStatus bool = false
		w.DefaultStatus = &defaultStatus
	}
	return w
}

func AccountUpsertToEntity(upsert *AccountUpsert) *Account {
	if upsert == nil {
		return nil
	}
	return &Account{
		ID:            upsert.ID,
		Name:          upsert.Name,
		No:            upsert.No,
		Remark:        upsert.Remark,
		Status:        upsert.Status,
		Sort:          upsert.Sort,
		DefaultStatus: upsert.DefaultStatus,
	}
}
