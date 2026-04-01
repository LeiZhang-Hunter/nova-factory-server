package resourcedao

import (
	"context"
	"nova-factory-server/app/business/iot/asset/resource/resourcemodels"
)

type IResourceFileDao interface {
	InsertResource(ctx context.Context, resource *resourcemodels.SysResourceFile) (*resourcemodels.SysResourceFile, error)
	UpdateResource(ctx context.Context, resource *resourcemodels.SysResourceFile) (*resourcemodels.SysResourceFile, error)
	List(ctx context.Context, query *resourcemodels.SysResourceFileDQL) (*resourcemodels.SysResourceFileList, error)
	Remove(ctx context.Context, ids []string) error
	CheckNameUnique(ctx context.Context, parentId int64, name string, resourceId int64) int64
	CheckChildren(ctx context.Context, resourceId int64) int64
}
