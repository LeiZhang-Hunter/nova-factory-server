package stockserviceimpl

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/utils/observer/integration/event"

	"github.com/gin-gonic/gin"
)

// StockServiceImpl 提供业务实现。
type StockServiceImpl struct {
	dao          stockdao.IStockDao
	productDao   masterdao.IProductDao
	warehouseDao masterdao.IWarehouseDao
	uniqueFields []stockUniqueField
}

// NewStockService 创建服务。
func NewStockService(dao stockdao.IStockDao, productDao masterdao.IProductDao, warehouseDao masterdao.IWarehouseDao) stockservice.IStockService {
	return &StockServiceImpl{
		dao:          dao,
		productDao:   productDao,
		warehouseDao: warehouseDao,
		uniqueFields: []stockUniqueField{},
	}
}

func (s *StockServiceImpl) create(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	stockTrimStringFields(req)
	if err := stockValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, 0); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *StockServiceImpl) update(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	id := stockGetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	stockTrimStringFields(req)
	if err := stockValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, id); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

func (s *StockServiceImpl) Create(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error) {
	return s.create(c, req)
}

func (s *StockServiceImpl) Update(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error) {
	return s.update(c, req)
}

func (s *StockServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *StockServiceImpl) GetByID(c *gin.Context, id int64) (*stockmodels.Stock, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	info, err := s.dao.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	productNameMap, warehouseNameMap, err := s.getStockRelatedNameMaps(c, []*stockmodels.Stock{info})
	if err != nil {
		return nil, err
	}
	info.ProductName = productNameMap[info.ProductID]
	info.WarehouseName = warehouseNameMap[info.WarehouseID]
	return info, nil
}

func (s *StockServiceImpl) ListPage(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error) {
	if req != nil {
		stockTrimStringFields(req)
	}
	result, err := s.dao.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	if err := s.fillStockNames(c, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *StockServiceImpl) validateUniqueFields(c *gin.Context, req *stockmodels.StockUpsert, currentID int64) error {
	if len(s.uniqueFields) == 0 {
		return nil
	}
	for _, field := range s.uniqueFields {
		value, ok := stockGetFieldValue(req, field.Field)
		if !ok {
			continue
		}
		normalized, empty := stockNormalizeValue(value)
		if empty {
			continue
		}
		exists, err := s.dao.GetByColumn(c, field.Column, normalized)
		if err != nil {
			return err
		}
		if exists == nil {
			continue
		}
		if stockGetIntField(exists, "ID") != currentID {
			label := strings.TrimSpace(field.Label)
			if label == "" {
				label = field.Column
			}
			return errors.New(label + "已存在")
		}
	}
	return nil
}

func (s *StockServiceImpl) List(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockListData{Rows: result.Rows, Total: result.Total}, nil
}

func (s *StockServiceImpl) fillStockNames(c *gin.Context, result *stockmodels.StockListData) error {
	if result == nil || len(result.Rows) == 0 {
		return nil
	}

	productNameMap, warehouseNameMap, err := s.getStockRelatedNameMaps(c, result.Rows)
	if err != nil {
		return err
	}
	for _, item := range result.Rows {
		if item == nil {
			continue
		}
		item.ProductName = productNameMap[item.ProductID]
		item.WarehouseName = warehouseNameMap[item.WarehouseID]
	}
	return nil
}

func (s *StockServiceImpl) getStockRelatedNameMaps(c *gin.Context, rows []*stockmodels.Stock) (map[int64]string, map[int64]string, error) {
	productIDs := make([]int64, 0)
	warehouseIDs := make([]int64, 0)
	productIDSet := make(map[int64]struct{})
	warehouseIDSet := make(map[int64]struct{})
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.ProductID > 0 {
			if _, exists := productIDSet[row.ProductID]; !exists {
				productIDSet[row.ProductID] = struct{}{}
				productIDs = append(productIDs, row.ProductID)
			}
		}
		if row.WarehouseID > 0 {
			if _, exists := warehouseIDSet[row.WarehouseID]; !exists {
				warehouseIDSet[row.WarehouseID] = struct{}{}
				warehouseIDs = append(warehouseIDs, row.WarehouseID)
			}
		}
	}

	productNameMap, err := s.batchGetProductNames(c, productIDs)
	if err != nil {
		return nil, nil, err
	}
	warehouseNameMap, err := s.batchGetWarehouseNames(c, warehouseIDs)
	if err != nil {
		return nil, nil, err
	}
	return productNameMap, warehouseNameMap, nil
}

