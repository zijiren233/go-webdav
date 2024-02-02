package gowebdav

import (
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

	SetCORS(config CORSConfig)

	FS() *webdav.Handler

	usersfunc
}

type Cli struct {
	readOnly bool
	usersfunc
	fs         *webdav.Handler
	group      *gin.RouterGroup
	pathPrefix string
}

// All client path prefix levels must match
func (server *webdavServer) DefaultClient(pathPrefix, filePath string) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.Dir(filePath),
		LockSystem: webdav.NewMemLS(),
	}
	client := Cli{pathPrefix: pathPrefix, fs: fs, usersfunc: newusers()}
	client.addMethod(server.ginengine, pathPrefix, Defaulthandle(&client))
	return &client
}

// Custom handler, All client path prefix levels must match
func (server *webdavServer) Client(pathPrefix, filePath string, handlerFunc HandlerFunc) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.Dir(filePath),
		LockSystem: webdav.NewMemLS(),
	}
	client := Cli{pathPrefix: pathPrefix, fs: fs, usersfunc: newusers()}
	client.addMethod(server.ginengine, pathPrefix, handlerFunc(&client))
	return &client
}

// All client path prefix levels must match
func (server *webdavServer) DefaultClientWithMemFS(pathPrefix string) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.NewMemFS(),
		LockSystem: webdav.NewMemLS(),
	}
	client := Cli{pathPrefix: pathPrefix, fs: fs, usersfunc: newusers()}
	client.addMethod(server.ginengine, pathPrefix, Defaulthandle(&client))
	return &client
}

// Custom handler, All client path prefix levels must match
func (server *webdavServer) ClientWithMemFS(pathPrefix string, handlerFunc HandlerFunc) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.NewMemFS(),
		LockSystem: webdav.NewMemLS(),
	}
	client := Cli{pathPrefix: pathPrefix, fs: fs, usersfunc: newusers()}
	client.addMethod(server.ginengine, pathPrefix, handlerFunc(&client))
	return &client
}

func (client *Cli) addMethod(ginengine *gin.Engine, pathPrefix string, handlerFunc gin.HandlerFunc) {
	group := ginengine.Group(pathPrefix)
	group.Use(client.webdavauth())
	group.Any("/*webdav", handlerFunc)
	for _, v := range missingMethods {
		group.Handle(v, "/*webdav", handlerFunc)
	}
	client.group = group
}

func (client *Cli) webdavauth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if client.UserNum() != 0 {
			username, pwd, ok := ctx.Request.BasicAuth()
			if !ok {
				authErr(ctx)
				return
			}
			user, ok := client.FindUser(username)
			if !ok || !user.ComparePassword(pwd) {
				authErr(ctx)
				return
			}
			if asd(ctx.Request.Method, user.Mode()) {
				methodNotAllowed(ctx)
				return
			}
		} else if client.readOnly && asd(ctx.Request.Method, O_RDONLY) {
			methodNotAllowed(ctx)
			return
		}
	}
}

// It only takes effect if no user is set
func (client *Cli) GlobalReadOnly() {
	if !client.readOnly {
		client.readOnly = true
	}
}

// It only takes effect if no user is set
func (client *Cli) UnSetReadOnly() {
	if client.readOnly {
		client.readOnly = false
	}
}

func (client *Cli) SetCORS(config CORSConfig) {
	client.group.Use(cors(config))
}

func (client *Cli) FS() *webdav.Handler {
	return client.fs
}

type HandlerFunc func(*Cli) gin.HandlerFunc

func Defaulthandle(client *Cli) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "GET" && client.handleDirList(client.fs.FileSystem, ctx) {
			return
		}
		client.fs.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
