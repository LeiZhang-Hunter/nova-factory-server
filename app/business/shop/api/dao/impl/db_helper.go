package impl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getCurrentDB 优先读取上下文中的事务对象，否则返回默认数据库连接。
func getCurrentDB(c *gin.Context, db *gorm.DB) *gorm.DB {
	if c != nil {
		if tx, ok := c.Get("db"); ok {
			if txDB, ok := tx.(*gorm.DB); ok && txDB != nil {
				return txDB
			}
		}
	}
	return db
}
