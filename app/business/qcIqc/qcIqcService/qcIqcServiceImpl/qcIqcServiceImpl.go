package qcIqcServiceImpl

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/qcIqc/qcIqcApi"
	"nova-factory-server/app/business/qcIqc/qcIqcDao"
	"nova-factory-server/app/business/qcIqc/qcIqcModel"
	"nova-factory-server/app/business/qcIqc/qcIqcService"
	"nova-factory-server/app/utils/snowflake"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QcIqcServiceImpl 来料检验单服务实现
type QcIqcServiceImpl struct {
	qcIqcDao qcIqcDao.IQcIqcDao
}

// NewQcIqcServiceImpl 创建来料检验单服务实现
func NewQcIqcServiceImpl(qcIqcDao qcIqcDao.IQcIqcDao) qcIqcService.IQcIqcService {
	return &QcIqcServiceImpl{
		qcIqcDao: qcIqcDao,
	}
}

// List 获取来料检验单列表
func (s *QcIqcServiceImpl) List(c *gin.Context, req *qcIqcApi.QcIqcQueryReq) (*qcIqcApi.QcIqcListRes, error) {
	// 设置默认分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	total, models, err := s.qcIqcDao.List(c, req)
	if err != nil {
		return nil, err
	}

	// 将 model 结构体转换为 API 结构体
	rows := make([]*qcIqcApi.QcIqcData, 0, len(models))
	for _, model := range models {
		row := &qcIqcApi.QcIqcData{
			IqcId:                  int64(model.ID),
			IqcCode:                model.IqcCode,
			IqcName:                model.IqcName,
			TemplateId:             model.TemplateId,
			SourceDocId:            model.SourceDocId,
			SourceDocType:          model.SourceDocType,
			SourceDocCode:          model.SourceDocCode,
			SourceLineId:           model.SourceLineId,
			VendorId:               model.VendorId,
			VendorCode:             model.VendorCode,
			VendorName:             model.VendorName,
			VendorNick:             model.VendorNick,
			VendorBatch:            model.VendorBatch,
			ItemId:                 model.ItemId,
			ItemCode:               model.ItemCode,
			ItemName:               model.ItemName,
			Specification:          model.Specification,
			UnitOfMeasure:          model.UnitOfMeasure,
			UnitName:               model.UnitName,
			QuantityMinCheck:       getIntValue(model.QuantityMinCheck),
			QuantityMaxUnqualified: getIntValue(model.QuantityMaxUnqualified),
			QuantityRecived:        model.QuantityRecived,
			QuantityCheck:          getIntValue(model.QuantityCheck),
			QuantityQualified:      getIntValue(model.QuantityQualified),
			QuantityUnqualified:    getIntValue(model.QuantityUnqualified),
			CrRate:                 getFloat64Value(model.CrRate),
			MajRate:                getFloat64Value(model.MajRate),
			MinRate:                getFloat64Value(model.MinRate),
			CrQuantity:             getIntValue(model.CrQuantity),
			MajQuantity:            getIntValue(model.MajQuantity),
			MinQuantity:            getIntValue(model.MinQuantity),
			CheckResult:            model.CheckResult,
			ReciveDate:             model.ReciveDate,
			InspectDate:            model.InspectDate,
			Inspector:              model.Inspector,
			Status:                 getStringValue(model.Status),
			Remark:                 model.Remark,
			Attr1:                  model.Attr1,
			Attr2:                  model.Attr2,
			Attr3:                  getIntValue(model.Attr3),
			Attr4:                  getIntValue(model.Attr4),
			CreateBy:               model.CreateBy,
			UpdateBy:               model.UpdateBy,
			CreateById:             model.CreateById,
			UpdateById:             model.UpdateById,
			CreateTime:             model.CreatedAt,
			UpdateTime:             model.UpdatedAt,
		}
		rows = append(rows, row)
	}

	return &qcIqcApi.QcIqcListRes{
		Rows:  rows,
		Total: total,
	}, nil
}

// Create 创建来料检验单
func (s *QcIqcServiceImpl) Create(c *gin.Context, req *qcIqcApi.QcIqcCreateReq) (*baize.EmptyResponse, error) {
	// 设置默认值
	if req.Status == nil {
		defaultStatus := "PREPARE"
		req.Status = &defaultStatus
	}
	if req.QuantityMinCheck == nil {
		defaultMinCheck := 1
		req.QuantityMinCheck = &defaultMinCheck
	}
	if req.QuantityMaxUnqualified == nil {
		defaultMaxUnqualified := 0
		req.QuantityMaxUnqualified = &defaultMaxUnqualified
	}
	if req.QuantityQualified == nil {
		defaultQualified := 0
		req.QuantityQualified = &defaultQualified
	}
	if req.QuantityUnqualified == nil {
		defaultUnqualified := 0
		req.QuantityUnqualified = &defaultUnqualified
	}
	if req.CrRate == nil {
		defaultCrRate := 0.0
		req.CrRate = &defaultCrRate
	}
	if req.MajRate == nil {
		defaultMajRate := 0.0
		req.MajRate = &defaultMajRate
	}
	if req.MinRate == nil {
		defaultMinRate := 0.0
		req.MinRate = &defaultMinRate
	}
	if req.CrQuantity == nil {
		defaultCrQuantity := 0
		req.CrQuantity = &defaultCrQuantity
	}
	if req.MajQuantity == nil {
		defaultMajQuantity := 0
		req.MajQuantity = &defaultMajQuantity
	}
	if req.MinQuantity == nil {
		defaultMinQuantity := 0
		req.MinQuantity = &defaultMinQuantity
	}
	if req.Attr3 == nil {
		defaultAttr3 := 0
		req.Attr3 = &defaultAttr3
	}
	if req.Attr4 == nil {
		defaultAttr4 := 0
		req.Attr4 = &defaultAttr4
	}

	return s.qcIqcDao.Create(c, req)
}

// Update 更新来料检验单
func (s *QcIqcServiceImpl) Update(c *gin.Context, req *qcIqcApi.QcIqcUpdateReq) (*baize.EmptyResponse, error) {
	return s.qcIqcDao.Update(c, req)
}

// Delete 批量删除来料检验单
func (s *QcIqcServiceImpl) Delete(c *gin.Context, iqcIds []int64) error {
	return s.qcIqcDao.Delete(c, iqcIds)
}

// GetById 根据ID获取来料检验单
func (s *QcIqcServiceImpl) GetById(c *gin.Context, iqcId int64) (*qcIqcModel.QcIqc, error) {
	return s.qcIqcDao.GetById(c, iqcId)
}

// generateIqcCode 生成检验单编号
func (s *QcIqcServiceImpl) generateIqcCode() string {
	return "IQC" + strconv.FormatInt(snowflake.GenID(), 10)
}

// getIntValue 安全获取 int 值
func getIntValue(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

// getFloat64Value 安全获取 float64 值
func getFloat64Value(value *float64) float64 {
	if value == nil {
		return 0.0
	}
	return *value
}

// getStringValue 安全获取 string 值
func getStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
