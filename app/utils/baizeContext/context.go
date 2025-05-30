package baizeContext

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/constant/dataScopeAspect"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/middlewares/session/sessionCache"
	"strconv"
)

func IsAdmin(c *gin.Context) bool {
	for _, id := range GetRoles(c) {
		if id == 1 {
			return true
		}
	}
	return false
}
func GetRoles(c *gin.Context) []int64 {
	get := GetSession(c).Get(c, sessionStatus.Role)
	roles := make([]int64, 0)
	err := json.Unmarshal([]byte(get), &roles)
	if err != nil {
		panic(err)
	}
	return roles
}
func GetRolesPerms(c *gin.Context) []string {
	get := GetSession(c).Get(c, sessionStatus.RolePerms)
	roles := make([]string, 0)
	err := json.Unmarshal([]byte(get), &roles)
	if err != nil {
		panic(err)
	}
	return roles
}

func GetSession(c *gin.Context) *sessionCache.Session {
	val, ok := c.Get(sessionStatus.SessionKey)
	if ok {
		return val.(*sessionCache.Session)
	}
	panic("不应该出现")
}
func GetPermission(c *gin.Context) []string {
	get := GetSession(c).Get(c, sessionStatus.Permission)
	permission := make([]string, 0)
	err := json.Unmarshal([]byte(get), &permission)
	if err != nil {
		panic(err)
	}
	return permission
}

func GetUserId(c *gin.Context) int64 {
	i, err := strconv.ParseInt(GetSession(c).Get(c, sessionStatus.UserId), 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
func GetDeptId(c *gin.Context) int64 {
	i, err := strconv.ParseInt(GetSession(c).Get(c, sessionStatus.DeptId), 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func GetDataScopeAspect(c *gin.Context) string {
	return GetSession(c).Get(c, sessionStatus.DataScopeAspect)
}

func GetUserName(c *gin.Context) string {
	return GetSession(c).Get(c, sessionStatus.UserName)
}
func GetAvatar(c *gin.Context) string {
	return GetSession(c).Get(c, sessionStatus.Avatar)
}
func GetDataScope(c *gin.Context, deptAlias string) string {

	var sqlString string
	switch GetDataScopeAspect(c) {
	case dataScopeAspect.DataScopeAll:
		sqlString = ""
	case dataScopeAspect.DataScopeCustom:
		sqlString += fmt.Sprintf(" %s.dept_id IN ( SELECT dept_id FROM sys_user_dept_scope WHERE user_id = %d ) ", deptAlias, GetUserId(c))
	case dataScopeAspect.DataScopeDept:
		sqlString += fmt.Sprintf(" %s.dept_id = %d ", deptAlias, GetDeptId(c))
	case dataScopeAspect.DataScopeDeptAndChild:
		sqlString += fmt.Sprintf(" %s.dept_id IN ( SELECT dept_id FROM sys_dept WHERE dept_id = %d or find_in_set( %d , ancestors ) ) ", deptAlias, GetDeptId(c), GetDeptId(c))
	default:
		sqlString += fmt.Sprintf(" 1=0")
	}

	return sqlString
}

func GetGormDataScope(c *gin.Context, db *gorm.DB) *gorm.DB {
	switch GetDataScopeAspect(c) {
	case dataScopeAspect.DataScopeAll:
	case dataScopeAspect.DataScopeCustom:
		db.Where(fmt.Sprintf(" %s.dept_id IN ( SELECT dept_id FROM sys_user_dept_scope WHERE user_id = ? ) ", db.Statement.Table), GetUserId(c))
	case dataScopeAspect.DataScopeDept:
		db.Where(fmt.Sprintf("  %s.dept_id = ? ", db.Statement.Table), GetDeptId(c))
	case dataScopeAspect.DataScopeDeptAndChild:
		db.Where(fmt.Sprintf(" %s.dept_id IN ( SELECT dept_id FROM sys_dept WHERE dept_id = ? or find_in_set( ? , ancestors ) ) ", db.Statement.Table), GetDeptId(c), GetDeptId(c))
	default:
		db.Where(" 1=0")
	}
	return db
}
