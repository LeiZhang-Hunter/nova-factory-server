package permissions

import "github.com/gin-gonic/gin"

type Permissions interface {
	GetPermission(c *gin.Context, userId int64) []string
}
