package gowebdav

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	Allow_Origin, Max_Age, Allow_Methods, Allow_Headers, Allow_Credentials string
}

var CORSDefault = CORSConfig{
	Allow_Origin:      "*",
	Max_Age:           "86400",
	Allow_Methods:     "*",
	Allow_Headers:     "*",
	Allow_Credentials: "true",
}

func cors(config CORSConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", config.Allow_Origin)
		ctx.Writer.Header().Set("Access-Control-Max-Age", config.Max_Age)
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", config.Allow_Methods)
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", config.Allow_Headers)
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", config.Allow_Credentials)

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}
