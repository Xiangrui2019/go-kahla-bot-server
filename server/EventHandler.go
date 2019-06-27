package server

import (
	"errors"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/dao"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/models"
	"github.com/xiangrui2019/go-kahla-bot-server/pusher"
	"log"
	"strconv"
)

type EventHandler struct {
	config *conf.Config
	client *kahla.Client
	friendRequestChan chan struct{}
}

func NewEventHandler(cilen *kahla.Client) *EventHandler {
	c, _ := conf.LoadConfigFromFile("./config.toml")

	return &EventHandler{
		config: c,
		client: cilen,
	}
}

func (h *EventHandler) NewMessageEvent(v *pusher.Pusher_NewMessageEvent) error {

	return nil
}

func (h *EventHandler) NewFriendRequestEvent(v *pusher.Pusher_NewFriendRequestEvent) error {
	log.Println("有一个新的用户请求加入公众号..")
	log.Printf("用户名: %s", v.Requester.NickName)
	h.AcceptFriendRequest()
	log.Println("已经同意此用户加入公众号.")
	return nil
}

func (h *EventHandler) WereDeletedEvent(v *pusher.Pusher_WereDeletedEvent) error {

	return nil
}

func (h *EventHandler) AcceptFriendRequest() {
	select {
	case h.friendRequestChan <- struct{}{}:
		go func() {
			err := h.acceptFriendRequest()

			if err != nil {
				log.Println(err)
			}

			<-h.friendRequestChan
		}()
	default:
	}
}

func (h *EventHandler) acceptFriendRequest() error {
	var err1 error

	response, _, err := h.client.Friendship.MyRequests()

	if err != nil {
		return err
	}

	if response.Code != enums.ResponseCodeOK {
		return errors.New(response.Message)
	}

	for _, v := range response.Items {
		if !v.Completed {
			response, _, err := h.client.Friendship.CompleteRequest(&kahla.Friendship_CompleteRequestRequest{
				Id: strconv.Itoa(int(v.Id)),
				Accept: true,
			})

			if err != nil {
				if err1 == nil {
					err1 = err
				}
				continue
			}

			if response.Code != enums.ResponseCodeOK {
				if err1 == nil {
					err1 = errors.New(response.Message)
				}
				continue
			}

			err = dao.CreateBotUser(&models.BotUser{
				Token: "",
				Nickname: v.Creator.NickName,
				KahlaUserId: v.Creator.Id,
			})

			if err != nil {
				err1 = err
			}
		}
	}

	return nil
}