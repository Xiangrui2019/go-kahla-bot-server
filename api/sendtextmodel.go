package api

type SendTextRequestModel struct {
	Token string `form:"token" binding:"Required"`
	Content string `form:"content" binding:"Required"`
}

type SendTextResponseModel struct {
	Code int `json:"code"`
	Message string `json:"message"`
}