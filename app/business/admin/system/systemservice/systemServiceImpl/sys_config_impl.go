package systemServiceImpl

import (
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/admin/system/systemmodels"
	"nova-factory-server/app/business/admin/system/systemservice"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/excel"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type ConfigService struct {
	cache cache.Cache
	cd    systemdao.IConfigDao
}

func NewConfigService(cd systemdao.IConfigDao,
	cache cache.Cache) systemservice.IConfigService {
	return &ConfigService{cache: cache, cd: cd}
}

func (cs *ConfigService) SelectConfigList(c *gin.Context, config *systemmodels.SysConfigDQL) (sysConfigList []*systemmodels.SysConfigVo, total int64) {
	return cs.cd.SelectConfigList(c, config)
}
func (cs *ConfigService) ConfigExport(c *gin.Context, config *systemmodels.SysConfigDQL) (data []byte) {
	list := cs.cd.SelectConfigListAll(c, config)
	toExcel, err := excel.SliceToExcel(list)
	if err != nil {
		panic(err)
	}
	buffer, err := toExcel.WriteToBuffer()
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func (cs *ConfigService) SelectConfigById(c *gin.Context, configId int64) (Config *systemmodels.SysConfigVo) {
	return cs.cd.SelectConfigById(c, configId)
}

func (cs *ConfigService) InsertConfig(c *gin.Context, config *systemmodels.SysConfigVo) {
	config.ConfigId = snowflake.GenID()
	cs.cd.InsertConfig(c, config)
}

func (cs *ConfigService) UpdateConfig(c *gin.Context, config *systemmodels.SysConfigVo) {
	cs.cd.UpdateConfig(c, config)
	cs.cache.Del(c, cs.getCacheKey(config.ConfigKey))
}

func (cs *ConfigService) DeleteConfigById(c *gin.Context, configId int64) {
	key := cs.cd.SelectConfigById(c, configId).ConfigKey
	cs.cd.DeleteConfigById(c, configId)
	cs.cache.Del(c, cs.getCacheKey(key))
}

func (cs *ConfigService) CheckConfigKeyUnique(c *gin.Context, configId int64, configKey string) bool {
	id := cs.cd.SelectConfigIdByConfigKey(c, configKey)
	if id == configId {
		return false
	}
	return true
}

func (cs *ConfigService) SelectConfigValueByKey(c *gin.Context, configKey string) string {
	v, err := cs.cache.Get(c, cs.getCacheKey(configKey))
	if err == nil {
		return v
	}

	value := cs.cd.SelectConfigValueByConfigKey(c, configKey)
	if value != "" {
		cs.cache.Set(c, cs.getCacheKey(configKey), value, 0)
	}
	return value
}

func (cs *ConfigService) getCacheKey(configKey string) string {
	return "sys_config:" + configKey

}
