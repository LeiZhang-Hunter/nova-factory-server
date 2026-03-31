package materialDaoImpl

import (
	"errors"
	"nova-factory-server/app/business/iot/asset/material/materialDao"
	materialModels2 "nova-factory-server/app/business/iot/asset/material/materialModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MaterialDaoImpl struct {
	ms            *gorm.DB
	tableName     string
	inboundTable  string
	outboundTable string
}

func NewMaterialDaoImpl(ms *gorm.DB) materialDao.IMaterialDao {
	return &MaterialDaoImpl{
		tableName:     "sys_material",
		inboundTable:  "sys_material_inbound",
		outboundTable: "sys_material_outbound",
		ms:            ms,
	}
}
func (m *MaterialDaoImpl) InsertMaterial(c *gin.Context, material *materialModels2.MaterialInfo) (*materialModels2.MaterialVO, error) {
	if material == nil {
		return nil, errors.New("material is nil")
	}
	vo := materialModels2.NewMaterialVO(material)
	vo.MaterialId = uint64(snowflake.GenID())
	vo.SetCreateBy(baizeContext.GetUserId(c))
	vo.DeptId = uint64(baizeContext.GetDeptId(c))
	vo.Total = 0
	ret := m.ms.Table(m.tableName).Create(vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return vo, nil
}
func (m *MaterialDaoImpl) UpdateMaterial(c *gin.Context, material *materialModels2.MaterialInfo) (*materialModels2.MaterialVO, error) {
	if material == nil {
		return nil, errors.New("material is nil")
	}
	//deptId := baizeContext.GetDeptId(c)
	ret := m.ms.Table(m.tableName).Where("material_id = ?", material.MaterialId).Updates(material)
	if ret.Error != nil {
		return nil, ret.Error
	}
	var vo materialModels2.MaterialVO
	ret = m.ms.Table(m.tableName).Where("material_id = ?", material.MaterialId).Find(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}
func (m *MaterialDaoImpl) GetMaterialGroupByName(c *gin.Context, name string) (*materialModels2.MaterialVO, error) {
	var vo materialModels2.MaterialVO
	ret := m.ms.Table(m.tableName).Where("name = ?", name).Where("state = ?", commonStatus.NORMAL).First(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}
func (m *MaterialDaoImpl) GetNoExitIdMaterialGroupByName(c *gin.Context, name string, id uint64) (*materialModels2.MaterialVO, error) {
	var vo materialModels2.MaterialVO
	ret := m.ms.Table(m.tableName).Where("material_id != ?", id).Where("name = ?", name).Where("state = ?", commonStatus.NORMAL).First(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}
func (m *MaterialDaoImpl) SelectMaterialList(c *gin.Context, req *materialModels2.MaterialListReq) (*materialModels2.MaterialInfoListData, error) {
	db := m.ms.Table(m.tableName)

	if req != nil && req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req != nil && req.Code != "" {
		db = db.Where("code = ?", req.Code)
	}
	if req != nil && req.Model != "" {
		db = db.Where("model = ?", req.Model)
	}
	if req != nil && req.Factory != "" {
		db = db.Where("factory LIKE ?", "%"+req.Factory+"%")
	}
	if req != nil && req.Address != "" {
		db = db.Where("address LIKE ?", "%"+req.Address+"%")
	}
	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*materialModels2.MaterialVO

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &materialModels2.MaterialInfoListData{
			Rows:  make([]*materialModels2.MaterialVO, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &materialModels2.MaterialInfoListData{
			Rows:  make([]*materialModels2.MaterialVO, 0),
			Total: 0,
		}, ret.Error
	}
	return &materialModels2.MaterialInfoListData{
		Rows:  dto,
		Total: total,
	}, nil
}
func (m *MaterialDaoImpl) DeleteByMaterialIds(c *gin.Context, ids []int64) error {
	ret := m.ms.Table(m.tableName).Where("material_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (m *MaterialDaoImpl) GetByMaterialId(c *gin.Context, id int64) (*materialModels2.MaterialVO, error) {
	var vo *materialModels2.MaterialVO
	ret := m.ms.Table(m.tableName).Where("material_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&vo)
	return vo, ret.Error
}

// Inbound 入库
func (m *MaterialDaoImpl) Inbound(c *gin.Context, info *materialModels2.InboundInfo) (*materialModels2.InboundVO, error) {
	tx := m.ms.Table(m.tableName).Begin()
	var vo *materialModels2.MaterialVO
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&vo, info.MaterialId).Error; err != nil {
		zap.L().Error("入库失败", zap.Error(err))
		tx.Rollback()
		return nil, errors.New("入库失败")
	}
	if vo == nil {
		tx.Rollback()
		return nil, errors.New("物料不存在")
	}
	vo.Total += info.Number
	vo.CurrentTotal += info.Number
	ret := tx.Where("material_id = ?", info.MaterialId).Updates(vo)
	if ret.Error != nil {
		zap.L().Error("入库失败", zap.Error(ret.Error))
		tx.Rollback()
		return nil, errors.New("入库失败")
	}
	value := &materialModels2.InboundVO{
		MaterialId: vo.MaterialId,
		InboundId:  uint64(snowflake.GenID()),
		DeptId:     baizeContext.GetDeptId(c),
		Number:     info.Number,
	}
	value.SetCreateBy(baizeContext.GetUserId(c))
	ret = tx.Table(m.inboundTable).Create(value)
	if ret.Error != nil {
		tx.Rollback()
		zap.L().Error("入库失败", zap.Error(ret.Error))
		return nil, errors.New("入库失败")
	}
	tx.Commit()
	return value, nil
}

// Outbound 出库
func (m *MaterialDaoImpl) Outbound(c *gin.Context, info *materialModels2.OutboundInfo) (*materialModels2.OutboundVO, error) {
	tx := m.ms.Table(m.tableName).Begin()
	var vo *materialModels2.MaterialVO
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&vo, info.MaterialId).Error; err != nil {
		zap.L().Error("入库失败", zap.Error(err))
		tx.Rollback()
		return nil, errors.New("出库失败")
	}
	if vo == nil {
		tx.Rollback()
		return nil, errors.New("物料不存在")
	}
	if vo.CurrentTotal < info.Number {
		tx.Rollback()
		return nil, errors.New("库存不足")
	}
	vo.CurrentTotal -= info.Number
	vo.Outbound += info.Number
	ret := tx.Where("material_id = ?", info.MaterialId).Updates(vo)
	if ret.Error != nil {
		tx.Rollback()
		zap.L().Error("出库失败", zap.Error(ret.Error))
		return nil, errors.New("出库失败")
	}
	value := &materialModels2.OutboundVO{
		MaterialId: vo.MaterialId,
		OutboundId: uint64(snowflake.GenID()),
		DeptId:     baizeContext.GetDeptId(c),
		Number:     info.Number,
		Reason:     info.Reason,
	}
	value.SetCreateBy(baizeContext.GetUserId(c))
	ret = tx.Table(m.outboundTable).Create(value)
	if ret.Error != nil {
		tx.Rollback()
		zap.L().Error("出库失败", zap.Error(ret.Error))
		return nil, errors.New("出库失败")
	}
	tx.Commit()
	return value, nil
}

// InboundList 入库列表
func (m *MaterialDaoImpl) InboundList(c *gin.Context, req *materialModels2.InboundListReq) (*materialModels2.InboundListData, error) {
	db := m.ms.Table(m.inboundTable).Debug().Select(
		"sys_material.name as name, sys_material.code as code, " +
			"sys_material.model as model," +
			"sys_material.factory as factory, sys_material.address as address, sys_material.price as price, " +
			"sys_material.total as total," + "sys_material.outbound as outbound, sys_material.unit," +
			"sys_material.current_total as current_total, sys_material_inbound.inbound_id as inbound_id, sys_material_inbound.material_id as material_id, sys_material_inbound.create_by as create_by, sys_material_inbound.create_time as create_time," +
			"sys_material_inbound.update_by as update_by, sys_material_inbound.update_time as update_time, sys_material_inbound.number as number").Joins("right join sys_material on sys_material_inbound.material_id = sys_material.material_id")

	if req != nil && req.Keyword != "" {
		db = db.Where("sys_material.name LIKE ?", "%"+req.Keyword+"%").Or(
			db.Or("sys_material.code LIKE ?", "%"+req.Keyword+"%"),
			db.Or("sys_material.model LIKE ?", "%"+req.Keyword+"%"),
			db.Or("sys_material.factory LIKE ?", "%"+req.Keyword+"%"),
			db.Or("sys_material.address LIKE ?", "%"+req.Keyword+"%"))

	}

	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*materialModels2.InboundData

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &materialModels2.InboundListData{
			Rows:  make([]*materialModels2.InboundData, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &materialModels2.InboundListData{
			Rows:  make([]*materialModels2.InboundData, 0),
			Total: 0,
		}, ret.Error
	}
	return &materialModels2.InboundListData{
		Rows:  dto,
		Total: total,
	}, nil
}

// OutboundList 出库列表
func (m *MaterialDaoImpl) OutboundList(c *gin.Context, req *materialModels2.OutboundListReq) (*materialModels2.OutboundListData, error) {
	db := m.ms.Table(m.outboundTable).Debug().Select(
		"sys_material.name as name, sys_material.code as code, " +
			"sys_material.model as model," +
			"sys_material.factory as factory, sys_material.address as address, sys_material.price as price, " +
			"sys_material.total as total," + "sys_material.outbound as outbound, sys_material.unit," +
			"sys_material.current_total as current_total, sys_material_outbound.outbound_id as outbound_id, sys_material_outbound.material_id as material_id, sys_material_outbound.create_by as create_by, sys_material_outbound.create_time as create_time," +
			"sys_material_outbound.update_by as update_by, sys_material_outbound.update_time as update_time," +
			"sys_material_outbound.reason as reason," + " sys_material_outbound.number as number").Joins("right join sys_material on sys_material_outbound.material_id = sys_material.material_id")

	if req != nil && req.Keyword != "" {
		db = db.Where("sys_material.name LIKE ?", "%"+req.Keyword+"%").Or(
			db.Or("sys_material.code LIKE ?", "%"+req.Keyword+"%"),
			db.Or("sys_material.model LIKE ?", "%"+req.Keyword+"%"),
			db.Or("sys_material.factory LIKE ?", "%"+req.Keyword+"%"),
			db.Or("sys_material.address LIKE ?", "%"+req.Keyword+"%"))
	}

	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*materialModels2.OutboundData

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &materialModels2.OutboundListData{
			Rows:  make([]*materialModels2.OutboundData, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Limit(size).Find(&dto).Order("create_time desc")
	if ret.Error != nil {
		return &materialModels2.OutboundListData{
			Rows:  make([]*materialModels2.OutboundData, 0),
			Total: 0,
		}, ret.Error
	}
	return &materialModels2.OutboundListData{
		Rows:  dto,
		Total: total,
	}, nil
}
