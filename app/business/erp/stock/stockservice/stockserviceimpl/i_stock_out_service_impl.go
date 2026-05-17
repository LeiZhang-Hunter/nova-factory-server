package stockserviceimpl

import (
	"errors"
	"reflect"
	"strings"

	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockOutServiceImpl 提供业务实现。
type StockOutServiceImpl struct {
	dao          stockdao.IStockOutDao
	uniqueFields []stockOutUniqueField
}

// NewStockOutService 创建服务。
func NewStockOutService(dao stockdao.IStockOutDao) stockservice.IStockOutService {
	return &StockOutServiceImpl{
		dao:          dao,
		uniqueFields: []stockOutUniqueField{{Field: "No", Column: "no", Label: "出库单号"}},
	}
}

func (s *StockOutServiceImpl) create(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	stockOutTrimStringFields(req)
	if err := stockOutValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, 0); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *StockOutServiceImpl) update(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	id := stockOutGetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	stockOutTrimStringFields(req)
	if err := stockOutValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, id); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

func (s *StockOutServiceImpl) Create(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error) {
	return s.create(c, req)
}

func (s *StockOutServiceImpl) Update(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error) {
	return s.update(c, req)
}

func (s *StockOutServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *StockOutServiceImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockOut, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *StockOutServiceImpl) ListPage(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error) {
	if req != nil {
		stockOutTrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *StockOutServiceImpl) validateUniqueFields(c *gin.Context, req *stockmodels.StockOutUpsert, currentID int64) error {
	if len(s.uniqueFields) == 0 {
		return nil
	}
	for _, field := range s.uniqueFields {
		value, ok := stockOutGetFieldValue(req, field.Field)
		if !ok {
			continue
		}
		normalized, empty := stockOutNormalizeValue(value)
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
		if stockOutGetIntField(exists, "ID") != currentID {
			label := strings.TrimSpace(field.Label)
			if label == "" {
				label = field.Column
			}
			return errors.New(label + "已存在")
		}
	}
	return nil
}

func (s *StockOutServiceImpl) List(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockOutListData{Rows: result.Rows, Total: result.Total}, nil
}

type stockOutUniqueField struct {
	Field  string
	Column string
	Label  string
}

func stockOutTrimStringFields(target any) {
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
	stockOutTrimStruct(value)
}

func stockOutValidateRequiredFields(target any) error {
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

func stockOutNormalizeValue(value reflect.Value) (any, bool) {
	if !value.IsValid() {
		return nil, true
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, true
		}
		return stockOutNormalizeValue(value.Elem())
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

func stockOutGetFieldValue(target any, name string) (reflect.Value, bool) {
	value := stockOutFieldValue(target, name)
	return value, value.IsValid()
}

func stockOutGetIntField(target any, name string) int64 {
	value := stockOutFieldValue(target, name)
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

func stockOutTrimStruct(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)
		if structField.PkgPath != "" {
			continue
		}
		if structField.Anonymous {
			if field.Kind() == reflect.Struct {
				stockOutTrimStruct(field)
			}
			continue
		}
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}

func stockOutFieldValue(target any, name string) reflect.Value {
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
