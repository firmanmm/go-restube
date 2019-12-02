package app

import (
	"github.com/firmanmm/go-restube/app/controller"
	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

func NewRestube() *gin.Engine {
	engine := gin.Default()
	fileStorageService := service.NewFileStorageService("/storage")
	downloaderService := service.NewDownloaderService(fileStorageService)
	downloaderController := controller.NewDownloaderController(downloaderService)
	engine.GET("/video/info", downloaderController.HandleGetVideoInfo)
	return engine
}
