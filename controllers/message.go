package controllers

import (
	"github.com/xiangrui2019/go-kahla-bot-server/api"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"github.com/xiangrui2019/go-kahla-bot-server/services"
	"gopkg.in/macaron.v1"
	"reflect"
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
		context.JSON(500, api.SendTextResponseModel{
			Code: enums.ResponseCodeSendMessageFailed,
			Message: err.Error(),
		})
		return
	}

	context.JSON(200, api.SendTextResponseModel{
		Code: enums.ResponseCodeOK,
		Message: "Successfully sent a message.",
	})
}

func (c *MessageController) SendImage(context *macaron.Context, model api.SendImageRequestModel) {
	err := c.messageService.SendImageMessageByToken(model.Token, model.Image)

	if err != nil {
		context.JSON(500, api.SendImageResponseModel{
			Code: enums.ImageError,
			Message: err.Error(),
		})
		return
	}

	context.JSON(200, api.SendImageResponseModel{
		Code: enums.ResponseCodeOK,
		Message: "Successfully sent a message.",
	})
}