package resourceservice

import (
	"context"
	"nova-factory-server/app/business/iot/asset/resource/resourcemodels"

	"github.com/gin-gonic/gin"
)

type IResourceFileService interface {
	InsertResource(c *gin.Context, resource *resourcemodels.SysResourceFileDML) (*resourcemodels.SysResourceFile, error)
	UpdateResource(c *gin.Context, resource *resourcemodels.SysResourceFileDML) (*resourcemodels.SysResourceFile, error)
	List(ctx context.Context, query *resourcemodels.SysResourceFileDQL) (*resourcemodels.SysResourceFileList, error)
	Remove(ctx context.Context, ids []string) error
	CheckNameUnique(ctx context.Context, parentId int64, name string, resourceId int64) int64
	CheckChildren(ctx context.Context, resourceId int64) int64
}
