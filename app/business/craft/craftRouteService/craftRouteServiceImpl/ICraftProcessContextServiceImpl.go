package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
)

type ICraftProcessContextServiceImpl struct {
	dao      craftRouteDao.IProcessContextDao
	iUserDao systemDao.IUserDao
}

func NewICraftProcessContextServiceImpl(dao craftRouteDao.IProcessContextDao, iUserDao systemDao.IUserDao) craftRouteService.ICraftProcessContextService {
	return &ICraftProcessContextServiceImpl{
		dao:      dao,
		iUserDao: iUserDao,
	}
}

func (i *ICraftProcessContextServiceImpl) Add(c *gin.Context, processContext *craftRouteModels.SysProSetProcessContent) (*craftRouteModels.SysProProcessContent, error) {
	return i.dao.Add(c, processContext)
}

func (i *ICraftProcessContextServiceImpl) Update(c *gin.Context, processContext *craftRouteModels.SysProSetProcessContent) (*craftRouteModels.SysProProcessContent, error) {
	return i.dao.Update(c, processContext)
}

func (i *ICraftProcessContextServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}

func (i *ICraftProcessContextServiceImpl) List(c *gin.Context, req *craftRouteModels.SysProProcessContextListReq) (*craftRouteModels.SysProProcessContextListData, error) {
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
	return list, nil
}
