package gowebdav

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

// Perfection of gin logger middleware
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(handlelogger)
}

var handlelogger = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = methodcolor(&param)
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func methodcolor(p *gin.LogFormatterParams) string {
	switch p.Method {
	case http.MethodGet, "PROPFIND":
		return blue
	case http.MethodPost, "PROPPATCH":
		return cyan
	case http.MethodPut, "MKCOL":
		return yellow
	case http.MethodDelete, "LOCK":
		return red
	case http.MethodPatch, "COPY":
		return green
	case http.MethodHead, "MOVE":
		return magenta
	case http.MethodOptions, "UNLOCK":
		return white
	default:
		return reset
	}
}
