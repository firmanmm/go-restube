package app

import (
	"github.com/firmanmm/go-restube/app/controller"
	"github.com/firmanmm/go-restube/app/middleware"
	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

func NewRestube() *gin.Engine {
	engine := gin.Default()
	fileStorageService := service.NewFileStorageService("storage")
	inMemoryAuthStorageService := service.NewInMemoryAuthentication()
	authService := service.NewAuthenticationService(inMemoryAuthStorageService)
	authController := controller.NewAuthenticationController(authService)
	downloaderService := service.NewDownloaderService(fileStorageService)
	downloaderController := controller.NewDownloaderController(downloaderService, authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	authGroup := engine.Group("/video", authMiddleware.Handle)
	authGroup.GET("/", downloaderController.HandleGetVideoInfo)
	authGroup.POST("/", downloaderController.HandleRequest)
	authGroup.GET("/:url", downloaderController.HandleDownload)
	engine.POST("/register", authController.HandleRegister)
	engine.POST("/login", authController.HandleAuthentication)
	engine.POST("/usage", authMiddleware.Handle, authController.HandleGetByteDownloaded)
	return engine
}
