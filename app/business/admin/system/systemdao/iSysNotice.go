package systemdao

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"
)

type ISysNoticeDao interface {
	SelectNoticeList(ctx context.Context, notice *systemmodels.NoticeDQL) (list []*systemmodels.SysNoticeVo, total int64)
	SelectNoticeById(ctx context.Context, id int64) *systemmodels.SysNoticeVo
	InsertNotice(ctx context.Context, notice *systemmodels.SysNoticeVo)
	DeleteNoticeById(ctx context.Context, id int64)
	BatchSysNoticeUsers(ctx context.Context, notice []*systemmodels.NoticeUser)
	SelectNewMessageCountByUserId(ctx context.Context, userId int64) int64
	SelectConsumptionNoticeById(ctx context.Context, userId, noticeId int64) *systemmodels.ConsumptionNoticeVo
	SelectConsumptionNoticeList(ctx context.Context, notice *systemmodels.ConsumptionNoticeDQL) (list []*systemmodels.ConsumptionNoticeVo, total int64)
	SelectNoticeStatusByNoticeIdAndUserId(ctx context.Context, noticeId, userId int64) int
	SelectNoticeStatusByNoticeIdsAndUserId(ctx context.Context, noticeId []int64, userId int64) int
	UpdateNoticeRead(ctx context.Context, noticeId int64, userId int64)
	UpdateNoticeReadAll(ctx context.Context, userId int64)
	DeleteConsumptionNotice(ctx context.Context, noticeId []int64, userId int64)
}
