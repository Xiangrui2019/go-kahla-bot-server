package middlewares

import (
	"gopkg.in/macaron.v1"
)

func CorsMiddleware(origin string) macaron.Handler {
	return func(context *macaron.Context) {
		context.Resp.Header().Set("Access-Control-Allow-Origin", origin)
		context.Resp.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Resp.Header().Set("Access-Control-Allow-Headers", "*")
		context.Resp.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if context.Req.Method == "OPTIONS" {
			context.Status(204)
			return
		}

		context.Next()
	}
}
