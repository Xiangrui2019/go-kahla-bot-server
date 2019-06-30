package api

import "mime/multipart"

// 定义图片发送输入模型
type SendImageRequestModel struct {
	Token string `form:"token" binding:"Required"`
	Image *multipart.FileHeader `form:"image" binding:"Required"`
}

// 定义图片发送输出模型
type SendImageResponseModel struct {
	Code int `json:"code"`
	Message string `json:"message"`
}