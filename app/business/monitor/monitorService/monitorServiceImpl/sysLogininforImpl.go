package monitorServiceImpl

import (
	"context"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/monitor/monitorDao"
	"nova-factory-server/app/business/monitor/monitorModels"
	"nova-factory-server/app/business/monitor/monitorService"
	"nova-factory-server/app/utils/excel"
	"nova-factory-server/app/utils/snowflake"
)

type LogininforService struct {
	ld monitorDao.ILogininforDao
}

func NewLogininforService(ld monitorDao.ILogininforDao) monitorService.ILogininforService {
	return &LogininforService{ld: ld}
}

func (ls *LogininforService) SelectLogininforList(c *gin.Context, logininfor *monitorModels.LogininforDQL) (list []*monitorModels.Logininfor, total int64) {
	return ls.ld.SelectLogininforList(c, logininfor)

}
func (ls *LogininforService) ExportLogininfor(c *gin.Context, logininfor *monitorModels.LogininforDQL) (data []byte) {
	list := ls.ld.SelectLogininforListAll(c, logininfor)
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

func (ls *LogininforService) InsertLogininfor(c context.Context, loginUser *monitorModels.Logininfor) {
	loginUser.InfoId = snowflake.GenID()
	ls.ld.InserLogininfor(c, loginUser)
}

func (ls *LogininforService) DeleteLogininforByIds(c *gin.Context, infoIds []int64) {
	ls.ld.DeleteLogininforByIds(c, infoIds)

}

func (ls *LogininforService) CleanLogininfor(c *gin.Context) {
	ls.ld.CleanLogininfor(c)

}
