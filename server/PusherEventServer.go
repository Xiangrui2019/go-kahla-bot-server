package server

import (
	"errors"
	"fmt"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/pusher"
	"log"
)

type PusherEventServer struct {
	client *kahla.Client
	config *conf.Config
	pushereventing *pusher.Pusher
	handler *EventHandler
}

func NewPusherServer() *PusherEventServer {
	var err error

	server := new(PusherEventServer)

	server.config, err = conf.LoadConfigFromFile("./config.toml")

	if err != nil {
		return nil
	}

	server.client = kahla.NewClient(server.config.BotConfig.KahlaServer, "https://oss.cdn.aiursoft.com")

	server.pushereventing = pusher.NewPusher("", server.EventHandler)

	server.handler = NewEventHandler()

	return server
}

func (server *PusherEventServer) login() error {
	response, _, err := server.client.Auth.AuthByPassword(&kahla.Auth_AuthByPasswordRequest{
		Email: server.config.BotConfig.Email,
		Password: server.config.BotConfig.Password,
	})

	if err != nil {
		return err
	}

	if response.Code != enums.ResponseCodeOK {
		return errors.New(response.Message)
	}

	return nil
}

func (server *PusherEventServer) initpusher() (*string, error) {
	response, _, err := server.client.Auth.InitPusher()

	if err != nil {
		return nil, err
	}

	if response.Code != enums.ResponseCodeOK {
		return nil, errors.New(response.Message)
	}

	return &response.ServerPath, nil
}

func (server *PusherEventServer) runWebsocket(interrupt chan struct{}) error {
	serverPath, err := server.initpusher()

	if err != nil {
		return err
	}

	server.pushereventing.Url = *serverPath

	err = server.pushereventing.Connect(interrupt)

	if err != nil {
		return err
	}

	return nil
}

func (server *PusherEventServer) EventHandler(i interface{}) {
	switch v := i.(type) {
	case *pusher.Pusher_NewMessageEvent:
		err := server.handler.NewMessageEvent(v)
		if err != nil {
			log.Println(err)
		}
	case *pusher.Pusher_NewFriendRequestEvent:
		err := server.handler.NewFriendRequestEvent(v)
		if err != nil {
			log.Println(err)
		}
	case *pusher.Pusher_WereDeletedEvent:
		err := server.handler.WereDeletedEvent(v)
		if err != nil {
			log.Println(err)
		}
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
	err := server.login()

	if err != nil {
		return err
	}

	err = server.runWebsocket(interrupt)

	if err != nil {
		return err
	}

	return nil
}