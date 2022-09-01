package gowebdav

import (
	"fmt"
	"sync"

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
	userfunc
}

type client struct {
	readOnly bool
	*users
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
	client := client{pathPrefix: pathPrefix, fs: fs, engine: server.ginengine, users: &users{usermap: make(map[string]*user), lock: &sync.RWMutex{}}}
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
	client := client{pathPrefix: pathPrefix, fs: fs, engine: server.ginengine, users: &users{usermap: make(map[string]*user), lock: &sync.RWMutex{}}}
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
		if len(client.users.usermap) != 0 {
			username, pwd, ok := ctx.Request.BasicAuth()
			if !ok {
				authErr(ctx)
				return
			}
			client.users.lock.RLock()
			v, ok := client.usermap[username]
			client.users.lock.RUnlock()
			if !ok {
				authErr(ctx)
				return
			} else {
				v.lock.RLock()
				if v.password != pwd {
					v.lock.RUnlock()
					authErr(ctx)
					return
				}
				v.lock.RUnlock()
			}
			v.lock.RLock()
			if v.mode == O_READONLY && readonle(ctx.Request.Method) {
				v.lock.RUnlock()
				return
			}
			v.lock.RUnlock()
		} else if client.readOnly && readonle(ctx.Request.Method) {
			return
		}
		if ctx.Request.Method == "GET" && client.handleDirList(client.fs.FileSystem, ctx) {
			return
		}
		client.fs.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
