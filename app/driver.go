package app

import (
	"crypto/tls"
	"flag"
	"log"

	"github.com/go-redis/redis"

	"github.com/firmanmm/go-restube/app/controller"
	"github.com/firmanmm/go-restube/app/middleware"
	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

func NewRestube() *gin.Engine {
	engine := gin.Default()

	var azureStorageUsername string
	var azureStorageKey string

	flag.StringVar(&azureStorageUsername, "azSU", "", "Azure Storage Username for remote file storage")
	flag.StringVar(&azureStorageKey, "azSK", "", "Azure Storage Key for remote file storage")

	var azureRedisKey string
	var azureRedisAddress string

	flag.StringVar(&azureRedisAddress, "azRA", "", "Azure Redis Address for session storage")
	flag.StringVar(&azureRedisKey, "azRK", "", "Azure Redis Key for session storage")

	flag.Parse()
	var fileStorageService service.IFileStorage
	var authStorageService service.IAuthStorage
	if len(azureStorageUsername) > 0 && len(azureStorageKey) > 0 {
		fileStorageService = service.NewFileBlobStorageService("storage", azureStorageUsername, azureStorageKey)
		log.Println("Running Storage Using Azure Blob Storage")
	} else {
		fileStorageService = service.NewFileStorageService("storage")
		log.Println("Running Storage Using Local File Storage")
	}
	if len(azureRedisKey) > 0 && len(azureRedisAddress) > 0 {
		authStorageService = service.NewRedisAuthentication(&redis.Options{
			Addr:      azureRedisAddress,
			Password:  azureRedisKey,
			TLSConfig: &tls.Config{},
		})
		log.Println("Running Authentication Using Azure Redis Cache")
	} else {
		authStorageService = service.NewInMemoryAuthentication()
		log.Println("Running Authentication Using InMemory Storage")
	}
	authService := service.NewAuthenticationService(authStorageService)
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
	engine.GET("/usage", authMiddleware.Handle, authController.HandleGetByteDownloaded)
	engine.GET("/usage/all", authMiddleware.Handle, authController.HandleListUsage)
	return engine
}
