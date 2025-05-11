package systemService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/monitor/monitorModels"
	"nova-factory-server/app/business/system/systemModels"
)

type ILoginService interface {
	Login(c *gin.Context, user *systemModels.User) string
	Register(c *gin.Context, user *systemModels.LoginBody)
	RecordLoginInfo(c *gin.Context, loginUser *monitorModels.Logininfor)
	GenerateCode(c *gin.Context) (m *systemModels.CaptchaVo)
	VerityCaptcha(c *gin.Context, id, base64 string) bool
	ForceLogout(c *gin.Context, token string)
	GetInfo(c *gin.Context) *systemModels.GetInfo
}
