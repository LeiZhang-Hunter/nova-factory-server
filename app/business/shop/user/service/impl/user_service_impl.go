package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"strconv"
	"strings"

	"nova-factory-server/app/utils/bCryptPasswordEncoder"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type ShopUserServiceImpl struct {
	dao dao.IShopUserDao
}

func NewShopUserService(dao dao.IShopUserDao) service.IShopUserService {
	return &ShopUserServiceImpl{dao: dao}
}

func (s *ShopUserServiceImpl) Create(c *gin.Context, req *models.UserUpsert) (*models.User, error) {
	if err := s.prepareUpsert(c, req, false); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *ShopUserServiceImpl) Update(c *gin.Context, req *models.UserUpsert) (*models.User, error) {
	if err := s.prepareUpsert(c, req, true); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

func (s *ShopUserServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopUserServiceImpl) GetByID(c *gin.Context, id int64) (*models.User, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopUserServiceImpl) List(c *gin.Context, req *models.UserQuery) (*models.UserListData, error) {
	return s.dao.List(c, req)
}

func (s *ShopUserServiceImpl) prepareUpsert(c *gin.Context, req *models.UserUpsert, isUpdate bool) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if isUpdate && req.ID == 0 {
		return errors.New("用户ID不能为空")
	}
	req.UserID = strings.TrimSpace(req.UserID)
	req.Username = strings.TrimSpace(req.Username)
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Mobile = strings.TrimSpace(req.Mobile)
	req.Email = strings.TrimSpace(req.Email)
	req.Avatar = strings.TrimSpace(req.Avatar)
	req.CompanyName = strings.TrimSpace(req.CompanyName)
	req.ContactName = strings.TrimSpace(req.ContactName)
	req.ContactPhone = strings.TrimSpace(req.ContactPhone)
	if req.Username == "" {
		return errors.New("用户名不能为空")
	}
	if req.UserType < 1 || req.UserType > 3 {
		return errors.New("用户类型仅支持 1代理商、2分销商、3工厂")
	}
	if req.Status == nil {
		req.Status = userStatusPtr(false)
	}
	if req.UserID == "" {
		req.UserID = strconv.FormatInt(snowflake.GenID(), 10)
	}
	if isUpdate {
		current, err := s.dao.GetByID(c, req.ID)
		if err != nil {
			return err
		}
		if current == nil {
			return errors.New("商城用户不存在")
		}
		if strings.TrimSpace(req.Password) == "" {
			req.Password = current.Password
		} else {
			req.Password = bCryptPasswordEncoder.HashPassword(strings.TrimSpace(req.Password))
		}
		return nil
	}
	if strings.TrimSpace(req.Password) != "" {
		req.Password = bCryptPasswordEncoder.HashPassword(strings.TrimSpace(req.Password))
	}
	return nil
}

func userStatusPtr(v bool) *bool {
	return &v
}
