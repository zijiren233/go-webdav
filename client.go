package gowebdav

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

type Client interface {
	// Global permissions
	GlobalReadOnly()
	UnSetReadOnly()

	// System
	FS() webdav.FileSystem
	LS() webdav.LockSystem

	// User
	AddUser(username, password string, mode int) Client
	ChangeUserMode(username string, mode int) Client
	ChangeUserPwd(username, password string) Client
	SetUserRights(username, password string, mode int) Client
}

type client struct {
	readOnly   bool
	userInfo   map[string]*user
	fs         *webdav.Handler
	pathPrefix string
	engine     *gin.Engine
}

// All client path prefix levels must match
func (server *webdavServer) NewClient(pathPrefix, filePath string) Client {
	fs := &webdav.Handler{
		Prefix:     pathPrefix,
		FileSystem: webdav.Dir(filePath),
		LockSystem: webdav.NewMemLS(),
	}
	client := client{pathPrefix: pathPrefix, fs: fs, engine: server.ginengine, userInfo: make(map[string]*user)}
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
	client := client{pathPrefix: pathPrefix, fs: fs, engine: server.ginengine, userInfo: make(map[string]*user)}
	server.ginengine.Any(fmt.Sprintf("%s/*webdav", pathPrefix), client.handleWebdav())
	return &client
}

// Users are authenticated individually without using global permissions
func (client *client) GlobalReadOnly() {
	if !client.readOnly {
		client.readOnly = true
	}
}

// Users are authenticated individually without using global permissions
func (client *client) UnSetReadOnly() {
	if client.readOnly {
		client.readOnly = false
	}
}

func (client *client) FS() webdav.FileSystem {
	return client.fs.FileSystem
}

func (client *client) LS() webdav.LockSystem {
	return client.fs.LockSystem
}

func (client *client) handleWebdav() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		client.userAuth(ctx)
		if ctx.Request.Method == "GET" && client.handleDirList(client.fs.FileSystem, ctx) {
			return
		}
		client.fs.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
