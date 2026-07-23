package systemServiceImpl

import (
	"nova-factory-server/app/baize"
	systemDao2 "nova-factory-server/app/business/admin/system/systemdao"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
	"nova-factory-server/app/business/admin/system/systemservice"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type SelectService struct {
	pd systemDao2.IPermissionDao
	dd systemDao2.IDeptDao
}

func NewSelectService(pd systemDao2.IPermissionDao, dd systemDao2.IDeptDao) systemservice.ISelectBoxService {
	return &SelectService{pd: pd, dd: dd}
}

func (cs *SelectService) SelectPermissionBox(c *gin.Context) (list []*modelresponse.SelectPermission) {
	return cs.pd.SelectPermissionListSelectBoxByPerm(c, baizeContext.GetPermission(c))
}

func (cs *SelectService) SelectDeptBox(c *gin.Context, be *baize.BaseEntityDQL) (list []*modelresponse.SelectDept) {
	return cs.dd.SelectDeptListSelectBox(c, be)
}
