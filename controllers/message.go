package controllers

import (
	"fmt"
	"github.com/xiangrui2019/go-kahla-bot-server/api"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/services"
	"gopkg.in/macaron.v1"
	"reflect"
	"strconv"
)

type MessageController struct {
	messageService *services.MessageService
	client *kahla.Client
}

func NewMessageController(macaronapp *macaron.Macaron, injector *injects.BasicInject) *MessageController {
	client := macaronapp.GetVal(reflect.TypeOf(injector.Client)).Interface().(*kahla.Client)

	return &MessageController{
		messageService: services.NewMessageService(client),
		client: client,
	}
}

func (c *MessageController) SendText(context *macaron.Context, model api.SendTextRequestModel) {
	err := c.messageService.SendMessageByToken(model.Token, model.Content)

	if err != nil {
		context.JSON(500, map[string]string{
			"code": strconv.Itoa(enums.ResponseCodeSendMessageFailed),
			"message": err.Error(),
		})
		return
	}

	context.JSON(200, map[string]string{
		"code": "0",
		"message": "Successfully sent a message.",
	})
}

func (c *MessageController) SendImage(context *macaron.Context, model api.SendImageRequestModel) {
	file, _ := model.Image.Open()
	response, _, _ := c.client.Files.UploadMedia(&kahla.Files_UploadMediaRequest{
		File: file,
		Name: model.Image.Filename,
	})

	fmt.Printf("%+v", *response)
}