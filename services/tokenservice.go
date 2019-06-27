package services

import (
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/utils"
)

type TokenService struct {
	messageService *MessageService
}

func NewTokenService(client *kahla.Client) *TokenService {
	return &TokenService{
		messageService: NewMessageService(client),
	}
}

func (s *TokenService) SendToken(conversationId uint32) (string, error) {
	token := utils.RandomString(32)

	// TODO:

	s.messageService.SendMessageByConversationId()

	return token, nil
}