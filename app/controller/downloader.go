package controller

import (
	"fmt"
	"net/http"

	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

type DownloaderController struct {
	downloader *service.DownloaderService
}

func (d *DownloaderController) HandleGetVideoInfo(ctx *gin.Context) {
	url := ctx.Query("url")
	if len(url) < 10 {
		ctx.String(http.StatusBadRequest, "Valid \"url\" parameter is required")
		return
	}
	res, err := d.downloader.GetVideoQuality(url)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	for key, val := range res {
		res[key] = fmt.Sprintf("[%d] %s", key, val)
	}
	ctx.JSON(http.StatusOK, res)
}

func (d *DownloaderController) HandleRequest(ctx *gin.Context) {
	//TODO : Make
}

func (d *DownloaderController) HandleDownload(ctx *gin.Context) {
	//TODO : Make
}

func NewDownloaderController(downloader *service.DownloaderService) *DownloaderController {
	instance := new(DownloaderController)
	instance.downloader = downloader
	return instance
}
