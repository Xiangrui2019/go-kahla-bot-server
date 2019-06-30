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
		context.Post("/sendtext", binding.Bind(api.SendTextRequestModel{}), messagecontroller.SendText)
		context.Post("/sendimage", binding.Bind(api.SendImageRequestModel{}), messagecontroller.SendImage)
		context.Post("/sendvoice", binding.Bind(api.SendVoiceRequestModel{}), messagecontroller.SendVoice)
		context.Post("/sendvideo", binding.Bind(api.SendVideoRequestModel{}), messagecontroller.SendVideo)
		context.Post("/sendfile", binding.Bind(api.SendFileRequestModel{}), messagecontroller.SendFile)
	})

	return nil
}
