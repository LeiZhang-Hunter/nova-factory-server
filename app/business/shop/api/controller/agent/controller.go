//go:build ai

package agent

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewConversations, NewMessage, NewVoice, wire.Struct(new(Controller), "*"))

type Controller struct {
	Conversations *Conversations
	Message       *Message
	Voice         *Voice
}
