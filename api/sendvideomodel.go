package api

import "mime/multipart"

// 发送视频输入模型
type SendVideoRequestModel struct {
	Token string                `form:"token" binding:"Required"`
	Video *multipart.FileHeader `form:"video" binding:"Required"`
}

// 发送视频输出模型
type SendVideoResponseModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
