package gowebdav

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	O_RDWR = iota
	O_READONLY
)

type Server interface {
	NewClient(pathPrefix, filePath string) Client
	Run(addr ...string)
}

type SingleServer interface {
	Client
	Run(addr ...string)
}

type webdavSingleServer struct {
	webdavServer
	Client
}

type webdavServer struct {
	ginengine *gin.Engine
}

// Cannot be used concurrently with NewSingleWebdav()
func NewWebdav() Server {
	webdavserver := webdavServer{}
	gin.SetMode(gin.ReleaseMode)

	webdavserver.ginengine = gin.Default()

	return &webdavserver
}

// Cannot be used concurrently with NewWebdav()
func NewSingleWebdav(filePath string) SingleServer {
	webdavserver := webdavSingleServer{}
	gin.SetMode(gin.ReleaseMode)

	webdavserver.ginengine = gin.Default()
	webdavserver.Client = webdavserver.newClient(filePath)

	return &webdavserver
}

func (webdavServer *webdavServer) Run(addr ...string) {
	fmt.Printf("Webdav run on port%s\n", resolveAddress(addr))
	webdavServer.ginengine.Run(addr...)
}
