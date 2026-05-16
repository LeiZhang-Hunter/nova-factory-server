package masterserviceimpl

import (
	"errors"
	"reflect"
	"strings"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// SupplierServiceImpl 提供业务实现。
type SupplierServiceImpl struct {
	dao          masterdao.ISupplierDao
	uniqueFields []supplierUniqueField
}

// NewSupplierService 创建服务。
func NewSupplierService(dao masterdao.ISupplierDao) masterservice.ISupplierService {
	return &SupplierServiceImpl{
		dao:          dao,
		uniqueFields: []supplierUniqueField{{Field: "Code", Column: "code", Label: "供应商编码"}},
	}
}

func (s *SupplierServiceImpl) create(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	supplierTrimStringFields(req)
	if err := supplierValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, 0); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *SupplierServiceImpl) update(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	id := supplierGetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	supplierTrimStringFields(req)
	if err := supplierValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, id); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

func (s *SupplierServiceImpl) Create(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error) {
	return s.create(c, req)
}

func (s *SupplierServiceImpl) Update(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error) {
	return s.update(c, req)
}

func (s *SupplierServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *SupplierServiceImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Supplier, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *SupplierServiceImpl) ListPage(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error) {
	if req != nil {
		supplierTrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *SupplierServiceImpl) validateUniqueFields(c *gin.Context, req *mastermodels.SupplierUpsert, currentID int64) error {
	if len(s.uniqueFields) == 0 {
		return nil
	}
	for _, field := range s.uniqueFields {
		value, ok := supplierGetFieldValue(req, field.Field)
		if !ok {
			continue
		}
		normalized, empty := supplierNormalizeValue(value)
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
		if supplierGetIntField(exists, "ID") != currentID {
			label := strings.TrimSpace(field.Label)
			if label == "" {
				label = field.Column
			}
			return errors.New(label + "已存在")
		}
	}
	return nil
}

func (s *SupplierServiceImpl) List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.SupplierListData{Rows: result.Rows, Total: result.Total}, nil
}

type supplierUniqueField struct {
	Field  string
	Column string
	Label  string
}

func supplierTrimStringFields(target any) {
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
	supplierTrimStruct(value)
}

func supplierValidateRequiredFields(target any) error {
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

func supplierNormalizeValue(value reflect.Value) (any, bool) {
	if !value.IsValid() {
		return nil, true
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, true
		}
		return supplierNormalizeValue(value.Elem())
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

func supplierGetFieldValue(target any, name string) (reflect.Value, bool) {
	value := supplierFieldValue(target, name)
	return value, value.IsValid()
}

func supplierGetIntField(target any, name string) int64 {
	value := supplierFieldValue(target, name)
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

func supplierTrimStruct(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)
		if structField.PkgPath != "" {
			continue
		}
		if structField.Anonymous {
			if field.Kind() == reflect.Struct {
				supplierTrimStruct(field)
			}
			continue
		}
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}

func supplierFieldValue(target any, name string) reflect.Value {
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
