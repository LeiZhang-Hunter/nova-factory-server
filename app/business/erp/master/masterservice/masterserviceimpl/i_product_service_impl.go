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

// ProductServiceImpl 提供业务实现。
type ProductServiceImpl struct {
	dao          masterdao.IProductDao
	categoryDao  masterdao.IProductCategoryDao
	uniqueFields []productUniqueField
}

// NewProductService 创建服务。
func NewProductService(dao masterdao.IProductDao, categoryDao masterdao.IProductCategoryDao) masterservice.IProductService {
	return &ProductServiceImpl{
		dao:          dao,
		categoryDao:  categoryDao,
		uniqueFields: []productUniqueField{},
	}
}

func (s *ProductServiceImpl) create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	productTrimStringFields(req)
	if err := productValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, 0); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *ProductServiceImpl) update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	id := productGetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	productTrimStringFields(req)
	if err := productValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, id); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

func (s *ProductServiceImpl) Create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error) {
	return s.create(c, req)
}

func (s *ProductServiceImpl) Update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error) {
	return s.update(c, req)
}

func (s *ProductServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ProductServiceImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Product, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *ProductServiceImpl) ListPage(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error) {
	if req != nil {
		productTrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *ProductServiceImpl) validateUniqueFields(c *gin.Context, req *mastermodels.ProductUpsert, currentID int64) error {
	if len(s.uniqueFields) == 0 {
		return nil
	}
	for _, field := range s.uniqueFields {
		value, ok := productGetFieldValue(req, field.Field)
		if !ok {
			continue
		}
		normalized, empty := productNormalizeValue(value)
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
		if productGetIntField(exists, "ID") != currentID {
			label := strings.TrimSpace(field.Label)
			if label == "" {
				label = field.Column
			}
			return errors.New(label + "已存在")
		}
	}
	return nil
}

func (s *ProductServiceImpl) List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	if len(result.Rows) == 0 {
		return result, nil
	}
	categoryIDSet := make(map[int64]struct{})
	for _, row := range result.Rows {
		if row.CategoryId > 0 {
			categoryIDSet[row.CategoryId] = struct{}{}
		}
	}
	if len(categoryIDSet) == 0 {
		return result, nil
	}
	ids := make([]int64, 0, len(categoryIDSet))
	for id := range categoryIDSet {
		ids = append(ids, id)
	}
	categories, err := s.categoryDao.GetByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	idToName := make(map[int64]string)
	for _, cat := range categories {
		idToName[cat.ID] = cat.Name
	}
	for _, row := range result.Rows {
		if row.CategoryId > 0 {
			if name, ok := idToName[row.CategoryId]; ok {
				row.CategoryName = name
			}
		}
	}
	return result, nil
}

type productUniqueField struct {
	Field  string
	Column string
	Label  string
}

func productTrimStringFields(target any) {
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
	productTrimStruct(value)
}

func productValidateRequiredFields(target any) error {
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

func productNormalizeValue(value reflect.Value) (any, bool) {
	if !value.IsValid() {
		return nil, true
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, true
		}
		return productNormalizeValue(value.Elem())
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

func productGetFieldValue(target any, name string) (reflect.Value, bool) {
	value := productFieldValue(target, name)
	return value, value.IsValid()
}

func productGetIntField(target any, name string) int64 {
	value := productFieldValue(target, name)
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

func productTrimStruct(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)
		if structField.PkgPath != "" {
			continue
		}
		if structField.Anonymous {
			if field.Kind() == reflect.Struct {
				productTrimStruct(field)
			}
			continue
		}
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}

func productFieldValue(target any, name string) reflect.Value {
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
