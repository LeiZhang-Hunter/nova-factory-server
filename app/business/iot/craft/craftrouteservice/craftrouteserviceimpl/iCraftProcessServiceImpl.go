package craftrouteserviceimpl

import (
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/admin/system/systemmodels"
	"nova-factory-server/app/business/iot/craft/craftroutedao"
	"nova-factory-server/app/business/iot/craft/craftroutemodels"
	"nova-factory-server/app/business/iot/craft/craftrouteservice"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ICraftProcessServiceImpl struct {
	dao      craftroutedao.IProcessDao
	iUserDao systemdao.IUserDao
}

func NewICraftProcessServiceImpl(dao craftroutedao.IProcessDao, iUserDao systemdao.IUserDao) craftrouteservice.ICraftProcessService {
	return &ICraftProcessServiceImpl{
		dao:      dao,
		iUserDao: iUserDao,
	}
}

func (i *ICraftProcessServiceImpl) Add(c *gin.Context, process *craftroutemodels.SysProSetProcessReq) (*craftroutemodels.SysProProcess, error) {
	data, err := i.dao.Add(c, process)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (i *ICraftProcessServiceImpl) Update(c *gin.Context, process *craftroutemodels.SysProSetProcessReq) (*craftroutemodels.SysProProcess, error) {
	data, err := i.dao.Update(c, process)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (i *ICraftProcessServiceImpl) Remove(c *gin.Context, processIds []int64) error {
	return i.dao.Remove(c, processIds)
}

func (i *ICraftProcessServiceImpl) List(c *gin.Context, req *craftroutemodels.SysProProcessListReq) (*craftroutemodels.SysProProcessListData, error) {
	list, err := i.dao.List(c, req)
	if err != nil {
		zap.L().Error("读取工序列表失败", zap.Error(err))
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
	userVoMap := make(map[int64]*systemmodels.SysUserDML)
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

func (i *ICraftProcessServiceImpl) GetById(c *gin.Context, id int64) (*craftroutemodels.SysProProcess, error) {
	return i.dao.GetById(c, id)
}
