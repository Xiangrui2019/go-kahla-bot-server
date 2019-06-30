package services

import (
	"errors"
	"fmt"
	"github.com/xiangrui2019/go-kahla-bot-server/cryptojs"
	"github.com/xiangrui2019/go-kahla-bot-server/dao"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/functions"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"mime/multipart"
	"net/http"
	"strconv"
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
	conversation, httpResponse, err := s.client.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: conversationId,
	})

	if err != nil {
		return err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	if conversation.Code != enums.ResponseCodeOK {
		return errors.New(conversation.Message)
	}

	err = s.SendRawMessage(conversationId, message, conversation.Value.AesKey)

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

	err = s.SendMessageByConversationId(user.ConversationId, message)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendImageMessageByToken(token string, file *multipart.FileHeader) error {
	user, err := dao.GetBotUserByToken(token)

	if err != nil {
		return err
	}

	err = s.SendImageMessageByConversationId(user.ConversationId, file)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendImageMessageByConversationId(conversationId uint32, file *multipart.FileHeader) error {
	conversation, httpResponse, err := s.client.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: conversationId,
	})

	if err != nil {
		return err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	if conversation.Code != enums.ResponseCodeOK {
		return errors.New(conversation.Message)
	}

	width, height, err := functions.GetImageSize(file)

	if err != nil {
		return err
	}

	imagefile, err := file.Open()

	if err != nil {
		return err
	}

	mediaresponse, httpResponse, _ := s.client.Files.UploadMedia(&kahla.Files_UploadMediaRequest{
		File: imagefile,
		Name: file.Filename,
	})

	if httpResponse.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	if mediaresponse.Code != enums.ResponseCodeOK {
		return errors.New(conversation.Message)
	}

	fileKey := functions.ParseFileKey(mediaresponse.DownloadPath)

	message := fmt.Sprintf("[img]%s-%s-%s-0", fileKey, strconv.Itoa(width), strconv.Itoa(height))

	err = s.SendRawMessage(conversationId, message, conversation.Value.AesKey)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendVoiceMessageByToken(token string, file *multipart.FileHeader) error {
	user, err := dao.GetBotUserByToken(token)

	if err != nil {
		return err
	}

	err = s.SendVoiceMessageByConversationId(user.ConversationId, file)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendVoiceMessageByConversationId(conversationId uint32, file *multipart.FileHeader) error {
	conversation, httpResponse, err := s.client.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: conversationId,
	})

	if err != nil {
		return err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	if conversation.Code != enums.ResponseCodeOK {
		return errors.New(conversation.Message)
	}

	voicefile, err := file.Open()

	if err != nil {
		return err
	}

	mediaresponse, httpResponse, _ := s.client.Files.UploadFile(&kahla.Files_UploadFileRequest{
		File:           voicefile,
		Name:           file.Filename,
		ConversationId: conversationId,
	})

	if httpResponse.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	if mediaresponse.Code != enums.ResponseCodeOK {
		return errors.New(conversation.Message)
	}

	message := fmt.Sprintf("[audio]%d", mediaresponse.FileKey)

	err = s.SendRawMessage(conversationId, message, conversation.Value.AesKey)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendRawMessage(conversationId uint32, message string, AesKey string) error {
	content, err := cryptojs.AesEncrypt(message, AesKey)

	if err != nil {
		return err
	}

	response, httpResponse, err := s.client.Conversation.SendMessage(&kahla.Conversation_SendMessageRequest{
		Id:      conversationId,
		Content: content,
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