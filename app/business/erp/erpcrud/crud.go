package erpcrud

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"nova-factory-server/app/baize"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UniqueField struct {
	Field  string
	Column string
	Label  string
}

type EntityConfig struct {
	Table        string
	OrderBy      string
	UniqueFields []UniqueField
}

type PageResult[T any] struct {
	Rows  []*T  `json:"rows"`
	Total int64 `json:"total"`
}

type CRUDStore[T any, U any, Q any] interface {
	Create(c *gin.Context, req *U) (*T, error)
	Update(c *gin.Context, req *U) (*T, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*T, error)
	GetByColumn(c *gin.Context, column string, value any) (*T, error)
	ListPage(c *gin.Context, req *Q) (*PageResult[T], error)
}

type CRUDDao[T any, U any, Q any] struct {
	db     *gorm.DB
	config EntityConfig
}

func NewCRUDDao[T any, U any, Q any](db *gorm.DB, config EntityConfig) *CRUDDao[T, U, Q] {
	return &CRUDDao[T, U, Q]{db: db, config: config}
}

func (d *CRUDDao[T, U, Q]) Create(c *gin.Context, req *U) (*T, error) {
	model := new(T)
	if err := copyStruct(model, req); err != nil {
		return nil, err
	}
	ensureID(model)
	ensureDeptID(model, baizeContext.GetDeptId(c))
	ensureState(model)
	setCreateAudit(model, baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table(d.config.Table).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *CRUDDao[T, U, Q]) Update(c *gin.Context, req *U) (*T, error) {
	id := getIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(T)
	if err := copyStruct(model, req); err != nil {
		return nil, err
	}
	setUpdateAudit(model, baizeContext.GetUserId(c))
	updates := buildUpdateMap(model)
	db := d.db.WithContext(c).Table(d.config.Table).Where("id = ?", id)
	if hasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *CRUDDao[T, U, Q]) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table(d.config.Table).Where("id IN ?", ids)
	if hasField(new(T), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *CRUDDao[T, U, Q]) GetByID(c *gin.Context, id int64) (*T, error) {
	item := new(T)
	db := d.db.WithContext(c).Table(d.config.Table).Where("id = ?", id)
	if hasField(item, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *CRUDDao[T, U, Q]) GetByColumn(c *gin.Context, column string, value any) (*T, error) {
	if column == "" {
		return nil, nil
	}
	item := new(T)
	db := d.db.WithContext(c).Table(d.config.Table).Where(fmt.Sprintf("%s = ?", column), value)
	if hasField(item, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *CRUDDao[T, U, Q]) ListPage(c *gin.Context, req *Q) (*PageResult[T], error) {
	if req == nil {
		req = new(Q)
	}
	db := d.db.WithContext(c).Table(d.config.Table)
	if hasField(new(T), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	applyFilters(db, req)
	db = applyFilters(db, req)
	page, size := getPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]T, 0)
	orderBy := d.config.OrderBy
	if strings.TrimSpace(orderBy) == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*T, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &PageResult[T]{Rows: result, Total: total}, nil
}

type CRUDService[T any, U any, Q any] struct {
	dao    CRUDStore[T, U, Q]
	config EntityConfig
}

func NewCRUDService[T any, U any, Q any](dao CRUDStore[T, U, Q], config EntityConfig) *CRUDService[T, U, Q] {
	return &CRUDService[T, U, Q]{dao: dao, config: config}
}

func (s *CRUDService[T, U, Q]) Create(c *gin.Context, req *U) (*T, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	trimStringFields(req)
	if err := validateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, 0); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *CRUDService[T, U, Q]) Update(c *gin.Context, req *U) (*T, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	id := getIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	trimStringFields(req)
	if err := validateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, id); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

func (s *CRUDService[T, U, Q]) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *CRUDService[T, U, Q]) GetByID(c *gin.Context, id int64) (*T, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *CRUDService[T, U, Q]) ListPage(c *gin.Context, req *Q) (*PageResult[T], error) {
	if req != nil {
		trimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *CRUDService[T, U, Q]) validateUniqueFields(c *gin.Context, req *U, currentID int64) error {
	if len(s.config.UniqueFields) == 0 {
		return nil
	}
	for _, field := range s.config.UniqueFields {
		value, ok := getFieldValue(req, field.Field)
		if !ok {
			continue
		}
		normalized, empty := normalizeValue(value)
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
		if getIntField(exists, "ID") != currentID {
			label := field.Label
			if strings.TrimSpace(label) == "" {
				label = field.Column
			}
			return errors.New(label + "已存在")
		}
	}
	return nil
}

func trimStringFields(target any) {
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
	trimStruct(value)
}

func trimStruct(value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)
		if structField.PkgPath != "" {
			continue
		}
		if structField.Anonymous {
			if field.Kind() == reflect.Struct {
				trimStruct(field)
			}
			continue
		}
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}

