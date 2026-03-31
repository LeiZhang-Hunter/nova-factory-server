package systemservice

import (
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	systemModels2 "nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type ILoginService interface {
	Login(c *gin.Context, user *systemModels2.User) string
	Register(c *gin.Context, user *systemModels2.LoginBody)
	RecordLoginInfo(c *gin.Context, loginUser *monitormodels.Logininfor)
	GenerateCode(c *gin.Context) (m *systemModels2.CaptchaVo)
	VerityCaptcha(c *gin.Context, id, base64 string) bool
	ForceLogout(c *gin.Context, token string)
	GetInfo(c *gin.Context) *systemModels2.GetInfo
}
