package api

import "mime/multipart"

type SendImageRequestModel struct {
	Token string `form:"token" binding:"Required"`
	Image *multipart.FileHeader `form:"image" binding:"Required"`
}

type SendImageResponseModel struct {
	Code int `json:"code"`
	Message string `json:"message"`
}