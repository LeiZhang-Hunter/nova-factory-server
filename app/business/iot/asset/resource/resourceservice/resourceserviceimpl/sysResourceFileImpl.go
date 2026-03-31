package resourceserviceimpl

import (
	"context"
	"nova-factory-server/app/business/iot/asset/resource/resourcedao"
	"nova-factory-server/app/business/iot/asset/resource/resourcemodels"
	"nova-factory-server/app/business/iot/asset/resource/resourceservice"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type sysResourceFileService struct {
	resourceDao resourcedao.IResourceFileDao
}

func NewSysResourceFileService(resourceDao resourcedao.IResourceFileDao) resourceservice.IResourceFileService {
	return &sysResourceFileService{resourceDao: resourceDao}
}

func (s *sysResourceFileService) InsertResource(c *gin.Context, resource *resourcemodels.SysResourceFileDML) (*resourcemodels.SysResourceFile, error) {
	value := resourcemodels.ToSysResourceFile(resource)
	value.SetCreateBy(baizeContext.GetUserId(c))
	value.ResourceID = snowflake.GenID()
	// Default status
	return s.resourceDao.InsertResource(c, value)
}

func (s *sysResourceFileService) UpdateResource(c *gin.Context, resource *resourcemodels.SysResourceFileDML) (*resourcemodels.SysResourceFile, error) {
	value := resourcemodels.ToSysResourceFile(resource)
	value.SetUpdateBy(baizeContext.GetUserId(c))
	return s.resourceDao.UpdateResource(c, value)
}

func (s *sysResourceFileService) List(ctx context.Context, query *resourcemodels.SysResourceFileDQL) (*resourcemodels.SysResourceFileList, error) {
	return s.resourceDao.List(ctx, query)
}

func (s *sysResourceFileService) Remove(ctx context.Context, ids []string) error {
	return s.resourceDao.Remove(ctx, ids)
}

func (s *sysResourceFileService) CheckNameUnique(ctx context.Context, parentId int64, name string, resourceId int64) int64 {
	return s.resourceDao.CheckNameUnique(ctx, parentId, name, resourceId)
}
func (s *sysResourceFileService) CheckChildren(ctx context.Context, resourceId int64) int64 {
	return s.resourceDao.CheckChildren(ctx, resourceId)
}
