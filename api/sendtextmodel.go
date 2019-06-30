package api

// 发送文本输入模型
type SendTextRequestModel struct {
	Token string `form:"token" binding:"Required"`
	Content string `form:"content" binding:"Required"`
}

// 发送文本输出模型
type SendTextResponseModel struct {
	Code int `json:"code"`
	Message string `json:"message"`
}