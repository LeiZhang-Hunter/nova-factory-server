package resourceDaoImpl

import (
	"context"
	"gorm.io/gorm"
	"nova-factory-server/app/business/asset/resource/resourceDao"
	"nova-factory-server/app/business/asset/resource/resourceModels"
	"nova-factory-server/app/constant/commonStatus"
)

type sysResourceFileDao struct {
	db    *gorm.DB
	table string
}

func NewSysResourceFileDao(db *gorm.DB) resourceDao.IResourceFileDao {
	return &sysResourceFileDao{
		db:    db,
		table: "sys_resource_file",
	}
}

func (dao *sysResourceFileDao) InsertResource(ctx context.Context, resource *resourceModels.SysResourceFile) (*resourceModels.SysResourceFile, error) {
	ret := dao.db.Table(dao.table).Create(&resource)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return resource, nil
}

func (dao *sysResourceFileDao) UpdateResource(ctx context.Context, resource *resourceModels.SysResourceFile) (*resourceModels.SysResourceFile, error) {
	ret := dao.db.Table(dao.table).Where("resource_id = ?", resource.ResourceID).Updates(&resource)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return resource, nil
}

func (dao *sysResourceFileDao) List(ctx context.Context, query *resourceModels.SysResourceFileDQL) (*resourceModels.SysResourceFileList, error) {
	result := &resourceModels.SysResourceFileList{
		List: make([]*resourceModels.SysResourceFileVo, 0),
	}

	dbQuery := dao.db.Table(dao.table)
	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Type != "" {
		dbQuery = dbQuery.Where("type = ?", query.Type)
	}
	if query.Category != "" {
		dbQuery = dbQuery.Where("category = ?", query.Category)
	}
	if query.ParentId != nil {
		dbQuery = dbQuery.Where("parent_id = ?", *query.ParentId)
	}

	if query.Status != nil {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	dbQuery = dbQuery.Where("state = ?", commonStatus.NORMAL)
	var total int64
	ret := dbQuery.Count(&total)
	if ret.Error != nil {
		return result, ret.Error
	}
	offset := 0
	if query.Page <= 0 {
		query.Page = 1
	} else {
		offset = int((query.Page - 1) * query.Size)
	}
	size := 0
	if query.Size <= 0 {
		size = 20
	} else {
		size = int(query.Size)
	}
	list := make([]*resourceModels.SysResourceFileVo, 0)
	if query.Type == "FOLDER" {
		ret = dbQuery.Where("state = ?", commonStatus.NORMAL).Offset(0).Order("create_time desc").Limit(1000).Find(&list)
		if ret.Error != nil {
			return result, ret.Error
		}
	} else {
		ret = dbQuery.Offset(offset).Order("create_time desc").Limit(size).Find(&list)
		if ret.Error != nil {
			return result, ret.Error
		}
	}

	result.List = list
	result.Total = total
	return result, ret.Error
}

func (dao *sysResourceFileDao) Remove(ctx context.Context, ids []string) error {
	ret := dao.db.Table(dao.table).Where("resource_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (dao *sysResourceFileDao) CheckNameUnique(ctx context.Context, parentId int64, name string, resourceId int64) int64 {
	query := dao.db.Table(dao.table)
	if resourceId > 0 {
		query = query.Where("resource_id != ?", resourceId)
	}
	query = query.Where("parent_id = ?", parentId)
	query = query.Where("name = ?", name).Where("state = ?", commonStatus.NORMAL)

	var count int64
	ret := query.Debug().Count(&count)
	if ret.Error != nil {
		return 0
	}
	return count
}

func (dao *sysResourceFileDao) CheckChildren(ctx context.Context, resourceId int64) int64 {
	query := dao.db.Table(dao.table)
	query = query.Where("parent_id = ?", resourceId).Where("state = ?", commonStatus.NORMAL)
	var count int64
	ret := query.Count(&count)
	if ret.Error != nil {
		return 0
	}
	return count
}
