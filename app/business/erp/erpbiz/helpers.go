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
