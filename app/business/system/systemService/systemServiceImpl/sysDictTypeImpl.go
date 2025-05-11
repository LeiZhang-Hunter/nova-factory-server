package systemServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/excel"
	"nova-factory-server/app/utils/snowflake"
)

type DictTypeService struct {
	cache       cache.Cache
	dictTypeDao systemDao.IDictTypeDao
	dictKey     string
}

func NewDictTypeService(dtd systemDao.IDictTypeDao,
	cache cache.Cache) systemService.IDictTypeService {
	return &DictTypeService{
		cache:       cache,
		dictTypeDao: dtd,
		dictKey:     "sys_dict:",
	}
}

func (dictTypeService *DictTypeService) SelectDictTypeList(c *gin.Context, dictType *systemModels.SysDictTypeDQL) (list []*systemModels.SysDictTypeVo, total int64) {
	return dictTypeService.dictTypeDao.SelectDictTypeList(c, dictType)

}
func (dictTypeService *DictTypeService) ExportDictType(c *gin.Context, dictType *systemModels.SysDictTypeDQL) (data []byte) {
	list := dictTypeService.dictTypeDao.SelectDictTypeAll(c, dictType)
	toExcel, err := excel.SliceToExcel(list)
	if err != nil {
		panic(err)
	}
	buffer, err := toExcel.WriteToBuffer()
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func (dictTypeService *DictTypeService) SelectDictTypeById(c *gin.Context, dictId int64) (dictType *systemModels.SysDictTypeVo) {
	return dictTypeService.dictTypeDao.SelectDictTypeById(c, dictId)

}
func (dictTypeService *DictTypeService) SelectDictTypeByIds(c *gin.Context, dictId []int64) (dictTypes []string) {
	return dictTypeService.dictTypeDao.SelectDictTypeByIds(c, dictId)
}

func (dictTypeService *DictTypeService) InsertDictType(c *gin.Context, dictType *systemModels.SysDictTypeVo) {
	dictType.DictId = snowflake.GenID()
	dictTypeService.dictTypeDao.InsertDictType(c, dictType)
}

func (dictTypeService *DictTypeService) UpdateDictType(c *gin.Context, dictType *systemModels.SysDictTypeVo) {
	dictTypeService.dictTypeDao.UpdateDictType(c, dictType)
}
func (dictTypeService *DictTypeService) DeleteDictTypeByIds(c *gin.Context, dictIds []int64) {
	dictTypeService.dictTypeDao.DeleteDictTypeByIds(c, dictIds)
}

func (dictTypeService *DictTypeService) CheckDictTypeUnique(c *gin.Context, id int64, dictType string) bool {
	dictId := dictTypeService.dictTypeDao.CheckDictTypeUnique(c, dictType)
	if dictId == id || dictId == 0 {
		return false
	}
	return true
}
func (dictTypeService *DictTypeService) DictTypeClearCache(c *gin.Context) {
	dictTypeService.cache.Del(c, dictTypeService.dictKey+"*")
}
func (dictTypeService *DictTypeService) SelectDictTypeAll(c *gin.Context) (list []*systemModels.SysDictTypeVo) {
	return dictTypeService.dictTypeDao.SelectDictTypeAll(c, new(systemModels.SysDictTypeDQL))
}
