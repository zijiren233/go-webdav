package gowebdav

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

type Client interface {
	AddUser(string, string, int) Client
	ChangeUserMode(string, int) Client
	ChangeUserPwd(string, string) Client
	SetUserRights(string, string, int) Client
}

type client struct {
	userInfo   map[string]*user
	fs         *webdav.Handler
	pathPrefix string
	mode       int
	engine     *gin.Engine
}

// All client path prefix levels must match
func (server *webdavServer) NewClient(pathPrefix, filePath string) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.Dir(filePath),
		LockSystem: webdav.NewMemLS(),
	}
	client := client{pathPrefix: pathPrefix, fs: fs, engine: server.ginengine, userInfo: make(map[string]*user), mode: O_RDWR}
	server.ginengine.Any(fmt.Sprintf("%s/*webdav", pathPrefix), client.handleWebdav())
	return &client
}

// All client path prefix levels must match
func (server *webdavServer) NewClientWithMemFS(pathPrefix string) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.NewMemFS(),
		LockSystem: webdav.NewMemLS(),
	}
	client := client{pathPrefix: pathPrefix, fs: fs, engine: server.ginengine, userInfo: make(map[string]*user), mode: O_RDWR}
	server.ginengine.Any(fmt.Sprintf("%s/*webdav", pathPrefix), client.handleWebdav())
	return &client
}

func (client *client) handleWebdav() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(client.userInfo) != 0 {
			user, pwd, ok := ctx.Request.BasicAuth()
			if !ok {
				ctx.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				ctx.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			if v, ok := client.userInfo[user]; !ok || v.password != pwd {
				ctx.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				ctx.Writer.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		switch client.mode {
		case O_RDWR:
		case O_READONLY:
			switch ctx.Request.Method {
			case "GET", "HEAD", "POST":
			default:
				return
			}
		default:
			return
		}
		if ctx.Request.Method == "GET" && client.handleDirList(client.fs.FileSystem, ctx) {
			return
		}
		client.fs.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
