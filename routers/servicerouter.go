package routers

import (
	"github.com/xiangrui2019/go-kahla-bot-server/controllers"
	"gopkg.in/macaron.v1"
)


func ConfigureServiceRouting(context *macaron.Macaron) error {
	maincontroller := controllers.NewMainController()

	context.Any("/", maincontroller.Index)

	return nil
}