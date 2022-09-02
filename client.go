package gowebdav

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

var (
	missingMethods = []string{
		"PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK",
	}
)

type Client interface {
	GlobalReadOnly()
	UnSetReadOnly()

	SetCORS(Allow_Origin, Max_Age, Allow_Methods, Allow_Headers, Allow_Credentials string)

	FS() webdav.FileSystem
	LS() webdav.LockSystem

	usersfunc
}

type client struct {
	readOnly bool
	usersfunc
	fs         *webdav.Handler
	group      *gin.RouterGroup
	pathPrefix string
}

// All client path prefix levels must match
func (server *webdavServer) NewClient(pathPrefix, filePath string) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.Dir(filePath),
		LockSystem: webdav.NewMemLS(),
	}
	client := client{pathPrefix: pathPrefix, fs: fs, usersfunc: newusers()}
	client.addMethod(server.ginengine, pathPrefix)
	return &client
}

// All client path prefix levels must match
func (server *webdavServer) NewClientWithMemFS(pathPrefix string) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.NewMemFS(),
		LockSystem: webdav.NewMemLS(),
	}
	client := client{pathPrefix: pathPrefix, fs: fs, usersfunc: newusers()}
	client.addMethod(server.ginengine, pathPrefix)
	return &client
}

func (client *client) addMethod(ginengine *gin.Engine, pathPrefix string) {
	group := ginengine.Group(pathPrefix)
	group.Use(corsDefaule(), client.webdavauth())
	group.Any("/*webdav", client.handleWebdav())
	for _, v := range missingMethods {
		group.Handle(v, "/*webdav", client.handleWebdav())
	}
	client.group = group
}

func (client *client) webdavauth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if client.UserNum() != 0 {
			username, pwd, ok := ctx.Request.BasicAuth()
			if !ok {
				authErr(ctx)
				return
			}
			user, ok := client.FindUser(username)
			if !ok || !user.comparePassword(pwd) {
				authErr(ctx)
				return
			}
			if user.Mode() == O_READONLY && readonle(ctx.Request.Method) {
				ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
				ctx.Abort()
				return
			}
		} else if client.readOnly && readonle(ctx.Request.Method) {
			ctx.Writer.WriteHeader(http.StatusMethodNotAllowed)
			ctx.Abort()
			return
		}
	}
}

// It only takes effect if no user is set
func (client *client) GlobalReadOnly() {
	if !client.readOnly {
		client.readOnly = true
	}
}

// It only takes effect if no user is set
func (client *client) UnSetReadOnly() {
	if client.readOnly {
		client.readOnly = false
	}
}

func (client *client) SetCORS(Allow_Origin, Max_Age, Allow_Methods, Allow_Headers, Allow_Credentials string) {
	client.group.Use(cors(Allow_Origin, Max_Age, Allow_Methods, Allow_Headers, Allow_Credentials))
}

func (client *client) FS() webdav.FileSystem {
	return client.fs.FileSystem
}

func (client *client) LS() webdav.LockSystem {
	return client.fs.LockSystem
}

func (client *client) handleWebdav() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "GET" && client.handleDirList(client.fs.FileSystem, ctx) {
			return
		}
		client.fs.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
