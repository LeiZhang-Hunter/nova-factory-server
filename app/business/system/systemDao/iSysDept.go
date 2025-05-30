package systemDao

import (
	"context"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/system/systemModels"
)

type IDeptDao interface {
	SelectDeptList(ctx context.Context, dept *systemModels.SysDeptDQL) (sysDeptList []*systemModels.SysDeptVo)
	SelectDeptListSelectBox(ctx context.Context, dept *baize.BaseEntityDQL) (list []*systemModels.SelectDept)
	SelectDeptById(ctx context.Context, deptId int64) (dept *systemModels.SysDeptVo)
	InsertDept(ctx context.Context, dept *systemModels.SysDeptVo)
	UpdateDept(ctx context.Context, dept *systemModels.SysDeptVo)
	DeleteDeptById(ctx context.Context, deptId int64)
	CheckDeptNameUnique(ctx context.Context, deptName string, parentId int64) int64
	HasChildByDeptId(ctx context.Context, deptId int64) int
	CheckDeptExistUser(ctx context.Context, deptId int64) int
}
