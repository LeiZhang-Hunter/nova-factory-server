package alertController

import "github.com/gin-gonic/gin"

type Alert struct {
}

func NewAlert() *Alert {
	return &Alert{}
}

func (a *Alert) GetAlertList(c *gin.Context) {
	return
}

func (a *Alert) SetAlert(c *gin.Context) {
	return
}

func (a *Alert) RemoveAlert(c *gin.Context) {
	return
}
