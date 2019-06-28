package controllers

import "gopkg.in/macaron.v1"

type MessageController struct {

}

func NewMessageController() *MessageController {
	return &MessageController{

	}
}

func (c *MessageController) SendText(context *macaron.Context) {
	
}