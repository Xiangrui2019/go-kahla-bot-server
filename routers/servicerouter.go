package routers

import (
	"github.com/go-macaron/binding"
	"github.com/xiangrui2019/go-kahla-bot-server/api"
	"github.com/xiangrui2019/go-kahla-bot-server/controllers"
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	"gopkg.in/macaron.v1"
)

// 配置控制器路由
func ConfigureServiceRouting(context *macaron.Macaron, injector *injects.BasicInject) error {
	maincontroller := controllers.NewMainController()
	messagecontroller := controllers.NewMessageController(context, injector)

	context.Any("/", maincontroller.Index)
	context.Any("/home", maincontroller.RedirectHome)

	context.Group("/message/", func() {
		context.Post("/sendtext", binding.Form(api.SendTextRequestModel{}), messagecontroller.SendText)
		context.Post("/sendimage", binding.MultipartForm(api.SendImageRequestModel{}), messagecontroller.SendImage)
		context.Post("/sendvoice", binding.MultipartForm(api.SendVoiceRequestModel{}), messagecontroller.SendVoice)
	})

	return nil
}