package importer

import (
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	"strconv"

	"nova-factory-server/app/constant/dataScopeAspect"
	"nova-factory-server/app/utils/snowflake"
)

func RowsToSysUserDMLList(rows [][]string, str string, failureNum int, dept map[string]int64, password string, userId int64) ([]*modelrequest.SysUserDML, string, int) {
	list := make([]*modelrequest.SysUserDML, 0, len(rows)-1)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if row[0] == "" {
			str += "<br/>第" + strconv.Itoa(i+1) + "行用户名为空"
			failureNum++
			continue
		}
		sysUser := new(modelrequest.SysUserDML)
		sysUser.UserId = snowflake.GenID()
		sysUser.UserName = row[0]
		sysUser.NickName = row[1]
		sysUser.DeptId = dept[row[2]]
		if sysUser.DeptId == 0 {
			str += "<br/>第" + strconv.Itoa(i+1) + "部门错误"
			failureNum++
			continue
		}
		sysUser.Email = row[3]
		sysUser.Phonenumber = row[4]
		sex := row[4]
		if sex == "男" {
			sysUser.Sex = "0"
		} else if sex == "女" {
			sysUser.Sex = "1"
		} else {
			sysUser.Sex = "2"
		}
		sysUser.Status = "0"
		sysUser.Password = password
		sysUser.DataScope = dataScopeAspect.NoDataScope
		sysUser.SetCreateBy(userId)
		list = append(list, sysUser)
	}
	return list, str, failureNum
}
