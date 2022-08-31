package gowebdav

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

type Client interface {
	SetAuth(string, string) Client
}

type client struct {
	userInfo   map[string]string
	fs         *webdav.Handler
	pathPrefix string
	mode       int
	engine     *gin.Engine
}

func (server *webdavServer) NewClient(pathPrefix, filePath string) Client {
	fs := &webdav.Handler{
		FileSystem: webdav.Dir(filePath),
		LockSystem: webdav.NewMemLS(),
	}
	client := client{pathPrefix: pathPrefix, fs: fs, engine: server.ginengine, userInfo: make(map[string]string), mode: O_RDWR}
	server.ginengine.Any(fmt.Sprintf("%s/*webdav", pathPrefix), client.handleWebdav())
	return &client
}

func (client *client) handleWebdav() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(client.userInfo) == 0 {
			return
		}
		user, pwd, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		if v, ok := client.userInfo[user]; !ok || v != pwd {
			ctx.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
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
		ctx.Request.URL.Path = ctx.Params.ByName("webdav")
		if ctx.Request.Method == "GET" && handleDirList(client.fs.FileSystem, ctx) {
			return
		}
		client.fs.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func (client *client) SetAuth(username, password string) Client {
	client.userInfo[username] = password
	return client
}

func (server *webdavSingleServer) newClient(filePath string) Client {
	fs := &webdav.Handler{
		FileSystem: webdav.Dir(filePath),
		LockSystem: webdav.NewMemLS(),
	}
	client := client{fs: fs, engine: server.ginengine, userInfo: make(map[string]string), mode: O_RDWR}
	server.ginengine.Any("/*webdav", client.handleWebdav())
	return &client
}