func validateRequiredFields(target any) error {
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

func copyStruct(dst any, src any) error {
	if dst == nil || src == nil {
		return nil
	}
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Pointer || dstValue.IsNil() {
		return errors.New("目标对象必须为指针")
	}
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() == reflect.Pointer {
		if srcValue.IsNil() {
			return nil
		}
		srcValue = srcValue.Elem()
	}
	if srcValue.Kind() != reflect.Struct {
		return nil
	}
	dstValue = dstValue.Elem()
	dstType := dstValue.Type()
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		srcFieldType := srcValue.Type().Field(i)
		if srcFieldType.PkgPath != "" || srcFieldType.Anonymous {
			continue
		}
		dstField := dstValue.FieldByName(srcFieldType.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}
		if srcField.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcField)
			continue
		}
		if dstField.Type() == reflect.TypeOf((*time.Time)(nil)) && srcField.Kind() == reflect.String {
			parsed, err := parseTimeString(strings.TrimSpace(srcField.String()), srcFieldType.Name)
			if err != nil {
				return err
			}
			if parsed == nil {
				dstField.Set(reflect.Zero(dstField.Type()))
			} else {
				dstField.Set(reflect.ValueOf(parsed))
			}
			continue
		}
		if dstField.Type().Kind() == reflect.Int64 && srcField.Kind() == reflect.Int64 {
			dstField.SetInt(srcField.Int())
			continue
		}
		_ = dstType
	}
	return nil
}

func parseTimeString(value string, fieldName string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	layouts := []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &parsed, nil
		}
	}
	return nil, errors.New(fieldName + "时间格式错误，要求: 2006-01-02 15:04:05")
}

func ensureID(target any) {
	value := fieldValue(target, "ID")
	if !value.IsValid() || !value.CanSet() || !value.IsZero() {
		return
	}
	switch value.Kind() {
	case reflect.Int64:
		value.SetInt(snowflake.GenID())
	case reflect.Uint64:
		value.SetUint(uint64(snowflake.GenID()))
	}
}

func ensureDeptID(target any, deptID int64) {
	value := fieldValue(target, "DeptID")
	if !value.IsValid() || !value.CanSet() || !value.IsZero() {
		return
	}
	if value.Kind() == reflect.Int64 {
		value.SetInt(deptID)
	}
}

func ensureState(target any) {
	value := fieldValue(target, "State")
	if !value.IsValid() || !value.CanSet() {
		return
	}
	if value.Kind() == reflect.Int32 && value.IsZero() {
		value.SetInt(commonStatus.NORMAL)
	}
}

func setCreateAudit(target any, userID int64) {
	base := fieldValue(target, "BaseEntity")
	if !base.IsValid() {
		return
	}
	if entity, ok := base.Addr().Interface().(*baize.BaseEntity); ok {
		entity.SetCreateBy(userID)
	}
}

func setUpdateAudit(target any, userID int64) {
	base := fieldValue(target, "BaseEntity")
	if !base.IsValid() {
		return
	}
	if entity, ok := base.Addr().Interface().(*baize.BaseEntity); ok {
		entity.SetUpdateBy(userID)
	}
}

func buildUpdateMap(target any) map[string]any {
	result := make(map[string]any)
	excluded := map[string]struct{}{
		"id":          {},
		"dept_id":     {},
		"state":       {},
		"create_by":   {},
		"create_time": {},
	}
	collectFields(reflect.Indirect(reflect.ValueOf(target)), result, excluded)
	return result
}

func collectFields(value reflect.Value, result map[string]any, excluded map[string]struct{}) {
	if !value.IsValid() || value.Kind() != reflect.Struct {
		return
	}
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := valueType.Field(i)
		if structField.PkgPath != "" {
			continue
		}
		if structField.Anonymous {
			collectFields(field, result, excluded)
			continue
		}
		column := columnName(structField)
		if _, ok := excluded[column]; ok {
			continue
		}
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				result[column] = nil
			} else {
				result[column] = field.Interface()
			}
			continue
		}
		result[column] = field.Interface()
	}
}

func applyFilters(db *gorm.DB, req any) *gorm.DB {
	value := reflect.ValueOf(req)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return db
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return db
	}
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := valueType.Field(i)
		filterTag := structField.Tag.Get("filter")
		if filterTag == "" {
			continue
		}
		parts := strings.Split(filterTag, ",")
		if len(parts) != 2 {
			continue
		}
		normalized, empty := normalizeValue(field)
		if empty {
			continue
		}
		switch parts[0] {
		case "like":
			db = db.Where(fmt.Sprintf("%s LIKE ?", parts[1]), "%"+fmt.Sprint(normalized)+"%")
		case "eq":
			db = db.Where(fmt.Sprintf("%s = ?", parts[1]), normalized)
		}
	}
	return db
}

func getPageSize(req any) (int64, int64) {
	page := getIntField(req, "Page")
	size := getIntField(req, "Size")
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 200 {
		size = 200
	}
	return page, size
}

func normalizeValue(value reflect.Value) (any, bool) {
	if !value.IsValid() {
		return nil, true
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, true
		}
		return normalizeValue(value.Elem())
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

func hasField(target any, name string) bool {
	return fieldValue(target, name).IsValid()
}

func getFieldValue(target any, name string) (reflect.Value, bool) {
	value := fieldValue(target, name)
	return value, value.IsValid()
}

func getIntField(target any, name string) int64 {
	value := fieldValue(target, name)
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

func fieldValue(target any, name string) reflect.Value {
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

func columnName(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")
	for _, item := range strings.Split(gormTag, ";") {
		item = strings.TrimSpace(item)
		if strings.HasPrefix(item, "column:") {
			return strings.TrimPrefix(item, "column:")
		}
	}
	dbTag := field.Tag.Get("db")
	if dbTag != "" {
		return dbTag
	}
	return toSnake(field.Name)
}

func toSnake(value string) string {
	if value == "" {
		return value
	}
	var builder strings.Builder
	for index, ch := range value {
		if ch >= 'A' && ch <= 'Z' {
			if index > 0 {
				builder.WriteByte('_')
			}
			builder.WriteRune(ch - 'A' + 'a')
			continue
		}
		builder.WriteRune(ch)
	}
	return builder.String()
}
