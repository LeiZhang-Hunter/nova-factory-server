package systemservice

import (
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type ISysNoticeService interface {
	SelectNoticeList(c *gin.Context, notice *modelquery.NoticeDQL) (list []*modelresponse.SysNoticeVo, total int64)
	SelectNoticeById(c *gin.Context, id int64) *modelresponse.SysNoticeVo
	InsertNotice(c *gin.Context, notice *modelrequest.SysNoticeDML)
	NewMessAge(c *gin.Context, userId int64) int64
	SelectConsumptionNoticeById(c *gin.Context, noticeId int64) *modelresponse.ConsumptionNoticeVo
	SelectConsumptionNoticeList(c *gin.Context, notice *modelquery.ConsumptionNoticeDQL) (list []*modelresponse.ConsumptionNoticeVo, total int64)
	UpdateNoticeRead(c *gin.Context, noticeId, userId int64)
	UpdateNoticeReadAll(c *gin.Context, userId int64)
	DeleteConsumptionNotice(c *gin.Context, noticeId []int64, userId int64)
}
