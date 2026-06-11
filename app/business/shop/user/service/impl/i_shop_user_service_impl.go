package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"strings"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/utils/bCryptPasswordEncoder"
	"strconv"
)

type ShopUserServiceImpl struct {
	dao   dao.IShopUserDao
	cache cache.Cache
}

func NewShopUserService(cache cache.Cache, dao dao.IShopUserDao) service.IShopUserService {
	return &ShopUserServiceImpl{dao: dao, cache: cache}
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
	user, err := s.dao.GetByID(c, id)
	if err != nil || user == nil {
		return user, err
	}
	user.IsOnline = s.isUserOnline(c, user.ID)
	return user, nil
}

func (s *ShopUserServiceImpl) List(c *gin.Context, req *models.UserQuery) (*models.UserListData, error) {
	data, err := s.dao.List(c, req)
	if err != nil || data == nil {
		return data, err
	}
	s.fillOnlineStatus(c, data.Rows)
	return data, nil
}

func (s *ShopUserServiceImpl) fillOnlineStatus(c *gin.Context, rows []*models.User) {
	if len(rows) == 0 {
		return
	}
	onlineUserIDs := s.onlineUserIDSet(c)
	for _, row := range rows {
		if row == nil {
			continue
		}
		_, row.IsOnline = onlineUserIDs[row.ID]
	}
}

func (s *ShopUserServiceImpl) isUserOnline(c *gin.Context, userID int64) bool {
	if userID == 0 {
		return false
	}
	_, ok := s.onlineUserIDSet(c)[userID]
	return ok
}

func (s *ShopUserServiceImpl) onlineUserIDSet(c *gin.Context) map[int64]struct{} {
	onlineUserIDs := make(map[int64]struct{})
	if s.cache == nil {
		return onlineUserIDs
	}
	manager := session.NewShopManager(s.cache)
	for _, sess := range manager.ScanSessions(c) {
		userID, err := strconv.ParseInt(sess.Get(c, sessionStatus.UserId), 10, 64)
		if err != nil || userID == 0 {
			continue
		}
		onlineUserIDs[userID] = struct{}{}
	}
	return onlineUserIDs
}

func (s *ShopUserServiceImpl) prepareUpsert(c *gin.Context, req *models.UserUpsert, isUpdate bool) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if isUpdate && req.ID == 0 {
		return errors.New("用户ID不能为空")
	}
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
