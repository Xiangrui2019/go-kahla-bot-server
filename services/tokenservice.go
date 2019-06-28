package services

import (
	"fmt"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/functions"
	"strings"
)

type TokenService struct {
	messageService *MessageService
	config *conf.Config
}

func NewTokenService(client *kahla.Client) *TokenService {
	c, _ := conf.LoadConfigFromFile("./config.toml")

	return &TokenService{
		messageService: NewMessageService(client),
		config: c,
	}
}

func (s *TokenService) SendToken(conversationId uint32) (*string, error) {
	token := functions.RandomString(s.config.TokenLength)

	if length := strings.Count(token, "") - 1; length > 48 {
		err := s.messageService.SendMessageByConversationId(conversationId, fmt.Sprintf("您被服务器分配的Token是: %s", token[1:48+1]))

		if err != nil {
			return nil, err
		}

		return &token, nil
	}

	err := s.messageService.SendMessageByConversationId(conversationId, fmt.Sprintf("您被服务器分配的令牌是: %s", token))

	if err != nil {
		return nil, err
	}

	return &token, nil
}