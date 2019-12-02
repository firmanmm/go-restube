package app

import (
	"github.com/firmanmm/go-restube/app/controller"
	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

func NewRestube() *gin.Engine {
	engine := gin.Default()
	fileStorageService := service.NewFileStorageService("storage")
	downloaderService := service.NewDownloaderService(fileStorageService)
	downloaderController := controller.NewDownloaderController(downloaderService)
	engine.GET("/info", downloaderController.HandleGetVideoInfo)
	engine.POST("/video", downloaderController.HandleRequest)
	engine.GET("/video/:url", downloaderController.HandleDownload)
	return engine
}
