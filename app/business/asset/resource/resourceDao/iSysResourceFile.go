package resourceDao

import (
	"context"
	"nova-factory-server/app/business/asset/resource/resourceModels"
)

type IResourceFileDao interface {
	InsertResource(ctx context.Context, resource *resourceModels.SysResourceFile) (*resourceModels.SysResourceFile, error)
	UpdateResource(ctx context.Context, resource *resourceModels.SysResourceFile) (*resourceModels.SysResourceFile, error)
	List(ctx context.Context, query *resourceModels.SysResourceFileDQL) (*resourceModels.SysResourceFileList, error)
	Remove(ctx context.Context, ids []string) error
	CheckNameUnique(ctx context.Context, parentId int64, name string, resourceId int64) int64
	CheckChildren(ctx context.Context, resourceId int64) int64
}
