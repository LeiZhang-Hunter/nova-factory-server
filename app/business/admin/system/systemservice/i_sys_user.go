package systemservice

import (
	"mime/multipart"
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	SelectUserByUserName(c *gin.Context, userName string) *modelentity.User
	SelectUserList(c *gin.Context, user *modelquery.SysUserDQL) (sysUserList []*modelresponse.SysUserVo, total int64)
	UserExport(c *gin.Context, user *modelquery.SysUserDQL) (data []byte)
	InsertUser(c *gin.Context, sysUser *modelrequest.SysUserDML)
	UpdateUser(c *gin.Context, sysUser *modelrequest.SysUserDML)

	UpdateUserDataScope(c *gin.Context, uds *modelrequest.SysUserDataScope)
	SelectUserDataScope(c *gin.Context, userId int64) *modelrequest.SysUserDataScope

	UpdateUserStatus(c *gin.Context, sysUser *modelrequest.EditUserStatus)
	ResetPwd(c *gin.Context, userId int64, password string)
	CheckUserNameUnique(c *gin.Context, userName string) bool
	CheckPhoneUnique(c *gin.Context, id int64, phonenumber string) bool
	CheckEmailUnique(c *gin.Context, id int64, email string) bool
	DeleteUserByIds(c *gin.Context, ids []int64)
	UserImportData(c *gin.Context, file *multipart.FileHeader) (msg string, failureNum int)
	UpdateUserAvatar(c *gin.Context, file *multipart.FileHeader) string
	ResetUserPwd(c *gin.Context, userId int64, password string)
	UpdateUserProfile(c *gin.Context, sysUser *modelrequest.SysUserDML)
	MatchesPassword(c *gin.Context, rawPassword string, userId int64) bool
	InsertUserAuth(c *gin.Context, userId int64, roleIds []int64)
	GetUserAuthRole(c *gin.Context, userId int64) *modelresponse.UserAndRoles
	SelectUserAndAccreditById(c *gin.Context, userId int64) (sysUser *modelresponse.UserAndAccredit)
	SelectAccredit(c *gin.Context) (sysUser *modelresponse.Accredit)
	ImportTemplate(c *gin.Context) (data []byte)
	GetUserProfile(c *gin.Context) *modelresponse.UserProfile
}
