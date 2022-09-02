package gowebdav

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func cors(Allow_Origin, Max_Age, Allow_Methods, Allow_Headers, Allow_Credentials string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", Allow_Origin)
		ctx.Writer.Header().Set("Access-Control-Max-Age", Max_Age)
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", Allow_Methods)
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", Allow_Headers)
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", Allow_Credentials)

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}

func corsDefaule() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}
