package systemdao

import (
	"context"
	systemModels2 "nova-factory-server/app/business/admin/system/systemmodels"
)

type IUserDao interface {
	SelectUserNameByUserName(ctx context.Context, userName []string) []string
	CheckUserNameUnique(ctx context.Context, userName string) int
	CheckPhoneUnique(ctx context.Context, phonenumber string) int64
	CheckEmailUnique(ctx context.Context, email string) int64
	InsertUser(ctx context.Context, sysUser *systemModels2.SysUserDML)
	BatchInsertUser(ctx context.Context, sysUser []*systemModels2.SysUserDML)
	UpdateUser(ctx context.Context, sysUser *systemModels2.SysUserDML)
	SelectUserByUserName(ctx context.Context, userName string) (loginUser *systemModels2.User)
	SelectUserById(ctx context.Context, userId int64) (sysUser *systemModels2.SysUserVo)
	SelectUserList(ctx context.Context, user *systemModels2.SysUserDQL) (sysUserList []*systemModels2.SysUserVo, total int64)
	SelectUserListAll(ctx context.Context, user *systemModels2.SysUserDQL) (list []*systemModels2.SysUserVo)
	DeleteUserByIds(ctx context.Context, ids []int64)
	UpdateUserAvatar(ctx context.Context, userId int64, avatar string)
	ResetUserPwd(ctx context.Context, userId int64, password string)
	SelectPasswordByUserId(ctx context.Context, userId int64) string
	SelectUserIdsByDeptIds(ctx context.Context, deptIds []int64) []int64
	SelectByUserIds(ctx context.Context, userIds []int64) []*systemModels2.SysUserDML
}
