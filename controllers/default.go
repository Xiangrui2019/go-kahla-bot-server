package controllers

import (
	"github.com/xiangrui2019/go-kahla-bot-server/controllers/api"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
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
		Code: enums.ResponseCodeOK,
		Message: "Welcome to the kahla robot server.",
	})
}