package systemServiceImpl

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/excel"
	"nova-factory-server/app/utils/response"
	"nova-factory-server/app/utils/snowflake"
	"time"
)

type DictDataService struct {
	cache       cache.Cache
	dictDataDao systemDao.IDictDataDao
	dictKey     string
	gzipNil     []byte
}

func NewDictDataService(ddd systemDao.IDictDataDao, cache cache.Cache) systemService.IDictDataService {
	return &DictDataService{
		cache:       cache,
		dictDataDao: ddd,
		dictKey:     "sys_dict:",
		gzipNil:     []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 170, 86, 74, 206, 79, 73, 85, 178, 50, 50, 48, 208, 81, 202, 45, 78, 87, 178, 82, 42, 46, 77, 78, 78, 45, 46, 86, 170, 5, 4, 0, 0, 255, 255, 166, 20, 213, 245, 28, 0, 0, 0},
	}
}

func (dictDataService *DictDataService) SelectDictDataByType(c *gin.Context, dictType string) (data []byte) {

	data = dictDataService.getDictCache(c, dictType)
	if data != nil && len(data) != 0 {
		return
	}
	sysDictDataList := dictDataService.dictDataDao.SelectDictDataByType(c, dictType)
	if len(sysDictDataList) != 0 {
		responseData := response.ResponseData{Code: response.Success, Msg: response.Success.Msg(), Data: sysDictDataList}
		marshal, err := json.Marshal(responseData)
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		if err != nil {
			panic(err)
		}
		if _, err = gz.Write(marshal); err != nil {
			panic(err)
		}
		if err = gz.Close(); err != nil {
			panic(err)
		}
		compressedData := buf.Bytes()
		go dictDataService.cache.Set(context.Background(), dictDataService.dictKey+dictType, string(compressedData), 0)
		return compressedData
	}
	return dictDataService.gzipNil
}
func (dictDataService *DictDataService) SelectDictDataList(c *gin.Context, dictData *systemModels.SysDictDataDQL) (list []*systemModels.SysDictDataVo, total int64) {
	return dictDataService.dictDataDao.SelectDictDataList(c, dictData)

}
func (dictDataService *DictDataService) ExportDictData(c *gin.Context, dictData *systemModels.SysDictDataDQL) (data []byte) {
	list, _ := dictDataService.dictDataDao.SelectDictDataList(c, dictData)
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
func (dictDataService *DictDataService) SelectDictDataById(c *gin.Context, dictCode int64) (dictData *systemModels.SysDictDataVo) {
	return dictDataService.dictDataDao.SelectDictDataById(c, dictCode)

}

func (dictDataService *DictDataService) InsertDictData(c *gin.Context, dictData *systemModels.SysDictDataVo) {
	dictData.DictCode = snowflake.GenID()
	dictDataService.dictDataDao.InsertDictData(c, dictData)
	dictDataService.deleteDictCache(dictData.DictType)

}

func (dictDataService *DictDataService) UpdateDictData(c *gin.Context, dictData *systemModels.SysDictDataVo) {
	dictDataService.dictDataDao.UpdateDictData(c, dictData)
	dictDataService.deleteDictCache(dictData.DictType)
}
func (dictDataService *DictDataService) DeleteDictDataByIds(c *gin.Context, dictCodes []int64) {

	codes := dictDataService.dictDataDao.SelectDictTypesByDictCodes(c, dictCodes)
	dictDataService.dictDataDao.DeleteDictDataByIds(c, dictCodes)
	for _, code := range codes {
		dictDataService.deleteDictCache(code)
	}

}
func (dictDataService *DictDataService) CheckDictDataByTypes(c *gin.Context, dictType []string) bool {
	return dictDataService.dictDataDao.CountDictDataByTypes(c, dictType) > 0

}
func (dictDataService *DictDataService) getDictCache(c context.Context, dictType string) (dictDataList []byte) {
	getString, _ := dictDataService.cache.Get(c, dictDataService.dictKey+dictType)
	return []byte(getString)
}

func (dictDataService *DictDataService) deleteDictCache(dictType string) {
	go func() {
		time.Sleep(time.Second * 3)
		dictDataService.cache.Del(context.Background(), dictDataService.dictKey+dictType)
	}()
}