func (s *StockServiceImpl) batchGetProductNames(c *gin.Context, ids []int64) (map[int64]string, error) {
	nameMap := make(map[int64]string, len(ids))
	if len(ids) == 0 {
		return nameMap, nil
	}

	rows, err := s.productDao.GetByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		nameMap[row.ID] = row.Name
	}
	return nameMap, nil
}

func (s *StockServiceImpl) batchGetWarehouseNames(c *gin.Context, ids []int64) (map[int64]string, error) {
	nameMap := make(map[int64]string, len(ids))
	if len(ids) == 0 {
		return nameMap, nil
	}

	rows, err := s.warehouseDao.GetByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		nameMap[row.ID] = row.Name
	}
	return nameMap, nil
}

type stockUniqueField struct {
	Field  string
	Column string
	Label  string
}

func stockTrimStringFields(target any) {
	if target == nil {
		return
	}
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return
	}
	stockTrimStruct(value)
}

func stockValidateRequiredFields(target any) error {
	if target == nil {
		return nil
	}
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := valueType.Field(i)
		if structField.PkgPath != "" || structField.Anonymous {
			continue
		}
		if !strings.Contains(structField.Tag.Get("binding"), "required") {
			continue
		}
		label := structField.Tag.Get("label")
		if label == "" {
			label = structField.Name
		}
		switch field.Kind() {
		case reflect.String:
			if strings.TrimSpace(field.String()) == "" {
				return errors.New(label + "不能为空")
			}
		case reflect.Pointer:
			if field.IsNil() {
				return errors.New(label + "不能为空")
			}
		default:
			if field.IsZero() {
				return errors.New(label + "不能为空")
			}
		}
	}
	return nil
}

func stockNormalizeValue(value reflect.Value) (any, bool) {
	if !value.IsValid() {
		return nil, true
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, true
		}
		return stockNormalizeValue(value.Elem())
	}
	switch value.Kind() {
	case reflect.String:
		trimmed := strings.TrimSpace(value.String())
		return trimmed, trimmed == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		current := value.Int()
		return current, current == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		current := value.Uint()
		return current, current == 0
	case reflect.Bool:
		return value.Bool(), false
	default:
		if value.IsZero() {
			return nil, true
		}
		return value.Interface(), false
	}
}

func stockGetFieldValue(target any, name string) (reflect.Value, bool) {
	value := stockFieldValue(target, name)
	return value, value.IsValid()
}

func stockGetIntField(target any, name string) int64 {
	value := stockFieldValue(target, name)
	if !value.IsValid() {
		return 0
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(value.Uint())
	}
	return 0
}

func stockTrimStruct(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)
		if structField.PkgPath != "" {
			continue
		}
		if structField.Anonymous {
			if field.Kind() == reflect.Struct {
				stockTrimStruct(field)
			}
			continue
		}
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}

func stockFieldValue(target any, name string) reflect.Value {
	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return reflect.Value{}
	}
	return value.FieldByName(name)
}

func (s *StockServiceImpl) SyncStock(stock event.StockEvent) error {
	if stock.GetDB() == nil {
		return errors.New("库存同步需要事务DB")
	}
	if len(stock.GetStocks()) == 0 {
		return nil
	}

	for _, e := range stock.GetStocks() {
		productID := e.SkuID()
		afterQty := e.AfterQty()
		if err := s.dao.UpdateStockByProductIDWithDB(stock.GetDB(), productID, afterQty); err != nil {
			return fmt.Errorf("更新ERP库存失败 productId=%d: %w", productID, err)
		}
	}

	return nil
}
