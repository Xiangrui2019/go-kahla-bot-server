package api

import "mime/multipart"

// 定义文件发送输入模型
type SendFileRequestModel struct {
	Token string `form:"token" binding:"Required"`
	File *multipart.FileHeader `form:"file" binding:"Required"`
}

// 定义文件发送输出模型
type SendFileResponseModel struct {
	Code int `json:"code"`
	Message string `json:"message"`
}