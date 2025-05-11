package systemServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
	"nova-factory-server/app/utils/baizeContext"
)

type SelectService struct {
	pd systemDao.IPermissionDao
	dd systemDao.IDeptDao
}

func NewSelectService(pd systemDao.IPermissionDao, dd systemDao.IDeptDao) systemService.ISelectBoxService {
	return &SelectService{pd: pd, dd: dd}
}

func (cs *SelectService) SelectPermissionBox(c *gin.Context) (list []*systemModels.SelectPermission) {
	return cs.pd.SelectPermissionListSelectBoxByPerm(c, baizeContext.GetPermission(c))
}

func (cs *SelectService) SelectDeptBox(c *gin.Context, be *baize.BaseEntityDQL) (list []*systemModels.SelectDept) {
	return cs.dd.SelectDeptListSelectBox(c, be)
}
