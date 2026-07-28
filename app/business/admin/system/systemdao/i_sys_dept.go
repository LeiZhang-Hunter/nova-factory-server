package systemdao

import (
	"context"
	"nova-factory-server/app/baize"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
)

type IDeptDao interface {
	SelectDeptList(ctx context.Context, dept *modelquery.SysDeptDQL) (sysDeptList []*modelresponse.SysDeptVo)
	SelectDeptListSelectBox(ctx context.Context, dept *baize.BaseEntityDQL) (list []*modelresponse.SelectDept)
	SelectDeptById(ctx context.Context, deptId int64) (dept *modelresponse.SysDeptVo)
	InsertDept(ctx context.Context, dept *modelrequest.SysDeptDML)
	UpdateDept(ctx context.Context, dept *modelrequest.SysDeptDML)
	DeleteDeptById(ctx context.Context, deptId int64)
	CheckDeptNameUnique(ctx context.Context, deptName string, parentId int64) int64
	HasChildByDeptId(ctx context.Context, deptId int64) int
	CheckDeptExistUser(ctx context.Context, deptId int64) int
}
