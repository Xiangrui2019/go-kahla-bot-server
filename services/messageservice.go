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
	conversation, err := s.getConversation(conversationId)

	if err != nil {
		return err
	}

	width, height, err := functions.GetImageSize(file)

	if err != nil {
		return err
	}

	fileKey, err := s.uploadMedia(file)

	if err != nil {
		return err
	}

	message := fmt.Sprintf("[img]%s-%s-%s-0", *fileKey, strconv.Itoa(width), strconv.Itoa(height))

	err = s.SendRawMessage(conversationId, message, conversation.Value.AesKey)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendVideoMessageByToken(token string, file *multipart.FileHeader) error {
	user, err := dao.GetBotUserByToken(token)

	if err != nil {
		return err
	}

	err = s.SendVideoMessageByConversationId(user.ConversationId, file)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendVideoMessageByConversationId(conversationId uint32, file *multipart.FileHeader) error {
	conversation, err := s.getConversation(conversationId)

	if err != nil {
		return err
	}

	fileKey, err := s.uploadMedia(file)

	if err != nil {
		return err
	}

	message := fmt.Sprintf("[video]%s", *fileKey)

	err = s.SendRawMessage(conversationId, message, conversation.Value.AesKey)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendFileMessageByToken(token string, file *multipart.FileHeader) error {
	user, err := dao.GetBotUserByToken(token)

	if err != nil {
		return err
	}

	err = s.SendFileMessageByConversationId(user.ConversationId, file)

	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) SendFileMessageByConversationId(conversationId uint32, file *multipart.FileHeader) error {
	conversation, err := s.getConversation(conversationId)

	if err != nil {
		return err
	}

	fileKey, err := s.uploadFile(file, conversationId)

	if err != nil {
		return err
	}

	fileSize, err := functions.CalcFileSize(file)

	if err != nil {
		return err
	}

	message := fmt.Sprintf("[file]%s-%s-%s", *fileKey, file.Filename, *fileSize)

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
	conversation, err := s.getConversation(conversationId)

	if err != nil {
		return err
	}

	fileKey, err := s.uploadFile(file, conversationId)

	if err != nil {
		return err
	}

	message := fmt.Sprintf("[audio]%s", *fileKey)

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

func (s *MessageService) getConversation(conversationId uint32) (*kahla.Conversation_ConversationDetailResponse, error) {
	conversation, httpResponse, err := s.client.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: conversationId,
	})

	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	if conversation.Code != enums.ResponseCodeOK {
		return nil, errors.New(conversation.Message)
	}

	return conversation, nil
}

func (s *MessageService) uploadMedia(file *multipart.FileHeader) (*string, error) {
	media, err := file.Open()

	defer media.Close()

	if err != nil {
		return nil, err
	}

	mediaresp, httpResponse, _ := s.client.Files.UploadMedia(&kahla.Files_UploadMediaRequest{
		File: media,
		Name: file.Filename,
	})

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	if mediaresp.Code != enums.ResponseCodeOK {
		return nil, errors.New(mediaresp.Message)
	}

	fileKey := functions.ParseFileKey(mediaresp.DownloadPath)

	return &fileKey, nil
}

func (s *MessageService) uploadFile(file *multipart.FileHeader, conversationId uint32) (*string, error) {
	filex, err := file.Open()

	defer filex.Close()

	if err != nil {
		return nil, err
	}

	fileresp, httpResponse, _ := s.client.Files.UploadFile(&kahla.Files_UploadFileRequest{
		File:           filex,
		Name:           file.Filename,
		ConversationId: conversationId,
	})

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New("status code not 200")
	}

	if fileresp.Code != enums.ResponseCodeOK {
		return nil, errors.New(fileresp.Message)
	}

	fileKey := strconv.Itoa(int(fileresp.FileKey))

	return &fileKey, nil
}