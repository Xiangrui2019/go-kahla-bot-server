package server

import (
	"errors"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/dao"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/models"
	"github.com/xiangrui2019/go-kahla-bot-server/pusher"
	"github.com/xiangrui2019/go-kahla-bot-server/services"
	"gopkg.in/macaron.v1"
	"log"
	"reflect"
	"strconv"
)

type EventHandler struct {
	config *conf.Config
	client *kahla.Client
	tokenService *services.TokenService
	friendRequestChan chan struct{}
}

func NewEventHandler(macaronapp *macaron.Macaron, injector *injects.BasicInject, cilen *kahla.Client) *EventHandler {
	handler := &EventHandler{
		config: macaronapp.GetVal(reflect.TypeOf(injector.Config)).Interface().(*conf.Config),
		client: cilen,
	}

	handler.tokenService = services.NewTokenService(macaronapp, injector, handler.client)

	return handler
}

func (h *EventHandler) NewMessageEvent(v *pusher.Pusher_NewMessageEvent) error {

	return nil
}

func (h *EventHandler) NewFriendRequestEvent(v *pusher.Pusher_NewFriendRequestEvent) error {
	log.Println("有一个新的用户请求加入公众号..")
	log.Printf("用户名: %s", v.Requester.NickName)
	log.Println("首先需要清理数据库...")
	err := h.UpdateConversation()
	log.Println("数据库清理成功...")
	if err != nil {
		log.Println(err)
		return err
	}
	err = h.AcceptFriendRequest()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("已经同意此用户加入公众号.")
	return nil
}

func (h *EventHandler) WereDeletedEvent(v *pusher.Pusher_WereDeletedEvent) error {
	log.Printf("有人从我们的公众号退出了: %s", v.Trigger.NickName)
	log.Printf("详细信息: ID: %s", v.Trigger.Id)
	err := h.UpdateConversation()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("退出事件处理完成.")
	return nil
}

func (h *EventHandler) UpdateConversation() error {
	response, _, err := h.client.Friendship.Mine()

	if err != nil {
		return err
	}

	allusers, err := dao.GetAllBotUser()

	if err != nil {
		return err
	}

	for _, v := range allusers {
		isinkahla := false
		for _, user := range response.Users {
			if v.KahlaUserId == user.Id {
				isinkahla = true
			}
		}

		if !isinkahla {
			log.Println("删除了当前不存在的好友信息")
			err := dao.DeleteBotUser(v.Id)
			return err
		}
	}

	return nil
}

func (h *EventHandler) AcceptFriendRequest() error {
	var err1 error

	response, _, err := h.client.Friendship.MyRequests()

	if err != nil {
		return err
	}

	if response.Code != enums.ResponseCodeOK {
		return errors.New(response.Message)
	}
CONTINUE:
	for _, v := range response.Items {
		if !v.Completed {
			response, _, err := h.client.Friendship.CompleteRequest(&kahla.Friendship_CompleteRequestRequest{
				Id: strconv.Itoa(int(v.Id)),
				Accept: true,
			})
			log.Printf("同意了一个好友请求: %d", v.Id)

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

			mines, _, err := h.client.Friendship.Mine()

			if err != nil {
				if err1 == nil {
					err1 = err
				}
				continue
			}

			for _, user := range mines.Users {
				if user.Id == v.CreatorId {
					response, _, err := h.client.Friendship.UserDetail(&kahla.Friendship_UserDetailRequest{
						Id: user.Id,
					})

					if err != nil {
						if err1 == nil {
							err1 = err
						}
						continue CONTINUE
					}

					if response.Code != enums.ResponseCodeOK {
						if err1 == nil {
							err1 = err
						}
						continue CONTINUE
					}

					if response.AreFriends != true {
						if err1 == nil {
							err1 = errors.New("your are not friends")
						}
						continue CONTINUE
					}

					token, err := h.tokenService.SendToken(response.ConversationId)

					if err != nil {
						if err1 == nil {
							err1 = err
						}
						continue CONTINUE
					}

					err = dao.CreateBotUser(&models.BotUser{
						Token: *token,
						Nickname: v.Creator.NickName,
						KahlaUserId: v.Creator.Id,
					})

					if err != nil {
						if err1 == nil {
							err1 = err
						}
						continue CONTINUE
					}

					err = h.UpdateConversation()

					if err != nil {
						if err1 == nil {
							err1 = err
						}
						continue CONTINUE
					}
				}
			}
		}
	}

	return err1
}