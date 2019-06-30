package controllers

import (
	"github.com/xiangrui2019/go-kahla-bot-server/api"
	"gopkg.in/macaron.v1"
)

type MainController struct {

}

func NewMainController() *MainController {
	return &MainController{

	}
}

func (c *MainController) Index(context *macaron.Context)  {
	context.JSON(200, api.MainIndexModel{
		Code: 200,
		Message: "Welcome to the kahla robot server.",
	})
}

func (c *MainController) RedirectHome(context *macaron.Context) {
	context.Redirect("https://github.com/xiangrui2019/go-kahla-bot-server")
}