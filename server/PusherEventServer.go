package server

import (
	"errors"
	"fmt"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/cryptojs"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/pusher"
	"log"
)

type PusherEventServer struct {
	client *kahla.Client
	config *conf.Config
}

func NewPusherServer() *PusherEventServer {
	var err error

	server := new(PusherEventServer)

	server.config, err = conf.LoadConfigFromFile("./config.toml")

	if err != nil {
		return nil
	}

	server.client = kahla.NewClient(server.config.BotConfig.KahlaServer, "https://oss.cdn.aiursoft.com")

	return server
}

func (server *PusherEventServer) login() error {
	response, _, err := server.client.Auth.AuthByPassword(&kahla.Auth_AuthByPasswordRequest{})

	if err != nil {
		return err
	}

	if response.Code == enums.ResponseCodeOK {
		return errors.New(response.Message)
	}

	return nil
}

func (server *PusherEventServer) initpusher() (*string, error) {
	response, _, err := server.client.Auth.InitPusher()

	if err != nil {
		return nil, err
	}

	if response.Code == enums.ResponseCodeOK {
		return nil, errors.New(response.Message)
	}

	return &response.ServerPath, nil
}

func (server *PusherEventServer) EventHandler(i interface{}) {
	switch v := i.(type) {
	case *pusher.Pusher_NewMessageEvent:
		content, err := cryptojs.AesDecrypt(v.Content, v.AesKey)
		if err != nil {
			log.Println(err)
		} else {
			title := v.Sender.NickName
			message := content
			log.Println(title, ":", message)
		}
	case *pusher.Pusher_NewFriendRequestEvent:
		title := "Friend request"
		message := "You have got a new friend request!"
		log.Println(title, ":", message, "nick name:", v.Requester.NickName, "id:", v.Requester.Id)
	case *pusher.Pusher_WereDeletedEvent:
		title := "Were deleted"
		message := "You were deleted by one of your friends from his friend list."
		log.Println(title, ":", message, "nick name:", v.Trigger.NickName, "id:", v.Trigger.Id)
	case *pusher.Pusher_FriendAcceptedEvent:
		title := "Friend request"
		message := "Your friend request was accepted!"
		log.Println(title, ":", message, "nick name:", v.Target.NickName, "id:", v.Target.Id)
	case *pusher.Pusher_TimerUpdatedEvent:
		title := "Self-destruct timer updated!"
		message := fmt.Sprintf("Your current message life time is: %d", v.NewTimer)
		log.Println(title, ":", message, "conversation id:", v.ConversationId)
	case *pusher.Pusher_NewMemberEvent:
		title := "New member"
		message := fmt.Sprintf("%s has join the group.", v.NewMember.NickName)
		log.Println(title, ":", message, "conversation id:", v.ConversationId)
	case *pusher.Pusher_SomeoneLeftEvent:
		title := "Someone left"
		message := fmt.Sprintf("%s has successfully left the group.", v.LeftUser.NickName)
		log.Println(title, ":", message, "conversation id:", v.ConversationId)
	case *pusher.Pusher_DissolveEvent:
		title := "Group Dissolved"
		message := "A group dissolved"
		log.Println(title, ":", message, "conversation id:", v.ConversationId)
	}
}

func (server *PusherEventServer) Run(interrupt chan struct{}) error {
	return nil
}