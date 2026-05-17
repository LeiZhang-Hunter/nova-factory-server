package masterserviceimpl

import (
	"errors"
	"reflect"
	"strings"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/utils/category"

	"github.com/gin-gonic/gin"
)

// ProductCategoryServiceImpl 提供业务实现。
type ProductCategoryServiceImpl struct {
	dao          masterdao.IProductCategoryDao
	uniqueFields []productCategoryUniqueField
}

// NewProductCategoryService 创建服务。
func NewProductCategoryService(dao masterdao.IProductCategoryDao) masterservice.IProductCategoryService {
	return &ProductCategoryServiceImpl{
		dao:          dao,
		uniqueFields: []productCategoryUniqueField{},
	}
}

func (s *ProductCategoryServiceImpl) create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	productCategoryTrimStringFields(req)
	if err := productCategoryValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, 0); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *ProductCategoryServiceImpl) update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	id := productCategoryGetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	productCategoryTrimStringFields(req)
	if err := productCategoryValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, id); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

// Create 新增 ERP 产品分类。
func (s *ProductCategoryServiceImpl) Create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	s.ensureCode(req)
	return s.create(c, req)
}

// Update 修改 ERP 产品分类。
func (s *ProductCategoryServiceImpl) Update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	s.ensureCode(req)
	return s.update(c, req)
}

func (s *ProductCategoryServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ProductCategoryServiceImpl) GetByID(c *gin.Context, id int64) (*mastermodels.ProductCategory, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	item, err := s.dao.GetByID(c, id)
	if err != nil || item == nil {
		return item, err
	}
	if item.ParentID > 0 {
		if parent, err := s.dao.GetByColumn(c, "id", item.ParentID); err == nil && parent != nil {
			item.ParentName = parent.Name
		}
	}
	return item, nil
}

func (s *ProductCategoryServiceImpl) ListPage(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error) {
	if req != nil {
		productCategoryTrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *ProductCategoryServiceImpl) validateUniqueFields(c *gin.Context, req *mastermodels.ProductCategoryUpsert, currentID int64) error {
	if len(s.uniqueFields) == 0 {
		return nil
	}
	for _, field := range s.uniqueFields {
		value, ok := productCategoryGetFieldValue(req, field.Field)
		if !ok {
			continue
		}
		normalized, empty := productCategoryNormalizeValue(value)
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
		if productCategoryGetIntField(exists, "ID") != currentID {
			label := strings.TrimSpace(field.Label)
			if label == "" {
				label = field.Column
			}
			return errors.New(label + "已存在")
		}
	}
	return nil
}

// List 查询 ERP 产品分类列表。
func (s *ProductCategoryServiceImpl) List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	if len(result.Rows) == 0 {
		return result, nil
	}
	parentIDSet := make(map[int64]struct{})
	for _, row := range result.Rows {
		if row.ParentID > 0 {
			parentIDSet[row.ParentID] = struct{}{}
		}
	}
	if len(parentIDSet) == 0 {
		return result, nil
	}
	ids := make([]int64, 0, len(parentIDSet))
	for id := range parentIDSet {
		ids = append(ids, id)
	}
	parents, err := s.dao.GetByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	idToName := make(map[int64]string)
	for _, parent := range parents {
		idToName[parent.ID] = parent.Name
	}
	for _, row := range result.Rows {
		if row.ParentID > 0 {
			if name, ok := idToName[row.ParentID]; ok {
				row.ParentName = name
			}
		}
	}
	return result, nil
}

func (s *ProductCategoryServiceImpl) ensureCode(req *mastermodels.ProductCategoryUpsert) {
	if req == nil {
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	if req.Code == "" {
		req.Code = category.GenerateProductCategoryCode()
	}
}

type productCategoryUniqueField struct {
	Field  string
	Column string
	Label  string
}

func productCategoryTrimStringFields(target any) {
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
	productCategoryTrimStruct(value)
}

func productCategoryValidateRequiredFields(target any) error {
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

func productCategoryNormalizeValue(value reflect.Value) (any, bool) {
	if !value.IsValid() {
		return nil, true
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, true
		}
		return productCategoryNormalizeValue(value.Elem())
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

func productCategoryGetFieldValue(target any, name string) (reflect.Value, bool) {
	value := productCategoryFieldValue(target, name)
	return value, value.IsValid()
}

func productCategoryGetIntField(target any, name string) int64 {
	value := productCategoryFieldValue(target, name)
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

func productCategoryTrimStruct(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)
		if structField.PkgPath != "" {
			continue
		}
		if structField.Anonymous {
			if field.Kind() == reflect.Struct {
				productCategoryTrimStruct(field)
			}
			continue
		}
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}

func productCategoryFieldValue(target any, name string) reflect.Value {
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
