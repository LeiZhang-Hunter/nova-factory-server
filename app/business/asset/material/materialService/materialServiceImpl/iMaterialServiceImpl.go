package materialServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/asset/material/materialDao"
	"nova-factory-server/app/business/asset/material/materialModels"
	"nova-factory-server/app/business/asset/material/materialService"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
)

type MaterialService struct {
	iMaterialDao materialDao.IMaterialDao
	iUserDao     systemDao.IUserDao
}

func NewMaterialService(iMaterialDao materialDao.IMaterialDao, iUserDao systemDao.IUserDao) materialService.IMaterialService {
	return &MaterialService{
		iMaterialDao: iMaterialDao,
		iUserDao:     iUserDao,
	}
}

func (m *MaterialService) InsertMaterial(c *gin.Context, value *materialModels.MaterialInfo) (*materialModels.MaterialVO, error) {
	info, err := m.iMaterialDao.GetMaterialGroupByName(c, value.Name)
	if info != nil {
		return nil, errors.New("物料已经存在")
	}
	if err != nil {
		zap.L().Error("读取物料名字失败", zap.Error(err))
	}

	vo, err := m.iMaterialDao.InsertMaterial(c, value)
	if err != nil {
		zap.L().Error("读取物料名字失败", zap.Error(err))
		return nil, errors.New("添加物料名字失败")
	}
	return vo, nil
}

func (m *MaterialService) UpdateMaterial(c *gin.Context, value *materialModels.MaterialInfo) (*materialModels.MaterialVO, error) {
	info, err := m.iMaterialDao.GetNoExitIdMaterialGroupByName(c, value.Name, value.MaterialId)
	if info != nil {
		return nil, errors.New("物料名已经存在")
	}
	if err != nil {
		zap.L().Error("读取物料名字失败", zap.Error(err))
	}
	return m.iMaterialDao.UpdateMaterial(c, value)
}

func (m *MaterialService) SelectMaterialList(c *gin.Context, req *materialModels.MaterialListReq) (*materialModels.MaterialInfoListValue, error) {
	list, err := m.iMaterialDao.SelectMaterialList(c, req)
	if err != nil {
		zap.L().Error("读取列表衰退", zap.Error(err))
		return &materialModels.MaterialInfoListValue{
			Rows:  make([]*materialModels.MaterialValue, 0),
			Total: 0,
		}, err
	}

	if len(list.Rows) == 0 {
		return &materialModels.MaterialInfoListValue{
			Rows:  make([]*materialModels.MaterialValue, 0),
			Total: 0,
		}, nil
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

	users := m.iUserDao.SelectByUserIds(c, userIds)
	userVoMap := make(map[int64]*systemModels.SysUserDML)
	for _, v := range users {
		userVoMap[v.UserId] = v
	}

	ret := make([]*materialModels.MaterialValue, 0)
	for _, v := range list.Rows {

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

		value := &materialModels.MaterialValue{
			MaterialId:     v.MaterialId,
			Name:           v.Name,
			Code:           v.Code,
			Model:          v.Model,
			Unit:           v.Unit,
			Factory:        v.Factory,
			Address:        v.Address,
			Price:          v.Price,
			Total:          v.Total,
			Outbound:       v.Outbound,
			CurrentTotal:   v.CurrentTotal,
			DeptId:         v.DeptId,
			State:          v.State,
			CreateUserName: createUserName,
			UpdateUserName: updateUserName,
			BaseEntity: baize.BaseEntity{
				CreateTime: v.CreateTime,
				UpdateTime: v.UpdateTime,
			},
		}

		ret = append(ret, value)
	}
	return &materialModels.MaterialInfoListValue{
		Rows:  ret,
		Total: list.Total,
	}, nil
}
func (m *MaterialService) DeleteByMaterialIds(c *gin.Context, ids []int64) error {
	return m.iMaterialDao.DeleteByMaterialIds(c, ids)
}

func (m *MaterialService) GetByMaterialId(c *gin.Context, id int64) (*materialModels.MaterialVO, error) {
	return m.iMaterialDao.GetByMaterialId(c, id)
}

func (m *MaterialService) Inbound(c *gin.Context, info *materialModels.InboundInfo) (*materialModels.InboundVO, error) {
	return m.iMaterialDao.Inbound(c, info)
}
func (m *MaterialService) Outbound(c *gin.Context, info *materialModels.OutboundInfo) (*materialModels.OutboundVO, error) {
	return m.iMaterialDao.Outbound(c, info)
}

func (m *MaterialService) InboundList(c *gin.Context, req *materialModels.InboundListReq) (*materialModels.InboundListData, error) {
	list, err := m.iMaterialDao.InboundList(c, req)
	if err != nil {
		return list, err
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

	users := m.iUserDao.SelectByUserIds(c, userIds)
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

	return list, err
}

func (m *MaterialService) OutboundList(c *gin.Context, req *materialModels.OutboundListReq) (*materialModels.OutboundListData, error) {
	list, err := m.iMaterialDao.OutboundList(c, req)
	if err != nil {
		return list, err
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

	users := m.iUserDao.SelectByUserIds(c, userIds)
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
	return list, err
}
