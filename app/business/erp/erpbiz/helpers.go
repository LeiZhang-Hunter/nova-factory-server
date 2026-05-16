package erpbiz

import (
	"errors"
	"fmt"
	"math"
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

const (
	AuditStatusProcess int32 = 10
	AuditStatusApprove int32 = 20

	PrefixStockIn        = "QTRK"
	PrefixStockOut       = "QCKD"
	PrefixStockMove      = "QCDB"
	PrefixStockCheck     = "QCPD"
	PrefixSaleOut        = "XSCK"
	PrefixSaleReturn     = "XSTH"
	PrefixPurchaseOrder  = "CGDD"
	PrefixPurchaseIn     = "CGRK"
	PrefixPurchaseReturn = "CGTH"
	PrefixFinancePayment = "FKD"
	PrefixFinanceReceipt = "SKD"
)

type UniqueField struct {
	Field  string
	Column string
	Label  string
}

type PageResult[T any] struct {
	Rows  []*T  `json:"rows"`
	Total int64 `json:"total"`
}

func GenerateNo(prefix string) string {
	suffix := snowflake.GenID()
	if suffix < 0 {
		suffix = -suffix
	}
	return strings.TrimSpace(prefix) + time.Now().Format("20060102") + fmt.Sprintf("%06d", suffix%1000000)
}

func ParseTime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	layouts := []string{
		time.DateTime,
		"2006-01-02 15:04",
		"2006-01-02",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &parsed, nil
		}
	}
	return nil, errors.New("时间格式不正确")
}

func RoundAmount(value float64) float64 {
	return math.Round(value*100) / 100
}

func CalculatePercentAmount(total, percent float64) float64 {
	if total == 0 || percent == 0 {
		return 0
	}
	return RoundAmount(total * percent / 100)
}

func PrepareCreate(model any, c *gin.Context) {
	setIntField(model, "ID", snowflake.GenID())
	setIntField(model, "DeptID", baizeContext.GetDeptId(c))
	setInt32Field(model, "State", commonStatus.NORMAL)
	setBaseCreate(model, baizeContext.GetUserId(c))
}

func PrepareUpdate(model any, c *gin.Context) error {
	if getInt64Field(model, "ID") <= 0 {
		return errors.New("id不能为空")
	}
	setBaseUpdate(model, baizeContext.GetUserId(c))
	return nil
}

func ReplaceChildren(tx *gorm.DB, c *gin.Context, table string, parentColumn string, parentID int64, children any) error {
	if parentID <= 0 {
		return errors.New("主表ID不能为空")
	}
	now := time.Now()
	if err := tx.WithContext(c).
		Table(table).
		Where(parentColumn+" = ? AND state = ?", parentID, commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": &now,
		}).Error; err != nil {
		return err
	}
	rv := reflect.ValueOf(children)
	if rv.Kind() != reflect.Slice || rv.Len() == 0 {
		return nil
	}
	return tx.WithContext(c).Table(table).Create(children).Error
}

func TrimStringFields(target any) {
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

func ValidateRequiredFields(target any) error {
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

func CopyStruct(dst any, src any) error {
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
	}
	return nil
}

func BuildUpdateMap(target any) map[string]any {
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

func ApplyFilters(db *gorm.DB, req any) *gorm.DB {
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
		normalized, empty := NormalizeValue(field)
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

func GetPageSize(req any) (int64, int64) {
	page := GetIntField(req, "Page")
	size := GetIntField(req, "Size")
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

func NormalizeValue(value reflect.Value) (any, bool) {
	if !value.IsValid() {
		return nil, true
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, true
		}
		return NormalizeValue(value.Elem())
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

func HasField(target any, name string) bool {
	return fieldValue(target, name).IsValid()
}

func GetFieldValue(target any, name string) (reflect.Value, bool) {
	value := fieldValue(target, name)
	return value, value.IsValid()
}

func GetIntField(target any, name string) int64 {
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

func setBaseCreate(model any, userID int64) {
	field := getBaseEntityField(model)
	if !field.IsValid() || !field.CanAddr() {
		return
	}
	base, ok := field.Addr().Interface().(*baize.BaseEntity)
	if !ok {
		return
	}
	base.SetCreateBy(userID)
}

func setBaseUpdate(model any, userID int64) {
	field := getBaseEntityField(model)
	if !field.IsValid() || !field.CanAddr() {
		return
	}
	base, ok := field.Addr().Interface().(*baize.BaseEntity)
	if !ok {
		return
	}
	base.SetUpdateBy(userID)
}

func getBaseEntityField(model any) reflect.Value {
	v := reflect.Indirect(reflect.ValueOf(model))
	if !v.IsValid() {
		return reflect.Value{}
	}
	return v.FieldByName("BaseEntity")
}

func setIntField(model any, name string, value int64) {
	field := reflect.Indirect(reflect.ValueOf(model)).FieldByName(name)
	if !field.IsValid() || !field.CanSet() {
		return
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() == 0 {
			field.SetInt(value)
		}
	}
}

func setInt32Field(model any, name string, value int32) {
	field := reflect.Indirect(reflect.ValueOf(model)).FieldByName(name)
	if !field.IsValid() || !field.CanSet() {
		return
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() == 0 {
			field.SetInt(int64(value))
		}
	}
}

func getInt64Field(model any, name string) int64 {
	field := reflect.Indirect(reflect.ValueOf(model)).FieldByName(name)
	if !field.IsValid() {
		return 0
	}
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int()
	}
	return 0
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
