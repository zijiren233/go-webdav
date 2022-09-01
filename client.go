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
	usersfunc
}

type client struct {
	readOnly bool
	usersfunc
	fs         *webdav.Handler
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
	client := client{pathPrefix: pathPrefix, fs: fs, usersfunc: newusers()}
	server.ginengine.Any(fmt.Sprintf("%s/*webdav", pathPrefix), client.handleWebdav())
	return &client
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

func (client *client) FS() webdav.FileSystem {
	return client.fs.FileSystem
}

func (client *client) LS() webdav.LockSystem {
	return client.fs.LockSystem
}

func (client *client) handleWebdav() gin.HandlerFunc {
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
				return
			}
		} else if client.readOnly && readonle(ctx.Request.Method) {
			return
		}
		if ctx.Request.Method == "GET" && client.handleDirList(client.fs.FileSystem, ctx) {
			return
		}
		client.fs.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
