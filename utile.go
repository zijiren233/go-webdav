package gowebdav

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func getsize(size int64) string {
	tmp := float64(size)
	if tmp < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if tmp = tmp / 1024; tmp < 1024 {
		return fmt.Sprintf("%.2f KB", tmp)
	} else if tmp = tmp / 1024; tmp < 1024 {
		return fmt.Sprintf("%.2f MB", tmp)
	} else if tmp = tmp / 1024; tmp < 1024 {
		return fmt.Sprintf("%.2f GB", tmp)
	} else {
		return fmt.Sprintf("%.2f TB", tmp/1024)
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func readonle(Method string) bool {
	switch Method {
	case "GET", "HEAD", "PROPFIND":
		return false
	default:
		return true
	}
}

func path2index(path string) string {
	s := strings.Split(path, "/")
	var tmp string
	for k, v := range s[1 : len(s)-1] {
		tmp += fmt.Sprintf("/<a href = \"%s\">%s</a>", strings.Repeat("../", len(s)-3-k), v)
	}
	return tmp
}

func authErr(ctx *gin.Context) {
	ctx.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	ctx.Writer.WriteHeader(http.StatusUnauthorized)
	ctx.Abort()
}
