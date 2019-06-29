package api

import "mime/multipart"

type SendVoiceRequestModel struct {
	Token string `form:"token" binding:"Required"`
	Voice *multipart.FileHeader `form:"voice" binding:"Required"`
}