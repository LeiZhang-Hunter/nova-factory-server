package resourceService

import (
	"context"
	"nova-factory-server/app/business/iot/asset/resource/resourceModels"

	"github.com/gin-gonic/gin"
)

type IResourceFileService interface {
	InsertResource(c *gin.Context, resource *resourceModels.SysResourceFileDML) (*resourceModels.SysResourceFile, error)
	UpdateResource(c *gin.Context, resource *resourceModels.SysResourceFileDML) (*resourceModels.SysResourceFile, error)
	List(ctx context.Context, query *resourceModels.SysResourceFileDQL) (*resourceModels.SysResourceFileList, error)
	Remove(ctx context.Context, ids []string) error
	CheckNameUnique(ctx context.Context, parentId int64, name string, resourceId int64) int64
	CheckChildren(ctx context.Context, resourceId int64) int64
}
