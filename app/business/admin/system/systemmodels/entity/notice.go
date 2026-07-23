package entity

type NoticeUser struct {
	NoticeId int64  `db:"notice_id,string"` //通知ID
	UserId   int64  `db:"user_id"`
	Status   string `db:"status"` //通知状态  1未读 2 已读
}
