package systemDao

import (
	"context"
	"nova-factory-server/app/business/system/systemModels"
)

type IUserDeptScopeDao interface {
	DeleteUserDeptScope(ctx context.Context, ids []int64)
	SelectUserDeptScopeDeptIdByUserId(ctx context.Context, id int64) []string
	DeleteUserDeptScopeByUserId(ctx context.Context, id int64)
	BatchUserDeptScope(ctx context.Context, list []*systemModels.SysUserDeptScope)
}
