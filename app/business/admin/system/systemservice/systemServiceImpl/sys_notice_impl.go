package systemServiceImpl

import (
	systemDao2 "nova-factory-server/app/business/admin/system/systemdao"
	systemModels2 "nova-factory-server/app/business/admin/system/systemmodels"
	systemService2 "nova-factory-server/app/business/admin/system/systemservice"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoticeService struct {
	nd  systemDao2.ISysNoticeDao
	sud systemDao2.IUserDao
	ss  systemService2.ISseService
	sss *systemModels2.SseType
}

func NewNoticeService(nd systemDao2.ISysNoticeDao,
	sud systemDao2.IUserDao, ss systemService2.ISseService) systemService2.ISysNoticeService {
	return &NoticeService{nd: nd, sud: sud, ss: ss, sss: &systemModels2.SseType{Key: "notice", Value: "1"}}
}

func (n *NoticeService) SelectNoticeList(c *gin.Context, notice *systemModels2.NoticeDQL) (list []*systemModels2.SysNoticeVo, total int64) {
	return n.nd.SelectNoticeList(c, notice)
}

func (n *NoticeService) SelectNoticeById(c *gin.Context, id int64) *systemModels2.SysNoticeVo {
	return n.nd.SelectNoticeById(c, id)
}

func (n *NoticeService) InsertNotice(c *gin.Context, notice *systemModels2.SysNoticeVo) {
	noticeId := snowflake.GenID()
	notice.Id = noticeId
	ids := notice.DeptIds
	notice.SetCreateBy(baizeContext.GetUserId(c))
	notice.DeptId = baizeContext.GetDeptId(c)
	notice.CreateName = baizeContext.GetUserName(c)
	deptIds := make([]int64, 0, len(ids.Data))
	for _, id := range ids.Data {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			continue
		}
		deptIds = append(deptIds, i)
	}
	userIds := n.sud.SelectUserIdsByDeptIds(c, deptIds)

	n.nd.InsertNotice(c, notice)
	if len(userIds) == 0 {
		return
	}
	users := make([]*systemModels2.NoticeUser, 0, len(userIds))
	for _, id := range userIds {
		s := new(systemModels2.NoticeUser)
		s.NoticeId = noticeId
		s.UserId = id
		s.Status = "1"
		users = append(users, s)
	}
	n.nd.BatchSysNoticeUsers(c, users)

}

func (n *NoticeService) NewMessAge(c *gin.Context, userId int64) int64 {
	return n.nd.SelectNewMessageCountByUserId(c, userId)
}

func (n *NoticeService) SelectConsumptionNoticeById(c *gin.Context, noticeId int64) *systemModels2.ConsumptionNoticeVo {
	userId := baizeContext.GetUserId(c)
	status := n.nd.SelectNoticeStatusByNoticeIdAndUserId(c, noticeId, userId)
	if status == 1 {
		n.nd.UpdateNoticeRead(c, noticeId, userId)
		go n.ss.SendNotification(c, &systemModels2.Sse{UserIds: []int64{userId}, Sse: n.sss})
	}
	return n.nd.SelectConsumptionNoticeById(c, userId, noticeId)
}

func (n *NoticeService) SelectConsumptionNoticeList(c *gin.Context, notice *systemModels2.ConsumptionNoticeDQL) (list []*systemModels2.ConsumptionNoticeVo, total int64) {
	return n.nd.SelectConsumptionNoticeList(c, notice)
}
func (n *NoticeService) UpdateNoticeRead(c *gin.Context, noticeId, userId int64) {
	status := n.nd.SelectNoticeStatusByNoticeIdAndUserId(c, noticeId, userId)
	if status == 0 {
		return
	}
	n.nd.UpdateNoticeRead(c, noticeId, userId)
	go n.ss.SendNotification(c, &systemModels2.Sse{UserIds: []int64{userId}, Sse: n.sss})
}
func (n *NoticeService) UpdateNoticeReadAll(c *gin.Context, userId int64) {
	n.nd.UpdateNoticeReadAll(c, userId)
	go n.ss.SendNotification(c, &systemModels2.Sse{UserIds: []int64{userId}, Sse: n.sss})
}
func (n *NoticeService) DeleteConsumptionNotice(c *gin.Context, noticeId []int64, userId int64) {
	i := n.nd.SelectNoticeStatusByNoticeIdsAndUserId(c, noticeId, userId)
	n.nd.DeleteConsumptionNotice(c, noticeId, userId)
	if i == 1 {
		go n.ss.SendNotification(c, &systemModels2.Sse{UserIds: []int64{userId}, Sse: n.sss})
	}
}
