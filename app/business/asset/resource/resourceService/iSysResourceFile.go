package resourceService

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/resource/resourceModels"
)

type IResourceFileService interface {
	InsertResource(c *gin.Context, resource *resourceModels.SysResourceFileDML) (*resourceModels.SysResourceFile, error)
	UpdateResource(c *gin.Context, resource *resourceModels.SysResourceFileDML) (*resourceModels.SysResourceFile, error)
	List(ctx context.Context, query *resourceModels.SysResourceFileDQL) (*resourceModels.SysResourceFileList, error)
	Remove(ctx context.Context, ids []string) error
}
