package craftRouteDaoImpl

import "nova-factory-server/app/business/craft/craftRouteDao"

type SysProRouteProductBomDaoImpl struct {
}

func NewSysProRouteProductBomDaoImpl() craftRouteDao.SysProRouteProductBomDao {
	return &SysProRouteProductBomDaoImpl{}
}
