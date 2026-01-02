package resourceServiceImpl

import (
	"context"
	"nova-factory-server/app/business/asset/resource/resourceDao"
	"nova-factory-server/app/business/asset/resource/resourceModels"
	"nova-factory-server/app/business/asset/resource/resourceService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type sysResourceFileService struct {
	resourceDao resourceDao.IResourceFileDao
}

func NewSysResourceFileService(resourceDao resourceDao.IResourceFileDao) resourceService.IResourceFileService {
	return &sysResourceFileService{resourceDao: resourceDao}
}

func (s *sysResourceFileService) InsertResource(c *gin.Context, resource *resourceModels.SysResourceFileDML) (*resourceModels.SysResourceFile, error) {
	value := resourceModels.ToSysResourceFile(resource)
	value.SetCreateBy(baizeContext.GetUserId(c))
	value.ResourceID = snowflake.GenID()
	// Default status
	return s.resourceDao.InsertResource(c, value)
}

func (s *sysResourceFileService) UpdateResource(c *gin.Context, resource *resourceModels.SysResourceFileDML) (*resourceModels.SysResourceFile, error) {
	value := resourceModels.ToSysResourceFile(resource)
	value.SetUpdateBy(baizeContext.GetUserId(c))
	return s.resourceDao.UpdateResource(c, value)
}

func (s *sysResourceFileService) List(ctx context.Context, query *resourceModels.SysResourceFileDQL) (*resourceModels.SysResourceFileList, error) {
	return s.resourceDao.List(ctx, query)
}

func (s *sysResourceFileService) Remove(ctx context.Context, ids []string) error {
	return s.resourceDao.Remove(ctx, ids)
}
