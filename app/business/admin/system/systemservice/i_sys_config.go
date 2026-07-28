package systemservice

import (
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type IConfigService interface {
	SelectConfigList(c *gin.Context, config *modelquery.SysConfigDQL) (sysConfigList []*modelresponse.SysConfigVo, total int64)
	ConfigExport(c *gin.Context, config *modelquery.SysConfigDQL) (data []byte)
	SelectConfigById(c *gin.Context, configId int64) (Config *modelresponse.SysConfigVo)
	InsertConfig(c *gin.Context, config *modelrequest.SysConfigDML)
	UpdateConfig(c *gin.Context, config *modelrequest.SysConfigDML)
	DeleteConfigById(c *gin.Context, configId int64)
	CheckConfigKeyUnique(c *gin.Context, configId int64, configKey string) bool
	SelectConfigValueByKey(c *gin.Context, configKey string) string
}
