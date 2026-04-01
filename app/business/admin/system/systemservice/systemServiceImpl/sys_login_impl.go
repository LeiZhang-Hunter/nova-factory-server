package systemServiceImpl

import (
	"nova-factory-server/app/business/admin/monitor/monitordao"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	systemDao2 "nova-factory-server/app/business/admin/system/systemdao"
	systemModels2 "nova-factory-server/app/business/admin/system/systemmodels"
	systemService2 "nova-factory-server/app/business/admin/system/systemservice"
	"nova-factory-server/app/constant/dataScopeAspect"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"nova-factory-server/app/utils/bCryptPasswordEncoder"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"time"

	"image/color"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type LoginService struct {
	cache       cache.Cache
	userDao     systemDao2.IUserDao
	pd          systemDao2.IPermissionDao
	roleDao     systemDao2.IRoleDao
	loginforDao monitordao.ILogininforDao
	driver      *base64Captcha.DriverMath
	store       base64Captcha.Store
	cs          systemService2.IConfigService
}

func NewLoginService(cache cache.Cache, ud systemDao2.IUserDao, pd systemDao2.IPermissionDao, rd systemDao2.IRoleDao, ld monitordao.ILogininforDao, cs systemService2.IConfigService) systemService2.ILoginService {
	return &LoginService{cache: cache, userDao: ud, pd: pd, roleDao: rd, loginforDao: ld, cs: cs,
		driver: base64Captcha.NewDriverMath(38, 106, 0, 0, &color.RGBA{0, 0, 0, 0}, nil, []string{"wqy-microhei.ttc"}),
		store:  base64Captcha.DefaultMemStore,
	}
}

func (loginService *LoginService) Login(c *gin.Context, user *systemModels2.User) string {
	manager := session.NewManger(loginService.cache)
	session, _ := manager.InitSession(c, user.UserId)
	session.Set(c, sessionStatus.Os, user.Os)
	session.Set(c, sessionStatus.Browser, user.Browser)
	session.Set(c, sessionStatus.UserName, user.UserName)
	session.Set(c, sessionStatus.Avatar, user.Avatar)
	return session.Id()
}

func (loginService *LoginService) Register(c *gin.Context, user *systemModels2.LoginBody) {
	u := new(systemModels2.SysUserDML)
	u.Password = bCryptPasswordEncoder.HashPassword(user.Password)
	u.DataScope = dataScopeAspect.NoDataScope
	u.UserId = snowflake.GenID()
	u.NickName = user.Username
	u.UserName = user.Username
	u.Status = "0"
	u.DeptId = 100
	u.SetCreateBy(u.UserId)
	loginService.userDao.InsertUser(c, u)
}

func (loginService *LoginService) RecordLoginInfo(c *gin.Context, loginUser *monitormodels.Logininfor) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("登录日志记录错误", zap.Any("error", err))
			}
		}()
		loginUser.InfoId = snowflake.GenID()
		loginUser.LoginTime = time.Now()
		loginService.loginforDao.InserLogininfor(c, loginUser)
	}()

}

func (loginService *LoginService) getPermission(c *gin.Context, userId int64) []string {
	perms := make([]string, 0)
	if baizeContext.IsAdmin(c) {
		perms = loginService.pd.SelectPermissionAll(c)
	} else {
		perms = loginService.pd.SelectPermissionByUserId(c, userId)
	}
	return perms
}

func (loginService *LoginService) GenerateCode(c *gin.Context) (m *systemModels2.CaptchaVo) {
	m = new(systemModels2.CaptchaVo)
	key := loginService.cs.SelectConfigValueByKey(c, "sys.account.captchaEnabled")
	if key != "false" {
		captcha := base64Captcha.NewCaptcha(loginService.driver, loginService.store)
		id, b64s, _, err := captcha.Generate()
		if err != nil {
			panic(err)
		}
		m.Id = id
		m.Img = b64s
		m.CaptchaEnabled = true
	}
	key = loginService.cs.SelectConfigValueByKey(c, "sys.account.registerUser")
	if key == "true" {
		m.RegisterEnabled = true
	}
	return m
}

func (loginService *LoginService) VerityCaptcha(c *gin.Context, id, base64 string) bool {
	return loginService.store.Verify(id, base64, true)
}

func (loginService *LoginService) ForceLogout(c *gin.Context, token string) {
	panic("等待补充")
}

func (loginService *LoginService) RolePermissionByRoles(roles []*systemModels2.SysRole) (loginRoles []int64) {
	loginRoles = make([]int64, 0, len(roles))
	for _, role := range roles {
		loginRoles = append(loginRoles, role.RoleId)
	}
	return
}
func (loginService *LoginService) GetInfo(c *gin.Context) *systemModels2.GetInfo {
	userId := baizeContext.GetUserId(c)
	roles := loginService.roleDao.SelectBasicRolesByUserId(c, userId)
	loginRoles := loginService.RolePermissionByRoles(roles)

	session := baizeContext.GetSession(c)
	session.Set(c, sessionStatus.Role, loginRoles)
	permission := loginService.getPermission(c, userId)
	session.Set(c, sessionStatus.Permission, permission)
	session.Set(c, sessionStatus.IpAddr, c.ClientIP())
	session.Set(c, sessionStatus.LoginTime, time.Now().Unix())
	user := loginService.userDao.SelectUserById(c, userId)
	session.Set(c, sessionStatus.UserName, user.UserName)
	session.Set(c, sessionStatus.Avatar, user.Avatar)
	session.Set(c, sessionStatus.DeptId, user.DeptId)
	session.Set(c, sessionStatus.DataScopeAspect, user.DataScope)
	getInfo := new(systemModels2.GetInfo)
	u := new(systemModels2.User)
	u.UserId = baizeContext.GetUserId(c)
	u.UserName = baizeContext.GetUserName(c)
	u.Avatar = baizeContext.GetAvatar(c)
	getInfo.User = u
	//getInfo.Roles = loginRoles
	getInfo.Permissions = permission
	return getInfo
}
