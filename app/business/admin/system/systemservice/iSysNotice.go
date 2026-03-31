package systemservice

import (
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type ISysNoticeService interface {
	SelectNoticeList(c *gin.Context, notice *systemmodels.NoticeDQL) (list []*systemmodels.SysNoticeVo, total int64)
	SelectNoticeById(c *gin.Context, id int64) *systemmodels.SysNoticeVo
	InsertNotice(c *gin.Context, notice *systemmodels.SysNoticeVo)
	NewMessAge(c *gin.Context, userId int64) int64
	SelectConsumptionNoticeById(c *gin.Context, noticeId int64) *systemmodels.ConsumptionNoticeVo
	SelectConsumptionNoticeList(c *gin.Context, notice *systemmodels.ConsumptionNoticeDQL) (list []*systemmodels.ConsumptionNoticeVo, total int64)
	UpdateNoticeRead(c *gin.Context, noticeId, userId int64)
	UpdateNoticeReadAll(c *gin.Context, userId int64)
	DeleteConsumptionNotice(c *gin.Context, noticeId []int64, userId int64)
}
