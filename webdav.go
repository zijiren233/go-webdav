package gowebdav

import (
	"fmt"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

type Server interface {
	NewClient(pathPrefix, filePath string) Client
	NewClientWithMemFS(pathPrefix string) Client
	Run(addr ...string) error
	RunTLS(addr string, certFile string, keyFile string) error
	RunAUTOTLS(domain ...string) error
	SSLRedirect(SSLHost string)
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

	webdavserver.ginengine.Use(Logger(), gin.Recovery())

	return &webdavserver
}

func NewWebdavWithGin(engine *gin.Engine) Server {
	webdavserver := webdavServer{}

	webdavserver.ginengine = engine

	return &webdavserver
}

func (webdavServer *webdavServer) Run(addr ...string) error {
	fmt.Printf("Webdav http run on port%s\n", resolveAddress(addr))
	return webdavServer.ginengine.Run(addr...)
}

func (webdavServer *webdavServer) RunTLS(addr string, certFile string, keyFile string) error {
	fmt.Printf("Webdav https run on port%s\n", addr)
	return webdavServer.ginengine.RunTLS(addr, certFile, keyFile)
}

// acme
func (webdavServer *webdavServer) RunAUTOTLS(domain ...string) error {
	return autotls.Run(webdavServer.ginengine, domain...)
}

func (webdavServer *webdavServer) SSLRedirect(SSLHost string) {
	webdavServer.ginengine.Use(func(ctx *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     SSLHost,
		})
		err := middleware.Process(ctx.Writer, ctx.Request)
		if err != nil {
			ctx.Abort()
			return
		}
	})
}

func (webdavServer *webdavServer) GetGinEngine() *gin.Engine {
	return webdavServer.ginengine
}
