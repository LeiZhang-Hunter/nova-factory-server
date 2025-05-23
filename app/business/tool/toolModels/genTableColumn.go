package toolModels

import (
	"nova-factory-server/app/baize"
	genUtils "nova-factory-server/app/business/tool/utils"
	"nova-factory-server/app/utils/snowflake"
	"nova-factory-server/app/utils/stringUtils"
	"strconv"
	"strings"
)

type GenTableColumnDML struct {
	ColumnId      int64  `json:"columnId,string" db:"column_id"`
	TableId       int64  `json:"tableId,string" db:"table_id"`
	ColumnName    string `json:"columnName" db:"column_name"`
	ColumnComment string `json:"columnComment" db:"column_comment"`
	ColumnType    string `json:"columnType" db:"column_type"`
	GoType        string `json:"goType" db:"go_type"`
	GoField       string `json:"goField" db:"go_field"`
	HtmlField     string `json:"htmlField" db:"html_field"`
	IsPk          string `json:"isPk" db:"is_pk"`
	IsRequired    string `json:"isRequired" db:"is_required"`
	IsInsert      string `json:"isInsert" db:"is_insert"`
	IsEdit        string `json:"isEdit" db:"is_edit"`
	IsList        string `json:"isList" db:"is_list"`
	IsQuery       string `json:"isQuery" db:"is_query"`
	QueryType     string `json:"queryType" db:"query_type"`
	HtmlType      string `json:"htmlType" db:"html_type"`
	DictType      string `json:"dictType" db:"dict_type"`
	Sort          int32  `json:"sort" db:"sort"`
	baize.BaseEntity
}

type GenTableColumnVo struct {
	ColumnId      int64  `json:"columnId,string" db:"column_id"`
	TableId       int64  `json:"tableId,string" db:"table_id"`
	ColumnName    string `json:"columnName" db:"column_name"`
	ColumnComment string `json:"columnComment" db:"column_comment"`
	ColumnType    string `json:"columnType" db:"column_type"`
	GoType        string `json:"goType" db:"go_type"`
	GoField       string `json:"goField" db:"go_field"`
	HtmlField     string `json:"htmlField" db:"html_field"`
	IsPk          string `json:"isPk" db:"is_pk"`
	IsRequired    string `json:"isRequired" db:"is_required"`
	IsInsert      string `json:"isInsert" db:"is_insert"`
	IsEdit        string `json:"isEdit" db:"is_edit"`
	IsList        string `json:"isList" db:"is_list"`
	IsQuery       string `json:"isQuery" db:"is_query"`
	IsEntity      string `json:"isEntity" db:"is_entity"`
	QueryType     string `json:"queryType" db:"query_type"`
	HtmlType      string `json:"htmlType" db:"html_type"`
	DictType      string `json:"dictType" db:"dict_type"`
	Sort          int32  `json:"remark" db:"sort"`
	baize.BaseEntity
}

type InformationSchemaColumn struct {
	ColumnName    string `db:"COLUMN_NAME"`
	ColumnComment string `db:"COLUMN_COMMENT"`
	ColumnType    string `db:"COLUMN_TYPE"`
	IsPk          string `db:"is_pk"`
	IsRequired    string `db:"is_required"`
	Sort          int32  `db:"sort"`
}

func GetGenTableColumnDML(column *InformationSchemaColumn, tableId int64, userId int64) *GenTableColumnDML {
	genTableColumn := new(GenTableColumnDML)
	dataType := column.ColumnType
	columnName := column.ColumnName
	genTableColumn.ColumnId = snowflake.GenID()
	genTableColumn.ColumnName = column.ColumnName
	genTableColumn.IsPk = column.IsPk
	genTableColumn.Sort = column.Sort
	genTableColumn.ColumnComment = column.ColumnComment
	genTableColumn.ColumnType = column.ColumnType
	genTableColumn.TableId = tableId
	genTableColumn.SetCreateBy(userId)
	//设置字段名
	genTableColumn.GoField = stringUtils.ConvertToBigCamelCase(columnName)
	genTableColumn.HtmlField = stringUtils.ConvertToLittleCamelCase(columnName)

	switch {
	case genUtils.IsTimeObject(dataType): //字段为时间类型
		genTableColumn.GoType = "Time"
		genTableColumn.HtmlType = "datetime"
	case genUtils.IsNumberObject(dataType): //字段为数字类型
		//字段为数字类型
		genTableColumn.HtmlType = "input"
		// 如果是浮点型
		tmp := genTableColumn.ColumnType
		if tmp == "float" || tmp == "double" {
			genTableColumn.GoType = "float64"
		} else {
			start := strings.Index(tmp, "(")
			end := strings.Index(tmp, ")")
			if end < 0 {
				genTableColumn.GoType = "int64"
			} else {
				arr := strings.Split(tmp[start+1:end], ",")
				i0, _ := strconv.Atoi(arr[0])
				if len(arr) == 2 && i0 > 0 {
					genTableColumn.GoType = "float64"
				} else {
					genTableColumn.GoType = "int64"
				}
			}
		}
	default:
		//字段为字符串类型
		genTableColumn.GoType = "string"
		if stringUtils.ReMatchingStr(dataType, "text|tinytext|mediumtext|longtext|longblob") {
			genTableColumn.HtmlType = "textarea"
		} else {
			columnLength := genUtils.GetColumnLength(column.ColumnType)
			if columnLength >= 500 {
				genTableColumn.HtmlType = "textarea"
			}
			genTableColumn.HtmlType = "input"
		}
	}

	//新增字段
	if genUtils.IsNotEntity(columnName) {
		genTableColumn.IsRequired = "1"
		if column.IsPk == "0" {
			genTableColumn.IsInsert = "1"
			genTableColumn.IsEdit = "1"
			genTableColumn.IsList = "1"
			genTableColumn.IsQuery = "1"
		} else {
			genTableColumn.IsInsert = "0"
			genTableColumn.IsEdit = "0"
			genTableColumn.IsList = "0"
			genTableColumn.IsQuery = "0"
		}
	} else {
		genTableColumn.IsRequired = "0"
		genTableColumn.IsInsert = "0"
		genTableColumn.IsEdit = "0"
		genTableColumn.IsList = "0"
		genTableColumn.IsQuery = "0"
	}

	// 查询字段类型
	if genUtils.CheckNameColumn(columnName) {
		genTableColumn.QueryType = "LIKE"
	} else {
		genTableColumn.QueryType = "EQ"
	}

	// 状态字段设置单选框
	if genUtils.CheckStatusColumn(columnName) {
		genTableColumn.HtmlType = "radio"
	} else if genUtils.CheckTypeColumn(columnName) || genUtils.CheckSexColumn(columnName) {
		// 类型&性别字段设置下拉框
		genTableColumn.HtmlType = "select"
	}
	return genTableColumn
}
