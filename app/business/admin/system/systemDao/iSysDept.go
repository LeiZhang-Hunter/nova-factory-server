package systemDao

import (
	"context"
	"nova-factory-server/app/baize"
	systemModels2 "nova-factory-server/app/business/admin/system/systemModels"
)

type IDeptDao interface {
	SelectDeptList(ctx context.Context, dept *systemModels2.SysDeptDQL) (sysDeptList []*systemModels2.SysDeptVo)
	SelectDeptListSelectBox(ctx context.Context, dept *baize.BaseEntityDQL) (list []*systemModels2.SelectDept)
	SelectDeptById(ctx context.Context, deptId int64) (dept *systemModels2.SysDeptVo)
	InsertDept(ctx context.Context, dept *systemModels2.SysDeptVo)
	UpdateDept(ctx context.Context, dept *systemModels2.SysDeptVo)
	DeleteDeptById(ctx context.Context, deptId int64)
	CheckDeptNameUnique(ctx context.Context, deptName string, parentId int64) int64
	HasChildByDeptId(ctx context.Context, deptId int64) int
	CheckDeptExistUser(ctx context.Context, deptId int64) int
}
