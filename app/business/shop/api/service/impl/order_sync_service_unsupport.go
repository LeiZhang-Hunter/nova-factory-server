//go:build !erp
// +build !erp

package impl

//
//import "github.com/gin-gonic/gin"
//
//type shopOrderSyncService struct{}
//
//func NewShopOrderSyncService() *shopOrderSyncService {
//	return &shopOrderSyncService{}
//}
//
//func (s *shopOrderSyncService) SyncCreatedOrder(c *gin.Context, tid string) error {
//	return nil
//}
//
//func truncateSyncMessage(message string) string {
//	if len(message) <= 500 {
//		return message
//	}
//	return message[:500]
//}
