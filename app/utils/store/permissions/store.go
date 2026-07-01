package permissions

import "github.com/gin-gonic/gin"

type emptyPermissions struct{}

func newEmptyPermissions() Permissions {
	return &emptyPermissions{}
}

func (e *emptyPermissions) GetPermission(c *gin.Context, userId int64) []string {
	return []string{}
}
