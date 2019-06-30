package api

import "mime/multipart"

// 发送语音输入模型
type SendVoiceRequestModel struct {
	Token string                `form:"token" binding:"Required"`
	Voice *multipart.FileHeader `form:"voice" binding:"Required"`
}

// 发送语音输出模型
type SendVoiceResponseModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
