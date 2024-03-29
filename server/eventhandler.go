package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"net/url"
	"strings"

	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/cryptojs"
	"github.com/xiangrui2019/go-kahla-bot-server/dao"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/models"
	"github.com/xiangrui2019/go-kahla-bot-server/pusher"
	"github.com/xiangrui2019/go-kahla-bot-server/services"
	"gopkg.in/macaron.v1"
)

type EventHandler struct {
	config            *conf.Config
	client            *kahla.Client
	tokenService      *services.TokenService
	httpclient        *http.Client
	friendRequestChan chan struct{}
}

func NewEventHandler(macaronapp *macaron.Macaron, injector *injects.BasicInject, cilen *kahla.Client) *EventHandler {
	handler := &EventHandler{
		config:     macaronapp.GetVal(reflect.TypeOf(injector.Config)).Interface().(*conf.Config),
		client:     cilen,
		httpclient: &http.Client{},
	}

	handler.tokenService = services.NewTokenService(macaronapp, injector, handler.client)

	return handler
}

func (h *EventHandler) NewMessageEvent(v *pusher.Pusher_NewMessageEvent) error {
	log.Println("成功接受了一条消息...")
	if v.Sender.NickName != h.config.BotConfig.Name {
		log.Println("消息验证成功, 准备处理消息")
		aesKey, err := h.getAesKey(v.ConversationId)
		if err != nil {
			return err
		}
		log.Println("获取AESKEY成功...")

		content, err := cryptojs.AesDecrypt(v.Content, *aesKey)
		if err != nil {
			return err
		}
		log.Printf("消息解密成功, 消息为: %s", content)

		log.Println("开始处理消息...")
		err = h.ProcessNewMessageEvent(content, v)
		log.Println("消息处理完成")
		if err != nil {
			return err
		}
	}
	log.Println("消息没有通过验证, 这条消息是机器自动发送的...")
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

func (h *EventHandler) ProcessNewMessageEvent(content string, v *pusher.Pusher_NewMessageEvent) error {
	message, messagetype, rawmessage := h.parseMessageContent(content)
	user, err := dao.GetBotUserByConversationId(v.ConversationId)
	if err != nil {
		return err
	}
	err = h.callbacktoServer(v.Sender.NickName, message, messagetype, rawmessage, user.Token, v.ConversationId)

	if err != nil {
		return err
	}

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
	users, err := dao.GetAllBotUser()

	if err != nil {
		return err
	}

	mines, err := h.getMines()

	if err != nil {
		return err
	}

	for _, v := range users {
		x := false
		for _, user := range *mines {
			if v.KahlaUserId == user.Id {
				x = true
			}
		}

		if !x {
			log.Println("删除了当前不存在的好友信息")
			err := dao.DeleteBotUser(v.Id)
			return err
		}
	}

	return nil
}

func (h *EventHandler) AcceptFriendRequest() error {
	var err1 error

	requests, err := h.getMyRequests()

	if err != nil {
		return err
	}

	for _, v := range *requests {
		if !v.Completed {
			err := h.acceptCompleteRequest(strconv.Itoa(int(v.Id)))

			if err != nil {
				if err1 == nil {
					err1 = err
				}
				continue
			}

			err = h.updateUser(&v)

			if err != nil {
				if err1 == nil {
					err1 = err
				}
				continue
			}

			err = h.UpdateConversation()

			if err != nil {
				if err1 == nil {
					err1 = err
				}
				continue
			}
		}
	}

	return err1
}

func (h *EventHandler) getConversationId(Id string) (*uint32, error) {
	response, httpResponse, err := h.client.Friendship.UserDetail(&kahla.Friendship_UserDetailRequest{
		Id: Id,
	})

	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	if response.Code != enums.ResponseCodeOK {
		return nil, errors.New(response.Message)
	}

	if response.AreFriends != true {
		return nil, errors.New("your are not friends")
	}

	return &response.ConversationId, nil
}

func (h *EventHandler) parseMessageContent(content string) (string, int, string) {
	switch {
	case strings.Contains(content, "[img]"):
		splited := strings.Split(content, "]")[1]
		filekey := strings.Split(splited, "-")[0]
		downloadurl := "https://oss.aiursoft.com/download/fromkey/" + filekey
		return downloadurl, enums.Image, content
	case strings.Contains(content, "[video]"):
		filekey := strings.Split(content, "]")[1]
		downloadurl := "https://oss.aiursoft.com/download/fromkey/" + filekey
		return downloadurl, enums.Video, content
	case strings.Contains(content, "[audio]"):
		splited := strings.Split(content, "]")[1]
		filekey, err := strconv.Atoi(splited)
		if err != nil {
			return "", 0, ""
		}
		result, err := h.getFileDownloadAddress(uint32(filekey))
		if err != nil {
			return "", 0, ""
		}
		downloadurl := strings.ReplaceAll(*result, "audio", "audio.ogg")
		return downloadurl, enums.Audio, content
	case strings.Contains(content, "[file]"):
		splited := strings.Split(content, "]")[1]
		filekey, err := strconv.Atoi(strings.Split(splited, "-")[0])
		if err != nil {
			return "", 0, ""
		}
		downloadurl, err := h.getFileDownloadAddress(uint32(filekey))
		if err != nil {
			return "", 0, ""
		}
		return *downloadurl, enums.File, content
	default:
		return content, enums.Text, ""
	}
}

func (h *EventHandler) buildcallbackParam(username string, content string, content_type int, rawcontent string, token string, conversationId uint32) url.Values {
	v := url.Values{}
	v.Add("username", username)
	v.Add("message", content)
	v.Add("messagetype", strconv.Itoa(content_type))
	v.Add("rawmessage", rawcontent)
	v.Add("token", token)
	v.Add("conversationId", strconv.Itoa(int(conversationId)))
	return v
}

func (h *EventHandler) callbacktoServer(username string, content string, content_type int, rawcontent string, token string, conversationId uint32) error {
	v := h.buildcallbackParam(username, content, content_type, rawcontent, token, conversationId)
	serverurl := fmt.Sprintf("%s%s", h.config.BotConfig.CallbackServer, h.config.BotConfig.MessageCallbackEndpoint)

	req, err := http.NewRequest("POST", serverurl, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := h.httpclient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	return nil
}

func (h *EventHandler) getFileDownloadAddress(fileKey uint32) (*string, error) {
	response, httpResponse, err := h.client.Files.FileDownloadAddress(&kahla.Files_FileDownloadAddressRequest{
		FileKey: fileKey,
	})

	if err != nil {
		return nil, err
	}

	if response.Code != enums.ResponseCodeOK {
		return nil, errors.New(response.Message)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	return &response.DownloadPath, nil
}

func (h *EventHandler) getAesKey(conversationId uint32) (*string, error) {
	response, httpResponse, err := h.client.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: conversationId,
	})

	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	if response.Code != enums.ResponseCodeOK {
		return nil, errors.New(response.Message)
	}

	return &response.Value.AesKey, nil
}

func (h *EventHandler) getMyRequests() (*[]kahla.Friendship_MyRequestsResponse_Item, error) {
	response, httpResponse, err := h.client.Friendship.MyRequests()

	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	if response.Code != enums.ResponseCodeOK {
		return nil, errors.New(response.Message)
	}

	return &response.Items, nil
}

func (h *EventHandler) getMines() (*[]kahla.Friendship_MineResponse_User, error) {
	response, httpResponse, err := h.client.Friendship.Mine()

	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	if response.Code != enums.ResponseCodeOK {
		return nil, errors.New(response.Message)
	}

	return &response.Users, nil
}

func (h *EventHandler) acceptCompleteRequest(id string) error {
	response, httpResponse, err := h.client.Friendship.CompleteRequest(&kahla.Friendship_CompleteRequestRequest{
		Id:     id,
		Accept: true,
	})

	if err != nil {
		return err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	if response.Code != enums.ResponseCodeOK {
		return errors.New(response.Message)
	}

	return nil
}

func (h *EventHandler) updateUser(v *kahla.Friendship_MyRequestsResponse_Item) error {
	conversationId, err := h.getConversationId(v.CreatorId)

	if err != nil {
		return err
	}

	token, err := h.tokenService.SendToken(*conversationId)

	if err != nil {
		return err
	}

	err = dao.CreateBotUser(&models.BotUser{
		Token:          *token,
		Nickname:       v.Creator.NickName,
		KahlaUserId:    v.Creator.Id,
		ConversationId: *conversationId,
	})

	if err != nil {
		return err
	}

	return nil
}
