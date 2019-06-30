package api

import "mime/multipart"

type SendVideoRequestModel struct {
	Token string `form:"token" binding:"Required"`
	Video *multipart.FileHeader `form:"video" binding:"Required"`
}

type SendVideoResponseModel struct {
	Code int `json:"code"`
	Message string `json:"message"`
}