package gowebdav

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server interface {
	NewClient(pathPrefix, filePath string) Client
	NewClientWithMemFS(pathPrefix string) Client
	Run(addr ...string)
	GetGinEngine() *gin.Engine
}

type webdavServer struct {
	ginengine *gin.Engine
}

// All client path prefix levels must match
func NewWebdav() Server {
	webdavserver := webdavServer{}
	gin.SetMode(gin.ReleaseMode)

	webdavserver.ginengine = gin.Default()

	return &webdavserver
}

func NewWebdavWithGin(engine *gin.Engine) Server {
	webdavserver := webdavServer{}
	gin.SetMode(gin.ReleaseMode)

	webdavserver.ginengine = engine

	return &webdavserver
}

func (webdavServer *webdavServer) Run(addr ...string) {
	fmt.Printf("Webdav run on port%s\n", resolveAddress(addr))
	webdavServer.ginengine.Run(addr...)
}

func (webdavServer *webdavServer) GetGinEngine() *gin.Engine {
	return webdavServer.ginengine
}
