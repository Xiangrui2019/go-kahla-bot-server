package services

import (
	"errors"
	"github.com/xiangrui2019/go-kahla-bot-server/cryptojs"
	"github.com/xiangrui2019/go-kahla-bot-server/dao"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"net/http"
)

type MessageService struct {
	client *kahla.Client
}

func NewMessageService(clien *kahla.Client) *MessageService {
	return &MessageService{
		client: clien,
	}
}

func (s *MessageService) SendMessageByConversationId(conversationId uint32, message string) error {
	response, httpResponse, err := s.client.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: conversationId,
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

	content, err := cryptojs.AesEncrypt(message, response.Value.AesKey)

	if err != nil {
		return err
	}

	err = s.SendRawMessage(conversationId, content)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendMessageByToken(token string, message string) error {
	user, err := dao.GetBotUserByToken(token)

	if err != nil {
		return err
	}

	response, httpResponse, err := s.client.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: user.ConversationId,
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

	content, err := cryptojs.AesEncrypt(message, response.Value.AesKey)

	if err != nil {
		return err
	}

	err = s.SendRawMessage(user.ConversationId, content)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendRawMessage(conversationId uint32, message string) error {
	response, httpResponse, err := s.client.Conversation.SendMessage(&kahla.Conversation_SendMessageRequest{
		Id:      conversationId,
		Content: message,
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