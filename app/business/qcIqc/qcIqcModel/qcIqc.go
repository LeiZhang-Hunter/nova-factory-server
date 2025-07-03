package qcIqcModel

import (
	"time"

	"gorm.io/gorm"
)

// QcIqc 来料检验单表
type QcIqc struct {
	gorm.Model
	IqcCode                string     `json:"iqcCode" gorm:"column:iqc_code;size:64;not null;comment:来料检验单编号"`
	IqcName                string     `json:"iqcName" gorm:"column:iqc_name;size:500;not null;comment:来料检验单名称"`
	TemplateId             int64      `json:"templateId" gorm:"column:template_id;not null;comment:检验模板ID"`
	SourceDocId            *int64     `json:"sourceDocId" gorm:"column:source_doc_id;comment:来源单据ID"`
	SourceDocType          *string    `json:"sourceDocType" gorm:"column:source_doc_type;size:64;comment:来源单据类型"`
	SourceDocCode          *string    `json:"sourceDocCode" gorm:"column:source_doc_code;size:64;comment:来源单据编号"`
	SourceLineId           *int64     `json:"sourceLineId" gorm:"column:source_line_id;comment:来源单据行ID"`
	VendorId               int64      `json:"vendorId" gorm:"column:vendor_id;not null;comment:供应商ID"`
	VendorCode             string     `json:"vendorCode" gorm:"column:vendor_code;size:64;not null;comment:供应商编码"`
	VendorName             string     `json:"vendorName" gorm:"column:vendor_name;size:255;not null;comment:供应商名称"`
	VendorNick             *string    `json:"vendorNick" gorm:"column:vendor_nick;size:255;comment:供应商简称"`
	VendorBatch            *string    `json:"vendorBatch" gorm:"column:vendor_batch;size:64;comment:供应商批次号"`
	ItemId                 int64      `json:"itemId" gorm:"column:item_id;not null;comment:产品物料ID"`
	ItemCode               *string    `json:"itemCode" gorm:"column:item_code;size:64;comment:产品物料编码"`
	ItemName               *string    `json:"itemName" gorm:"column:item_name;size:255;comment:产品物料名称"`
	Specification          *string    `json:"specification" gorm:"column:specification;size:500;comment:规格型号"`
	UnitOfMeasure          *string    `json:"unitOfMeasure" gorm:"column:unit_of_measure;size:64;comment:单位"`
	UnitName               *string    `json:"unitName" gorm:"column:unit_name;size:64;comment:单位名称"`
	QuantityMinCheck       *int       `json:"quantityMinCheck" gorm:"column:quantity_min_check;default:1;comment:最低检测数"`
	QuantityMaxUnqualified *int       `json:"quantityMaxUnqualified" gorm:"column:quantity_max_unqualified;default:0;comment:最大不合格数"`
	QuantityRecived        float64    `json:"quantityRecived" gorm:"column:quantity_recived;not null;comment:本次接收数量"`
	QuantityCheck          *int       `json:"quantityCheck" gorm:"column:quantity_check;comment:本次检测数量"`
	QuantityQualified      *int       `json:"quantityQualified" gorm:"column:quantity_qualified;default:0;comment:合格数"`
	QuantityUnqualified    *int       `json:"quantityUnqualified" gorm:"column:quantity_unqualified;default:0;comment:不合格数"`
	CrRate                 *float64   `json:"crRate" gorm:"column:cr_rate;default:0.00;comment:致命缺陷率"`
	MajRate                *float64   `json:"majRate" gorm:"column:maj_rate;default:0.00;comment:严重缺陷率"`
	MinRate                *float64   `json:"minRate" gorm:"column:min_rate;default:0.00;comment:轻微缺陷率"`
	CrQuantity             *int       `json:"crQuantity" gorm:"column:cr_quantity;default:0;comment:致命缺陷数量"`
	MajQuantity            *int       `json:"majQuantity" gorm:"column:maj_quantity;default:0;comment:严重缺陷数量"`
	MinQuantity            *int       `json:"minQuantity" gorm:"column:min_quantity;default:0;comment:轻微缺陷数量"`
	CheckResult            *string    `json:"checkResult" gorm:"column:check_result;size:64;comment:检测结果"`
	ReciveDate             *time.Time `json:"reciveDate" gorm:"column:recive_date;comment:来料日期"`
	InspectDate            *time.Time `json:"inspectDate" gorm:"column:inspect_date;comment:检测日期"`
	Inspector              *string    `json:"inspector" gorm:"column:inspector;size:64;comment:检测人员"`
	Status                 *string    `json:"status" gorm:"column:status;size:64;comment:单据状态"`
	Remark                 *string    `json:"remark" gorm:"column:remark;size:500;comment:备注"`
	Attr1                  *string    `json:"attr1" gorm:"column:attr1;size:64;comment:预留字段1"`
	Attr2                  *string    `json:"attr2" gorm:"column:attr2;size:255;comment:预留字段2"`
	Attr3                  *int       `json:"attr3" gorm:"column:attr3;default:0;comment:预留字段3"`
	Attr4                  *int       `json:"attr4" gorm:"column:attr4;default:0;comment:预留字段4"`
	CreateBy               string     `json:"createBy" gorm:"column:create_by;size:64;comment:创建者"`
	UpdateBy               string     `json:"updateBy" gorm:"column:update_by;size:64;comment:更新者"`
	CreateById             int64      `json:"createById" gorm:"column:create_by_id;comment:创建者ID"`
	UpdateById             int64      `json:"updateById" gorm:"column:update_by_id;comment:更新者ID"`
}

func (QcIqc) TableName() string {
	return "qc_iqc"
}
