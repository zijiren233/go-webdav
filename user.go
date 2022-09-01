package gowebdav

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	name, password string
	mode           int
}

const (
	O_RDWR = iota
	O_READONLY
)

func (client *client) AddUser(username, password string, mode int) Client {
	client.userInfo[username] = &user{name: username, password: password, mode: mode}
	return client
}

func (client *client) ChangeUserMode(username string, mode int) Client {
	if v, ok := client.userInfo[username]; ok {
		v.mode = mode
	}
	return client
}

func (client *client) ChangeUserPwd(username, password string) Client {
	if v, ok := client.userInfo[username]; ok {
		v.password = password
	}
	return client
}

func (client *client) SetUserRights(username, password string, mode int) Client {
	if v, ok := client.userInfo[username]; ok {
		v.password = password
		v.mode = mode
	}
	return client
}

func (client *client) userAuth(ctx *gin.Context) {
	if len(client.userInfo) != 0 {
		user, pwd, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			ctx.Abort()
		}
		v, ok := client.userInfo[user]
		if !ok || v.password != pwd {
			ctx.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			ctx.Abort()
		}
		v.userAuthentication(ctx)
		return
	}
	if client.readOnly {
		readonle(ctx)
	}
}

func (user *user) userAuthentication(ctx *gin.Context) {
	switch user.mode {
	case O_RDWR:
	case O_READONLY:
		readonle(ctx)
	default:
		ctx.Abort()
	}
}
