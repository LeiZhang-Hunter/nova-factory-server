//go:build !ai

package agent

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewConversations, NewMessage, NewVoice, NewASR, wire.Struct(new(Controller), "*"))

type Controller struct {
	Conversations *Conversations
	Message       *Message
	Voice         *Voice
	ASR           *ASR
}

type Conversations struct{}

func NewConversations() *Conversations {
	return &Conversations{}
}

func (*Conversations) PrivateRoutes(_ *gin.RouterGroup) {}

func (*Conversations) ConfigRoutes(_ *gin.RouterGroup) {}

func (*Conversations) WsChatRegister(_ *gin.RouterGroup) {}

type Message struct{}

func NewMessage() *Message {
	return &Message{}
}

func (*Message) PrivateRoutes(_ *gin.RouterGroup) {}
