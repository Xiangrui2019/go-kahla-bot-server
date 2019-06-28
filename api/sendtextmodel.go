package api

type SendTextRequestModel struct {
	Token string `form:"token" binding:"Required"`
	Content string `form:"content" binding:"Required"`
}