package systemservice

import (
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type IConfigService interface {
	SelectConfigList(c *gin.Context, config *systemmodels.SysConfigDQL) (sysConfigList []*systemmodels.SysConfigVo, total int64)
	ConfigExport(c *gin.Context, config *systemmodels.SysConfigDQL) (data []byte)
	SelectConfigById(c *gin.Context, configId int64) (Config *systemmodels.SysConfigVo)
	InsertConfig(c *gin.Context, config *systemmodels.SysConfigVo)
	UpdateConfig(c *gin.Context, config *systemmodels.SysConfigVo)
	DeleteConfigById(c *gin.Context, configId int64)
	CheckConfigKeyUnique(c *gin.Context, configId int64, configKey string) bool
	SelectConfigValueByKey(c *gin.Context, configKey string) string
}
