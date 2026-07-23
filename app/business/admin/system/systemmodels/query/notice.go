package query

import "nova-factory-server/app/baize"

type NoticeDQL struct {
	NoticeTitle string `form:"noticeTitle" db:"notice_title"`
	CreateBy    string `form:"createBy" db:"create_by"`
	NoticeType  string `form:"noticeType" db:"notice_type"`
	baize.BaseEntityDQL
}

type ConsumptionNoticeDQL struct {
	Status string `form:"status" db:"status"` //未读消息1 已读2 全部不填
	Title  string `form:"title" db:"title"`
	Type   string `form:"type" db:"type" ` //消息类型
	UserId int64  `db:"user_id"  swaggerignore:"true"`
	baize.BaseEntityDQL
}
