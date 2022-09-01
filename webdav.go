package gowebdav

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server interface {
	NewClient(pathPrefix, filePath string) Client
	NewClientWithMemFS(pathPrefix string) Client
	Run(addr ...string) error
	GetGinEngine() *gin.Engine
}

type webdavServer struct {
	ginengine *gin.Engine
}

// All client path prefix levels must match
func NewWebdav() Server {
	webdavserver := webdavServer{}
	gin.SetMode(gin.ReleaseMode)

	webdavserver.ginengine = gin.New()

	webdavserver.ginengine.Use(Formatter(), gin.Recovery())

	return &webdavserver
}

func NewWebdavWithGin(engine *gin.Engine) Server {
	webdavserver := webdavServer{}

	webdavserver.ginengine = engine

	return &webdavserver
}

func (webdavServer *webdavServer) Run(addr ...string) error {
	fmt.Printf("Webdav run on port%s\n", resolveAddress(addr))
	return webdavServer.ginengine.Run(addr...)
}

func (webdavServer *webdavServer) GetGinEngine() *gin.Engine {
	return webdavServer.ginengine
}
