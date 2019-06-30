package controllers

import (
	"github.com/xiangrui2019/go-kahla-bot-server/api"
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
	// 通过反射获取卡拉client服务
	client := macaronapp.GetVal(reflect.TypeOf(injector.Client)).Interface().(*kahla.Client)

	return &MessageController{
		messageService: services.NewMessageService(client),
		client: client,
	}
}

func (c *MessageController) SendText(context *macaron.Context, model api.SendTextRequestModel) {
	err := c.messageService.SendMessageByToken(model.Token, model.Content)

	if err != nil {
		context.JSON(200, api.SendTextResponseModel{
			Code: 500,
			Message: err.Error(),
		})
		return
	}

	context.JSON(200, api.SendTextResponseModel{
		Code: 200,
		Message: "Successfully sent a message.",
	})
}

func (c *MessageController) SendImage(context *macaron.Context, model api.SendImageRequestModel) {
	err := c.messageService.SendImageMessageByToken(model.Token, model.Image)

	if err != nil {
		context.JSON(200, api.SendImageResponseModel{
			Code: 500,
			Message: err.Error(),
		})
		return
	}

	context.JSON(200, api.SendImageResponseModel{
		Code: 200,
		Message: "Successfully sent a message.",
	})
}

func (c *MessageController) SendVideo(context *macaron.Context, model api.SendVideoRequestModel) {
	err := c.messageService.SendVideoMessageByToken(model.Token, model.Video)

	if err != nil {
		context.JSON(200, api.SendVideoResponseModel{
			Code: 500,
			Message: err.Error(),
		})
		return
	}

	context.JSON(200, api.SendVideoResponseModel{
		Code: 200,
		Message: "Successfully sent a message.",
	})
}

func (c *MessageController) SendFile(context *macaron.Context) {

}

func (c *MessageController) SendVoice(context *macaron.Context, model api.SendVoiceRequestModel) {
	err := c.messageService.SendVoiceMessageByToken(model.Token, model.Voice)

	if err != nil {
		context.JSON(200, api.SendImageResponseModel{
			Code: 500,
			Message: err.Error(),
		})
		return
	}

	context.JSON(200, api.SendImageResponseModel{
		Code: 200,
		Message: "Successfully sent a message.",
	})
}