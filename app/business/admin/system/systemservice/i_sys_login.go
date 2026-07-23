package systemservice

import (
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type ILoginService interface {
	Login(c *gin.Context, user *modelentity.User) *modelresponse.LoginResp
	Register(c *gin.Context, user *modelrequest.LoginBody)
	RecordLoginInfo(c *gin.Context, loginUser *monitormodels.Logininfor)
	GenerateCode(c *gin.Context) (m *modelresponse.CaptchaVo)
	VerityCaptcha(c *gin.Context, id, base64 string) bool
	ForceLogout(c *gin.Context, token string)
	GetInfo(c *gin.Context) *modelresponse.GetInfo
}
