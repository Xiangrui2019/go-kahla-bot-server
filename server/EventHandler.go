package server

import (
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/pusher"
)

type EventHandler struct {
	config *conf.Config
}

func NewEventHandler() *EventHandler {
	c, _ := conf.LoadConfigFromFile("./config.toml")

	return &EventHandler{
		config: c,
	}
}

func (h *EventHandler) NewMessageEvent(v *pusher.Pusher_NewMessageEvent) error {

	return nil
}

func (h *EventHandler) NewFriendRequestEvent(v *pusher.Pusher_NewFriendRequestEvent) error {

	return nil
}

func (h *EventHandler) WereDeletedEvent(v *pusher.Pusher_WereDeletedEvent) error {

	return nil
}
