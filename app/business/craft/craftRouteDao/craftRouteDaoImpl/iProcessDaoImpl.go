package craftRouteDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IProcessDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIProcessDaoImpl(db *gorm.DB) craftRouteDao.IProcessDao {
	return &IProcessDaoImpl{
		db:        db,
		tableName: "sys_pro_process",
	}
}

func (i *IProcessDaoImpl) Add(c *gin.Context, process *craftRouteModels.SysProProcess) (*craftRouteModels.SysProProcess, error) {
	process.ProcessID = snowflake.GenID()
	process.SetCreateBy(baizeContext.GetUserId(c))
	ret := i.db.Table(i.tableName).Create(process)
	return process, ret.Error
}

func (i *IProcessDaoImpl) Update(c *gin.Context, process *craftRouteModels.SysProProcess) (*craftRouteModels.SysProProcess, error) {
	process.SetUpdateBy(baizeContext.GetUserId(c))
	ret := i.db.Table(i.tableName).Where("process_id = ?", process.ProcessID).Updates(process)
	return process, ret.Error
}

func (i *IProcessDaoImpl) Remove(c *gin.Context, processIds []int64) error {
	ret := i.db.Table(i.tableName).Where("process_id in (?)", processIds).Update("state", commonStatus.DELETE)
	return ret.Error
}
