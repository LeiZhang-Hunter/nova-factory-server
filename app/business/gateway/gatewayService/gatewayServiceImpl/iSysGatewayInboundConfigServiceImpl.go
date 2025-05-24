package gatewayServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/gateway/gatewayDao"
	"nova-factory-server/app/business/gateway/gatewayModels"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
)

type ISysGatewayInboundConfigServiceImpl struct {
	dao      gatewayDao.ISysGatewayInboundConfigDao
	iUserDao systemDao.IUserDao
}

func NewISysGatewayInboundConfigServiceImpl(dao gatewayDao.ISysGatewayInboundConfigDao, iUserDao systemDao.IUserDao) *ISysGatewayInboundConfigServiceImpl {
	return &ISysGatewayInboundConfigServiceImpl{
		dao:      dao,
		iUserDao: iUserDao,
	}
}

func (i *ISysGatewayInboundConfigServiceImpl) Add(c *gin.Context, config *gatewayModels.SysSetGatewayInboundConfig) (*gatewayModels.SysGatewayInboundConfig, error) {
	return i.dao.Add(c, config)
}

func (i *ISysGatewayInboundConfigServiceImpl) Update(c *gin.Context, config *gatewayModels.SysSetGatewayInboundConfig) (*gatewayModels.SysGatewayInboundConfig, error) {
	return i.dao.Update(c, config)
}

func (i *ISysGatewayInboundConfigServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}

func (i *ISysGatewayInboundConfigServiceImpl) List(c *gin.Context, req *gatewayModels.SysSetGatewayInboundConfigReq) (*gatewayModels.SysSetGatewayInboundConfigList, error) {
	list, err := i.dao.List(c, req)
	if err != nil {
		return list, err
	}
	if len(list.Rows) == 0 {
		return list, nil
	}

	//  读取用户id集合
	userIdMap := make(map[int64]bool)
	for _, v := range list.Rows {
		if v.CreateBy > 0 {
			userIdMap[v.CreateBy] = true
		}

		if v.UpdateBy > 0 {
			userIdMap[v.UpdateBy] = true
		}
	}

	// 格式化服务id
	userIds := make([]int64, 0)
	for k, _ := range userIdMap {
		if k > 0 {
			userIds = append(userIds, k)
		}
	}

	users := i.iUserDao.SelectByUserIds(c, userIds)
	userVoMap := make(map[int64]*systemModels.SysUserDML)
	for _, v := range users {
		userVoMap[v.UserId] = v
	}

	for k, v := range list.Rows {
		var createUserName string
		var updateUserName string
		userVo, ok := userVoMap[v.CreateBy]
		if ok {
			createUserName = userVo.UserName
		}

		userVo, ok = userVoMap[v.UpdateBy]
		if ok {
			updateUserName = userVo.UserName
		}

		list.Rows[k].CreateUserName = createUserName
		list.Rows[k].UpdateUserName = updateUserName
	}
	return nil, nil
}
