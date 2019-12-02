package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/firmanmm/go-restube/app/service"
	"github.com/gin-gonic/gin"
)

type DownloaderController struct {
	downloader  *service.DownloaderService
	authService *service.AuthenticationService
}

func (d *DownloaderController) HandleGetVideoInfo(ctx *gin.Context) {
	url := ctx.Query("url")
	if len(url) < 10 {
		ctx.String(http.StatusBadRequest, "Valid \"url\" query is required")
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
	url := ctx.PostForm("url")
	if len(url) < 10 {
		ctx.String(http.StatusBadRequest, "Valid \"url\" post parameter is required")
		return
	}
	modeS := ctx.PostForm("mode")
	mode := 0
	if len(modeS) <= 0 {
		ctx.String(http.StatusBadRequest, "Valid \"mode\" post parameter is required")
		return
	} else {
		mode64, err := strconv.ParseInt(modeS, 10, 0)
		if err != nil {
			ctx.String(http.StatusBadRequest, "Failed to parse mode")
			return
		}
		mode = int(mode64)
	}
	res, err := d.downloader.Request(url, mode)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Please use \"url\" payload to access your requested video at /video/{url} and replace {url} to given url",
		"url":     res,
	})
}

func (d *DownloaderController) HandleDownload(ctx *gin.Context) {
	url := ctx.Param("url")
	if len(url) < 10 {
		ctx.String(http.StatusBadRequest, "Valid \"url\" parameter is required")
		return
	}
	data, err := d.downloader.Download(url)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	sessionID := ctx.GetHeader("Authorization")
	err = d.authService.AddByteDownloaded(sessionID, uint(len(data)))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	contentType := http.DetectContentType(data)
	ctx.Data(http.StatusOK, contentType, data)
}

func NewDownloaderController(downloader *service.DownloaderService, authService *service.AuthenticationService) *DownloaderController {
	instance := new(DownloaderController)
	instance.downloader = downloader
	instance.authService = authService
	return instance
}
