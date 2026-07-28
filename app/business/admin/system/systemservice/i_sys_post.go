package systemservice

import (
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type IPostService interface {
	PostExport(c *gin.Context, role *modelquery.SysPostDQL) (data []byte)
	SelectPostList(c *gin.Context, post *modelquery.SysPostDQL) (list []*modelresponse.SysPostVo, total int64)
	SelectPostById(c *gin.Context, postId int64) (Post *modelresponse.SysPostVo)
	InsertPost(c *gin.Context, post *modelrequest.SysPostDML)
	UpdatePost(c *gin.Context, post *modelrequest.SysPostDML)
	DeletePostByIds(c *gin.Context, postId []int64)
}
