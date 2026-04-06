package aiDataSetDaoImpl

import (
	"errors"
	"strings"
	"time"

	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IDatasetRolePermissionDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewIDatasetRolePermissionDaoImpl 知识库/文档-角色权限DAO构造函数。
func NewIDatasetRolePermissionDaoImpl(db *gorm.DB) aidatasetdao.IDatasetRolePermissionDao {
	return &IDatasetRolePermissionDaoImpl{
		db:    db,
		table: "sys_dataset_role_permission",
	}
}

// Create 创建知识库/文档-角色权限记录。
func (i *IDatasetRolePermissionDaoImpl) Create(c *gin.Context, req *aidatasetmodels.SetDatasetRolePermission) (*aidatasetmodels.DatasetRolePermission, error) {
	data := &aidatasetmodels.DatasetRolePermission{
		ID:         snowflake.GenID(),
		RoleID:     req.RoleID,
		DatasetID:  req.DatasetID,
		DocumentID: req.DocumentID,
		Permission: req.Permission,
		DeptID:     baizeContext.GetDeptId(c),
		State:      commonStatus.NORMAL,
	}
	data.SetCreateBy(baizeContext.GetUserId(c))
	if err := i.db.WithContext(c).Table(i.table).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update 更新知识库/文档-角色权限记录。
func (i *IDatasetRolePermissionDaoImpl) Update(c *gin.Context, req *aidatasetmodels.SetDatasetRolePermission) (*aidatasetmodels.DatasetRolePermission, error) {
	if req.ID == 0 {
		return nil, errors.New("id不能为空")
	}
	data := &aidatasetmodels.DatasetRolePermission{}
	if err := i.db.WithContext(c).Table(i.table).
		Where("id = ?", req.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(data).Error; err != nil {
		return nil, err
	}
	data.RoleID = req.RoleID
	data.DatasetID = req.DatasetID
	data.DocumentID = req.DocumentID
	data.Permission = req.Permission
	data.SetUpdateBy(baizeContext.GetUserId(c))
	if err := i.db.WithContext(c).Table(i.table).
		Where("id = ?", req.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"role_id":     data.RoleID,
			"dataset_id":  data.DatasetID,
			"document_id": data.DocumentID,
			"permission":  data.Permission,
			"update_by":   data.UpdateBy,
			"update_time": data.UpdateTime,
		}).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// List 查询知识库/文档-角色权限列表。
func (i *IDatasetRolePermissionDaoImpl) List(c *gin.Context, req *aidatasetmodels.DatasetRolePermissionQuery) (*aidatasetmodels.DatasetRolePermissionListData, error) {
	db := i.db.WithContext(c).Table(i.table).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if req.RoleID > 0 {
		db = db.Where("role_id = ?", req.RoleID)
	}
	if req.DatasetID > 0 {
		db = db.Where("dataset_id = ?", req.DatasetID)
	}
	if req.DocumentID > 0 {
		db = db.Where("document_id = ?", req.DocumentID)
	}
	if strings.TrimSpace(req.Permission) != "" {
		db = db.Where("permission = ?", strings.TrimSpace(req.Permission))
	}
	page := req.GetPage()
	size := req.GetSize()
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*aidatasetmodels.DatasetRolePermission, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &aidatasetmodels.DatasetRolePermissionListData{
		Rows:  rows,
		Total: total,
	}, nil
}

// Remove 删除知识库/文档-角色权限记录（软删除）。
func (i *IDatasetRolePermissionDaoImpl) Remove(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return i.db.WithContext(c).Table(i.table).
		Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}
