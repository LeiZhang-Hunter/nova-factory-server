package systemservice

import (
	"mime/multipart"
	systemModels2 "nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	SelectUserByUserName(c *gin.Context, userName string) *systemModels2.User
	SelectUserList(c *gin.Context, user *systemModels2.SysUserDQL) (sysUserList []*systemModels2.SysUserVo, total int64)
	UserExport(c *gin.Context, user *systemModels2.SysUserDQL) (data []byte)
	InsertUser(c *gin.Context, sysUser *systemModels2.SysUserDML)
	UpdateUser(c *gin.Context, sysUser *systemModels2.SysUserDML)

	UpdateUserDataScope(c *gin.Context, uds *systemModels2.SysUserDataScope)
	SelectUserDataScope(c *gin.Context, userId int64) *systemModels2.SysUserDataScope

	UpdateUserStatus(c *gin.Context, sysUser *systemModels2.EditUserStatus)
	ResetPwd(c *gin.Context, userId int64, password string)
	CheckUserNameUnique(c *gin.Context, userName string) bool
	CheckPhoneUnique(c *gin.Context, id int64, phonenumber string) bool
	CheckEmailUnique(c *gin.Context, id int64, email string) bool
	DeleteUserByIds(c *gin.Context, ids []int64)
	UserImportData(c *gin.Context, file *multipart.FileHeader) (msg string, failureNum int)
	UpdateUserAvatar(c *gin.Context, file *multipart.FileHeader) string
	ResetUserPwd(c *gin.Context, userId int64, password string)
	UpdateUserProfile(c *gin.Context, sysUser *systemModels2.SysUserDML)
	MatchesPassword(c *gin.Context, rawPassword string, userId int64) bool
	InsertUserAuth(c *gin.Context, userId int64, roleIds []int64)
	GetUserAuthRole(c *gin.Context, userId int64) *systemModels2.UserAndRoles
	SelectUserAndAccreditById(c *gin.Context, userId int64) (sysUser *systemModels2.UserAndAccredit)
	SelectAccredit(c *gin.Context) (sysUser *systemModels2.Accredit)
	ImportTemplate(c *gin.Context) (data []byte)
	GetUserProfile(c *gin.Context) *systemModels2.UserProfile
}
