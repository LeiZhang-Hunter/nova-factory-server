package masterserviceimpl

import (
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

type StockServiceImpl struct {
	dao masterdao.IStockDao
}

func NewStockService(dao masterdao.IStockDao) masterservice.IStockService {
	return &StockServiceImpl{
		dao: dao,
	}
}

func (s StockServiceImpl) Create(c *gin.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s StockServiceImpl) Update(c *gin.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s StockServiceImpl) UpdateStock(c *gin.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s StockServiceImpl) GetStock(c *gin.Context) error {
	//TODO implement me
	panic("implement me")
}
