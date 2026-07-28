package systemServiceImpl

import (
	"nova-factory-server/app/business/admin/system/systemdao"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
	"nova-factory-server/app/business/admin/system/systemservice"
	"nova-factory-server/app/utils/excel"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type PostService struct {
	postDao systemdao.IPostDao
}

func NewPostService(pd systemdao.IPostDao) systemservice.IPostService {
	return &PostService{
		postDao: pd,
	}
}

func (postService *PostService) SelectPostList(c *gin.Context, post *modelquery.SysPostDQL) (list []*modelresponse.SysPostVo, total int64) {
	return postService.postDao.SelectPostList(c, post)

}
func (postService *PostService) PostExport(c *gin.Context, post *modelquery.SysPostDQL) (data []byte) {
	list := postService.postDao.SelectPostListAll(c, post)
	toExcel, err := excel.SliceToExcel(list)
	if err != nil {
		panic(err)
	}
	buffer, err := toExcel.WriteToBuffer()
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func (postService *PostService) SelectPostById(c *gin.Context, postId int64) (Post *modelresponse.SysPostVo) {
	return postService.postDao.SelectPostById(c, postId)

}

func (postService *PostService) InsertPost(c *gin.Context, post *modelrequest.SysPostDML) {
	post.PostId = snowflake.GenID()
	postService.postDao.InsertPost(c, post)
}

func (postService *PostService) UpdatePost(c *gin.Context, post *modelrequest.SysPostDML) {
	postService.postDao.UpdatePost(c, post)
}
func (postService *PostService) DeletePostByIds(c *gin.Context, postId []int64) {
	postService.postDao.DeletePostByIds(c, postId)
	return
}
