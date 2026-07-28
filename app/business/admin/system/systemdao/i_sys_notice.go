package systemdao

import (
	"context"
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
)

type ISysNoticeDao interface {
	SelectNoticeList(ctx context.Context, notice *modelquery.NoticeDQL) (list []*modelresponse.SysNoticeVo, total int64)
	SelectNoticeById(ctx context.Context, id int64) *modelresponse.SysNoticeVo
	InsertNotice(ctx context.Context, notice *modelrequest.SysNoticeDML)
	DeleteNoticeById(ctx context.Context, id int64)
	BatchSysNoticeUsers(ctx context.Context, notice []*modelentity.NoticeUser)
	SelectNewMessageCountByUserId(ctx context.Context, userId int64) int64
	SelectConsumptionNoticeById(ctx context.Context, userId, noticeId int64) *modelresponse.ConsumptionNoticeVo
	SelectConsumptionNoticeList(ctx context.Context, notice *modelquery.ConsumptionNoticeDQL) (list []*modelresponse.ConsumptionNoticeVo, total int64)
	SelectNoticeStatusByNoticeIdAndUserId(ctx context.Context, noticeId, userId int64) int
	SelectNoticeStatusByNoticeIdsAndUserId(ctx context.Context, noticeId []int64, userId int64) int
	UpdateNoticeRead(ctx context.Context, noticeId int64, userId int64)
	UpdateNoticeReadAll(ctx context.Context, userId int64)
	DeleteConsumptionNotice(ctx context.Context, noticeId []int64, userId int64)
}
