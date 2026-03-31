package productServiceImpl

import (
	"nova-factory-server/app/business/admin/product/productdao"
	"nova-factory-server/app/business/admin/product/productmodels"
	"nova-factory-server/app/business/admin/product/productservice"
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type SysProductLaboratoryService struct {
	dao      productdao.ISysProductLaboratoryDao
	iUserDao systemdao.IUserDao
}

func NewSysProductLaboratoryService(dao productdao.ISysProductLaboratoryDao, iUserDao systemdao.IUserDao) productservice.ISysProductLaboratoryService {
	return &SysProductLaboratoryService{
		dao:      dao,
		iUserDao: iUserDao,
	}
}

func (s *SysProductLaboratoryService) SelectLaboratoryList(c *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (*productmodels.SysProductLaboratoryList, error) {
	list, err := s.dao.SelectLaboratoryList(c, dql)
	if err != nil {
		return nil, err
	}

	if list == nil {
		return &productmodels.SysProductLaboratoryList{
			Rows: make([]*productmodels.SysProductLaboratory, 0),
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
	users := s.iUserDao.SelectByUserIds(c, userIds)
	userVoMap := make(map[int64]*systemmodels.SysUserDML)
	for _, v := range users {
		userVoMap[v.UserId] = v
	}

	for k, v := range list.Rows {
		var createUserName string
		var updateUserName string
		userVo, ok := userVoMap[v.CreateBy]
		if ok {
			createUserName = userVo.NickName
		}

		userVo, ok = userVoMap[v.UpdateBy]
		if ok {
			updateUserName = userVo.NickName
		}
		list.Rows[k].CreateUserName = createUserName
		list.Rows[k].UpdateUserName = updateUserName
	}

	return list, err
}

func (s *SysProductLaboratoryService) SelectLaboratoryById(c *gin.Context, id int64) (*productmodels.SysProductLaboratoryVo, error) {
	return s.dao.SelectLaboratoryById(c, id)
}

func (s *SysProductLaboratoryService) Set(c *gin.Context, data *productmodels.SysProductLaboratoryVo) (*productmodels.SysProductLaboratory, error) {
	return s.dao.Set(c, data)
}

func (s *SysProductLaboratoryService) DeleteLaboratoryByIds(c *gin.Context, ids []int64) error {
	err := s.dao.DeleteLaboratoryByIds(c, ids)
	return err
}

// SelectUserLaboratoryList 读取用户化验单
func (s *SysProductLaboratoryService) SelectUserLaboratoryList(ctx *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (list *productmodels.SysProductLaboratoryList, err error) {
	return s.dao.SelectUserLaboratoryList(ctx, dql)
}

// FirstLaboratoryInfo 读取用户化验单
func (s *SysProductLaboratoryService) FirstLaboratoryInfo(ctx *gin.Context, req *productmodels.SysProductLaboratoryInfoDQL) (*productmodels.SysProductLaboratory, error) {
	return s.dao.FirstLaboratoryInfo(ctx, req)
}

func (s *SysProductLaboratoryService) FirstLaboratoryList(ctx *gin.Context, dql *productmodels.SysProductLaboratoryDQL) (*productmodels.SysProductLaboratoryList, error) {
	list, err := s.dao.FirstLaboratoryList(ctx, dql)
	if err != nil {
		return nil, err
	}
	if list == nil {
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
	users := s.iUserDao.SelectByUserIds(ctx, userIds)
	userVoMap := make(map[int64]*systemmodels.SysUserDML)
	for _, v := range users {
		userVoMap[v.UserId] = v
	}

	for k, v := range list.Rows {
		var createUserName string
		var updateUserName string
		userVo, ok := userVoMap[v.CreateBy]
		if ok {
			createUserName = userVo.NickName
		}

		userVo, ok = userVoMap[v.UpdateBy]
		if ok {
			updateUserName = userVo.NickName
		}
		list.Rows[k].CreateUserName = createUserName
		list.Rows[k].UpdateUserName = updateUserName
	}
	return list, nil
}
