package systemservice

import (
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type IPostService interface {
	PostExport(c *gin.Context, role *systemmodels.SysPostDQL) (data []byte)
	SelectPostList(c *gin.Context, post *systemmodels.SysPostDQL) (list []*systemmodels.SysPostVo, total int64)
	SelectPostById(c *gin.Context, postId int64) (Post *systemmodels.SysPostVo)
	InsertPost(c *gin.Context, post *systemmodels.SysPostVo)
	UpdatePost(c *gin.Context, post *systemmodels.SysPostVo)
	DeletePostByIds(c *gin.Context, postId []int64)
}
