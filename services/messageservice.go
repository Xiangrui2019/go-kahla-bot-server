package services

type MessageService struct {

}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) SendMessageByConversationId() error {
	return nil
}

func (s *MessageService) SendMessageByToken() error {
	return nil
}
