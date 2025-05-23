package systemServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
	"nova-factory-server/app/utils/snowflake"
	"strconv"
)

type DeptService struct {
	deptDao systemDao.IDeptDao
	roleDao systemDao.IRoleDao
}

func NewDeptService(dd systemDao.IDeptDao, rd systemDao.IRoleDao) systemService.IDeptService {
	return &DeptService{deptDao: dd, roleDao: rd}
}

func (ds *DeptService) SelectDeptList(c *gin.Context, dept *systemModels.SysDeptDQL) (list []*systemModels.SysDeptVo) {
	return ds.deptDao.SelectDeptList(c, dept)

}

func (ds *DeptService) SelectDeptById(c *gin.Context, deptId int64) (dept *systemModels.SysDeptVo) {
	return ds.deptDao.SelectDeptById(c, deptId)

}

func (ds *DeptService) InsertDept(c *gin.Context, dept *systemModels.SysDeptVo) {
	//获取上级部门祖籍信息
	parentDept := ds.SelectDeptById(c, dept.ParentId)
	dept.Ancestors = parentDept.Ancestors + "," + strconv.FormatInt(dept.ParentId, 10)
	//执行添加
	dept.DeptId = snowflake.GenID()
	ds.deptDao.InsertDept(c, dept)
	return
}

func (ds *DeptService) UpdateDept(c *gin.Context, dept *systemModels.SysDeptVo) {
	ds.deptDao.UpdateDept(c, dept)
}
func (ds *DeptService) DeleteDeptById(c *gin.Context, dept int64) {
	ds.deptDao.DeleteDeptById(c, dept)
	return
}
func (ds *DeptService) CheckDeptNameUnique(c *gin.Context, id, parentId int64, deptName string) bool {
	deptId := ds.deptDao.CheckDeptNameUnique(c, deptName, parentId)
	if deptId == id || deptId == 0 {
		return false
	}
	return true
}
func (ds *DeptService) HasChildByDeptId(c *gin.Context, deptId int64) bool {
	return ds.deptDao.HasChildByDeptId(c, deptId) > 0
}

func (ds *DeptService) CheckDeptExistUser(c *gin.Context, deptId int64) bool {
	return ds.deptDao.CheckDeptExistUser(c, deptId) > 0
}
