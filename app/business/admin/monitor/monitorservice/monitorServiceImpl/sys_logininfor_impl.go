package monitorServiceImpl

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitordao"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	"nova-factory-server/app/business/admin/monitor/monitorservice"
	"nova-factory-server/app/utils/excel"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type LogininforService struct {
	ld monitordao.ILogininforDao
}

func NewLogininforService(ld monitordao.ILogininforDao) monitorservice.ILogininforService {
	return &LogininforService{ld: ld}
}

func (ls *LogininforService) SelectLogininforList(c *gin.Context, logininfor *monitormodels.LogininforDQL) (list []*monitormodels.Logininfor, total int64) {
	return ls.ld.SelectLogininforList(c, logininfor)

}
func (ls *LogininforService) ExportLogininfor(c *gin.Context, logininfor *monitormodels.LogininforDQL) (data []byte) {
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

func (ls *LogininforService) InsertLogininfor(c context.Context, loginUser *monitormodels.Logininfor) {
	loginUser.InfoId = snowflake.GenID()
	ls.ld.InserLogininfor(c, loginUser)
}

func (ls *LogininforService) DeleteLogininforByIds(c *gin.Context, infoIds []int64) {
	ls.ld.DeleteLogininforByIds(c, infoIds)

}

func (ls *LogininforService) CleanLogininfor(c *gin.Context) {
	ls.ld.CleanLogininfor(c)

}
