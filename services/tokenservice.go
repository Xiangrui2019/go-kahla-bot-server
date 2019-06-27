package services

import (
	"fmt"
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

func (s *TokenService) SendToken(conversationId uint32) (*string, error) {
	token := utils.RandomString(32)

	err := s.messageService.SendMessageByConversationId(conversationId, fmt.Sprintf("您被服务器分配的令牌是: %s", token))

	if err != nil {
		return nil, nil
	}

	return &token, nil
}